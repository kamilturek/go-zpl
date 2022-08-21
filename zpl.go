package zpl

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

type converter struct {
	Input        io.Reader
	output       io.Writer
	outputFormat string
	density      int // dpmm
	width        int // inch
	height       int // inch
}

type option func(c *converter) error

func WithInput(input io.Reader) option {
	return func(c *converter) error {
		if input == nil {
			return ErrNilInput
		}

		c.Input = input

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

		c.Input = f

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
		Input:        os.Stdin,
		output:       os.Stdout,
		outputFormat: PNG,
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

func (conv *converter) Convert() error {
	converted, err := conv.doRequest()
	if err != nil {
		return err
	}

	defer converted.Close()

	if _, err := io.Copy(conv.output, converted); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

func (conv *converter) doRequest() (io.ReadCloser, error) {
	url := fmt.Sprintf("http://api.labelary.com/v1/printers/%ddpmm/labels/%dx%d/0/", conv.density, conv.width, conv.height)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, conv.Input)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	contentTypes := map[string]string{
		PDF: "application/pdf",
		PNG: "image/png",
	}

	req.Header.Set("Accept", contentTypes[conv.outputFormat])

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return res.Body, nil
}

func Convert(opts ...option) error {
	conv, err := NewConverter(opts...)
	if err != nil {
		return err
	}

	if err := conv.Convert(); err != nil {
		return fmt.Errorf("failed to convert: %w", err)
	}

	return nil
}

func ConvertBytes(content []byte, opts ...option) ([]byte, error) {
	output := &bytes.Buffer{}

	options := []option{
		WithInput(bytes.NewBuffer(content)),
		WithOutput(output),
	}
	options = append(options, opts...)

	if err := Convert(options...); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func ToPNG(content []byte, opts ...option) ([]byte, error) {
	options := []option{WithOutputFormat(PNG)}
	options = append(options, opts...)

	return ConvertBytes(content, options...)
}

func ToPDF(content []byte, opts ...option) ([]byte, error) {
	options := []option{WithOutputFormat(PDF)}
	options = append(options, opts...)

	return ConvertBytes(content, options...)
}
