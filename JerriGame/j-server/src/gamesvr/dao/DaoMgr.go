package dao

import (
	// "errors"
	"fmt"
	// "net"
	// "os"
	// "syscall"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

var (
	user     = "root"
	passward = "123456"
	host     = "127.0.0.1"
	port     = "3306"
	dbname   = "j_server_game_db"
	charset  = "utf8mb4"
)

var (
	AccountTable = "JAccount"
	PlayerTable  = "JPlayer"
	GameDB       = "j_server_game_db"
)

type DBAccData struct {
	UserId   string `gorm:"column:UserId"`
	Platform string `gorm:"column:Platform"`
	AccId    uint64 `gorm:"column:AccId"`
	CreateAt uint64 `gorm:"column:CreateAt"`
}

type DaoMgr struct {
	db *gorm.DB
}

// 自定义 io.Writer
type customWriter struct{}

func (cw *customWriter) Write(p []byte) (n int, err error) {
	// 这里你可以自定义日志处理逻辑
	fmt.Printf("Custom Log: %s", p) // 输出到控制台
	return len(p), nil
}

func (cw *customWriter) Printf(format string, args ...interface{}) {
	// 实现 Printf 方法
	fmt.Printf(format, args...)
}

var (
	daoMgr *DaoMgr

	newLogger glog.Interface = glog.New(
		&customWriter{}, // io writer
		glog.Config{
			SlowThreshold: 200 * time.Millisecond, // 慢 SQL 阈值
			Colorful:      false,                  // 禁用彩色打印
			LogLevel:      glog.LogLevel(4),       // Log level
		},
	)
)

func init() {
	fmt.Println("dao init")
	daoMgr = new(DaoMgr)
}

func InitDaoMgr() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", user, passward, host, port, dbname, charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(100)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Duration(3600))
	} else {
		panic(err)
	}

	daoMgr.db = db
}

func (d *DaoMgr) GetDB() *gorm.DB {
	return d.db
}

func (d *DaoMgr) Close() {
	sqlDB, err := d.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// func IsNeedReInitDB(err error) bool {
// 	mysqlErr, ok := err.(*mysql.MySQLError)
// 	if ok {
// 		//mysql 错误,无需重连mysql服务器, 根据_.Number 可以判断具体错误类型,当前无需细化到具体类型
// 		return false
// 	}
// 	switch {
// 	case errors.Is(err, gorm.ErrInvalidDB):
// 		//gorm里只有ErrInvalidDB这种错误重连数据库是有意义的
// 		return true
// 	}
// 	//判断网络连接错误
// 	netErr, ok := err.(net.Error)
// 	if !ok {
// 		//不是网络类的错误，不重连
// 		return false
// 	}

// 	if netErr.Timeout() {
// 		return true
// 	}

// 	opErr, ok := netErr.(*net.OpError)
// 	if !ok {
// 		return false
// 	}

// 	switch t := opErr.Err.(type) {
// 	case *net.DNSError:
// 		return true
// 	case *os.SyscallError:
// 		if errno, ok := t.Err.(syscall.Errno); ok {
// 			switch errno {
// 			case syscall.ECONNREFUSED:
// 				return true
// 			case syscall.ETIMEDOUT:
// 				return true
// 			}
// 		}
// 	}

// 	return false
// }

func reInitDb() {
	InitDaoMgr()
}

func GetTableData(dbName, tablePrefix string, primary uint64, out interface{}, args ...interface{}) error {
	db := daoMgr.GetDB()

	if db == nil {
		return nil
	}

	tableName := tablePrefix

	argc := len(args)
	// for i := 0; i < 2; i++ {
	var err error
	if argc < 1 {
		// 扫描全表
		err = db.Table(tableName).Find(out).Error
	} else if argc == 1 {
		// 有 where
		if args[0] == nil {
			err = db.Table(tableName).Find(out).Error
		} else {
			err = db.Table(tableName).Where(args[0]).Find(out).Error
		}
	} else if argc == 2 {
		// 有 where offset
		if args[0] == nil {
			err = db.Table(tableName).Offset(args[1].(int)).Find(out).Error
		} else {
			err = db.Table(tableName).Where(args[0]).Offset(args[1].(int)).Find(out).Error
		}
	} else if argc > 2 {
		// 有 where offset limit
		if args[0] == nil {
			err = db.Table(tableName).Offset(args[1].(int)).Limit(args[2].(int)).Find(out).Error
		} else {
			err = db.Table(tableName).Where(args[0]).Offset(args[1].(int)).Limit(args[2].(int)).Find(out).Error
		}
	}

	// 数据库没有错误，则不需要重试
	if err == nil || err == gorm.ErrRecordNotFound {
		return nil
	}

	fmt.Printf("GetDataByTable failed. DBName: %s, Primary: %d, TableName: %s, Error: %v\n", dbName, primary, tableName, err)

	// 数据库有错误则尝试重连之后再 select
	// 	if IsNeedReInitDB(err) {
	// 		InitDaoMgr()
	// 		db = daoMgr.GetDB()
	// 		if db == nil {
	// 			return fmt.Errorf(fmt.Sprintf("reInitHandler gorm.DB failed. DBName: %s, Primary: %d", dbName, primary))
	// 		}
	// 	}
	// }
	return fmt.Errorf("unknown failed. DBName: %s, Primary: %d", dbName, primary)
}

func InsertData(dbName, tablePrefix string, primary uint64, iModel, out interface{}) error {
	db := daoMgr.GetDB()
	if db == nil {
		return fmt.Errorf("get gorm.DB failed. DBName: %s, Primary: %d", dbName, primary)
	}

	tableName := tablePrefix

	for i := 0; i < 2; i++ {
		if err := db.Table(tableName).Model(iModel).Create(out).Error; err != nil {
			fmt.Printf("InsertData failed. DBName: %s, Primary: %d, TableName: %s, Error: %v\n", dbName, primary, tableName, err)
			// 重试
			reInitDb()
			db = daoMgr.GetDB()
		} else {
			// 数据库没有错误，不需要重试
			return nil
		}
	}

	return fmt.Errorf("unknown failed. DBName: %s, Primary: %d", dbName, primary)
}
