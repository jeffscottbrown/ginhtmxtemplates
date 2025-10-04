package ginhtmx

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Htmx provides functionality to render HTML templates with optional layout decoration.
type Htmx struct {
	template *template.Template
	config   HtmxConfig
}

// HtmxConfig holds configuration options for the Htmx instance.
type HtmxConfig struct {
	// LayoutTemplateName is the name of the layout template that templates will be wrapped in
	LayoutTemplateName string

	// ContentVariableName is the name of the variable in the layout template that will hold the body content
	ContentVariableName string
}

// NewHtmxWithConfig creates a new instance of Htmx with the provided HTML templates and configuration.
func NewHtmxWithConfig(template *template.Template, config HtmxConfig) *Htmx {
	return &Htmx{
		config:   config,
		template: template,
	}
}

// NewHtmx creates a new instance of Htmx with the provided HTML templates and
// configuration. The default configuration uses "layout" as the layout
// template name and "Content" as the body variable name.
func NewHtmx(template *template.Template) *Htmx {
	return &Htmx{
		config: HtmxConfig{
			LayoutTemplateName:  "layout",
			ContentVariableName: "Content",
		},
		template: template,
	}
}

// RenderTemplateWithStatus renders the specified template with the provided data and status code.
// If the request does not inlcude the "Hx-Request" header indicating this is an HTMX request
// then the template will be wrapped in the layout page.
func (htmx *Htmx) RenderTemplateWithStatus(ginContext *gin.Context, templateName string, data gin.H, status int) {
	ginContext.Status(status)
	isHTMX := ginContext.GetHeader("HX-Request") != ""

	if isHTMX {
		_ = htmx.template.ExecuteTemplate(ginContext.Writer, templateName, data)
	} else {
		_ = htmx.template.ExecuteTemplate(ginContext.Writer, htmx.config.LayoutTemplateName, gin.H{
			//nolint:gosec
			htmx.config.ContentVariableName: template.HTML(htmx.renderTemplateToString(templateName, data)),
		})
	}
}

// RenderTemplate renders the specified template with the provided data and an HTTP 200 status code.
// If the request does not inlcude the "Hx-Request" header indicating this is an HTMX request
// then the template will be wrapped in the layout page.
func (htmx *Htmx) RenderTemplate(c *gin.Context, templateName string, data gin.H) {
	htmx.RenderTemplateWithStatus(c, templateName, data, http.StatusOK)
}

func (htmx *Htmx) renderTemplateToString(name string, data any) string {
	var buf []byte

	writer := &buffer{&buf}
	_ = htmx.template.ExecuteTemplate(writer, name, data)

	return string(*writer.buf)
}

type buffer struct {
	buf *[]byte
}

func (w *buffer) Write(p []byte) (int, error) {
	*w.buf = append(*w.buf, p...)

	return len(p), nil
}
