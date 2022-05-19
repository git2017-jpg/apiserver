package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"monitor-apiserver/pkg/config"
	"time"
)

// var _ IDataSource = new(*defaultMysqlDataSource) 也可
var _ IDataSource = (*defaultMysqlDataSource)(nil)

// IDataSource 定义数据库数据源接口，按照业务需求可以返回主库链接Master和从库链接Slave
type IDataSource interface {
	Master() *gorm.DB
	Slave() *gorm.DB
	Close()
}

// defaultMysqlDataSource 默认mysql数据源实现
type defaultMysqlDataSource struct {
	master *gorm.DB // 定义私有属性，用来持有主库链接，防止每次创建，创建后直接返回该变量。
	slave  *gorm.DB // 同上，从库链接
}

func (d *defaultMysqlDataSource) Master() *gorm.DB {
	if d.master == nil {
		panic("The [master] connection is nil, Please initialize it first.")
	}
	return d.master
}

func (d *defaultMysqlDataSource) Slave() *gorm.DB {
	if d.master == nil {
		panic("The [slave] connection is nil, Please initialize it first.")
	}
	return d.slave
}

func (d *defaultMysqlDataSource) Close() {
	// 关闭主库链接
	if d.master != nil {
		m, err := d.master.DB()
		if err != nil {
			m.Close()
		}
	}
	// 关闭从库链接
	if d.slave != nil {
		s, err := d.slave.DB()
		if err != nil {
			s.Close()
		}
	}
}

func NewDefaultMysql(c config.DBConfig) *defaultMysqlDataSource {
	return &defaultMysqlDataSource{
		master: connect(
			c.Username,
			c.Password,
			c.Host,
			c.Port,
			c.Dbname,
			c.MaximumPoolSize,
			c.MaximumIdleSize),
	}
}

func connect(user, password, host, port, dbname string, maxPoolSize, maxIdle int) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:                              true, // 缓存每一条sql语句，提高执行速度
		DisableForeignKeyConstraintWhenMigrating: true, // 禁止外键约束, 生产环境不建议使用外键约束
	})
	if err != nil {
		panic(err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}

	// set for db connection
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	// 设置连接池大小
	sqlDb.SetMaxOpenConns(maxPoolSize)
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDb.SetMaxIdleConns(maxIdle)
	sqlDb.SetConnMaxLifetime(time.Hour)
	return db
}
