package handlers

import (
	"Diplom/back/dCounter"
	"Diplom/back/db"
	"database/sql"
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
		nowTime, err = time.Parse(dCounter.DateFormat, now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	newDate, err := dCounter.NextDate(nowTime, date, rule)
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

func writeJSON(w http.ResponseWriter, data any) {
	respBytes, err := json.Marshal(data)
	if err != nil {
		writeError(w, "Ошибка сериализации ответа: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

func checkDate(task *db.Task) error {
	nowRaw := time.Now()
	now := time.Date(nowRaw.Year(), nowRaw.Month(), nowRaw.Day(), 0, 0, 0, 0, time.UTC)
	if task.Date == "" {
		task.Date = now.Format(dCounter.DateFormat)
	}
	date, err := time.Parse(dCounter.DateFormat, task.Date)
	if err != nil {
		return errors.New("invalid date format")
	}
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	var next string
	if len(task.Repeat) > 0 {
		next, err = dCounter.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return errors.New("invalid repeat rule: " + err.Error())
		}
	}
	if dCounter.AfterNow(date, now) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(dCounter.DateFormat)
		} else {
			task.Date = next
		}
	}
	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
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
		writeError(w, "Ошибка добавления задачи в БД: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := IdResponse{
		ID: strconv.Itoa(id),
	}
	writeJSON(w, resp)
}

type Tasks struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("search")
	var tasks []*db.Task
	var err error
	if key != "" {
		tasks, err = db.SearchTaskByKey(key)
	} else {
		tasks, err = db.Tasks(20)
	}
	if err != nil {
		writeError(w, "Ошибка получения записи "+err.Error(), http.StatusInternalServerError)
		return
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}
	resp := Tasks{
		Tasks: tasks,
	}
	writeJSON(w, resp)
}

func taskIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "Не указан id", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, "Запись не найдена", http.StatusBadRequest)
			return
		}
		writeError(w, "server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, "Ошибка чтения запроса "+err.Error(), http.StatusBadRequest)
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
	err = db.UpdateTask(&task)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, struct{}{})
}
func Init(r *chi.Mux) {
	r.Get("/api/nextdate", nextDayHandler)
	r.Post("/api/task", addTaskHandler)
	r.Get("/api/tasks", tasksHandler)
	r.Get("/api/task", taskIdHandler)
	r.Put("/api/task", updateTaskHandler)
}
