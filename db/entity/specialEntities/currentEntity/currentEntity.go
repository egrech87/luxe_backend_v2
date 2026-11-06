package currentEntity

import (
	"Quote/utils/anyUtils"
	"database/sql/driver"
)

type CurrentEntity struct {
	LoggedInUserId *uint64 `gorm:"type:uint;null"` //@Deprecated
}

func New() *CurrentEntity {
	return new(CurrentEntity)
}

func (e *CurrentEntity) GetLoggedInUserId() *uint64 {
	return e.LoggedInUserId
}
func (e *CurrentEntity) SetLoggedInUserId(value *uint64) {
	e.LoggedInUserId = value
}

func (*CurrentEntity) TableName() string {
	return TableName()
}
func TableName() string {
	return "current"
}

func (e *CurrentEntity) Scan(gormValue interface{}) error {
	if gormValue == nil {
		return nil
	}

	value, err := anyUtils.ToUint64(gormValue)
	if err != nil {
		return err
	}
	e.LoggedInUserId = &value

	return nil
}

// Value implements the driver Valuer interface.
func (e *CurrentEntity) Value() (driver.Value, error) {
	if e.LoggedInUserId == nil {
		return 0, nil
	}
	return *e.LoggedInUserId, nil
}
