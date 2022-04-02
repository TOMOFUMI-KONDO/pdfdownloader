package pdfdownloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"

	"golang.org/x/net/html"
)

type Downloader struct {
	url        *url.URL
	dir        string
	user, pass string
}

func NewDownloader(rawURL, user, pass string) (*Downloader, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url '%s': %w", rawURL, err)
	}

	return &Downloader{
		url:  u,
		dir:  path.Dir(u.Path),
		user: user,
		pass: pass,
	}, nil
}

func (d *Downloader) Download() ([]*PDF, error) {
	req, err := d.newRequestWithBasicAuth("GET", nil, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request '%s': %w", req.URL, err)
	}
	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse responsed html: %w", err)
	}

	var pdfs []*PDF
	for _, l := range sessionPageLinks(doc) {
		result, err := d.downloadPDFs(l)
		if err != nil {
			return nil, fmt.Errorf("failed to download PDFs '%s': %w", l, err)
		}
		pdfs = append(pdfs, result...)
	}

	return pdfs, nil
}

func (d *Downloader) downloadPDFs(sessionPageLink string) ([]*PDF, error) {
	req, err := d.newRequestWithBasicAuth("GET", nil, sessionPageLink)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	fmt.Println(req.URL)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request '%s': %w", req.URL, err)
	}
	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse responsed html: %w", err)
	}

	var pdfs []*PDF
	mu := &sync.Mutex{}
	var wg sync.WaitGroup
	for _, l := range pdfLinks(doc) {
		wg.Add(1)
		go func(pdfLink, title string) {
			d.downloadPDF(pdfLink, title, &pdfs, mu)
			fmt.Println(title)
			wg.Done()
		}(path.Join(path.Dir(sessionPageLink), l.link), l.title)
	}
	wg.Wait()

	return pdfs, nil
}

func (d *Downloader) downloadPDF(pdfLink, title string, out *[]*PDF, mu *sync.Mutex) {
	req, err := d.newRequestWithBasicAuth("GET", nil, pdfLink)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create request: %v\n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to request: '%s': %v\n", req.URL, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "faiiled to read response: %v\n", err)
	}

	mu.Lock()
	*out = append(*out, NewPDF(path.Base(pdfLink), title, string(body)))
	mu.Unlock()
}

func (d *Downloader) newRequestWithBasicAuth(method string, body io.Reader, link string) (*http.Request, error) {
	var url_ string

	if link != "" {
		path_ := path.Join(d.dir, link)
		url_ = fmt.Sprintf("%s://%s/%s", d.url.Scheme, d.url.Host, path_)
	} else {
		url_ = d.url.String()
	}

	req, err := http.NewRequest(method, url_, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(d.user, d.pass)

	return req, nil
}
