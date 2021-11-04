package configs

var AppConfig appConfig

// 应用配置
type appConfig struct {
	LogPath string		// 日志路径
	LogFile string		// 日志名
}

var DBConfig dbConfig

// 数据库配置
type dbConfig struct {
	Database string
	Type string
	Host string
	User string
	Password string
	MaxIdleConns int
	MaxOpenConns int
}