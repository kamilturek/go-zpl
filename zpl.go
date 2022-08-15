package zpl

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func RunCLI() {
}

type converter struct {
	input   io.Reader
	density int // dpmm
	width   int // inch
	height  int // inch
}

type option func(c *converter) error

func WithInput(input io.Reader) option {
	return func(c *converter) error {
		if input == nil {
			return errors.New("nil input reader")
		}

		c.input = input

		return nil
	}
}

func WithInputFromArgs(args []string) option {
	return func(c *converter) error {
		if len(args) == 0 {
			return nil
		}

		f, err := os.Open(args[0])
		if err != nil {
			return err
		}

		c.input = f

		return nil
	}
}

func NewConverter(opts ...option) (*converter, error) {
	c := &converter{
		input:   os.Stdin,
		width:   4,
		height:  6,
		density: 8,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *converter) ToZPL() (io.Reader, error) {
	return c.input, nil
}

func (c *converter) ToPNG() (io.Reader, error) {
	return c.doRequest("image/png")
}

func (c *converter) ToPDF() (io.Reader, error) {
	return c.doRequest("application/pdf")
}

const templateURL = "http://api.labelary.com/v1/printers/%ddpmm/labels/%dx%d/0/"

func (c *converter) doRequest(contentType string) (io.Reader, error) {
	url := fmt.Sprintf(templateURL, c.density, c.width, c.height)

	req, err := http.NewRequest(http.MethodPost, url, c.input)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", contentType)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
