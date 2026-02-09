package infrastructure

import (
	"bytes"
	"html/template"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type WkhtmlPdfGenerator struct{}

func NewWkhtmlPdfGenerator() *WkhtmlPdfGenerator {
	return &WkhtmlPdfGenerator{}
}

func (g *WkhtmlPdfGenerator) Generate(
	templatePath string,
	data any,
) ([]byte, error) {

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	var htmlBuf bytes.Buffer
	if err := tpl.Execute(&htmlBuf, data); err != nil {
		return nil, err
	}

	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	page := wkhtml.NewPageReader(&htmlBuf)
	page.EnableLocalFileAccess.Set(true)

	pdfg.AddPage(page)
	pdfg.PageSize.Set(wkhtml.PageSizeA4)

	if err := pdfg.Create(); err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}
