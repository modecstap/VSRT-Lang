package user

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func encodePNG(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode: %v", err)
	}
	return buf.Bytes()
}

func encodeJPEG(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatalf("jpeg.Encode: %v", err)
	}
	return buf.Bytes()
}

func solidRGBA(w, h int, c color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	return img
}

func decodePNG(t *testing.T, raw []byte) image.Image {
	t.Helper()
	img, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("png.Decode: %v", err)
	}
	return img
}

func TestSetAvatar_SquareFormatsBecomePNG(t *testing.T) {
	t.Parallel()

	red := color.RGBA{R: 255, A: 255}
	cases := []struct {
		name string
		raw  []byte
		side int
	}{
		{name: "png", raw: encodePNG(t, solidRGBA(16, 16, red)), side: 16},
		{name: "jpeg", raw: encodeJPEG(t, solidRGBA(16, 16, red)), side: 16},
		{name: "webp", raw: readTestdata(t, "square.webp"), side: 8},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			u := &User{}
			err := u.SetAvatar(c.raw)
			if err != nil {
				t.Fatalf("SetAvatar: %v", err)
			}
			got := u.Avatar
			if got.MediaType != "image/png" {
				t.Fatalf("MediaType = %q, want image/png", got.MediaType)
			}
			bounds := decodePNG(t, got.Bytes).Bounds()
			if bounds.Dx() != c.side || bounds.Dy() != c.side {
				t.Fatalf("bounds = %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), c.side, c.side)
			}
		})
	}
}

func TestSetAvatar_CenterCropsNonSquare(t *testing.T) {
	t.Parallel()

	src := image.NewRGBA(image.Rect(0, 0, 500, 400))
	for y := 0; y < 400; y++ {
		for x := 0; x < 500; x++ {
			switch {
			case x < 50:
				src.SetRGBA(x, y, color.RGBA{R: 255, A: 255})
			case x >= 450:
				src.SetRGBA(x, y, color.RGBA{B: 255, A: 255})
			default:
				src.SetRGBA(x, y, color.RGBA{G: 255, A: 255})
			}
		}
	}

	u := &User{}
	err := u.SetAvatar(encodePNG(t, src))
	if err != nil {
		t.Fatalf("SetAvatar: %v", err)
	}
	img := decodePNG(t, u.Avatar.Bytes)
	if img.Bounds().Dx() != 400 || img.Bounds().Dy() != 400 {
		t.Fatalf("bounds = %dx%d, want 400x400", img.Bounds().Dx(), img.Bounds().Dy())
	}

	min := img.Bounds().Min
	left := color.RGBAModel.Convert(img.At(min.X, min.Y)).(color.RGBA)
	right := color.RGBAModel.Convert(img.At(min.X+399, min.Y)).(color.RGBA)
	if left.R != 0 || left.G != 255 {
		t.Fatalf("left-edge pixel after crop = %+v, want green (red strip gone)", left)
	}
	if right.B != 0 || right.G != 255 {
		t.Fatalf("right-edge pixel after crop = %+v, want green (blue strip gone)", right)
	}
}

func TestSetAvatar_Square512IsNotCropped(t *testing.T) {
	t.Parallel()

	u := &User{}
	err := u.SetAvatar(encodePNG(t, solidRGBA(512, 512, color.RGBA{A: 255})))
	if err != nil {
		t.Fatalf("SetAvatar: %v", err)
	}
	bounds := decodePNG(t, u.Avatar.Bytes).Bounds()
	if bounds.Dx() != 512 || bounds.Dy() != 512 {
		t.Fatalf("bounds = %dx%d, want 512x512", bounds.Dx(), bounds.Dy())
	}
}

func TestSetAvatar_Rejects(t *testing.T) {
	t.Parallel()

	var gifBuf bytes.Buffer
	if err := gif.Encode(&gifBuf, solidRGBA(8, 8, color.RGBA{A: 255}), nil); err != nil {
		t.Fatalf("gif.Encode: %v", err)
	}

	cases := []struct {
		name string
		raw  []byte
		want error
	}{
		{name: "empty", raw: []byte{}, want: ErrInvalidAvatarType},
		{name: "oversize", raw: make([]byte, MaxAvatarBytes+1), want: ErrAvatarTooLarge},
		{name: "gif", raw: gifBuf.Bytes(), want: ErrInvalidAvatarType},
		{name: "svg", raw: []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="8" height="8"></svg>`), want: ErrInvalidAvatarType},
		{name: "513x1", raw: encodePNG(t, solidRGBA(513, 1, color.RGBA{A: 255})), want: ErrAvatarDimensions},
		{name: "800x400", raw: encodePNG(t, solidRGBA(800, 400, color.RGBA{A: 255})), want: ErrAvatarDimensions},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			u := &User{}
			err := u.SetAvatar(c.raw)
			if err != c.want {
				t.Fatalf("error = %v, want %v", err, c.want)
			}
		})
	}
}

func readTestdata(t *testing.T, name string) []byte {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read testdata %s: %v", name, err)
	}
	return raw
}
