package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"

	"go-tex2pdf/converter"
)

type Server struct {
	converter *converter.Converter
	version   string
}

func New(converter *converter.Converter, version string) *Server {
	return &Server{
		converter: converter,
		version:   version,
	}
}

func (s *Server) ProcessRequest(w http.ResponseWriter, r *http.Request) {
	// Read raw LaTeX content directly from request body
	texContent, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(texContent) == 0 {
		http.Error(w, "No LaTeX content provided", http.StatusBadRequest)
		return
	}

	pdf, err := s.converter.Convert(string(texContent))
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
		fmt.Fprintf(w, "%s (Go: %s)\n", s.version, runtime.Version())
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Welcome to the LaTeX to PDF converter API v%s.\n", s.version)
		fmt.Fprintf(w, "Send POST requests to /convert with raw LaTeX content in the request body\n")
		fmt.Fprintf(w, "Check /version for version information.\n")
	})

	return mux
}
