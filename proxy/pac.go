package proxy

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func serveFile(w http.ResponseWriter, _ *http.Request, file string) {
	f, err := os.Open(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("file not found: %s", file), http.StatusNotFound)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		http.Error(w, fmt.Sprintf("file stat error: %s", file), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	w.Header().Set("Accept-Ranges", "bytes")
	if _, err := io.Copy(w, f); err != nil {
		println(err)
	}
}

func (p *Proxy) servePacFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
	serveFile(w, r, p.Options.PacFile)
}
