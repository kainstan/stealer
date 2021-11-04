package main

import (
	"stealer/configs"
	"stealer/internal/infra/database"
	"stealer/internal/task"
)

func init()  {
	configs.Load()
	database.SetUp()
}

func main() {
	task.CronManger()
}
