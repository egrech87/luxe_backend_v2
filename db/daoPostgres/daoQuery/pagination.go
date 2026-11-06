package daoQuery

import (
	"Quote/utils/mapUtils"
	"encoding/json"
	"math"

	"gorm.io/gorm"
)

const (
	MAX_LIMIT     int = 100
	DEFAULT_LIMIT int = 10
	DEFAULT_PAGE  int = 1
)

type Meta struct {
	Limit      int    `json:"pagesize,omitempty"`
	Page       int    `json:"page,omitempty"`
	Sort       string `json:"sort"`
	TotalItems int    `json:"total_items"`
	TotalPages int    `json:"total_pages"`
}

type Pagination struct {
	Meta    Meta `json:"pagination"`
	Data    any  `json:"data"`
	Applies bool `json:"-"` // When false, pagination is not applied to the query
}

func NewPagination() *Pagination {
	pagination := new(Pagination)
	pagination.SetLimit(DEFAULT_LIMIT)
	pagination.SetPage(DEFAULT_PAGE)
	pagination.Applies = true
	return pagination
}

func (p *Pagination) UnmarshalJSON(text []byte) error {
	type Alias Pagination
	aux := Alias{
		Applies: true, // set the default value before parsing JSON
	}
	if err := json.Unmarshal(text, &aux); err != nil {
		return err
	}
	*p = Pagination(aux)
	return nil
}

func (p *Pagination) SetLimit(limit int) {
	if limit < 1 {
		p.Meta.Limit = DEFAULT_LIMIT
	} else if limit > MAX_LIMIT {
		p.Meta.Limit = MAX_LIMIT
	} else {
		p.Meta.Limit = limit
	}
}

func (p *Pagination) SetPage(page int) {
	if page < 1 {
		p.Meta.Page = DEFAULT_PAGE
	} else {
		p.Meta.Page = page
	}
}

func (p *Pagination) SetSort(sort string) {
	p.Meta.Sort = sort
}

func (p *Pagination) GetOffset() int {
	return (p.Meta.Page - 1) * p.Meta.Limit
}

func (p *Pagination) CalculateTotalPages() int {
	return int(math.Ceil(float64(p.Meta.TotalItems) / float64(p.Meta.Limit)))
}

func (p *Pagination) FillFromJson(m map[string]any) error {
	return mapUtils.ToStruct_usingJsonTags(m, p)
}

func (*Pagination) Paginate(pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagination.Applies {
			return db.Offset(pagination.GetOffset()).Limit(pagination.Meta.Limit).Order(pagination.Meta.Sort)
		} else {
			return db.Order(pagination.Meta.Sort)
		}
	}
}
