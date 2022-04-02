# pdfdownloader

This is a pdf-downloader for papers of IPSJ2021.
This does web scraping IPSJ web site to get pdf links of papers, and fetch that links.

## Example

```bash
go run ./cmd/main.go -url=XXX -user=XXX -pass=pass
```

## Parameters

### -url

URL of IPSJ website that has pdf links of papers.

### -user

username for Basic Authentication

### -pass

password for Basic Authentication
