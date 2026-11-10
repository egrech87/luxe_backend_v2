package termsAndConditionsEntity

import "Quote/db/entity/baseEntity"

// TermsAndConditionsEntity points to the PDF (stored in the files table) that currently holds the terms and conditions.
// At most one non-deleted row is expected to exist.
type TermsAndConditionsEntity struct {
	baseEntity.BaseEntity `json:"-" structs:"-"`
	FileID                uint64 `json:"file_id" structs:"file_id" gorm:"column:file_id; not null"`
}

func New() *TermsAndConditionsEntity {
	return new(TermsAndConditionsEntity)
}

func (*TermsAndConditionsEntity) TableName() string {
	return TableName()
}
func TableName() string {
	return "terms_and_conditions"
}
