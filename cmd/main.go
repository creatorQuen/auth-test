package main

import (
	"auth-test/config"
	"auth-test/migrations"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	bin "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	var cnf config.Config
	err := envconfig.Process("", &cnf)
	if err != nil {
		log.Fatal(err)
		return
	}

	migrateDB(&cnf)

	_ = connectDB(&cnf)
	fmt.Println("Successful migrate")
}

func connectDB(conf *config.Config) (db *sql.DB) {
	psqlConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.ConfigDataBase.Host, conf.ConfigDataBase.Port, conf.ConfigDataBase.User, conf.ConfigDataBase.Password, conf.ConfigDataBase.NameDataBase)
	db, err := sql.Open("postgres", psqlConnStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func migrateDB(conf *config.Config) {
	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require",
		conf.ConfigDataBase.User,
		conf.ConfigDataBase.Password,
		conf.ConfigDataBase.Host,
		conf.ConfigDataBase.Port,
		conf.ConfigDataBase.NameDataBase,
	)

	source := bin.Resource(migrations.AssetNames(), migrations.Asset)
	driver, err := bin.WithInstance(source)
	if err != nil {
		log.Fatal(err)
	}
	migration, err := migrate.NewWithSourceInstance("go-bindata", driver, databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err = migration.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println(err)
		} else {
			log.Fatal(err)
		}
	}
}
