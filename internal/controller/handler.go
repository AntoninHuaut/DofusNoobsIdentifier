package controller

import (
	"DofusNoobsIdentifier/domain"
	"DofusNoobsIdentifier/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type HttpHandler interface {
	HandleQuests(c *gin.Context)
}

type httpHandler struct {
	quests usecase.Quests
}

func NewHttpHandler(quests usecase.Quests) HttpHandler {
	return &httpHandler{quests: quests}
}

func (h *httpHandler) HandleQuests(c *gin.Context) {
	idRaw := c.Param("id")
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest id"})
		return
	}

	shouldRedirect := c.DefaultQuery(domain.RedirectParam, "") == "1"

	location, err := h.quests.HandleQuests(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if shouldRedirect {
		c.Redirect(http.StatusFound, *location)
		return
	}

	c.JSON(http.StatusOK, gin.H{"location": location})
}
