package handler

import (
	"github.com/gin-gonic/gin"
	"infinite-canvas-server/model"
	"infinite-canvas-server/service"
)

type CaptchaHandler struct {
	captchaService *service.CaptchaService
}

func NewCaptchaHandler(cs *service.CaptchaService) *CaptchaHandler {
	return &CaptchaHandler{captchaService: cs}
}

func (h *CaptchaHandler) Generate(c *gin.Context) {
	id, svg := h.captchaService.Generate()
	model.OK(c, gin.H{
		"captcha_id": id,
		"svg":        svg,
	})
}
