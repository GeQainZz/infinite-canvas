package repository

import (
	"errors"
	"infinite-canvas-server/model"

	"gorm.io/gorm"
)

type GenerationRecordRepo struct {
	db *gorm.DB
}

func NewGenerationRecordRepo(db *gorm.DB) *GenerationRecordRepo {
	return &GenerationRecordRepo{db: db}
}

func (r *GenerationRecordRepo) Upsert(record *model.GenerationRecord) error {
	var existing model.GenerationRecord
	err := r.db.Where("tenant_id = ? AND user_id = ? AND record_id = ? AND type = ?", record.TenantID, record.UserID, record.RecordID, record.Type).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(record).Error
	}
	if err != nil {
		return err
	}
	return r.db.Model(&existing).Updates(map[string]any{
		"status":  record.Status,
		"payload": record.Payload,
	}).Error
}

func (r *GenerationRecordRepo) ListByUser(tenantID uint, userID uint, recordType string) ([]model.GenerationRecord, error) {
	var records []model.GenerationRecord
	err := r.db.Where("tenant_id = ? AND user_id = ? AND type = ?", tenantID, userID, recordType).Order("updated_at DESC").Find(&records).Error
	return records, err
}

func (r *GenerationRecordRepo) Delete(tenantID uint, userID uint, recordType string, recordID string) error {
	return r.db.Where("tenant_id = ? AND user_id = ? AND type = ? AND record_id = ?", tenantID, userID, recordType, recordID).Delete(&model.GenerationRecord{}).Error
}

func (r *GenerationRecordRepo) DeleteBatch(tenantID uint, userID uint, recordType string, recordIDs []string) error {
	return r.db.Where("tenant_id = ? AND user_id = ? AND type = ? AND record_id IN ?", tenantID, userID, recordType, recordIDs).Delete(&model.GenerationRecord{}).Error
}
