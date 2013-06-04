package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BenLubar/forum/resource"
)

func init() {
	bootstrapLenGz := strconv.Itoa(len(resource.BootstrapCssGz))
	bootstrapLen := strconv.Itoa(len(resource.BootstrapCss))

	http.HandleFunc("/css/bootstrap.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Content-Length", bootstrapLenGz)
			w.Write(resource.BootstrapCssGz)
		} else {
			w.Header().Set("Content-Length", bootstrapLen)
			w.Write(resource.BootstrapCss)
		}
	})

	fontawesomeLenGz := strconv.Itoa(len(resource.FontawesomeCssGz))
	fontawesomeLen := strconv.Itoa(len(resource.FontawesomeCss))

	http.HandleFunc("/css/fontawesome.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Content-Length", fontawesomeLenGz)
			w.Write(resource.FontawesomeCssGz)
		} else {
			w.Header().Set("Content-Length", fontawesomeLen)
			w.Write(resource.FontawesomeCss)
		}
	})

	fonts := map[string][]byte{
		"/font/fontawesome-webfont.eot":  resource.FontawesomeWebfontEot,
		"/font/fontawesome-webfont.svg":  resource.FontawesomeWebfontSvg,
		"/font/fontawesome-webfont.ttf":  resource.FontawesomeWebfontTtf,
		"/font/fontawesome-webfont.woff": resource.FontawesomeWebfontWoff,
	}
	fontLen := make(map[string]string, len(fonts))
	for f, font := range fonts {
		fontLen[f] = strconv.Itoa(len(font))
	}
	http.HandleFunc("/font/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))

		if font, ok := fonts[r.URL.Path]; ok {
			w.Header().Set("Content-Length", fontLen[r.URL.Path])
			w.Write(font)
		} else {
			http.NotFound(w, r)
		}
	})
}
