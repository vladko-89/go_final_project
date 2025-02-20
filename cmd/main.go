package main

import (
	"fmt"
	"net/http"

	"final-project/internal/date"
	"final-project/internal/static"
	"final-project/internal/tasks"
	"final-project/pkg/configs"
	"final-project/pkg/db"
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

	fmt.Printf("Server is running on port %s\n", conf.Port.Port)
	server.ListenAndServe()
}
