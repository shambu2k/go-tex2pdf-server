package converter

import (
	"github.com/rwestlund/gotex"
)

// Converter handles the conversion of LaTeX to PDF.
type Converter struct {
	pdfLatexPath string
	runs         int
	texInputs    string
}

// New creates a new Converter instance.
func New(pdfLatexPath string, runs int, texInputs string) *Converter {
	return &Converter{
		pdfLatexPath: pdfLatexPath,
		runs:         runs,
		texInputs:    texInputs,
	}
}

// Convert transforms LaTeX content into PDF.
func (c *Converter) Convert(tex string) ([]byte, error) {
	pdf, err := gotex.Render(tex, gotex.Options{
		Command:   c.pdfLatexPath,
		Runs:      c.runs,
		Texinputs: c.texInputs,
	})
	if err != nil {
		return nil, err
	}
	return pdf, nil
}
