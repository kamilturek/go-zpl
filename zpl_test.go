package zpl_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kamilturek/go-zpl"
)

func TestWithInput(t *testing.T) {
	t.Parallel()

	input := bytes.NewBufferString("^xa^xz")

	c, err := zpl.NewConverter(
		zpl.WithInput(input),
	)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte("^xa^xz")

	output, err := c.Convert()
	if err != nil {
		t.Fatal(err)
	}

	got, err := io.ReadAll(output)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestWithInputFromArgs(t *testing.T) {
	t.Parallel()

	args := []string{"testdata/hello.zpl"}

	c, err := zpl.NewConverter(
		zpl.WithInputFromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte("^xa^cfa,50^fo100,100^fdHello World^fs^xz")

	output, err := c.Convert()
	if err != nil {
		t.Fatal(err)
	}

	got, err := io.ReadAll(output)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestInputFromArgsNoArgs(t *testing.T) {
	t.Parallel()

	input := bytes.NewBufferString("^xa^xz")
	args := []string{}

	c, err := zpl.NewConverter(
		zpl.WithInput(input),
		zpl.WithInputFromArgs(args),
	)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte("^xa^xz")

	output, err := c.Convert()
	if err != nil {
		t.Fatal(err)
	}

	got, err := io.ReadAll(output)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}
