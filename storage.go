package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s Storage) InitDatabase() {
	querySQL := `
		CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY,
		date VARCHAR(8) NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat VARCHAR(128));
	    CREATE INDEX IF NOT EXISTS indexdate ON scheduler (date);`

	log.Println("[INFO] Creating new table...")
	_, err := s.db.Exec(querySQL)
	checkError(err, "table")
}

func checkError(err error, s string) {
	if err != nil {
		log.Println("[Error] Failed: " + s)
		log.Fatal(err)
	}
	log.Println("[Info] Success: " + s)
}

func (s Storage) InsertTask(task Task) (int, error) {
	querySQL := `INSERT INTO scheduler (date, title, comment, repeat) 
	             VALUES (:date, :title, :comment, :repeat)`
	res, err := s.db.Exec(querySQL,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		return 0, fmt.Errorf("add query error: %w", err)
	}
	log.Println("[Info] Success: add query executed ")

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insertion id error: %w", err)
	}

	return int(id), nil
}

func (s Storage) SelectTasks() ([]Task, error) {
	var tasks []Task
	querySQL := `SELECT id, date, title, comment, repeat 
				 FROM scheduler 
				 ORDER BY date ASC
				 LIMIT 10`
	rows, err := s.db.Query(querySQL)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s Storage) SelectById(id int) (Task, error) {
	var task Task
	querySQL := `SELECT id, date, title, comment, repeat
 				 FROM scheduler
 				 WHERE id = :id`

	row := s.db.QueryRow(querySQL, sql.Named("id", id))

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}

	return task, nil
}

func (s Storage) UpdateTask(task Task) error {
	querySQL := `UPDATE scheduler 
				 SET date = :date, title = :title, comment = :comment, repeat = :repeat 
				 WHERE id = :id`
	res, err := s.db.Exec(querySQL,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		return fmt.Errorf("[UPDATE] failed: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("[UPDATE] RowsAffection failed: %w", err)
	}

	if n == 0 {
		return fmt.Errorf("update failed")
	}
	return nil
}

func (s Storage) DeleteTask(id int) error {
	querySQL := `DELETE FROM scheduler 
				 WHERE id = :id`
	res, err := s.db.Exec(querySQL, sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("deleteTask failed: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting RowsAffected failed: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("delete failed")
	}
	return nil
}
