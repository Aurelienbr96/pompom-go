package model

// swagger:model Tag
type Tag struct {
	ID int64 `db:"id" json:"id"`
	TagToCreate
}

type TagToCreate struct {
	Name   string `db:"name" json:"name" validate:"required,min=1"`
	Color  string `db:"color" json:"color" validate:"required,min=1"`
	UserId int    `db:"userid" json:"userId" validate:"required"`
}

type TagToUpdate struct {
	Id    int    `json:"id" db:"id"`
	Color string `json:"color" db:"color"`
	Name  string `json:"name" db:"name"`
}
