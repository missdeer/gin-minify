package minify

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	minifier "github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
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

	for _, setter := range options {
		setter(handler.Options)
	}

	if !handler.Options.IgnoreCSS {
		handler.m.AddFuncRegexp(regexp.MustCompile("^text/css"), css.Minify)
	}
	if !handler.Options.IgnoreHTML {
		handler.m.AddRegexp(regexp.MustCompile("^text/html"), &html.Minifier{
			KeepQuotes: true,
		})
	}
	if !handler.Options.IgnoreSVG {
		handler.m.AddFuncRegexp(regexp.MustCompile("^image/svg+xml"), svg.Minify)
	}
	if !handler.Options.IgnoreJS {
		handler.m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	}
	if !handler.Options.IgnoreJSON {
		handler.m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	}
	if !handler.Options.IgnoreXML {
		handler.m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	}

	return handler
}

func (g *minifyHandler) Handle(c *gin.Context) {
	c.Writer = &minifyWriter{
		ResponseWriter: c.Writer,
		m:              g.m,
	}
	defer func() {
		c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
	}()
	c.Next()
}
