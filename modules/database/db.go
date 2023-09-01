package database

import (
	"database/sql"
	"fmt"

	"biqx.com.br/acgm_agent/modules/config"
	"github.com/rs/zerolog/log"
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

	log.Info().Str("package", "Database").Msgf("Database at: %s:%d/%s", conf.Database.Host, conf.Database.Port, conf.Database.Database)

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
		log.Error().Str("package", "Database").Err(err).Msg("Failed to open database connection")
		return err
	}

	d.Conn, err = d.Db.DB()
	if err != nil {
		log.Error().Str("package", "Database").Err(err).Msg("Failed to get database handler")
		return err
	}

	err = d.Conn.Ping()
	if err != nil {
		log.Error().Str("package", "Database").Err(err).Msg("Failed to ping the database")
		return err
	}

	d.connected = true
	log.Info().Str("package", "Database").Err(err).Msg("Connected to the database")

	d.Tx, err = d.Begin()
	if err != nil {
		return err
	}

	return nil
}

func (d *Db) Begin() (*gorm.DB, error) {

	if !d.connected {
		return nil, nil
	}

	tx := d.Db.Begin()
	if d.Db.Error != nil {
		log.Error().Str("package", "Database").Err(d.Db.Error).Msg("Failed to start database transaction")
		return nil, d.Db.Error
	}

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
		log.Error().Str("package", "Database").Err(err).Msg("Failed to commit database transaction")
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
		log.Error().Str("package", "Database").Err(err).Msg("Failed to rollback database transaction")
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

	log.Info().Str("package", "Database").Msg("Disconnected from the database")

	return nil
}
