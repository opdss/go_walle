package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Order struct {
	Field     string `json:"field" form:"field" query:"field"`
	OrderDesc string `json:"order" form:"order" query:"order"`
}

func (o Order) OrderQuery() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if o.Field == "" {
			return db
		}
		return db.Order(clause.OrderByColumn{Column: clause.Column{Name: o.Field}, Desc: o.OrderDesc == "descend"})
	}
}
