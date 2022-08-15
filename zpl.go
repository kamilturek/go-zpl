package zpl

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const PDF = "pdf"
const PNG = "png"
const ZPL = "zpl"

var allowedDensities = []int{6, 8, 12, 24}
var allowedOutputFormats = []string{PDF, PNG, ZPL}

var contentTypes = map[string]string{
	PDF: "application/pdf",
	PNG: "image/png",
}

type converter struct {
	input        io.Reader
	density      int // dpmm
	width        int // inch
	height       int // inch
	output       io.Writer
	outputFormat string
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

func WithOutput(output io.Writer) option {
	return func(c *converter) error {
		if output == nil {
			return errors.New("nil output writer")
		}

		c.output = output

		return nil
	}
}

func WithOutputPath(path string) option {
	return func(c *converter) error {
		if path == "" {
			return nil
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}

		c.output = f

		return nil
	}
}

func WithOutputFormat(outputFormat string) option {
	return func(c *converter) error {
		for _, allowedOutputFormat := range allowedOutputFormats {
			if outputFormat == allowedOutputFormat {
				c.outputFormat = outputFormat

				return nil
			}
		}

		return fmt.Errorf("invalid output format: must be one of %v", allowedOutputFormats)
	}
}

func WithDensity(density int) option {
	return func(c *converter) error {
		for _, allowedDensity := range allowedDensities {
			if density == allowedDensity {
				c.density = density

				return nil
			}
		}

		return fmt.Errorf("invalid density: must be one of %v", allowedDensities)
	}
}

func WithWidth(width int) option {
	return func(c *converter) error {
		c.width = width

		return nil
	}
}

func WithHeight(height int) option {
	return func(c *converter) error {
		c.height = height

		return nil
	}
}

func NewConverter(opts ...option) (*converter, error) {
	c := &converter{
		input:        os.Stdin,
		output:       os.Stdout,
		outputFormat: ZPL,
		density:      8,
		width:        4,
		height:       6,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *converter) Convert() (io.Reader, error) {
	if c.outputFormat == ZPL {
		return c.input, nil
	}

	return c.doRequest()
}

func (c *converter) ConvertAndWrite() error {
	output, err := c.Convert()
	if err != nil {
		return err
	}

	_, err = io.Copy(c.output, output)
	return err
}

func (c *converter) doRequest() (io.Reader, error) {
	url := fmt.Sprintf("http://api.labelary.com/v1/printers/%ddpmm/labels/%dx%d/0/", c.density, c.width, c.height)

	req, err := http.NewRequest(http.MethodPost, url, c.input)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", contentTypes[c.outputFormat])

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
