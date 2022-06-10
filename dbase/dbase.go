package dbase

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Opt struct {
	AliasName         string
	DriverName        string
	DriverTyp         orm.DriverType
	Host              string //host <must>
	Port              string //default 3306
	User              string //username <must>
	Password          string
	DBName            string         //connect DB name <must>
	SslMode           string         //default disable
	TimeZone          *time.Location //default local
	MaxIdleConnes     int            //default 10
	MaxOpenConnes     int            //default 30
	MaxLifeTimeConnes time.Duration  //default 3600
	SyncDB            bool           //is need orm auto sync DB struct
}

func (opt *Opt) getDriverName() string {
	if opt.DriverName == "" {
		return "mysql"
	}
	return opt.DriverName
}

func (opt *Opt) getDriverTyp() orm.DriverType {
	if opt.DriverTyp == 0 {
		return orm.DRMySQL
	}
	return opt.DriverTyp
}

func (opt *Opt) getAliasName() string {
	if opt.AliasName == "" {
		return "default"
	}
	return opt.AliasName
}

func (opt *Opt) getPort() string {
	if opt.Port == "" {
		return "3306"
	}
	return opt.Port
}

func (opt *Opt) getMaxIdleConnes() int {
	if opt.MaxOpenConnes == 0 {
		return 10
	}
	return opt.MaxOpenConnes
}

func (opt *Opt) getMaxOpenConnes() int {
	if opt.MaxOpenConnes == 0 {
		return 30
	}
	return opt.MaxOpenConnes
}

func (opt *Opt) getMaxLifeTimeConnes() time.Duration {
	if opt.MaxLifeTimeConnes == 0 {
		return 3600
	}
	return opt.MaxLifeTimeConnes
}

func (opt *Opt) getTimeZone() *time.Location {
	if opt.TimeZone == nil {
		return time.Local
	}
	return opt.TimeZone
}

func (opt *Opt) getSslMode() string {
	if opt.SslMode == "" {
		return "disable"
	}
	return opt.SslMode
}

// Init 初始化orm数据库连接池
func Init(opts ...Opt) error {
	//TODO init orm pool
	for _, opt := range opts {
		err := orm.RegisterDriver(opt.getDriverName(), opt.getDriverTyp())
		if err != nil {
			return err
		}
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", opt.User, opt.Password, opt.Host, opt.getPort(), opt.DBName)
		aliasName := opt.getAliasName()
		err = orm.RegisterDataBase(aliasName, opt.getDriverName(), dsn,
			orm.MaxIdleConnections(opt.getMaxIdleConnes()),
			orm.MaxOpenConnections(opt.getMaxOpenConnes()),
			orm.ConnMaxLifetime(opt.getMaxLifeTimeConnes()))
		if err != nil {
			return err
		}
		err = orm.SetDataBaseTZ(aliasName, opt.getTimeZone())
		if err != nil {
			return err
		}
		if opt.SyncDB {
			err = orm.RunSyncdb(opt.getAliasName(), false, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type DB struct {
	orm.Ormer
}

// Orm 获取一个orm，orm本身有缓存
func Orm(aliasName ...string) *DB {
	name := "default"
	if len(aliasName) != 0 {
		name = aliasName[0]
	}
	return &DB{orm.NewOrmUsingDB(name)}
}
