package handlers

import (
	counter "Diplom/back/Counter"
	"net/http"
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

func Init(r *chi.Mux) {
	r.Get("/api/nextdate", nextDayHandler)
}
