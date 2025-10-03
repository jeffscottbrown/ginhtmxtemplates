// Package ginhtmx provides utilities for rendering templates which use HTMX
// in Gin web applications.
//
// From your gin handler, you can use the RenderTemplate or RenderTemplateWithStatus
// methods to render your templates. If the request includes the "Hx-Request" header
// indicating that it is an HTMX request, then the template will be rendered as-is.
// If the request does not include that header, then the template will be wrapped in
// a layout template.
//
// By default, the layout template is named "layout" and the body content is placed
// in a variable named "Content". You can customize these names by using the
// NewHtmxWithConfig function to create your Htmx instance.
//
// Here is an example of how to use ginhtmx in your Gin application:
//
// . package server
//
//	 import (
//		 "embed"
//		 "html/template"
//
//	 	"github.com/gin-gonic/gin"
//		 "github.com/jeffscottbrown/ginhtmxtemplates/ginhtmx"
//	 )
//
//	 //go:embed templates/*.html
//	 var embeddedHTMLFiles embed.FS
//
//	 func Run() {
//		 router := gin.Default()
//
//		 tmpl := template.Must(template.New("").Funcs(router.FuncMap).ParseFS(embeddedHTMLFiles, "templates/*.html"))
//		 htmx := ginhtmx.NewHtmx(tmpl)
//
//		 handlers := &demoHandlers{htmx: htmx}
//
//		 router.GET("/", handlers.home)
//		 router.GET("/about", handlers.about)
//
//		 router.Run(":8080")
//	 }
//
//	 func (handlers *demoHandlers) home(c *gin.Context) {
//		 // Your handler logic here before rendering template
//		 handlers.htmx.RenderTemplate(c, "home", gin.H{})
//	 }
//
//	 func (handlers *demoHandlers) about(c *gin.Context) {
//	 	// Your handler logic here before rendering template
//	 	handlers.htmx.RenderTemplate(c, "about", gin.H{})
//	 }
//
//	 type demoHandlers struct {
//	 	htmx *ginhtmx.Htmx
//	 }
//
// That server runs a simple hello word application with two routes, "/" and "/about".
// The templates are embedded in the binary using the embed package.
//
// Here is an example of a layout template named "layout" defined in the templates/layout.html file:
//
//	{{define "layout"}}
//
// . <!DOCTYPE html>
//
//	<html lang="en">
//	  <head>
//	    <script src="https://unpkg.com/htmx.org@2.0.4"></script>
//	    <title>My Website</title>
//	  </head>
//	  <body>
//	    <header>
//	       <h1>Welcome</h1>
//	       <nav>
//	           <a hx-get="/" hx-push-url="true" hx-target="#content">Home</a> |
//	           <a hx-get="/about" hx-push-url="true" hx-target="#content">About</a>
//	       </nav>
//	    </header>
//	    <main id="content">
//	       {{ .Content }}
//	    </main>
//	    <footer>
//	       <p>&copy; 2025 My Website</p>
//	    </footer>
//	 </body>
//	</html>
//
// {{end}}
//
// The layout template includes the HTMX script and defines a basic HTML structure
// with a header, navigation links, a main content area, and a footer. The main content
// area uses the "Content" variable to include the body of the page.
//
// A simple "home" template defined in the templates/home.html file:
//
//	{{define "home"}}
//	 <h1>Home Page</h1>
//	{{end}}
//
// A simple "about" template defined in the templates/about.html file:
//
//	{{define "about"}}
//	<h1>About Page</h1>
//	{{end}}
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

	// BodyVariableName is the name of the variable in the layout template that will hold the body content
	BodyVariableName string
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
			LayoutTemplateName: "layout",
			BodyVariableName:   "Content",
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
			htmx.config.BodyVariableName: template.HTML(htmx.renderTemplateToString(templateName, data)),
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
