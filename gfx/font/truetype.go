package font

import (
	xfont "golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"
)

// http://www.freetype.org/freetype2/docs/tutorial/step2.html

// LoadTrueType loads a truetype fontAtlas from the given stream and
// applies the given fontAtlas scale in points.
//
// The low and high values determine the lower and upper rune limits
// we should load for this Font. For standard ASCII this would be:32, 127.
func LoadTrueType(r io.Reader, lc TTFConfig) (*fontAtlas, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Read the TrueType Font
	ttf, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}

	// Create an image(512*512) to store all requested glyphs.
	iw, ih := 1024, 1024
	_, fg := image.Black, image.White
	rect := image.Rect(0, 0, iw, ih)
	img := image.NewRGBA(rect)
	padding := fixed.I(2)

	// Use an OpenType face to do the drawing.
	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    float64(lc.FontSize()),
		DPI:     72,
		Hinting: xfont.HintingNone,
	})
	if err != nil {
		return nil, err
	}

	// new font-atlas
	f := &fontAtlas{glyphs: make(map[rune]Glyph)}

	// Iterate over all relevant glyphs in the truetype fontAtlas and
	// draw them all to the image buffer.
	gb, err := ttf.Bounds(nil, fixed.I(lc.FontSize()), xfont.HintingNone)
	if err != nil {
		return nil, err
	}
	m := face.Metrics()

	var (
		gw = gb.Max.X - gb.Min.X
		gh = gb.Max.Y - gb.Min.Y
		hh = m.Ascent + m.Descent
		ah = m.Ascent
	)

	width := fixed.I(iw)
	dot := fixed.Point26_6{X: padding, Y: padding}
	adjust := padding / 2

	//for ch := low; ch <= high; ch++ {
	for _, ch := range lc.Runes() {
		bb, advance, ok := face.GlyphBounds(ch)
		if !ok {
			continue
		}
		gw := bb.Max.X - bb.Min.X
		gh := bb.Max.Y - bb.Min.Y

		// draw to canvas
		if x := dot.X + gw; x > width {
			dot.X = padding
			dot.Y += hh + padding
		}

		d := dot
		d.Y -= bb.Min.Y
		d.X -= bb.Min.X
		dr, mask, mp, _, ok := face.Glyph(d, ch)
		if !ok {
			continue
		}
		draw.DrawMask(img, dr, fg, image.Point{}, mask, mp, draw.Over)

		// record glyph region
		g := Glyph{
			Rune:    ch,
			Advance: int(advance.Floor()),
			X:       fixed2f32(dot.X - adjust),
			Y:       fixed2f32(dot.Y - adjust),
			Width:   fixed2f32(gw + padding),
			Height:  fixed2f32(gh + padding),

			XOffset: fixed2f32(bb.Min.X),
			YOffset: fixed2f32(ah + bb.Min.Y - padding),
		}

		// add padding
		dot.X += gw + padding

		// add to glyph-array
		f.addGlyphs(ch, g)
	}

	// set bounds
	f.gWidth = fixed2f32(gw)
	f.gHeight = fixed2f32(gh)

	// load image
	err = f.loadTex(img)
	if err != nil {
		return nil, err
	}
	// save baked fontAtlas-image
	//savePng(img)
	return f, nil
}

func fixed2f32(fixed fixed.Int26_6) float32 {
	return float32(fixed) / (1 << 6)
}

// debug only.
func savePng(img image.Image) {
	f, err := os.Create("ttf.png")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		log.Println(err)
	}
}
