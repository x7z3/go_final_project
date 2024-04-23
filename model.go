package main

import (
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID      int
	Date    string
	Title   string
	Comment string
	Repeat  string
}

type DTO struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type ResponseTasks struct {
	Dtos []DTO `json:"tasks"`
}

func DtoToTask(dto DTO) (Task, error) {
	if dto.Title == "" {
		return Task{}, fmt.Errorf("empty title")
	}

	y, m, d := time.Now().Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	date := today

	var err error

	if dto.Date != "" {
		date, err = time.Parse("20060102", dto.Date)
		if err != nil {
			return Task{}, fmt.Errorf("invalid date format")
		}
	}

	if date.Before(today) {
		if dto.Repeat == "" {
			date = today
		} else {
			date, err = CalculateNextDate(date, today, dto.Repeat)
			if err != nil {
				return Task{}, fmt.Errorf("can't get next date: %w", err)
			}
		}
	}

	id := 0
	if dto.ID != "" {
		id, _ = strconv.Atoi(dto.ID)
	}

	return Task{
		ID:      int(id),
		Date:    date.Format("20060102"),
		Title:   dto.Title,
		Comment: dto.Comment,
		Repeat:  dto.Repeat,
	}, nil
}

func TasksToDto(tasks []Task) []DTO {
	dtos := make([]DTO, 0, len(tasks))

	for _, task := range tasks {
		dtos = append(dtos, DTO{
			ID:      strconv.Itoa(int(task.ID)),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}
	return dtos
}

func TaskToDto(task Task) DTO {

	return DTO{
		ID:      strconv.Itoa(int(task.ID)),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}
