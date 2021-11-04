package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HtmlHandler struct {}

func NewHtmlHandler() *HtmlHandler {
	return &HtmlHandler{}
}

func (h *HtmlHandler) RedirectIndex(c *gin.Context) {
	c.Redirect(http.StatusFound, "/ui")
	return
}

func (h *HtmlHandler) Index(c *gin.Context) {
	c.Header("content-type", "text/html;charset=utf-8")
	c.String(200, string(Html))
	return
}