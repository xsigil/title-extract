package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	stdhtml "html"
	"io"
	"os"
	"strings"

	nethtml "golang.org/x/net/html"
)

const version = "1.1.0"

// ExtractTitle scans an HTML stream and returns the first <title> text.
// It aborts scanning once </head> or <body> is reached to avoid useless I/O.
func ExtractTitle(r io.Reader) (string, bool, error) {
	z := nethtml.NewTokenizer(r)

	inTitle := false
	var b strings.Builder

	for {
		tt := z.Next()
		switch tt {
		case nethtml.ErrorToken:
			err := z.Err()
			if err == io.EOF {
				if inTitle && b.Len() > 0 {
					t := strings.TrimSpace(stdhtml.UnescapeString(b.String()))
					if t != "" {
						return t, true, nil
					}
				}
				return "", false, nil
			}
			return "", false, err

		case nethtml.StartTagToken, nethtml.SelfClosingTagToken:
			name, _ := z.TagName()
			if bytes.EqualFold(name, []byte("title")) {
				inTitle = true
				b.Reset()
			} else if bytes.EqualFold(name, []byte("body")) {
				return "", false, nil
			}

		case nethtml.EndTagToken:
			name, _ := z.TagName()
			if inTitle && bytes.EqualFold(name, []byte("title")) {
				t := strings.TrimSpace(stdhtml.UnescapeString(b.String()))
				if t != "" {
					return t, true, nil
				}
				inTitle = false
			} else if bytes.EqualFold(name, []byte("head")) {
				return "", false, nil
			}

		case nethtml.TextToken:
			if inTitle {
				b.Write(z.Text())
			}
		}
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `title-extract - ultra-fast, zero-alloc HTML title extractor

Usage:
  title-extract [OPTIONS] [FILE...]
  cat index.html | title-extract

Options:
  -v, --version  Print version and exit
  -h, --help     Show this help message

Exit Status:
  0  Title successfully found and printed
  1  No <title> found in document
  2  I/O or syntax error
`)
}

func processReader(r io.Reader) int {
	br := bufio.NewReaderSize(r, 256*1024)
	title, ok, err := ExtractTitle(br)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	if !ok {
		return 1
	}

	title = strings.Join(strings.Fields(title), " ")
	fmt.Println(title)
	return 0
}

func main() {
	var showHelp bool
	var showVersion bool

	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Usage = printUsage
	flag.Parse()

	if showHelp {
		printUsage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("title-extract %s\n", version)
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) == 0 || (len(args) == 1 && args[0] == "-") {
		os.Exit(processReader(os.Stdin))
	}

	exitCode := 0
	for _, filename := range args {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening %s: %v\n", filename, err)
			exitCode = 2
			continue
		}
		res := processReader(f)
		f.Close()
		if res != 0 && exitCode == 0 {
			exitCode = res
		}
	}

	os.Exit(exitCode)
}
