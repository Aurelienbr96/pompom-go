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

func (c *TagDb) GetAllTags(userId int) ([]model.Tag, error) {
	tags := []model.Tag{}
	query := "SELECT * FROM tag WHERE userId = $1"
	err := c.DB.Select(&tags, query, userId)
	if err != nil {
		log.Printf("Error during GetAllTags query: %v", err)
		log.Printf("Tags error: %+v", tags)
		return nil, err
	}
	return tags, err
}

func (c *TagDb) CreateNewTag(tag dto.Tag) error {
	sqlStatement := `INSERT INTO tag (name, color) VALUES (:name, :color) RETURNING *`
	namedStmt, err := c.DB.PrepareNamed(sqlStatement)
	if err != nil {
		log.Printf("Failed to prepare named statement: %s", err)
		return err
	}
	defer namedStmt.Close()
	namedStmt.Exec(tag)
	return nil
}

/* func (c *TagDb) DeleteTag(id int) error {
	sqlStatement := `DELETE ...`
} */

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

	for _, tag := range tags {
		createdTag, err := stmt.Exec(tag)
		if err != nil {
			tx.Rollback()
			return err
		}
		rowsAffected, err := createdTag.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}
		log.Printf("New tag  %+v created, rows affected: %d", tag, rowsAffected)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
