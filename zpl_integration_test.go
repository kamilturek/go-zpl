//go:build integration

package zpl_test

import (
	"os"
	"testing"

	"github.com/kamilturek/go-zpl"
	"github.com/kamilturek/go-zpl/utils"
)

func TestToPNG(t *testing.T) {
	t.Parallel()

	args := []string{"testdata/hello.zpl"}

	c, err := zpl.NewConverter(
		zpl.WithInputFromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}

	want, err := os.Open("testdata/hello.png")
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ToPNG()
	if err != nil {
		t.Fatal(err)
	}

	if err := utils.CompareImages(want, got); err != nil {
		t.Fatal(err)
	}
}
