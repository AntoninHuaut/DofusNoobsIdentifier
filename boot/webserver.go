package boot

import (
	"DofusNoobsIdentifier/internal/controller"
	"DofusNoobsIdentifier/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadWebserver() *gin.Engine {
	r := gin.Default()

	handler := controller.NewHttpHandler(usecase.NewQuests(httpClient, Sitemap))

	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
		c.Writer.WriteHeaderNow()
	})

	r.GET("/quests/:id", handler.HandleQuests)

	return r
}
