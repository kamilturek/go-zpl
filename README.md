# go-zpl

A CLI tool & Go package for conversion of ZPL files.
A wrapper around [Labelary ZPL Web Service](http://labelary.com/service.html).

## Usage

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
