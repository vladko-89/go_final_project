package main

import (
	"final-project/internal/date"
	"final-project/internal/static"
	"final-project/internal/tasks"
	"final-project/pgk/configs"
	"final-project/pgk/db"
	"fmt"
	"net/http"
)

func main() {

	conf := configs.LoadConfig()

	db, err := db.InitDb(conf.Db.Path)
	if err != nil {
		fmt.Println("Ошибка инициализации базы данных:", err)
		return
	}

	defer db.Close()

	router := http.NewServeMux()

	taskRepository := tasks.NewTaskRepository(db)

	tasks.NewTaskHandler(
		router,
		tasks.TaskHandlerDeps{TaskRepository: taskRepository},
	)
	static.NewStaticHandler(router)

	date.NewDateHandler(router)

	server := http.Server{
		Addr:    conf.Port.Port,
		Handler: router,
	}

	fmt.Println("Server is running on port 8081")
	server.ListenAndServe()
}
