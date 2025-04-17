package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/rwestlund/gotex"
)

const Version = "0.1.0"

type TexRequest struct {
	Tex string `json:"tex"`
}

type Converter struct {
	pdfLatexPath string
	runs         int
	texInputs    string
}

func NewConverter(pdfLatexPath string, runs int, texInputs string) *Converter {
	return &Converter{
		pdfLatexPath: pdfLatexPath,
		runs:         runs,
		texInputs:    texInputs,
	}
}

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

type Server struct {
	converter *Converter
}

func NewServer(converter *Converter) *Server {
	return &Server{
		converter: converter,
	}
}

func (s *Server) ProcessRequest(w http.ResponseWriter, r *http.Request) {
	var req TexRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Tex == "" {
		http.Error(w, "No LaTeX content provided", http.StatusBadRequest)
		return
	}

	pdf, err := s.converter.Convert(req.Tex)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		http.Error(w, "Failed to generate PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=document.pdf")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdf)))

	if _, err := w.Write(pdf); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (s *Server) HandleRequests() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
			return
		}
		s.ProcessRequest(w, r)
	})

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s (Go: %s)\n", Version, runtime.Version())
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Welcome to the LaTeX to PDF converter API v%s.\n", Version)
		fmt.Fprintf(w, "Send POST requests to /convert with {\"tex\": \"your LaTeX content\"}\n")
		fmt.Fprintf(w, "Check /version for version information.\n")
	})

	return mux
}

func main() {
	converter := NewConverter(
		"/usr/bin/pdflatex",
		1,
		"/my/asset/dir:/my/other/asset/dir",
	)

	server := NewServer(converter)

	port := "8080"
	log.Printf("Starting TeX to PDF converter server v%s on port %s", Version, port)
	log.Fatal(http.ListenAndServe(":"+port, server.HandleRequests()))
}
