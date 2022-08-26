# go-zpl

A CLI tool & Go package for conversion of ZPL files.
A wrapper around [Labelary ZPL Web Service](http://labelary.com/service.html).

## Installation

### CLI Tool

Using Go:

```bash
go install github.com/kamilturek/go-zpl/cmd/go-zpl@latest
```

Using Homebrew:

```bash
brew install kamilturek/tap/go-zpl
```

### Go package

```bash
go get github.com/kamilturek/go-zpl
```

## Usage

### CLI Tool

`go-zpl` provides a help message describing all available options along with
their default values.

```bash
$ go-zpl --help              
Usage of go-zpl:
  -d int
        input label density [dpmm] (default 8)
  -f string
        output file format (default "png")
  -h int
        input label height [inch] (default 6)
  -o string
        output file path
  -w int
        input label width [inch] (default 4)
```

By default, `go-zpl` reads ZPL data from the standard input and writes result
to the standard output.

```bash
$ cat hello.zpl | go-zpl -f pdf | head -n 1
%PDF-1.4
```

It is possible to read ZPL data from a file and write the result to another one.

```bash
$ go-zpl -f pdf -o hello.pdf hello.zpl
$ head -n 1 hello.pdf                   
%PDF-1.4
```

### Go package

```go
package main

import "github.com/kamilturek/go-zpl"

func main() {
      // Using convenience wrapper
      zpl.ToPNG(
            []byte("^xa^cfa,50^fo100,100^fdHello World^fs^xz"),
            // Optional extra configuration
            zpl.WithWidth(4),
            zpl.WithHeight(6),
            zpl.WithDensity(8),
      )

      // Or configuring the converter from scratch
      f, _ := os.Open("hello.zpl")
      zpl.Convert(
            zpl.WithInput(f),
            zpl.WithOutput(os.Stdout),
            zpl.WithOutputFormat(zpl.PNG),
            zpl.WithDensity(12),
      )
}
```

## License

See [LICENSE](LICENSE.md).

## Contributing

Please create a GitHub issue for any feedback, bugs, requests or issues.
PRs are also welcome.

### Running tests

- General tests

    ```bash
    go test
    ```

- Integration tests

    ```bash
    go test -tags=integration
    ```
