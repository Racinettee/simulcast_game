package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/Racinettee/asefile"
	"github.com/jessevdk/go-flags"
)

type outFileTemplate struct {
	FileName string
	// For each frame (key) list of points (value)
	Points map[int][]int
}

type outputTemplate struct {
	// The package to insert into
	Package string
	// The list of files with respective frames -> shapes
	Files []outFileTemplate
}

var outputTemplateStr string = `package {{.Package}}
{{ range .Files }}
var {{.FileName}}Shapes = map[int][]int{ {{ range $frame, $points := .Points }}
	{{ $frame }}: { {{ range $index, $elem := $points }}{{ if $index }},{{ end }} {{ $elem }}{{ end }} },{{ end }}
}
{{ end }}`

type Pointi struct {
	X, Y int
}
type ByPointX []Pointi

func (a ByPointX) Len() int           { return len(a) }
func (a ByPointX) Less(i, j int) bool { return a[i].X < a[j].X }
func (a ByPointX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func main() {
	var opts struct {
		InFiles string `short:"i"`
		OutFile string `short:"o"`
		Package string `short:"p"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}
	mainResult := outputTemplate{
		Package: opts.Package,
		Files:   make([]outFileTemplate, 0),
	}

	inFiles := strings.Split(opts.InFiles, ",")

	// For each file in inFiles we parse the ase file looking for a collision layer
	for _, file := range inFiles {
		outResult := outFileTemplate{
			FileName: transformTitle(file),
			Points:   make(map[int][]int),
		}
		var aseFile asefile.AsepriteFile
		err = aseFile.DecodeFile(file)
		if err != nil {
			panic(err)
		}
		if len(aseFile.Frames) == 0 {
			println("No frames in file")
			continue
		}
		collisLayer := -1
		for x, layer := range aseFile.Frames[0].Layers {
			if layer.LayerName == "<Collision>" {
				collisLayer = x
			}
		}
		if collisLayer == -1 {
			println("No collision layers - use <Collision> to define collision layers")
			continue
		}
		// For each frame we want to scan the pixels and find all the ones of a certain color
		for framei, frame := range aseFile.Frames {

			cel := frame.Cels[collisLayer]
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
		mainResult.Files = append(mainResult.Files, outResult)
	}

	t, err := template.New("shapes").Parse(outputTemplateStr)

	if err != nil {
		panic(err)
	}

	outFile, err := os.Create(opts.OutFile)

	if err != nil {
		panic(err)
	}

	t.Execute(outFile, &mainResult)
}

func transformTitle(title string) string {
	str := strings.Title(strings.TrimSuffix(filepath.Base(title), filepath.Ext(title)))
	return strings.Map(func(r rune) rune {
		switch {
		case unicode.IsDigit(r), unicode.IsLetter(r):
			return r
		}
		return -1
	}, str)
}
