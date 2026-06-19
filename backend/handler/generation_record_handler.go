package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"infinite-canvas-server/model"
	"infinite-canvas-server/repository"
)

type GenerationRecordHandler struct {
	repo *repository.GenerationRecordRepo
}

func NewGenerationRecordHandler(repo *repository.GenerationRecordRepo) *GenerationRecordHandler {
	return &GenerationRecordHandler{repo: repo}
}

type GenerationRecordSaveRequest struct {
	RecordID string      `json:"record_id"`
	Type     string      `json:"type"`
	Status   string      `json:"status"`
	Payload  interface{} `json:"payload"`
}

type GenerationRecordDeleteBatchRequest struct {
	Type string   `json:"type"`
	IDs  []string `json:"ids"`
}

func (h *GenerationRecordHandler) Save(c *gin.Context) {
	var req GenerationRecordSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.Fail(c, 400, "无效的请求参数")
		return
	}
	if req.RecordID == "" || (req.Type != "image" && req.Type != "video") || req.Payload == nil {
		model.Fail(c, 400, "缺少必要参数")
		return
	}

	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		model.Fail(c, 400, "记录内容无法序列化")
		return
	}

	record := &model.GenerationRecord{
		TenantID: c.GetUint("tenant_id"),
		UserID:   c.GetUint("user_id"),
		RecordID: req.RecordID,
		Type:     req.Type,
		Status:   req.Status,
		Payload:  string(payloadBytes),
	}
	if err := h.repo.Upsert(record); err != nil {
		model.Fail(c, 500, err.Error())
		return
	}
	model.OK(c, gin.H{"record_id": req.RecordID})
}

func (h *GenerationRecordHandler) List(c *gin.Context) {
	recordType := c.Query("type")
	if recordType != "image" && recordType != "video" {
		model.Fail(c, 400, "无效的记录类型")
		return
	}

	records, err := h.repo.ListByUser(c.GetUint("tenant_id"), c.GetUint("user_id"), recordType)
	if err != nil {
		model.Fail(c, 500, err.Error())
		return
	}

	items := make([]gin.H, 0, len(records))
	for _, record := range records {
		var payload interface{}
		if err := json.Unmarshal([]byte(record.Payload), &payload); err != nil {
			payload = gin.H{}
		}
		items = append(items, gin.H{
			"record_id":   record.RecordID,
			"type":        record.Type,
			"status":      record.Status,
			"payload":     payload,
			"created_at":  record.CreatedAt,
			"updated_at":  record.UpdatedAt,
		})
	}
	model.OK(c, items)
}

func (h *GenerationRecordHandler) Delete(c *gin.Context) {
	recordType := c.Query("type")
	recordID := c.Param("id")
	if recordID == "" || (recordType != "image" && recordType != "video") {
		model.Fail(c, 400, "无效的请求参数")
		return
	}
	if err := h.repo.Delete(c.GetUint("tenant_id"), c.GetUint("user_id"), recordType, recordID); err != nil {
		model.Fail(c, 500, err.Error())
		return
	}
	model.OK(c, nil)
}

func (h *GenerationRecordHandler) DeleteBatch(c *gin.Context) {
	var req GenerationRecordDeleteBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.Fail(c, 400, "无效的请求参数")
		return
	}
	if (req.Type != "image" && req.Type != "video") || len(req.IDs) == 0 {
		model.Fail(c, 400, "缺少必要参数")
		return
	}
	if err := h.repo.DeleteBatch(c.GetUint("tenant_id"), c.GetUint("user_id"), req.Type, req.IDs); err != nil {
		model.Fail(c, 500, err.Error())
		return
	}
	model.OK(c, nil)
}
