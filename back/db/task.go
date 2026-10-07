package db

import (
	"Diplom/back/dCounter"
	"database/sql"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES(:date, :title, :comment, :repeat)`

	res, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func Tasks(n int) ([]*Task, error) {
	query := `SELECT * FROM scheduler ORDER BY date LIMIT :queue`
	res, err := DB.Query(query,
		sql.Named("queue", n),
	)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	tasks := make([]*Task, 0)
	for res.Next() {
		task := &Task{}
		if err := res.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func SearchTaskByKey(key string) ([]*Task, error) {
	var query string
	var searchParam sql.NamedArg
	date, err := time.Parse("02.01.2006", key)
	if err == nil {
		query = `SELECT * FROM scheduler WHERE date=:date ORDER BY date `
		searchKey := date.Format(dCounter.DateFormat)
		searchParam = sql.Named("date", searchKey)

	} else {
		query = `SELECT * FROM scheduler WHERE title LIKE :key OR comment LIKE :key ORDER BY date`
		searchKey := "%" + key + "%"
		searchParam = sql.Named("key", searchKey)

	}
	res, err := DB.Query(query, searchParam)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	tasks := make([]*Task, 0)
	for res.Next() {
		task := &Task{}
		if err := res.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
