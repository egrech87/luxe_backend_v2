package profileModel

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/profileEntity"
	"Quote/model/baseModel"
	"context"
)

type ProfileModel struct {
	*profileEntity.ProfileEntity
}

func New() *ProfileModel {
	return new(ProfileModel)
}

func (m *ProfileModel) GetEntity() baseEntity.BaseEntityI {
	return m.ProfileEntity
}
func (m *ProfileModel) SetEntity(value baseEntity.BaseEntityI) {
	m.ProfileEntity = value.(*profileEntity.ProfileEntity)
}

func (m *ProfileModel) ToJson() map[string]any {
	return m.ProfileEntity.ToJson()
}

func (m *ProfileModel) FillFromJson(json map[string]any) error {
	if m.ProfileEntity == nil {
		m.ProfileEntity = profileEntity.New()
	}

	return m.ProfileEntity.FillFromJson(json)
}

func (m *ProfileModel) Clone() (baseModel.BaseModelI, error) {
	clone := New()
	clonedEntity, err := m.ProfileEntity.Clone()
	if err != nil {
		return nil, err
	}

	clone.ProfileEntity = clonedEntity.(*profileEntity.ProfileEntity)
	return clone, nil
}

func (m *ProfileModel) Validate(ctx context.Context) error {
	if err := m.GetEntity().Validate(); err != nil {
		return err
	}
	return nil
}

func (m *ProfileModel) GetUpdates(oldModel baseModel.BaseModelI) map[string]any {
	return m.GetEntity().GetUpdates(oldModel.GetEntity())
}
