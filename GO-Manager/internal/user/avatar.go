package user

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"

	_ "golang.org/x/image/webp"
)

const (
	MaxAvatarBytes = 2 * 1024 * 1024
	MaxAvatarSide  = 512
)

var (
	ErrAvatarTooLarge    = errors.New("avatar too large")
	ErrInvalidAvatarType = errors.New("invalid avatar type")
	ErrAvatarDimensions  = errors.New("avatar dimensions invalid")
)

type Avatar struct {
	Bytes     []byte
	MediaType string
}

func prepareAvatar(raw []byte) (Avatar, error) {
	if len(raw) == 0 {
		return Avatar{}, ErrInvalidAvatarType
	}
	if len(raw) > MaxAvatarBytes {
		return Avatar{}, ErrAvatarTooLarge
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(raw))
	if err != nil {
		return Avatar{}, ErrInvalidAvatarType
	}
	switch format {
	case "jpeg", "png", "webp":
	default:
		return Avatar{}, ErrInvalidAvatarType
	}

	if cfg.Width > MaxAvatarSide || cfg.Height > MaxAvatarSide {
		return Avatar{}, ErrAvatarDimensions
	}

	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return Avatar{}, ErrInvalidAvatarType
	}

	if cfg.Width != cfg.Height {
		side := cfg.Width
		if cfg.Height < side {
			side = cfg.Height
		}
		origin := image.Pt((cfg.Width-side)/2, (cfg.Height-side)/2)
		dst := image.NewRGBA(image.Rect(0, 0, side, side))
		draw.Draw(dst, dst.Bounds(), img, origin, draw.Src)
		img = dst
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return Avatar{}, err
	}
	return Avatar{Bytes: buf.Bytes(), MediaType: "image/png"}, nil
}
