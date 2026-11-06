package daoQuery

import (
	"Quote/db/daoPostgres/daoQuery/enums/daoOperator"
	"Quote/utils/mapUtils"
	"fmt"
)

type Filter struct {
	Field    string               `json:"field"`
	Operator daoOperator.Operator `json:"operator"`
	Value    any                  `json:"value"`
}

func (f *Filter) FillFromJson(m map[string]any) error {
	_, err := mapUtils.FieldAnyToType(m, "field", fmt.Sprintf("%T", f.Field))
	if err != nil {
		return err
	}

	operatorText := m["operator"].(string)
	m["operator"] = daoOperator.Operator(operatorText)

	return mapUtils.ToStruct_usingJsonTags(m, f)
}
