package date

import (
	"encoding/json"
	"net/http"
	"time"

	consts "final-project/pkg/consts"
	"final-project/pkg/date_pkg"
)

type DateHandler struct{}

func NewDateHandler(router *http.ServeMux) {
	handler := &DateHandler{}

	router.HandleFunc("/api/nextdate", handler.getNextDate)
}
func (handler *DateHandler) getNextDate(res http.ResponseWriter, req *http.Request) {
	nowStr := req.FormValue("now")
	dateStr := req.FormValue("date")
	repeat := req.FormValue("repeat")

	now, err := time.Parse(consts.FormatDate, nowStr)
	if err != nil {
		http.Error(res, "invalid 'now' parameter", http.StatusBadRequest)
		return
	}

	nextDate, err := date_pkg.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(map[string]string{
			"error": "",
		})
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(nextDate))
}
