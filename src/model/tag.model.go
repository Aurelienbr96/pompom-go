package model

import (
	"pompom/go/src/dto"
)

// swagger:model Tag
type Tag struct {
	// ID of Tag
	// in: int64
	ID int64 `db:"id" json:"id"`
	// Name of Tag
	// in: string
	Name string `db:"name" json:"name"`
	// Color of Tag
	// in: string
	Color  string `db:"color" json:"color"`
	UserId string `db:"userid" json:"userId"`
}

type TagService interface {
	GetAllTags(userId int) ([]Tag, error)
	CreateManyTags(tags []dto.Tag) error
	CreateNewTag(tag dto.Tag) error
}
