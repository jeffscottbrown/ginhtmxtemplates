// Package ginhtmx provides utilities for rendering templates which use HTMX
// in Gin web applications.
//
// From your gin handler, you can use the RenderTemplate or RenderTemplateWithStatus
// methods to render templates. If the request includes the "Hx-Request" header
// indicating that it is an HTMX request, then the template will be rendered as-is.
// If the request does not include that header, then the template will be wrapped in
// a layout template.  This is useful for providing a consistent layout for
// non-HTMX requests while still allowing HTMX requests to receive only the
// content defined in the requested template.
//
// By default, the layout template is named "layout" and the body content is placed
// in a variable named "Content". These names may be customized by using the
// NewHtmxWithConfig function to create your Htmx instance.  The
// NewHtmxWithConfig function takes a HtmxConfig struct which allows you to
// specify the layout template name and body variable name.
//
// Here is an example of using ginhtmx in a simple Gin application:
//
//	package server
//
//	import (
//	  "embed"
//	  "html/template"
//
//	  "github.com/gin-gonic/gin"
//	  "github.com/jeffscottbrown/ginhtmxtemplates/ginhtmx"
//	)
//
//	//go:embed templates/*.html
//	var embeddedHTMLFiles embed.FS
//
//	func Run() {
//	  router := gin.Default()
//
//	  tmpl := template.Must(template.New("").Funcs(router.FuncMap).ParseFS(embeddedHTMLFiles, "templates/*.html"))
//	  htmx := ginhtmx.NewHtmx(tmpl)
//
//	  handler := &demoHandler{htmx: htmx}
//
//	  router.GET("/", handler.home)
//	  router.GET("/about", handler.about)
//
//	  router.Run(":8080")
//	}
//
//	func (handler *demoHandler) home(c *gin.Context) {
//	  // Your handler logic here before rendering template
//	  handler.htmx.RenderTemplate(c, "home", gin.H{})
//	}
//
//	func (handler *demoHandler) about(c *gin.Context) {
//	  // Your handler logic here before rendering template
//	  handler.htmx.RenderTemplate(c, "about", gin.H{})
//	}
//
//	type demoHandler struct {
//	  htmx *ginhtmx.Htmx
//	}
//
// That server runs a simple hello word application with two routes, "/" and "/about".
// The templates are embedded in the binary using the embed package.
//
// Here is an example of a layout template named "layout" defined in the templates/layout.html file:
//
//	{{define "layout"}}
//
//	<!DOCTYPE html>
//
//	<html lang="en">
//	  <head>
//	    <script src="https://unpkg.com/htmx.org@2.0.4"></script>
//	    <title>My Website</title>
//	  </head>
//	  <body>
//	    <header>
//	      <h1>Welcome</h1>
//	      <nav>
//	        <a hx-get="/" hx-push-url="true" hx-target="#content">Home</a> |
//	        <a hx-get="/about" hx-push-url="true" hx-target="#content">About</a>
//	      </nav>
//	    </header>
//	    <main id="content">
//	      {{ .Content }}
//	    </main>
//	    <footer>
//	      <p>&copy; 2025 My Website</p>
//	    </footer>
//	  </body>
//	</html>
//	{{end}}
//
// The layout template includes the HTMX script and defines a basic HTML structure
// with a header, navigation links, a main content area, and a footer. The main content
// area uses the "Content" variable to include the body of the page.
//
// A simple "home" template defined in the templates/home.html file:
//
//	{{define "home"}}
//	<h1>Home Page</h1>
//	{{end}}
//
// A simple "about" template defined in the templates/about.html file:
//
//	{{define "about"}}
//	<h1>About Page</h1>
//	{{end}}
package ginhtmx
