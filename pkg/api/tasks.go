package api

import (
	"go_final_project/pkg/db"
	"log"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {

	search := r.FormValue("search")

	var tasks []*db.Task
	var err error

	if search == "" {
		tasks, err = db.Tasks(50)
	} else {
		tasks, err = db.SearchTasks(50, search)
	}

	if err != nil {
		log.Printf("Ошибка получения списка задач: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{

			"error": "Внутренняя ошибка сервера",
		})
		return
	}
	writeJSON(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})

}
