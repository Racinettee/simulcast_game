package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/Racinettee/asefile"
	"github.com/Racinettee/generics"
	"github.com/jessevdk/go-flags"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type outputShape struct {
	FileName string
	// For each frame (key) list of points (value)
	Points map[int][]int
}

type outputTemplate struct {
	// The package to insert into
	Package string
	// The list of files with respective frames -> shapes
	Files []outputShape
}

var outputTemplateStr string = `package {{.Package}}
// This code is generated, do not edit
{{ range .Files }}
var {{.FileName}}Shape = map[int][]int{ {{ range $frame, $points := .Points }}
	{{ $frame }}: { {{ range $index, $elem := $points }}{{ if $index }},{{ end }} {{ $elem }}{{ end }} },{{ end }}
}
{{ end }}`

type Pointi struct{ X, Y int }
type ByPointX []Pointi

func (a ByPointX) Len() int           { return len(a) }
func (a ByPointX) Less(i, j int) bool { return a[i].X < a[j].X }
func (a ByPointX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type options struct {
	Config  bool   `short:"c" description:"Use config file instead of the other flags"`
	InFiles string `short:"i" description:"List of semi-colon seperated file paths"`
	OutFile string `short:"o" description:"The name of the output file to create"`
	Package string `short:"p" description:"The package name to write into the out file"`
}

func main() {
	var opts options
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}
	inFiles, outFile, pkg := config(opts)
	mainResult := outputTemplate{
		Package: pkg,
		Files:   make([]outputShape, 0),
	}
	// For each file in inFiles we parse the ase file looking for a collision layer
	for _, file := range inFiles {
		mainResult.Files = append(mainResult.Files, collectCollisionData(file)...)
	}

	t, err := template.New("shapes").Parse(outputTemplateStr)
	if err != nil {
		panic(err)
	}
	of, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	t.Execute(of, &mainResult)
}

func transformTitle(title string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			return r
		}
		return -1
	}, cases.Title(language.English).String(strings.TrimSuffix(filepath.Base(title), filepath.Ext(title))))
}

func collectCollisionData(file string) []outputShape {
	results := generics.NewList[outputShape](0)
	// Load the ase file and quickly parse it for any layers that match our naming scheme
	var aseFile asefile.AsepriteFile
	if err := aseFile.DecodeFile(file); err != nil {
		panic(err)
	}
	if len(aseFile.Frames) == 0 {
		log.Printf("No frames in %v\n", file)
		return results
	}
	// Figure out which layers we care about before parsing the frames
	collisLayers := generics.NewList[generics.Pair[int, string]](0)
	for x, layer := range aseFile.Frames[0].Layers {
		if layer.LayerName == "<HurtBox>" {
			collisLayers.Push(generics.NewPair(x, "_HurtBox"))
		} else if layer.LayerName == "<HitBox>" {
			collisLayers.Push(generics.NewPair(x, "_HitBox"))
		}
	}
	// If there was no matching layers then return from this function to move on
	if collisLayers.Len() == 0 {
		log.Printf("no collision layers in %v - use <HitBox> or <HurtBox> to define collision layers\n", file)
		return results
	}
	// For each collision layer we need to scan the frames and collect all the points
	for _, layer := range collisLayers {
		results.Push(parseFrame(aseFile, file, layer))
	}
	return results
}

func parseFrame(aseFile asefile.AsepriteFile, file string, layer generics.Pair[int, string]) outputShape {
	outResult := outputShape{
		FileName: transformTitle(file) + layer.Right,
		Points:   make(map[int][]int),
	}
	// For each frame we want to scan the pixels and find all the ones of a certain color
	for framei, frame := range aseFile.Frames {
		cel := frame.Cels[layer.Left]
		dat := cel.RawCelData
		if len(dat) == 0 {
			continue
		}
		points := make([]Pointi, 0)
		// For now, We'll just go with 255,0,255,255 as our color
		w, h := cel.WidthInPix, cel.HeightInPix
		offset := 0
		for y := 0; y < int(h); y += 1 {
			for x := 0; x < int(w); x, offset = x+1, offset+4 {
				r, g, b, a := dat[offset], dat[offset+1], dat[offset+2], dat[offset+3]
				if r == 255 && g == 0 && b == 255 && a == 255 {
					points = append(points, Pointi{x, y})
				}
			}
		}
		// Sort the points first by X order
		sort.Sort(ByPointX(points))
		// Sort the highest and lowest points by the middle of the median y value
		highs := make([]Pointi, 0)
		lows := make([]Pointi, 0)

		middle := int(h / 2)
		for _, point := range points {
			if point.Y <= middle {
				highs = append(highs, point)
			} else {
				lows = append(lows, point)
			}
		}
		result := make([]int, 0)
		// With all high values (points near the top have lower y value)
		// iterate them forward and generate the coordinates
		for _, point := range highs {
			// Add the point - but including the offset of the cel for the final data
			result = append(result, point.X+int(cel.X), point.Y+int(cel.Y))
		}
		for i := len(lows) - 1; i >= 0; i-- {
			result = append(result, lows[i].X+int(cel.X), lows[i].Y+int(cel.Y))
		}
		outResult.Points[framei] = result
	}
	return outResult
}

type assetConfig struct {
	Package   string
	AssetPath string
	OutFile   string
	Files     []string
}

func config(opts options) (inFiles []string, outFile, pkg string) {
	if opts.Config {
		var assetConf assetConfig
		toml.DecodeFile("asset-conf.toml", &assetConf)
		for x := range assetConf.Files {
			assetConf.Files[x] = assetConf.AssetPath + "/" + assetConf.Files[x]
		}
		inFiles = assetConf.Files
		outFile = assetConf.OutFile
		pkg = assetConf.Package
	} else {
		inFiles = strings.Split(opts.InFiles, ";")
		outFile = opts.OutFile
		pkg = opts.Package
	}
	return
}
