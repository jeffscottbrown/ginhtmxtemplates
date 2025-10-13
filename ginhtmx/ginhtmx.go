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

	// ModelDecorator is an optional interface that can be implemented to modify the model.
	// If provided, the DecorateModel method will be called before rendering any templates.
	ModelDecorator ModelDecorator
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
			ModelDecorator:      nil,
		},
		template: template,
	}
}

// RenderWithStatus renders the specified templates with the provided data, concatenates the
// results and then writes that to the response with the provided status code.
// The templates are rendered and concatenated together in the order they are provided.
// If the request does not inlcude the "Hx-Request" header indicating this is an HTMX request
// then the contents will be wrapped in the layout page.
func (htmx *Htmx) RenderWithStatus(ginContext *gin.Context, data gin.H, status int, templateNames ...string) {
	ginContext.Status(status)
	isHTMX := ginContext.GetHeader("HX-Request") != ""

	if htmx.config.ModelDecorator != nil {
		htmx.config.ModelDecorator.DecorateModel(ginContext, &data)
	}

	// Concatenate the rendered templates
	var content string
	for _, name := range templateNames {
		content += htmx.renderTemplateToString(name, data)
	}

	if isHTMX {
		ginContext.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
	} else {
		//nolint:gosec
		data[htmx.config.ContentVariableName] = template.HTML(content)
		_ = htmx.template.ExecuteTemplate(ginContext.Writer, htmx.config.LayoutTemplateName, data)
	}
}

// Render renders the specified templates with the provided data, concatenates the
// results and then writes that to the response with a 200 status code.
// The templates are rendered and concatenated together in the order they are provided.
// If the request does not inlcude the "Hx-Request" header indicating this is an HTMX request
// then the contents will be wrapped in the layout page.
func (htmx *Htmx) Render(c *gin.Context, data gin.H, templateNames ...string) {
	htmx.RenderWithStatus(c, data, http.StatusOK, templateNames...)
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
