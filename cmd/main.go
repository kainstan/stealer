package main

import (
	"tiktok-uploader/configs"
	"tiktok-uploader/internal/infra/database"
	"tiktok-uploader/internal/task"
)

func init()  {
	configs.Load()
	database.SetUp()
}

func main() {
	task.CronManger()
}
