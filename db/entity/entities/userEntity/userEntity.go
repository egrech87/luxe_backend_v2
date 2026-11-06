package userEntity

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/components"
	"Quote/utils/mapUtils"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
)

type UserEntity struct {
	baseEntity.BaseEntity `json:"-" structs:"-"`
	components.Fullname   `json:"-" structs:"-"`
	FriendlyName          string     `json:"friendly_name" structs:"friendly_name"`
	LastLogin             *time.Time `json:"last_login" structs:"last_login"`
}

func New() *UserEntity {
	return new(UserEntity)
}

func (e *UserEntity) GetFriendlyName() string {
	return e.FriendlyName
}
func (e *UserEntity) GetLastLogin() *time.Time {
	return e.LastLogin
}
func (e *UserEntity) SetFriendlyName(value string) {
	e.FriendlyName = value
}
func (e *UserEntity) SetLastLogin(value *time.Time) {
	e.LastLogin = value
}

func (*UserEntity) TableName() string {
	return TableName()
}
func TableName() string {
	return "users"
}
func (e *UserEntity) FieldNames() []string {
	result := e.BaseEntity.FieldNames()
	return append(result, []string{"name", "surname", "friendly_name", "last_login"}...)
}

func (e *UserEntity) ToJson() map[string]any {
	json := structs.Map(e)

	for key, element := range e.BaseEntity.ToJson() {
		json[key] = element
	}
	for key, element := range e.Fullname.ToJson() {
		json[key] = element
	}

	return json
}

func (e *UserEntity) FillFromJson(json map[string]any) error {
	if err := e.BaseEntity.FillFromJson(json); err != nil {
		return err
	}

	if err := e.Fullname.FillFromJson(json); err != nil {
		return err
	}

	_, err := mapUtils.FieldAnyToType(json, "last_login", fmt.Sprintf("%T", e.LastLogin))
	if err != nil {
		return err
	}

	return mapUtils.ToStruct_usingJsonTags(json, e)
}

func (e *UserEntity) FillFromDb(dbMap map[string]any) error {
	itemMap := map[string]any{}
	keyPrefix := e.TableName() + "_"
	for key, value := range dbMap {
		if strings.HasPrefix(key, keyPrefix) {
			newKey := strings.TrimPrefix(key, keyPrefix)
			itemMap[newKey] = value
		}
	}

	if err := e.FillFromJson(itemMap); err != nil {
		return err
	}

	return nil
}

func (e *UserEntity) Clone() (baseEntity.BaseEntityI, error) {
	clone := New()
	err := clone.FillFromJson(e.ToJson())
	return clone, err
}

func (e *UserEntity) Validate() error {
	if err := e.Fullname.Validate(); err != nil {
		return err
	}

	if e.FriendlyName == "" {
		return errors.Errorf("Friendly name cannot be empty")
	}

	return nil
}

func (e *UserEntity) GetUpdates(oldItem baseEntity.BaseEntityI) map[string]any {
	oldUser := oldItem.(*UserEntity)
	updates := map[string]any{}

	for key, element := range e.Fullname.GetUpdates(oldUser.Fullname) {
		updates[key] = element
	}

	if e.FriendlyName != oldUser.FriendlyName {
		updates["friendly_name"] = e.FriendlyName
	}
	if e.LastLogin != oldUser.LastLogin {
		updates["last_login"] = e.LastLogin
	}

	return updates
}
