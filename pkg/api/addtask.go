package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан заголовок задачи",
		})
		return
	}

	err = checkDate(&task)

	if err != nil {

		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	id, err := db.AddTask(&task)

	if err != nil {
		log.Printf("Ошибка сохранения задачи: %v", err)

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Внутренняя ошибка сервера",
		})
		return
	}

	idStr := strconv.FormatInt(id, 10)

	writeJSON(w, http.StatusOK, map[string]string{
		"id": idStr,
	})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {

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

	writeJSON(w, http.StatusOK, task)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")

	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	err := db.DeleteTask(id)

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

}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан заголовок задачи",
		})
		return
	}

	err = checkDate(&task)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = db.UpdateTask(&task)

	if err != nil {

		if errors.Is(err, db.ErrTaskNotFound) {

			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "Задача не найдена",
			})
			return

		}

		log.Printf("Ошибка обновления задачи: %v", err)

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Внутренняя ошибка сервера",
		})
		return

	}

	writeJSON(w, http.StatusOK, map[string]string{})
}

func checkDate(task *db.Task) error {

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)

	if err != nil {
		return err
	}

	var next string

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)

		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = next
		}
	}
	return nil

}
