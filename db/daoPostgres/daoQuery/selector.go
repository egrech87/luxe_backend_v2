package daoQuery

import (
	"Quote/utils/mapUtils"
	"fmt"
)

type Selector struct {
	ExcludeFields []string `json:"exclude_fields"`
	SubItems      []string `json:"sub_items"`
	Filters       []Filter `json:"filters"`
}

func (s *Selector) FillFromJson(json map[string]any) error {
	_, err := mapUtils.FieldAnyToType(json, "exclude_fields", fmt.Sprintf("%T", s.ExcludeFields))
	if err != nil {
		return err
	}

	_, err = mapUtils.FieldAnyToType(json, "sub_items", fmt.Sprintf("%T", s.SubItems))
	if err != nil {
		return err
	}

	if json["filters"] != nil {
		filtersAny := json["filters"].([]any)

		filters := []Filter{}
		for _, filterMap := range filtersAny {
			filter := Filter{}
			if err := filter.FillFromJson(filterMap.(map[string]any)); err != nil {
				return err
			}
			filters = append(filters, filter)
		}
		json["filters"] = filters
	}

	return mapUtils.ToStruct_usingJsonTags(json, s)
}
