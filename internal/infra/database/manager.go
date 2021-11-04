package database

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	idColumn = "id"
	deleteColumn = "deleted"
	createColumn = "create_time"
	updateColumn = "update_time"
)

var (
	availableCond = deleteColumn + " = " + fmt.Sprintf("%d", Available)
	deletedCond = deleteColumn + " = " + fmt.Sprintf("%d", Deleted)

	notFoundError = errors.New("record not found")
)

// Dao 基础dao
type Dao struct {
}

// Session 获取数据库连接
func (d *Dao) Session() *gorm.DB {
	session := mysqlDB.Session(&gorm.Session{
		SkipDefaultTransaction: true,
		PrepareStmt: true,
	})
	return session
}

// LongSession 获取长保持连接
func (d *Dao) LongSession() (*gorm.DB, context.CancelFunc) {
	session := d.Session()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return session.WithContext(ctx), cancel
}

// Paginate 分页封装
func (d *Dao) Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// Create 执行对象保存
func (d *Dao) Create(obj interface{}) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		return db.Create(obj)
	}
}

// Update 更新操作
func (d *Dao) Update(data interface{}) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		return db.Omit(idColumn, deleteColumn, createColumn).Save(data)
	}
}

// Delete 删除操作
func (d *Dao) Delete(data interface{}) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		return db.Model(data).
			Select(deleteColumn, updateColumn).
			Updates(data)
	}
}

// BatchDelete 批量删除操作
func (d *Dao) BatchDelete(table string, ids *[]uint, time *time.Time) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		return db.Table(table).
			Where(availableCond+ " AND id IN (?)", *ids).
			Updates(map[string]interface{}{
			deleteColumn: Deleted,
			updateColumn: *time,
			});
	}
}

// DoCreate 执行对象保存
func (d *Dao) DoCreate(obj interface{}) error {
	session := d.Session()
	tx := session.Scopes(d.Create(obj))
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 执行更新
func (d *Dao) DoUpdate(obj interface{}, exclude ...string) error {
	session := d.Session()
	tx := session.Omit(exclude...).Scopes(d.Update(obj))
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 执行删除
func (d *Dao) DoDelete(obj interface{}) error {
	session := d.Session()
	tx := session.Scopes(d.Delete(obj))
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 执行id查询
func (d *Dao) FindById(id uint, data interface{}) error {
	rs := d.Session().
		Where(idColumn+ " = ? AND " +availableCond, id).
		First(data)
	if d.IsError(rs) {
		return rs.Error
	}
	return nil
}

// 执行条件查询
func (d *Dao) FindByCond(cond, data interface{}) error {
	rs := d.Session().Where(cond).Find(data)
	if d.IsError(rs) {
		return rs.Error
	}
	return nil
}

func (d *Dao) IsError(tx *gorm.DB) bool {
	if tx.Error == nil {
		return false
	}
	if tx.Error.Error() == notFoundError.Error() {
		return false
	}
	return true
}

// 执行事务
func (d *Dao) TxExecute() func(...func(*gorm.DB) *gorm.DB) error {
	return func(fl ...func(*gorm.DB) *gorm.DB) error {
		// 开始事务
		tx := mysqlDB.Begin()
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		var db *gorm.DB
		for _, f := range fl {
			db = f(tx)
			if db.Error != nil {
				// 遇到错误时回滚事务
				tx.Rollback()
				return db.Error
			}
		}
		// 否则，提交事务
		tx.Commit()
		return nil
	}
}

// 事务保存
func (d *Dao) Transaction(fl ...func(*gorm.DB) *gorm.DB) error {
	// 开始事务
	tx := mysqlDB.Begin()
	// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
	var db *gorm.DB
	for _, f := range fl {
		db = f(tx)
		if db.Error != nil {
			// 遇到错误时回滚事务
			tx.Rollback()
			return db.Error
		}
	}
	// 否则，提交事务
	tx.Commit()
	return nil
}

// GetBranchInsertSql 获取批量添加数据sql语句
//func GetBranchInsertSql(tableName, column string, objs []interface{}) string {
//	if len(objs) == 0 {
//		return ""
//	}
//
//	buf := new(bytes.Buffer)
//	buf.WriteByte('(')
//
//	var valueList []string
//	for _, obj := range objs {
//		objV := reflect.ValueOf(obj)
//		v := "("
//		for index, i := range valueTypeList {
//			if index == fieldNum-1 {
//				v += GetFormatField(objV, index, i, "")
//			} else {
//				v += GetFormatField(objV, index, i, ",")
//			}
//		}
//		valueList = append(valueList, v)
//	}
//	buf.WriteByte(')')
//	insertSql := fmt.Sprintf("insert into `%s` (%s) values %s", tableName, column,
//		strings.Join(valueList, ",")+";")
//	return insertSql
//}

//func GetFieldValue(objV reflect.Value, index int) string {
//	field := objV.Field(index)
//	isNum := objV.Type().Bits() & int(types.IsNumeric) != 0
//
//	objV.
//	v := ""
//	if t == "string" {
//		v += fmt.Sprintf("'%s'%s", objV.Field(index).String(), sep)
//	} else if t == "uint" {
//		v += fmt.Sprintf("%d%s", objV.Field(index).Uint(), sep)
//	} else if t == "int" {
//		v += fmt.Sprintf("%d%s", objV.Field(index).Int(), sep)
//	}
//	return v
//
//}

//func BatchCreate(tx *gorm.DB, dataList []interface{}, tableName string) (err error) {
//	if len(dataList) == 0 {
//		return
//	}
//
//	var bills = make([]interface{}, 0)
//	if a == page {
//		bills = dataList[(a-1)*size:]
//	} else {
//		bills = dataList[(a-1)*size : a*size]
//	}
//	sql := GetBranchInsertSql(bills, tableName)
//	if err = tx.Exec(sql).Error; err != nil {
//		fmt.Println(fmt.Sprintf("batch create data error: %v, sql: %s, tableName: %s", err, sql, tableName))
//		return
//	}
//	return
//}
