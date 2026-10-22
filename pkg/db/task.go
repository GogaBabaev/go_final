package db

import (
	"database/sql"
	"errors"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`

	res, err := db.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

func Tasks(limit int) ([]*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date
		LIMIT ?
		`

	rows, err := db.Query(query, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {

		task := &Task{}

		err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func SearchTasks(limit int, search string) ([]*Task, error) {

	date, err := time.Parse("02.01.2006", search)

	var rows *sql.Rows

	if err == nil {

		dateStr := date.Format("20060102")

		query := `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE date = ?
			ORDER BY date
			LIMIT ?
		`

		rows, err = db.Query(query, dateStr, limit)

	} else {
		query := `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE title LIKE ? OR comment LIKE ?
			ORDER BY date
			LIMIT ?
		`

		search = "%" + search + "%"

		rows, err = db.Query(query, search, search, limit)

	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {

		task := &Task{}

		err := rows.Scan(

			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`

	row := db.QueryRow(query, id)

	task := &Task{}

	err := row.Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	query := `
		UPDATE scheduler
		SET date = ?,
			title = ?,
			comment = ?,
			repeat = ?
		WHERE id = ?
	`

	res, err := db.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {

		return ErrTaskNotFound

	}

	return nil
}

func DeleteTask(id string) error {
	query := `
	DELETE FROM scheduler
	WHERE id = ?
	`

	res, err := db.Exec(query, id)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func UpdateDate(next string, id string) error {
	query := `
		UPDATE scheduler
		SET date = ?
		WHERE id = ?
	`

	res, err := db.Exec(query, next, id)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return ErrTaskNotFound
	}

	return nil
}
