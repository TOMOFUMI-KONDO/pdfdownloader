package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/TOMOFUMI-KONDO/pdfdownloader"
)

var (
	url, user, pass string
)

func init() {
	flag.StringVar(&url, "url", "", "URL of web site from which you download PDF")
	flag.StringVar(&user, "user", "", "user name for Basic Authentication")
	flag.StringVar(&pass, "pass", "", "password for Basic Authentication")
	flag.Parse()
}

func main() {
	d, err := pdfdownloader.NewDownloader(url, user, pass)
	if err != nil {
		panic(err)
	}

	pdfs, err := d.Download()
	if err != nil {
		panic(err)
	}

	for _, p := range pdfs {
		if err := writePDF(p); err != nil {
			panic(err)
		}
	}
}

func writePDF(p *pdfdownloader.PDF) error {
	sessNum := p.FileName[:strings.LastIndexByte(p.FileName, '.')]
	fileName := fmt.Sprintf("%s_%s.pdf", sessNum, strings.ReplaceAll(p.Title, "/", "-"))

	f, err := os.Create(path.Join("pdfs", fileName))
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", fileName, err)
	}
	defer f.Close()

	if _, err = f.WriteString(p.Body); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}
