//go:build integration

package zpl_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/kamilturek/go-zpl"
	zplutils "github.com/kamilturek/go-zpl/utils"
)

func TestConvert(t *testing.T) {
	t.Parallel()

	args := []string{"testdata/hello.zpl"}
	output := &bytes.Buffer{}

	if err := zpl.Convert(
		zpl.WithInputFromArgs(args),
		zpl.WithOutputFormat(zpl.PNG),
		zpl.WithOutput(output),
	); err != nil {
		t.Fatal(err)
	}

	want, err := os.Open("testdata/hello.png")
	if err != nil {
		t.Fatal(err)
	}

	got := output
	if err := zplutils.CompareImages(want, got); err != nil {
		t.Fatal(err)
	}
}

func TestConvertBytes(t *testing.T) {
	t.Parallel()

	res, err := zpl.ConvertBytes(
		[]byte("^xa^cfa,50^fo100,100^fdHello World^fs^xz"),
		zpl.WithOutputFormat(zpl.PNG),
	)
	if err != nil {
		t.Fatal(err)
	}

	want, err := os.Open("testdata/hello.png")
	if err != nil {
		t.Fatal(err)
	}

	got := bytes.NewBuffer(res)
	if err := zplutils.CompareImages(want, got); err != nil {
		t.Fatal(err)
	}
}

func TestToPNG(t *testing.T) {
	t.Parallel()

	res, err := zpl.ToPNG([]byte("^xa^cfa,50^fo100,100^fdHello World^fs^xz"))
	if err != nil {
		t.Fatal(err)
	}

	want, err := os.Open("testdata/hello.png")
	if err != nil {
		t.Fatal(err)
	}

	got := bytes.NewBuffer(res)
	if err := zplutils.CompareImages(want, got); err != nil {
		t.Fatal(err)
	}
}
