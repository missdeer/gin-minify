package minify

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	minifier "github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

type minifyHandler struct {
	*Options
	m *minifier.M
}

func newMinifyHandler(options ...Option) *minifyHandler {
	handler := &minifyHandler{
		Options: DefaultOptions,
		m:       minifier.New(),
	}
	handler.m.AddFuncRegexp(regexp.MustCompile("^text/css"), css.Minify)
	// handler.m.AddRegexp(regexp.MustCompile("^text/html"), &html.Minifier{
	// 	KeepQuotes: true,
	// })
	handler.m.AddFuncRegexp(regexp.MustCompile("^image/svg+xml"), svg.Minify)
	handler.m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	handler.m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	handler.m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	for _, setter := range options {
		setter(handler.Options)
	}
	return handler
}

func (g *minifyHandler) Handle(c *gin.Context) {
	c.Writer.Header().Set("Transfer-Encoding", "identity")
	c.Writer = &minifyWriter{
		ResponseWriter: c.Writer,
		m:              g.m,
	}
	defer func() {
		c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
	}()
	c.Next()
}
