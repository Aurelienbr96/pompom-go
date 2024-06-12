package services

import (
	model "pompom/go/src/model"

	"github.com/jmoiron/sqlx"
)

type db struct {
	DB *sqlx.DB
}

type StatisticService interface {
	GetStatistic(userId int) ([]model.Statistic, error)
}

func NewStatistic(sqlxDb *sqlx.DB) StatisticService {
	return &db{DB: sqlxDb}
}

func (c *db) GetStatistic(userId int) ([]model.Statistic, error) {
	query := `
	SELECT 
		tag.id as id, 
		tag.name, 
		tag.color,
		SUM(task.duration) as Total_duration 
		FROM 
		task 
		LEFT JOIN 
				tagtotask 
		ON 
				task.id = tagtotask.taskid
		LEFT JOIN
				tag
		ON
			tagtotask.tagid = tag.id 
			WHERE task.userid = $1
		GROUP BY tag.id, tag.name
    `

	rows, err := c.DB.Queryx(query, userId)
	if err != nil {
		return []model.Statistic{}, err
	}
	defer rows.Close()

	var statistics []model.Statistic

	for rows.Next() {
		var statistic model.Statistic

		err := rows.StructScan(&statistic)
		if err != nil {
			return []model.Statistic{}, err
		}
		statistics = append(statistics, statistic)

	}

	return statistics, nil
}
