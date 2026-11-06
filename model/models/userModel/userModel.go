package userModel

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/userEntity"
	"Quote/model/baseModel"
	"context"
)

type UserModel struct {
	*userEntity.UserEntity
}

func New() *UserModel {
	return new(UserModel)
}

func (m *UserModel) GetEntity() baseEntity.BaseEntityI {
	return m.UserEntity
}
func (m *UserModel) SetEntity(value baseEntity.BaseEntityI) {
	m.UserEntity = value.(*userEntity.UserEntity)
}

func (m *UserModel) ToJson() map[string]any {
	return m.UserEntity.ToJson()
}

func (m *UserModel) FillFromJson(json map[string]any) error {
	if m.UserEntity == nil {
		m.UserEntity = userEntity.New()
	}

	return m.UserEntity.FillFromJson(json)
}

func (m *UserModel) Clone() (baseModel.BaseModelI, error) {
	clone := New()
	clonedEntity, err := m.UserEntity.Clone()
	if err != nil {
		return nil, err
	}

	clone.UserEntity = clonedEntity.(*userEntity.UserEntity)
	return clone, nil
}

func (m *UserModel) Validate(ctx context.Context) error {
	if err := m.GetEntity().Validate(); err != nil {
		return err
	}

	return nil
}

func (m *UserModel) GetUpdates(oldModel baseModel.BaseModelI) map[string]any {
	return m.GetEntity().GetUpdates(oldModel.GetEntity())
}
