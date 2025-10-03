package ginhtmx_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/jeffscottbrown/ginhtmxtemplates/ginhtmx"
	"github.com/stretchr/testify/suite"
)

func (suite *GinHtmxTestSuite) TestPageIsDecoratedWithLayout() {
	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)

	testContext.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	suite.htmx.RenderTemplate(testContext, "hello", gin.H{
		"Name": "Jerry",
	})

	suite.Equal(http.StatusOK, recorder.Code, "Expected status 200")

	doc, err := goquery.NewDocumentFromReader(recorder.Body)
	suite.Require().NoError(err, "Expected no error parsing HTML")

	greeting := doc.Find("#greeting").Text()

	suite.Equal("Hello, Jerry!", greeting)

	suite.Equal(3, doc.Find("body > div").Length())

	suite.Equal("Menu Bar Here", doc.Find("body > div").First().Text())
	suite.Equal("Footer Here", doc.Find("body > div").Last().Text())
}

func (suite *GinHtmxTestSuite) TestPageIsNotDecoratedWithLayoutForHtmxRequest() {
	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)

	testContext.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	testContext.Request.Header.Set("Hx-Request", "true")

	suite.htmx.RenderTemplate(testContext, "hello", gin.H{
		"Name": "Jerry",
	})

	suite.Equal(http.StatusOK, recorder.Code)

	doc, err := goquery.NewDocumentFromReader(recorder.Body)
	suite.Require().NoError(err, "Expected no error parsing HTML")

	suite.Equal("Hello, Jerry!", doc.Find("#greeting").Text())
	suite.Equal(0, doc.Find("body > div").Length())
}

func (suite *GinHtmxTestSuite) SetupSuite() {
	templateContent := `
 {{define "layout"}}
<html>
<body>
  <div>Menu Bar Here</div>
  <div>
	{{.Content}}
  </div>
  <div>Footer Here</div>
</body>
</html>
{{end}}

{{define "hello"}}
<h1 id="greeting">Hello, {{.Name}}!</h1>
{{end}}

{{ define "error" }}
{{ if .Message }}
<div class="alert alert-danger" role="alert">
    {{ .Message }}
</div>
{{ end }}
{{ end }}
`
	tmpl := template.Must(template.New("").Parse(templateContent))
	suite.htmx = ginhtmx.NewHtmx(tmpl)
}

func TestGinHtmxTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GinHtmxTestSuite))
}

type GinHtmxTestSuite struct {
	suite.Suite

	htmx *ginhtmx.Htmx
}
