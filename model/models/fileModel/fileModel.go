package fileModel

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/fileEntity"
	"Quote/model/baseModel"
	"Quote/model/enums/fileTypes"
	"context"

	"github.com/pkg/errors"
)

type FileModel struct {
	*fileEntity.FileEntity
}

func New() *FileModel {
	return new(FileModel)
}

func (m *FileModel) GetEntity() baseEntity.BaseEntityI {
	return m.FileEntity
}
func (m *FileModel) SetEntity(value baseEntity.BaseEntityI) {
	m.FileEntity = value.(*fileEntity.FileEntity)
}

func (m *FileModel) GetType() *fileTypes.FileType {
	key := m.FileEntity.GetType()
	return fileTypes.FromInt(key)
}

func (m *FileModel) SetType(value fileTypes.FileType) {
	m.Type = int(value)
}

func (m *FileModel) ToJson() map[string]any {
	return m.FileEntity.ToJson()
}

func (m *FileModel) FillFromJson(json map[string]any) error {
	if m.FileEntity == nil {
		m.FileEntity = fileEntity.New()
	}

	return m.FileEntity.FillFromJson(json)
}

func (m *FileModel) Clone() (baseModel.BaseModelI, error) {
	clone := New()
	clonedEntity, err := m.FileEntity.Clone()
	if err != nil {
		return nil, err
	}

	clone.FileEntity = clonedEntity.(*fileEntity.FileEntity)
	return clone, nil
}

func (m *FileModel) ValidateType(ctx context.Context) error {
	key := m.FileEntity.GetType()
	if !fileTypes.Exists(key) {
		return errors.Errorf("Type '%d' is invalid", key)
	}

	return nil
}

func (m *FileModel) Validate(ctx context.Context) error {
	if err := m.GetEntity().Validate(); err != nil {
		return err
	}
	// Validate the Type
	{
		if err := m.ValidateType(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (m *FileModel) GetUpdates(oldModel baseModel.BaseModelI) map[string]any {
	return m.GetEntity().GetUpdates(oldModel.GetEntity())
}
