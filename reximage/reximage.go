// A package for decoding and handling .xp files produced by Kyzrati's fabulous
// REXPaint program, the gold-standard in ASCII art drawing programs. It can be
// found at www.gridsagegames.com/rexpaint.
//
// reximage is part of the BURL-E engine by Benjamin Nicholls, but feel free to
// use it as a standalone package!
package reximage

import (
	"compress/gzip"
	"encoding/binary"
	"errors"
	"os"
	"strings"
)

// ImageData is the struct holding the decoded and exported image data.
type ImageData struct {
	Width  int
	Height int
	Cells  []CellData //will have Width*Height Elements
}

// GetCell returns the CellData at coordinate (x, y) of the decoded image, with (0,0)
// at the top-left of the image.
func (id ImageData) GetCell(x, y int) (cd CellData, err error) {
	if id.Cells == nil || len(id.Cells) == 0 {
		return CellData{}, errors.New("Image has no data.")
	}

	if x >= id.Width || y >= id.Height || x+y*id.Width > len(id.Cells) {
		return CellData{}, errors.New("x, y coordinates out of bounds.")
	}

	cd = id.Cells[x+y*id.Width]

	return
}

// CellData holds the decoded data for a single cell. Colours are split into uint8
// components so the user can combine them into whatever colour format they need.
// Some popular colour format conversion functions are provided as well.
type CellData struct {
	Glyph uint32 // ASCII code for glyph
	R_f   uint8  //
	G_f   uint8  // Foreground Colour Elements
	B_f   uint8  //
	R_b   uint8  //
	G_b   uint8  // Background Colour Elements
	B_b   uint8  //
}

// ARGB returns the foreground and background colours of the cell in ARGB format.
// Alpha in this case is always set to maximum (255)
func (cd CellData) ARGB() (fore, back uint32) {
	fore = uint32(0xFF << 24) //set alpha to 255
	fore |= uint32(cd.R_f) << 16
	fore |= uint32(cd.G_f) << 8
	fore |= uint32(cd.B_f)

	back = uint32(0xFF << 24) //set alpha to 255
	back |= uint32(cd.R_b) << 16
	back |= uint32(cd.G_b) << 8
	back |= uint32(cd.B_b)

	return
}

// RGBA returns the foreground and background colours of the cell in RGBA format.
// Alpha in this case is always set to maximum (255)
func (cd CellData) RGBA() (fore, back uint32) {
	fore = uint32(cd.R_f) << 24
	fore |= uint32(cd.G_f) << 16
	fore |= uint32(cd.B_f) << 8
	fore |= 0xFF //set alpha to 255

	back = uint32(cd.R_b) << 24
	back |= uint32(cd.G_b) << 16
	back |= uint32(cd.B_b) << 8
	back |= 0xFF //set alpha to 255

	return
}

//Import imports an image from the xp file at the provided path. Returns the Imagedata and an error.
//If an error is present, ImageData will be no good.
func Import(path string) (image ImageData, err error) {
	image = ImageData{}

	if !strings.HasSuffix(path, ".xp") {
		err = errors.New("File is not an XP image.")
		return
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	//xp data is gzipped
	data, err := gzip.NewReader(f)
	if err != nil {
		return
	}

	//read rexpaint version num and the number of layers
	var version int32
	var numLayers uint32
	err = binary.Read(data, binary.LittleEndian, &version)
	err = binary.Read(data, binary.LittleEndian, &numLayers)
	if err != nil {
		return
	}

	//read into the first layer so we can get the image dimensions and initialize
	//cell data
	var w, h uint32
	err = binary.Read(data, binary.LittleEndian, &w)
	err = binary.Read(data, binary.LittleEndian, &h)
	if err != nil {
		return
	}

	image.Width, image.Height = int(w), int(h)
	image.Cells = make([]CellData, image.Width*image.Height)

	//read layers, painting from lowest layer to highest
	for layer := 0; layer < int(numLayers); layer++ {
		if layer != 0 {
			//if reading subsequent layers, throw away the dimension
			//bytes since we've already read them before
			err = binary.Read(data, binary.LittleEndian, &w)
			err = binary.Read(data, binary.LittleEndian, &h)
		}

		for i := 0; i < image.Width*image.Height; i++ {
			//read bytes for each cell.
			c := CellData{}
			err = binary.Read(data, binary.LittleEndian, &c)
			if err != nil {
				return
			}

			//check for undrawn cell, identified by bgcolour = (255, 0, 255)
			if c.R_b == 255 && c.G_b == 0 && c.B_b == 255 {
				continue
			} else {
				//xp images are encoded in the totally insane column-major order for some reason,
				//we correct that here (sorry Kyzrati, gotta put my foot down on this one)
				x, y := i/image.Height, i%image.Height
				image.Cells[y*image.Width+x] = c
			}
		}
	}

	return
}
