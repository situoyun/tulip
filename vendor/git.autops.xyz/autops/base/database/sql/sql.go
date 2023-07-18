package sql

import (
	"fmt"
	"time"

	"git.autops.xyz/autops/base/database/zapgorm2"
	"git.autops.xyz/autops/base/logs"

	"gorm.io/gorm"
	// 在这里添加默认的driver
	"gorm.io/driver/mysql"
	gormlogger "gorm.io/gorm/logger"
)

var (
	defaultDBName string
	dbs           = make(map[string]*gorm.DB)
)

type Model struct {
	ID            uint      `gorm:"type:int(11) auto_increment;primary_key" json:"id,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	DeletedStatus uint8     `gorm:"column:delete_status;type:tinyint(2);default:0" json:"delete_status"`
}

type SqlAddr struct {
	Account  string
	Password string
	Addr     string
	Database string
}

//Config sql config
type Config struct {
	Name string //default poplar
	// 数据库类型 mysql,postgres,sqlite,mssql
	Driver string
	// 连接字符串
	Source string

	Addr *SqlAddr

	LogLevel gormlogger.LogLevel //打印等级
}

//NewSQL init db sql
func NewSQL(configs []*Config, defaultSQLDB string, slowThreshold int) error {
	for _, config := range configs {
		dsn := config.Source
		if dsn == "" {
			if config.Addr == nil || config.Addr.Account == "" {
				err := fmt.Errorf("invald account or source")
				logs.Errorf("open db for %s error %s", config.Source, err.Error())
				return err
			}
			//root:cmstop@tcp(mysql.services:3306)/poplar?charset=utf8&parseTime=True&loc=Local
			config.Source = config.Addr.Account + ":" + config.Addr.Password +
				"@tcp(" + config.Addr.Addr + ")/" + config.Addr.Database +
				"?charset=utf8mb4&parseTime=True&loc=Local"
			dsn = config.Source
		}

		logLevel := gormlogger.Info
		if config.LogLevel != 0 {
			logLevel = config.LogLevel
		}
		d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: zapgorm2.New(logs.GetLogger(), gormlogger.Config{
				Colorful:                  true,
				IgnoreRecordNotFoundError: false,
				SlowThreshold:             time.Duration(slowThreshold) * time.Millisecond,
				LogLevel:                  logLevel,
			}),
		})
		if err != nil {
			logs.Errorf("open db for %s error %s", config.Source, err.Error())
			return err
		}

		sqlDB, err := d.DB()
		if err != nil {
			logs.Errorf("getdb db for %s error %s", config.Source, err.Error())
			return err
		}
		sqlDB.SetMaxOpenConns(255)
		sqlDB.SetMaxIdleConns(255)
		sqlDB.SetConnMaxLifetime(10 * time.Second)

		if config.Name == "" {
			config.Name = defaultSQLDB
		}

		dbs[config.Name] = d
	}

	defaultDBName = defaultSQLDB
	return nil
}

func AddSQL(configs []*Config, slowThreshold int) error {
	for _, config := range configs {
		dsn := config.Source
		if dsn == "" {
			if config.Addr == nil || config.Addr.Account == "" {
				err := fmt.Errorf("invald account or source")
				logs.Errorf("open db for %s error %s", config.Source, err.Error())
				return err
			}
			//root:cmstop@tcp(mysql.services:3306)/poplar?charset=utf8&parseTime=True&loc=Local
			config.Source = config.Addr.Account + ":" + config.Addr.Password +
				"@tcp(" + config.Addr.Addr + ")/" + config.Addr.Database +
				"?charset=utf8mb4&parseTime=True&loc=Local"
			dsn = config.Source
		}

		if config.Name == "" {
			err := fmt.Errorf("no name")
			logs.Errorf("open db for %s error %s", config.Source, err.Error())
			return err
		}

		logLevel := gormlogger.Info
		if config.LogLevel == 0 {
			logLevel = config.LogLevel
		}
		d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: zapgorm2.New(logs.GetLogger(), gormlogger.Config{
				Colorful:                  true,
				IgnoreRecordNotFoundError: false,
				SlowThreshold:             time.Duration(slowThreshold) * time.Millisecond,
				LogLevel:                  logLevel,
			}),
		})
		if err != nil {
			logs.Errorf("open db for %s error %s", config.Source, err.Error())
			return err
		}

		sqlDB, err := d.DB()
		if err != nil {
			logs.Errorf("getdb db for %s error %s", config.Source, err.Error())
			return err
		}
		sqlDB.SetMaxOpenConns(255)
		sqlDB.SetMaxIdleConns(255)
		sqlDB.SetConnMaxLifetime(10 * time.Second)

		dbs[config.Name] = d
	}

	return nil
}

//CloseDB close db
func CloseDB() {
	for _, db := range dbs {
		sqlDB, err := db.DB()
		if err != nil {
			continue
		}
		sqlDB.Close()
	}
}

//GetDB get db
func GetDB(name ...string) *gorm.DB {
	if len(name) == 0 {
		return dbs[defaultDBName]
	}

	return dbs[name[0]]
}
