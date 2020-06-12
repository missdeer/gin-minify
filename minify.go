package minify

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	minifier "github.com/tdewolff/minify/v2"
)

func Minify(options ...Option) gin.HandlerFunc {
	return newMinifyHandler(options...).Handle
}

type minifyWriter struct {
	gin.ResponseWriter
	m *minifier.M
}

func (g *minifyWriter) WriteString(s string) (int, error) {
	ct := strings.Split(g.Header().Get("Content-Type"), ";")
	b, e := g.m.Bytes(ct[0], []byte(s))
	if e != nil || len(b) <= 0 {
		log.Println(e, len(b), g.Header().Get("Content-Type"), s)
		return g.ResponseWriter.Write([]byte(s))
	}
	return g.ResponseWriter.Write(b)
}

func (g *minifyWriter) Write(data []byte) (int, error) {
	ct := strings.Split(g.Header().Get("Content-Type"), ";")
	b, e := g.m.Bytes(ct[0], data)
	if e != nil || len(b) <= 0 {
		return g.ResponseWriter.Write(data)
	}
	return g.ResponseWriter.Write(b)
}

// Fix: https://github.com/mholt/caddy/issues/38
func (g *minifyWriter) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}
