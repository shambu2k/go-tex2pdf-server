# Docker file for go serverless

# Use the official Golang image as a base image
FROM golang:1.23

# Install TeX Live for pdflatex and other dependencies
RUN apt-get update && apt-get install -y \
    texlive-latex-base \
    texlive-fonts-recommended \
    texlive-latex-recommended \
    texlive-latex-extra \
    texlive-font-utils \
    texlive-fonts-extra \
    perl \
    wget \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install additional LaTeX packages using tlmgr (TeX Live Package Manager)
RUN wget http://mirror.ctan.org/systems/texlive/tlnet/install-tl-unx.tar.gz && \
    tar -xzf install-tl-unx.tar.gz && \
    cd install-tl-* && \
    echo "selected_scheme scheme-basic" > texlive.profile && \
    echo "tlpdbopt_install_docfiles 0" >> texlive.profile && \
    echo "tlpdbopt_install_srcfiles 0" >> texlive.profile && \
    ./install-tl --profile=texlive.profile && \
    ln -s /usr/local/texlive/*/bin/* /usr/local/bin/texlive && \
    export PATH="/usr/local/bin/texlive:$PATH" && \
    cd .. && \
    rm -rf install-tl-* install-tl-unx.tar.gz

# Add TeX Live binaries to PATH permanently
ENV PATH="/usr/local/bin/texlive:$PATH"

# Install the required LaTeX packages
RUN tlmgr init-usertree && tlmgr update --self

RUN tlmgr install fontawesome5 enumitem marvosym framed titlesec preprint

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Build the application
RUN go build -o tex2pdf-server .

# Create directories for TeX inputs if needed
RUN mkdir -p /my/asset/dir /my/other/asset/dir

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./tex2pdf-server"]
