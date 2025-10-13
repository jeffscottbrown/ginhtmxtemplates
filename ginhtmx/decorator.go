package ginhtmx

import "github.com/gin-gonic/gin"

// ModelDecorator is an interface that can be implemented to modify the model before rendering.
// If provided, the DecorateModel method will be called before rendering any templates.
type ModelDecorator interface {
	DecorateModel(ginContext *gin.Context, model *gin.H)
}
