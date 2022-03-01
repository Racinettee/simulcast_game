package main

import (
	"github.com/Racinettee/asefile"
	"github.com/jessevdk/go-flags"
)

type Pointi struct {
	X, Y int
}

func main() {
	var opts struct {
		InFile  string `short:"i"`
		OutFile string `short:"o"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}
	var aseFile asefile.AsepriteFile
	err = aseFile.DecodeFile(opts.InFile)
	if err != nil {
		panic(err)
	}
	if len(aseFile.Frames) == 0 {
		println("No frames in file")
		return
	}
	collisLayer := -1
	for x, layer := range aseFile.Frames[0].Layers {
		if layer.LayerName == "<Collision>" {
			collisLayer = x
		}
	}
	if collisLayer == -1 {
		println("No collision layers - use <Collision> to define collision layers")
		return
	}

	// For each frame we want to scan the pixels and find all the ones of a certain color
	for framei, frame := range aseFile.Frames {
		points := make([]Pointi, 0)
		if len(points) == 0 {

		}
		// Find the highest and lowest point

		cel := frame.Cels[collisLayer]
		dat := cel.RawCelData
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
	}
}
