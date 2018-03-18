# reximage

## A Go package For importing REXPaint's .xp images

REXPaint (www.gridsagegames.com/rexpaint) is a fabulous program for creating ASCII art, made by ultra-famous roguelike developer Kyzrati. This package allows you to decode and import image data from the .xp file format produced by REXPaint for use in your project. It is part of the larger burl-E engine, but is fine to be used standalone!

## Documentation

Obtain this package in the usual way:

`go get github.com/bennicholls/burl-E/reximage`

To use, import the package and call

```Go
image, err := reximage.Import(pathname)
```

ensuring that `pathname` is a string describing a path to an .xp file. It will throw an error if the pathname is invalid or the file cannot be read for some other reason.

Access the cell data using:

```Go
cell, err := image.GetCell(x, y)
```

where x and y are coordinates to a cell in the image (bounded by `image.Width, image.Height`). It will throw an error if (x,y) is not in bounds or if the image hasn't been imported yet. My engine `burl-E` uses SDL, so `reximage` follows SDL's convention of setting the (0,0) coordinate in the top-left of the image.

`cell` consists of a `Glyph` (ASCII codepoint) and RGB components for both foreground and background colours. Each component is 8bits (0-255). You can extract 32bit colours using some helper functions:

```Go
foregroundRGBA, backgroundRGBA := cell.RGBA()
foregroundARGB, backgroundARGB := cell.ARGB()
```

I imagine these are the most popular colour formats but I don't have the research to back that up. Of course, if your program uses a different colour format you can access the individual components and form the colour yourself.

Complete(-ish) documentation can be found on [GoDoc.org](https://godoc.org/github.com/BenNicholls/burl-E/reximage).

## Future

This is all my engine needs at the moment so I'm not sure what else would be helpful here. On the off-chance someone else is using this, let me know if there's anything you think it could use. Conceivably, the package could be expanded to write changes to .xp files I suppose?

## License

reximage is licensed under the MIT license (see LICENSE file).