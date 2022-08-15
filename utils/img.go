//go:build integration

package utils

import (
	"errors"
	"image"
	"io"

	"github.com/google/go-cmp/cmp"

	_ "image/png"
)

func CompareImages(a io.Reader, b io.Reader) error {
	ia, _, err := image.Decode(a)
	if err != nil {
		return err
	}

	ib, _, err := image.Decode(b)
	if err != nil {
		return err
	}

	if !cmp.Equal(ia.Bounds(), ib.Bounds()) {
		return errors.New("bounds not equal")
	}

	for x := ia.Bounds().Min.X; x < ia.Bounds().Max.X; x++ {
		for y := ia.Bounds().Min.Y; y < ia.Bounds().Max.Y; y++ {
			ra, ga, ba, aa := ia.At(x, y).RGBA()
			rb, gb, bb, ab := ia.At(x, y).RGBA()

			if ra != rb || ga != gb || ba != bb || aa != ab {
				return errors.New("colors not equal")
			}
		}
	}

	return nil
}
