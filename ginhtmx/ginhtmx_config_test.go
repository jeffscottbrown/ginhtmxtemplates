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

func (suite *GinHtmxConfigTestSuite) TestPageIsDecoratedWithNonDefaultLayout() {
	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)

	testContext.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	suite.htmx.Render(testContext, gin.H{
		"Name": "Jerry",
	}, "hello")

	suite.Equal(http.StatusOK, recorder.Code, "Expected status 200")

	doc, err := goquery.NewDocumentFromReader(recorder.Body)
	suite.Require().NoError(err, "Expected no error parsing HTML")

	greeting := doc.Find("#greeting").Text()

	suite.Equal("Hello, Jerry!", greeting)

	suite.Equal(3, doc.Find("body > div").Length())

	suite.Equal("Menu Bar Here", doc.Find("body > div").First().Text())
	suite.Equal("Footer Here", doc.Find("body > div").Last().Text())
}

func (suite *GinHtmxConfigTestSuite) SetupSuite() {
	templateContent := `
 {{define "customlayout"}}
<html>
<body>
  <div>Menu Bar Here</div>
  <div>
	{{.CustomBody}}
  </div>
  <div>Footer Here</div>
</body>
</html>
{{end}}

{{define "hello"}}
<h1 id="greeting">Hello, {{.Name}}!</h1>
{{end}}
`
	tmpl := template.Must(template.New("").Parse(templateContent))
	suite.htmx = ginhtmx.NewHtmxWithConfig(tmpl, ginhtmx.HtmxConfig{
		LayoutTemplateName:  "customlayout",
		ContentVariableName: "CustomBody",
		ModelDecorator:      nil,
	})
}

func TestGinHtmxConfigTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GinHtmxConfigTestSuite))
}

type GinHtmxConfigTestSuite struct {
	suite.Suite

	htmx *ginhtmx.Htmx
}
