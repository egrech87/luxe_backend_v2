package fileEntity

import (
	"Quote/db/entity/baseEntity"
	"Quote/utils/mapUtils"
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/structs"
)

type FileEntity struct {
	baseEntity.BaseEntity `json:"-" structs:"-"`
	Type                  int    `json:"type" structs:"type"`
	Data                  []byte `json:"data" structs:"data"`
}

func New() *FileEntity {
	entity := new(FileEntity)
	entity.init()
	return entity
}

func (e *FileEntity) init() {
	e.Type = -1 // Intentionally an invalid value
	e.Data = []byte("")
}

func (e *FileEntity) GetType() int {
	return e.Type
}
func (e *FileEntity) GetData() []byte {
	return e.Data
}

func (e *FileEntity) SetType(value int) {
	e.Type = value
}
func (e *FileEntity) SetData(value []byte) {
	e.Data = value
}

func (*FileEntity) TableName() string {
	return TableName()
}
func TableName() string {
	return "files"
}
func (e *FileEntity) FieldNames() []string {
	result := e.BaseEntity.FieldNames()
	return append(result, []string{"type", "data"}...)
}

func (e *FileEntity) ToJson() map[string]any {
	json := structs.Map(e)

	for key, element := range e.BaseEntity.ToJson() {
		json[key] = element
	}

	return json
}

func (e *FileEntity) FillFromJson(json map[string]any) error {
	if err := e.BaseEntity.FillFromJson(json); err != nil {
		return err
	}

	_, err := mapUtils.FieldAnyToType(json, "type", fmt.Sprintf("%T", e.Type))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "data", fmt.Sprintf("%T", e.Data))
	if err != nil {
		return err
	}

	return mapUtils.ToStruct_usingJsonTags(json, e)
}

func (e *FileEntity) FillFromDb(dbMap map[string]any) error {
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

func (e *FileEntity) Clone() (baseEntity.BaseEntityI, error) {
	clone := New()
	err := clone.FillFromJson(e.ToJson())
	return clone, err
}

func (e *FileEntity) Validate() error {
	return nil
}

func (e *FileEntity) GetUpdates(oldItem baseEntity.BaseEntityI) map[string]any {
	oldFile := oldItem.(*FileEntity)
	updates := map[string]any{}

	if e.Type != oldFile.Type {
		updates["type"] = e.Type
	}
	if !bytes.Equal(e.Data, oldFile.Data) {
		updates["data"] = e.Data
	}

	return updates
}
