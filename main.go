package main

import (
	"log"
	"net/http"

	"go-tex2pdf/converter"
	"go-tex2pdf/server"
)

const Version = "0.1.0"

func main() {
	texConverter := converter.New(
		"/usr//bin/pdflatex",
		1,
		"/my/asset/dir:/my/other/asset/dir",
	)

	httpServer := server.New(texConverter, Version)

	port := "8080"
	log.Printf("Starting TeX to PDF converter server v%s on port %s", Version, port)
	log.Fatal(http.ListenAndServe(":"+port, httpServer.HandleRequests()))
}
