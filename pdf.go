package pdfdownloader

type PDF struct {
	FileName, Title, Body string
}

func NewPDF(fileName, title, body string) *PDF {
	return &PDF{
		FileName: fileName,
		Title:    title,
		Body:     body,
	}
}
