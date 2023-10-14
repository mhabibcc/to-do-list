package main

import (
	"database/sql"
	"fmt"
	"to-do-list/internal/config"
	handler_http "to-do-list/internal/handler/http/task"
	repo "to-do-list/internal/repo/task"
	usecase "to-do-list/internal/usecase/task"
	redis_client "to-do-list/pkg/redis"

	_ "github.com/lib/pq"
)

func startApp(cfg *config.Config) error {

	db_credential := fmt.Sprintf(cfg.Database.Credential, cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DbName)

	db, err := sql.Open(cfg.Database.Driver, db_credential)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	redis := redis_client.NewRedisClient(cfg.Redis.Host, cfg.Redis.Password)

	taskRepo := repo.NewTaskRepository(db, redis)

	taskUseCase := usecase.NewUseCase(taskRepo)

	taskHandler := handler_http.NewHandler(taskUseCase)

	router := newRoutes(taskHandler)

	return startServer(router, cfg)
}
