package handler

import (
	"github.com/gin-gonic/gin"
	"infinite-canvas-server/crypto"
	"infinite-canvas-server/config"
	"infinite-canvas-server/model"
	"infinite-canvas-server/repository"
	"infinite-canvas-server/service"
)

type ApiConfigHandler struct {
	apiConfigRepo *repository.ApiConfigRepo
	cfg           *config.Config
}

func NewApiConfigHandler(apiConfigRepo *repository.ApiConfigRepo, cfg *config.Config) *ApiConfigHandler {
	return &ApiConfigHandler{apiConfigRepo: apiConfigRepo, cfg: cfg}
}

type SaveApiConfigInput struct {
	BaseUrl string `json:"base_url"`
	ApiKey  string `json:"api_key"`
}

func (h *ApiConfigHandler) Get(c *gin.Context) {
	claims := c.MustGet("claims").(*service.Claims)
	cfg, err := h.apiConfigRepo.FindByTenant(claims.TenantID)
	if err != nil {
		model.Fail(c, 404, "未配置 API")
		return
	}
	model.OK(c, gin.H{
		"base_url": cfg.BaseUrl,
		"has_key":  len(cfg.ApiKey) > 0,
	})
}

func (h *ApiConfigHandler) Save(c *gin.Context) {
	claims := c.MustGet("claims").(*service.Claims)
	var input SaveApiConfigInput
	if err := c.ShouldBindJSON(&input); err != nil {
		model.Fail(c, 400, "无效的请求参数")
		return
	}

	encryptedKey, err := crypto.Encrypt(h.cfg.ApiKeyEncryptKey, input.ApiKey)
	if err != nil {
		model.Fail(c, 500, "加密 API Key 失败")
		return
	}

	cfg := &model.TenantApiConfig{
		TenantID: claims.TenantID,
		BaseUrl:  input.BaseUrl,
		ApiKey:   encryptedKey,
	}
	if err := h.apiConfigRepo.Save(cfg); err != nil {
		model.Fail(c, 500, err.Error())
		return
	}
	model.OK(c, gin.H{"saved": true})
}
