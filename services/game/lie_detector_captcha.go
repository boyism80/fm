package game

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math/rand/v2"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	captchaLength   = 6
	captchaAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	captchaWidth    = 194
	captchaHeight   = 43
	captchaScale    = 3
	captchaLines    = 6
	captchaDots     = 200
)

type Captcha struct {
	Image  []byte
	Answer string
}

func NewCaptcha() (*Captcha, error) {
	answer := make([]byte, captchaLength)
	for i := range answer {
		answer[i] = captchaAlphabet[rand.IntN(len(captchaAlphabet))]
	}

	img := image.NewRGBA(image.Rect(0, 0, captchaWidth, captchaHeight))
	background := color.RGBA{uint8(225 + rand.IntN(31)), uint8(225 + rand.IntN(31)), uint8(225 + rand.IntN(31)), 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: background}, image.Point{}, draw.Src)

	face := basicfont.Face7x13
	glyph := image.NewAlpha(image.Rect(0, 0, face.Advance, face.Height))
	for i, c := range answer {
		draw.Draw(glyph, glyph.Bounds(), image.Transparent, image.Point{}, draw.Src)
		drawer := &font.Drawer{Dst: glyph, Src: image.Opaque, Face: face, Dot: fixed.P(0, face.Ascent)}
		drawer.DrawString(string(c))

		ink := color.RGBA{uint8(rand.IntN(120)), uint8(rand.IntN(120)), uint8(rand.IntN(120)), 255}
		left := 12 + i*29 + rand.IntN(5)
		top := rand.IntN(captchaHeight - face.Height*captchaScale + 1)
		for y := 0; y < face.Height*captchaScale; y++ {
			for x := 0; x < face.Advance*captchaScale; x++ {
				if glyph.AlphaAt(x/captchaScale, y/captchaScale).A != 0 {
					img.Set(left+x, top+y, ink)
				}
			}
		}
	}

	for range captchaLines {
		ink := color.RGBA{uint8(rand.IntN(200)), uint8(rand.IntN(200)), uint8(rand.IntN(200)), 255}
		x0, y0 := rand.IntN(captchaWidth), rand.IntN(captchaHeight)
		x1, y1 := rand.IntN(captchaWidth), rand.IntN(captchaHeight)
		steps := max(x1-x0, x0-x1, y1-y0, y0-y1, 1)
		for s := 0; s <= steps; s++ {
			img.Set(x0+(x1-x0)*s/steps, y0+(y1-y0)*s/steps, ink)
		}
	}
	for range captchaDots {
		ink := color.RGBA{uint8(rand.IntN(256)), uint8(rand.IntN(256)), uint8(rand.IntN(256)), 255}
		img.Set(rand.IntN(captchaWidth), rand.IntN(captchaHeight), ink)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
		return nil, err
	}
	return &Captcha{Image: buf.Bytes(), Answer: string(answer)}, nil
}
