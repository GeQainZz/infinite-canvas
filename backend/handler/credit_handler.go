package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"infinite-canvas-server/model"
	"infinite-canvas-server/repository"
	"infinite-canvas-server/service"
)

type CreditHandler struct {
	creditService *service.CreditService
	creditRepo    *repository.CreditRepo
}

func NewCreditHandler(creditService *service.CreditService, creditRepo *repository.CreditRepo) *CreditHandler {
	return &CreditHandler{creditService: creditService, creditRepo: creditRepo}
}

func (h *CreditHandler) GetBalance(c *gin.Context) {
	claims := c.MustGet("claims").(*service.Claims)
	account, err := h.creditService.GetOrCreateAccount(claims.TenantID, claims.UserID)
	if err != nil {
		model.Fail(c, 500, err.Error())
		return
	}
	model.OK(c, gin.H{
		"balance":      account.Balance,
		"total_earned": account.TotalEarned,
		"total_spent":  account.TotalSpent,
	})
}

func (h *CreditHandler) GetTransactions(c *gin.Context) {
	claims := c.MustGet("claims").(*service.Claims)
	account, err := h.creditService.GetOrCreateAccount(claims.TenantID, claims.UserID)
	if err != nil {
		model.Fail(c, 500, err.Error())
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	txs, total, err := h.creditRepo.ListTransactions(account.ID, page, pageSize)
	if err != nil {
		model.Fail(c, 500, err.Error())
		return
	}
	model.OKPage(c, txs, total, page, pageSize)
}

func (h *CreditHandler) ListPricing(c *gin.Context) {
	claims := c.MustGet("claims").(*service.Claims)
	items, err := h.creditRepo.ListPricing(claims.TenantID)
	if err != nil {
		model.Fail(c, 500, err.Error())
		return
	}
	model.OK(c, items)
}

func (h *CreditHandler) SavePricing(c *gin.Context) {
	claims := c.MustGet("claims").(*service.Claims)
	var pricing model.CreditPricing
	if err := c.ShouldBindJSON(&pricing); err != nil {
		model.Fail(c, 400, "invalid request body")
		return
	}
	pricing.TenantID = claims.TenantID
	if err := h.creditRepo.SavePricing(&pricing); err != nil {
		model.Fail(c, 500, err.Error())
		return
	}
	model.OK(c, pricing)
}

func (h *CreditHandler) DeletePricing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		model.Fail(c, 400, "invalid id")
		return
	}
	if err := h.creditRepo.DeletePricing(uint(id)); err != nil {
		model.Fail(c, 500, err.Error())
		return
	}
	model.OK(c, nil)
}

func (h *CreditHandler) EstimateCost(c *gin.Context) {
	claims := c.MustGet("claims").(*service.Claims)
	modelName := c.Query("model")
	if modelName == "" {
		model.Fail(c, 400, "model is required")
		return
	}
	pricing, err := h.creditRepo.FindPricing(claims.TenantID, modelName)
	if err != nil {
		model.OK(c, gin.H{"credits_per_unit": 0, "note": "no pricing configured"})
		return
	}
	model.OK(c, gin.H{
		"model":            pricing.Model,
		"credits_per_unit": pricing.CreditsPerUnit,
		"unit_type":        pricing.UnitType,
	})
}
