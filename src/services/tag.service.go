package services

import (
	"fmt"
	"log"
	dto "pompom/go/src/dto"
	model "pompom/go/src/model"

	"github.com/jmoiron/sqlx"
)

type TagService interface {
	GetAllTags(userId int) ([]model.Tag, error)
	CreateManyTags(tags []dto.Tag) error
	CreateNewTag(tag model.TagToCreate) error
	DeleteTag(tagId int) (int, error)
	UpdateTag(tagToUpdate model.TagToUpdate) (model.Tag, error)
}

type TagDb struct {
	DB *sqlx.DB
}

func NewTagService(db *sqlx.DB) TagService {
	return &TagDb{DB: db}
}

func (c *TagDb) DeleteTag(tagId int) (int, error) {
	res, err := c.DB.Exec("DELETE FROM tag WHERE id = $1", tagId)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	return int(count), err
}

func (c *TagDb) UpdateTag(tagToUpdate model.TagToUpdate) (model.Tag, error) {
	var updatedTag model.Tag

	query := `
        UPDATE tag 
        SET name = :name, color = :color 
        WHERE id = :id
        RETURNING id, name, color, userid
    `

	params := map[string]interface{}{
		"name":  tagToUpdate.Name,
		"color": tagToUpdate.Color,
		"id":    tagToUpdate.Id,
	}
	rows, err := c.DB.NamedQuery(query, params)
	if err != nil {
		return model.Tag{}, err
	}
	defer rows.Close()

	// Fetch the updated row
	if rows.Next() {
		err = rows.StructScan(&updatedTag)
		if err != nil {
			return model.Tag{}, err
		}
	} else {
		return model.Tag{}, fmt.Errorf("no tag found with id %d", tagToUpdate.Id)
	}

	return updatedTag, nil
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

func (c *TagDb) CreateNewTag(tag model.TagToCreate) error {
	sqlStatement := `INSERT INTO tag (name, color, userid) VALUES (:name, :color, :userid) RETURNING *`
	namedStmt, err := c.DB.PrepareNamed(sqlStatement)
	if err != nil {
		log.Printf("Failed to prepare named statement: %s", err)
		return err
	}
	defer namedStmt.Close()
	namedStmt.Exec(tag)
	return nil
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
