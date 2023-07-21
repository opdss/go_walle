package db

import (
	"gorm.io/gorm"
)

const defPage = 1
const defPageSize = 100
const maxPageSize = 5000

type Paginator struct {
	Page     int `form:"page" json:"page" query:"page" validate:"omitempty"`                //页码
	PageSize int `form:"page_size" json:"page_size" query:"page_size" validate:"omitempty"` //分页数量
	infinite bool
}

func (p *Paginator) GetPage() int {
	if p.Page <= 0 {
		p.Page = defPage
	}
	return p.Page
}

func (p *Paginator) GetPageSize() int {
	if p.PageSize <= 0 || p.PageSize > maxPageSize {
		p.PageSize = defPageSize
	}
	return p.PageSize
}

func (p *Paginator) PageQuery() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.infinite {
			return db
		}
		offset := (p.GetPage() - 1) * p.GetPageSize()
		return db.Offset(offset).Limit(p.GetPageSize())
	}
}

// SetInfinite 不限制分页，即获取所有结果
func (p *Paginator) SetInfinite() {
	p.infinite = true
}
