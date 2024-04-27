package services

import (
	"log"
	dto "pompom/go/src/dto"
	model "pompom/go/src/model"

	"github.com/jmoiron/sqlx"
)

type TagDb struct {
	DB *sqlx.DB
}

func NewTagService(db *sqlx.DB) model.TagService {
	return &TagDb{DB: db}
}

func (c *TagDb) GetAllTags() ([]model.Tag, error) {
	tags := []model.Tag{}
	err := c.DB.Select(&tags, "SELECT * FROM tag")
	if err != nil {
		log.Printf("Error during GetAllTags query: %v", err)
		log.Printf("Tags error: %+v", tags)
		return nil, err
	}
	return tags, err
}

func (c *TagDb) CreateManyTags(tags []dto.Tag) error {
	tx, err := c.DB.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed("INSERT INTO tag (name, color) VALUES (:name, :color)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, task := range tags {
		createdTag, err := stmt.Exec(task)
		if err != nil {
			tx.Rollback()
			return err
		}
		rowsAffected, err := createdTag.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}
		log.Printf("New task  %+v created, rows affected: %d", task, rowsAffected)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
