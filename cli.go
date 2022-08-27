package zpl

import (
	"flag"
	"fmt"
	"os"
)

var version string

func RunCLI() {
	showVersion := flag.Bool("v", false, "version for go-zpl")
	showVersionLong := flag.Bool("version", false, "version for go-zpl")

	density := flag.Int("d", 8, "input label density [dpmm]")
	width := flag.Int("w", 4, "input label width [inch]")
	height := flag.Int("h", 6, "input label height [inch]")
	outputPath := flag.String("o", "", "output file path")
	outputFormat := flag.String("f", "png", "output file format")
	flag.Parse()

	if *showVersion || *showVersionLong {
		fmt.Fprintln(os.Stdout, version)
		os.Exit(0)
	}

	if err := Convert(
		WithInputFromArgs(flag.Args()),
		WithOutputPath(*outputPath),
		WithOutputFormat(*outputFormat),
		WithDensity(*density),
		WithWidth(*width),
		WithHeight(*height),
	); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
