//go:build integration

package zpl_test

import (
	"os"
	"testing"

	"github.com/kamilturek/go-zpl"
	zplutils "github.com/kamilturek/go-zpl/utils"
)

func TestConvert(t *testing.T) {
	t.Parallel()

	args := []string{"testdata/hello.zpl"}

	c, err := zpl.NewConverter(
		zpl.WithInputFromArgs(args),
		zpl.WithOutputFormat(zpl.PNG),
	)
	if err != nil {
		t.Fatal(err)
	}

	want, err := os.Open("testdata/hello.png")
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.Convert()
	if err != nil {
		t.Fatal(err)
	}

	if err := zplutils.CompareImages(want, got); err != nil {
		t.Fatal(err)
	}
}
