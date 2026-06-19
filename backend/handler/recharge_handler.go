package handler

import (
	"strconv"

	"infinite-canvas-server/model"
	"infinite-canvas-server/repository"
	"infinite-canvas-server/service"

	"github.com/gin-gonic/gin"
)

type RechargeHandler struct {
	rechargeRepo  *repository.RechargeRepo
	paymentSvc    service.PaymentGateway
	creditService *service.CreditService
}

func NewRechargeHandler(rechargeRepo *repository.RechargeRepo, paymentSvc service.PaymentGateway, creditService *service.CreditService) *RechargeHandler {
	return &RechargeHandler{rechargeRepo: rechargeRepo, paymentSvc: paymentSvc, creditService: creditService}
}

type CreateOrderInput struct {
	PayoutID string `json:"payout_id" binding:"required"`
}

func (h *RechargeHandler) CreateOrder(c *gin.Context) {
	model.Fail(c, 403, "当前暂不开放在线购买，请联系管理员")
}

func (h *RechargeHandler) ListMyOrders(c *gin.Context) {
	claims := c.MustGet("claims").(*service.Claims)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, total, err := h.rechargeRepo.ListByUser(claims.UserID, page, pageSize)
	if err != nil {
		model.Fail(c, 500, err.Error())
		return
	}
	model.OKPage(c, orders, total, page, pageSize)
}

func (h *RechargeHandler) ListPayouts(c *gin.Context) {
	payouts := service.GetDefaultPayouts()
	model.OK(c, payouts)
}
