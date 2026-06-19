package model

type GenerationRecord struct {
	BaseModel
	TenantID uint   `gorm:"index;uniqueIndex:idx_generation_records_user_record_type;not null" json:"tenant_id"`
	UserID   uint   `gorm:"index;uniqueIndex:idx_generation_records_user_record_type;not null" json:"user_id"`
	RecordID string `gorm:"size:64;not null;uniqueIndex:idx_generation_records_user_record_type" json:"record_id"`
	Type     string `gorm:"size:20;not null;index;uniqueIndex:idx_generation_records_user_record_type" json:"type"`
	Status   string `gorm:"size:20" json:"status"`
	Payload  string `gorm:"type:longtext;not null" json:"payload"`
}

func (GenerationRecord) TableName() string { return "generation_records" }
