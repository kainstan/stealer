package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"stealer/configs"
	"time"
)

var (
	mysqlDB *gorm.DB
)

// 创建连接池
func mysqlConn() *gorm.DB {
	// user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	connArgs := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Asia%%2FShanghai",
		configs.DBConfig.User, configs.DBConfig.Password, configs.DBConfig.Host, configs.DBConfig.Database)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: connArgs, // data source name
		DefaultStringSize: 256, // string 类型字段的默认长度
		DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex: true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn: true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		PrepareStmt: true,
		//ConnPool: pool,
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	})

	if err != nil {
		panic("failed to connect database, " + err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to connect database, " + err.Error())
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(configs.DBConfig.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(configs.DBConfig.MaxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间
	//sqlDB.SetConnMaxLifetime(3 * time.Minute)
	// SetConnMaxIdleTime 设置了连接空闲的最大时间
	sqlDB.SetConnMaxIdleTime(3 * time.Minute)

	return db
}

func SetUp() {
	mysqlDB = mysqlConn()
}