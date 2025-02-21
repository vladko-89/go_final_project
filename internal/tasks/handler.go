package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"final-project/pkg/consts"
	"final-project/pkg/date_pkg"
	response "final-project/pkg/response"
)

type TaskHandlerDeps struct {
	TaskRepository *TaskRepository
}

type TaskHandler struct {
	TaskRepository *TaskRepository
}

func NewTaskHandler(
	router *http.ServeMux,
	deps TaskHandlerDeps,
) {
	handler := &TaskHandler{
		TaskRepository: deps.TaskRepository,
	}

	router.HandleFunc("/api/task", handler.GetById())
	router.HandleFunc("POST /api/task", handler.Create())
	router.HandleFunc("DELETE /api/task", handler.Delete())
	router.HandleFunc("PUT /api/task", handler.Update())
	router.HandleFunc("POST /api/task/done", handler.Done())
	router.HandleFunc("/api/tasks", handler.List())
}

func (handler *TaskHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			response.SendError(w, http.StatusBadRequest, err)
			return
		}

		if task.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Поле 'title' является обязательным",
			})
			return
		}

		var err error

		if strings.ToLower(task.Date) == consts.TODAY {
			fmt.Println("today")
			task.Date = time.Now().Format(consts.FormatDate)
		}

		if task.Date == strings.ToLower(consts.TODAY) {
			task.Date = time.Now().Add(24 * time.Hour).Format(consts.FormatDate)
		}

		// Если повторение не указано
		if task.Repeat == "" {
			// Парсим дату задачи, если она указана
			var taskDate time.Time
			if task.Date != "" {
				taskDate, err = time.Parse(consts.FormatDate, task.Date)
				if err != nil {
					response.SendError(w, http.StatusBadRequest, err)
					return
				}
				task.Date = time.Now().Format(consts.FormatDate)
			}
			if taskDate.Before(time.Now()) {
				task.Date = time.Now().Format(consts.FormatDate)
			}

		} else {
			var taskDate time.Time
			taskDate, err = time.Parse(consts.FormatDate, task.Date)
			if err != nil {
				response.SendError(w, http.StatusBadRequest, err)
				return
			}
			date, err := date_pkg.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				response.SendError(w, http.StatusBadRequest, err)
				return

			}
			if taskDate.Before(time.Now()) {
				task.Date = date
			}
		}

		id, err := handler.TaskRepository.Create(&task)
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id": id,
		})
	}
}

func (handler *TaskHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := handler.TaskRepository.GetList()
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, err)
			return
		}

		var tasks []map[string]string
		for rows.Next() {
			var id, date, title, comment, repeat string
			if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
				response.SendError(w, http.StatusInternalServerError, err)
				return
			}
			tasks = append(tasks, map[string]string{
				"id":      id,
				"date":    date,
				"title":   title,
				"comment": comment,
				"repeat":  repeat,
			})
		}

		if err = rows.Err(); err != nil {
			response.SendError(w, http.StatusInternalServerError, err)
			return
		}

		if len(tasks) == 0 {
			tasks = []map[string]string{} // Возвращаем пустой список вместо nil
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"tasks": tasks,
		})
	}
}

func (handler *TaskHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Не указан идентификатор задачи",
			})
			return
		}

		var task TaskRow

		err := handler.TaskRepository.Get(id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			if err == sql.ErrNoRows {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Задача не найдена",
				})
				return
			}
			response.SendError(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)

	}
}

func (handler *TaskHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Не указан идентификатор задачи",
			})
			return
		}

		var task TaskRow

		err := handler.TaskRepository.Get(id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			if err == sql.ErrNoRows {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Задача не найдена",
				})
				return
			}
			response.SendError(w, http.StatusInternalServerError, err)
			return
		}

		err = handler.TaskRepository.Delete(id)
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, err)
			return
		}
		response.SendOk(w)
	}
}

func (handler *TaskHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task TaskRow
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Неверный формат JSON",
			})
			return
		}

		if task.ID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Поле 'id' является обязательным",
			})
			return
		}

		if task.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Поле 'title' является обязательным",
			})
			return
		}

		err := handler.TaskRepository.Update(&task)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}
		response.SendOk(w)
	}
}

func (handler *TaskHandler) Done() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Метод не поддерживается",
			})
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Не указан идентификатор задачи",
			})
			return
		}

		var task TaskRow

		err := handler.TaskRepository.Get(id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			if err == sql.ErrNoRows {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Задача не найдена",
				})
				return
			}
			response.SendError(w, http.StatusInternalServerError, err)
			return
		}

		if task.Repeat == "" {
			err = handler.TaskRepository.Delete(id)
			if err != nil {
				response.SendError(w, http.StatusInternalServerError, err)
				return
			}
			response.SendOk(w)
		}

		now := time.Now().AddDate(0, 0, 1)
		nextDate, err := date_pkg.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, err)
			return
		}

		dateFormated, err := time.Parse(consts.FormatDate, nextDate)
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, err)
			return
		}

		err = handler.TaskRepository.UpdateDate(id, dateFormated.Format(consts.FormatDate))
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, err)
			return
		}
		response.SendOk(w)
	}
}
