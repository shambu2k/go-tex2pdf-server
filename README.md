# TeX to PDF Converter Server

A lightweight web service that converts LaTeX documents to PDF files using pdflatex. Supports fontawesome, preprint, enumitem and more.

## Features

- Simple RESTful API to convert LaTeX to PDF
- Accepts raw LaTeX content in request body
- Returns generated PDF as binary response
- Docker support for easy deployment

## Quick Start

### Using Docker

```bash
# Pull from Docker Hub
docker pull shambu2k/go-tex2pdf-server:latest

# Or pull from GitHub Container Registry
docker pull ghcr.io/shambu2k/go-tex2pdf-server:latest

# Run the container
docker run -p 8080:8080 shambu2k/go-tex2pdf-server:latest
```

### Running locally

1. Ensure you have Go installed (version 1.23 or higher)
2. Install pdflatex and required LaTeX packages:
   ```bash
   # For Ubuntu/Debian
   apt-get install texlive-latex-base texlive-fonts-recommended texlive-latex-extra
   
   # For macOS using Homebrew
   brew install --cask mactex
   
   # For macOS using tlmgr
   tlmgr install fontawesome5 enumitem marvosym framed titlesec preprint fullpage
   ```
3. Build and run the server:
   ```bash
   go build -o tex2pdf-server
   ./tex2pdf-server
   ```

## API Usage

### Convert LaTeX to PDF

**Endpoint**: `POST /convert`

**Request Body**: Raw LaTeX content

**Response**: Binary PDF file

**Example**:
```bash
curl -X POST http://localhost:8080/convert \
  -H "Content-Type: text/plain" \
  -d '\documentclass{article}\begin{document}Hello World!\end{document}' \
  --output document.pdf
```

### Get Version Information

**Endpoint**: `GET /version`

**Response**: Version string

### Root Path

**Endpoint**: `GET /`

**Response**: Welcome message with basic API usage instructions

## License

[MIT License](LICENSE)

## Acknowledgements

This project uses [gotex](https://github.com/rwestlund/gotex) for LaTeX rendering.
