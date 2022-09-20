package db

import (
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"habiko-go/pkg/models"
)

type DB struct {
	Connection          *gorm.DB
	TimescaleConnection *gorm.DB
}

func jsonMarshal(value interface{}) datatypes.JSON {
	body, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return body
}

func dbConn(t models.Transaction, dbConn *gorm.DB) *gorm.DB {
	if t == nil {
		return dbConn
	}
	return t.(*Transaction).GetTx()
}

func getConnectionString(dbType string) string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		viper.GetString(dbType+".host"),
		viper.GetString(dbType+".port"),
		viper.GetString(dbType+".name"),
		viper.GetString(dbType+".user"),
		viper.GetString(dbType+".password"),
	)
}

func executeMigrations(conn *gorm.DB) {
	err := conn.AutoMigrate(
		&models.User{},
		&models.Tracker{},
		&models.TrackerEvent{},
		&models.SampleTracker{},
	)
	if err != nil {
		panic(fmt.Sprintf("Could not migrate: %v", err))
	}

	// execute all custom migrations
	customMigrations := append(
		[]*gormigrate.Migration{},
		getSimpleTrackerMigrations()...,
	)
	m := gormigrate.New(conn, gormigrate.DefaultOptions, customMigrations)
	if err = m.Migrate(); err != nil {
		panic(fmt.Sprintf("Could not migrate: %v", err))
	}
}

func Connect() *DB {
	conn, err := gorm.Open(postgres.Open(getConnectionString("db")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	executeMigrations(conn)

	db := &DB{Connection: conn}

	if viper.GetBool("timescale_db.is_available") {
		timescaleConn, err := gorm.Open(postgres.Open(getConnectionString("timescale_db")), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		db.TimescaleConnection = timescaleConn
	}

	return db
}
