package handlers

import (
	"Diplom/back/counter"
	"Diplom/back/db"
	"encoding/json/v2"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	qry := r.URL.Query()
	now := qry.Get("now")
	date := qry.Get("date")
	rule := qry.Get("repeat")
	var nowTime time.Time
	if now == "" {
		nowTime = time.Now()
	} else {
		var err error
		nowTime, err = time.Parse(counter.DateFormat, now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	newDate, err := counter.NextDate(nowTime, date, rule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(newDate))
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type IdResponse struct {
	ID string `json:"id"`
}

func writeError(w http.ResponseWriter, errText string, errCode int) {

	resp := ErrorResponse{Error: errText}
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errCode)
	w.Write(data)
}

func checkDate(task *db.Task) error {
	nowRaw := time.Now()
	now := time.Date(nowRaw.Year(), nowRaw.Month(), nowRaw.Day(), 0, 0, 0, 0, time.UTC)
	if task.Date == "" {
		task.Date = now.Format(counter.DateFormat)
	}
	t, err := time.Parse(counter.DateFormat, task.Date)
	if err != nil {
		return errors.New("invalid date format")
	}
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	var next string
	if len(task.Repeat) > 0 {
		next, err = counter.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return errors.New("invalid repeat rule: " + err.Error())
		}
	}
	if counter.AfterNow(t, now) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(counter.DateFormat)
		} else {
			task.Date = next
		}
	}
	return nil
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, "Ошибка чтения запроса: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var task db.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeError(w, "Ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		writeError(w, "Заголовок не может быть пустым", http.StatusBadRequest)
		return
	}
	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "Ошибка обавления задачи в БД: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := IdResponse{
		ID: strconv.Itoa(id),
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		writeError(w, "Ошибка сериализации ответа: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

func Init(r *chi.Mux) {
	r.Get("/api/nextdate", nextDayHandler)
	r.Post("/api/task", AddTaskHandler)
}
