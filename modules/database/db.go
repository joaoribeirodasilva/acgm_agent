package database

import (
	"database/sql"
	"fmt"

	"biqx.com.br/acgm_agent/modules/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Db struct {
	dsn       string
	connected bool
	Conn      *sql.DB
	Db        *gorm.DB
	Tx        *gorm.DB
}

func New(conf *config.Config) *Db {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=%t",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Database,
		conf.Database.ParseTime,
	)
	if conf.Database.Charset != "" {
		dsn = fmt.Sprintf(
			"%s&charset=%s",
			dsn,
			conf.Database.Charset,
		)
	}

	return &Db{
		dsn:  dsn,
		Db:   nil,
		Conn: nil,
		Tx:   nil,
	}
}

func (d *Db) Connect() error {

	var err error

	d.Db, err = gorm.Open(mysql.Open(d.dsn), &gorm.Config{})
	if err != nil {
		// New error
		return err
	}

	d.Conn, err = d.Db.DB()
	if err != nil {
		// New error
		return err
	}

	err = d.Conn.Ping()
	if err != nil {
		// New error
		return err
	}

	d.connected = true

	d.Tx, err = d.Begin()
	if err != nil {
		// New error
		return err
	}

	return nil
}

func (d *Db) Begin() (*gorm.DB, error) {
	if !d.connected {
		return nil, nil
	}

	tx := d.Db.Begin()

	return tx, nil
}

func (d *Db) Commit(tx *gorm.DB) error {

	var err error
	if tx == nil {
		d.Tx.Commit()
	} else {
		tx.Commit()
	}

	err = d.Tx.Error
	if err != nil {
		return err
	}

	return nil
}

func (d *Db) Rollback(tx *gorm.DB) error {

	var err error
	if tx == nil {
		d.Tx.Rollback()
	} else {
		tx.Rollback()
	}

	err = d.Tx.Error
	if err != nil {
		return err
	}

	return nil
}

func (d *Db) Disconnect() error {

	db, err := d.Db.DB()
	if err != nil {
		return err
	}
	db.Close()
	d.connected = false
	return nil
}
