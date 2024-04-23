package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TaskService struct {
	storage *Storage
}

func NewTaskService(storage *Storage) TaskService {
	return TaskService{storage: storage}
}

func (t TaskService) TasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	tasksFromDB, _ := t.storage.SelectTasks()
	dtos := TasksToDto(tasksFromDB)
	response := ResponseTasks{Dtos: dtos}
	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("json Marshal:", err)

		return
	}
	log.Println("[Info] Success: tasks from DB are given")
	w.Write(responseBody)
}

func (t TaskService) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		t.addTask(w, r)
	case http.MethodGet:
		t.getTask(w, r)
	case http.MethodPut:
		t.editTask(w, r)
	case http.MethodDelete:
		t.removeTask(w, r)
	}
}

func (t TaskService) removeTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, `{"error":"wrong id"}`, http.StatusBadRequest)
		log.Println("Atoi:", err)

		return
	}

	err = t.storage.DeleteTask(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)

		return
	}

	w.Write([]byte(`{}`))
}

func (t TaskService) editTask(w http.ResponseWriter, r *http.Request) {

	var inDTO DTO
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewDecoder(r.Body).Decode(&inDTO); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("json Decoder:", err)

		return
	}

	task, err := DtoToTask(inDTO)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("DtoToTask:", err)

		return
	}

	if task.ID == 0 {
		http.Error(w, `{"error":"wrong id"}`, http.StatusBadRequest)

		return
	}

	err = t.storage.UpdateTask(task)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("UpdateTask:", err)

		return
	}

	w.Write([]byte(`{}`))
}

func (t TaskService) getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	num := r.FormValue("id")
	if num == "" {
		http.Error(w, `{"error":"true"}`, http.StatusBadRequest)
		log.Println("[WARN] empty param")
		return
	}

	id, err := strconv.Atoi(num)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("[WARN] Failed convertation:", err)
		return
	}

	task, err := t.storage.SelectById(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("[WARN] Failed convertation:", err)
		return
	}

	dto := TaskToDto(task)

	responseBody, err := json.Marshal(dto)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("json Marshal:", err)

		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(responseBody); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
	}
}

func (t TaskService) addTask(w http.ResponseWriter, r *http.Request) {
	var inDTO DTO

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&inDTO)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("[WARN] Failed json decoding:", err)
		return
	}

	task, err := DtoToTask(inDTO)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		log.Println("[WARN] Failed Dto-to-Task convertation:", err)
		return
	}

	log.Println("[DTO ] : " + inDTO.Date + inDTO.Title + inDTO.Comment + inDTO.Repeat)
	log.Println("[TASK] : " + task.Date + task.Title + task.Comment + task.Repeat)

	id, err := t.storage.InsertTask(task)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		log.Println("[WARN] Failed to add a task:", err)
		return
	}

	log.Println("[Info] Success: Task added with id = " + strconv.Itoa(id))
	w.Write([]byte(fmt.Sprintf(`{"id":"%d"}`, id)))
}

func (t TaskService) DoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	num := r.FormValue("id")
	id, err := strconv.Atoi(num)
	if err != nil {
		http.Error(w, `{"error":"wrong id"}`, http.StatusBadRequest)
		log.Println("Atoi:", err)

		return
	}

	err = t.getTaskDone(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)

		return
	}

	w.Write([]byte(`{}`))
}

func (t TaskService) getTaskDone(id int) error {
	task, err := t.storage.SelectById(id)
	if err != nil {
		return fmt.Errorf(`{"error":"%s"}`, err.Error())
	}

	if task.Repeat == "" {
		return t.storage.DeleteTask(id)
	}

	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format in db")
	}

	nextDate, _ := CalculateNextDate(date, time.Now(), task.Repeat)

	task.Date = nextDate.Format("20060102")

	err = t.storage.UpdateTask(task)
	if err != nil {
		return err
	}

	return nil
}
