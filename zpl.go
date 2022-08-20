package zpl

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

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
			return ErrNilInput
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
			return fmt.Errorf("failed to open input file: %w", err)
		}

		c.input = f

		return nil
	}
}

func WithOutput(output io.Writer) option {
	return func(c *converter) error {
		if output == nil {
			return ErrNilOutput
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
			return fmt.Errorf("failed to open output file: %w", err)
		}

		c.output = f

		return nil
	}
}

func WithOutputFormat(outputFormat string) option {
	return func(c *converter) error {
		for _, allowedOutputFormat := range allowedOutputFormats() {
			if outputFormat == allowedOutputFormat {
				c.outputFormat = outputFormat

				return nil
			}
		}

		return ErrInvalidOutputFormat
	}
}

func WithDensity(density int) option {
	return func(c *converter) error {
		for _, allowedDensity := range allowedDensities() {
			if density == allowedDensity {
				c.density = density

				return nil
			}
		}

		return ErrInvalidDensity
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

func (c *converter) Convert() (io.ReadCloser, error) {
	if c.outputFormat == ZPL {
		return io.NopCloser(c.input), nil
	}

	return c.doRequest()
}

func (c *converter) ConvertAndWrite() error {
	output, err := c.Convert()
	if err != nil {
		return err
	}

	defer output.Close()

	_, err = io.Copy(c.output, output)

	return fmt.Errorf("failed to write output: %w", err)
}

func (c *converter) doRequest() (io.ReadCloser, error) {
	url := fmt.Sprintf("http://api.labelary.com/v1/printers/%ddpmm/labels/%dx%d/0/", c.density, c.width, c.height)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, c.input)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	contentTypes := map[string]string{
		PDF: "application/pdf",
		PNG: "image/png",
	}

	req.Header.Set("Accept", contentTypes[c.outputFormat])

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return res.Body, nil
}
