package api

import (
	"net/http"
	"time"

	"github.com/AlexFox86/scheduler/internal/service/tasks"
)

const dateFmt = "20060102"

// NextDayHandler processes the 'api/nextdate' request
func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	repeatForm := r.FormValue("repeat")
	nowForm := r.FormValue("now")
	dateForm := r.FormValue("date")

	var err error
	var now time.Time

	if nowForm == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFmt, nowForm)
		if err != nil {
			http.Error(w, "invalid 'now' parameter", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := tasks.NextDate(now, dateForm, repeatForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}
