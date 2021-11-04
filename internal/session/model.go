package session

import (
	"tiktok-uploader/internal/infra/database"
	"time"
)

type User struct {
	database.BaseModel

	account string
	token string
	category int
	postTime *time.Time
	proxy    Proxy
}

type Proxy struct {
	ip string
	port int
	username string
	password string
}