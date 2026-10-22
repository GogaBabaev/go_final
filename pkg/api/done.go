package api

import (
	"database/sql"
	"errors"
	"go_final_project/pkg/db"
	"log"
	"net/http"
	"time"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")

	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "Задача не найдена",
			})
			return
		}

		log.Printf("Ошибка получения задачи: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Внутренняя ошибка сервера",
		})
		return

	}

	if task.Repeat == "" {

		err = db.DeleteTask(id)
		if err != nil {
			if errors.Is(err, db.ErrTaskNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{
					"error": "Задача не найдена",
				})
				return
			}
			log.Printf("Ошибка удаления задачи: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Внутренняя ошибка сервера",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{})
		return

	}

	now := time.Now()

	next, err := NextDate(now, task.Date, task.Repeat)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = db.UpdateDate(next, id)

	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {

			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "Задача не найдена",
			})
			return

		}

		log.Printf("Ошибка обновления даты задачи: %v", err)

		writeJSON(w, http.StatusInternalServerError, map[string]string{

			"error": "Внутренняя ошибка сервера",
		})

		return
	}
	writeJSON(w, http.StatusOK, map[string]string{})

}
