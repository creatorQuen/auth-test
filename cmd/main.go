package main

import (
	"auth-test/config"
	"auth-test/infrastructure/handlers"
	"auth-test/infrastructure/repository"
	"auth-test/infrastructure/services"
	"auth-test/migrations"
	"auth-test/pkg/logging"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	bin "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	log := logging.GetLogger()

	var cnf config.Config
	err := envconfig.Process("", &cnf)
	if err != nil {
		log.Fatal(err)
		return
	}

	migrateDB(&cnf, log)
	db := connectDB(&cnf, log)

	repoAuthUser := repository.NewAuthUserRepositoryDb(db, log)
	serviceAuthUser := services.NewAuthUserService(repoAuthUser, log, cnf.Salt)
	handlerAuthUser := handlers.NewAuthUserHandler(serviceAuthUser, log)

	e := echo.New()
	e.Validator = &handlers.CustomValidator{Validator: validator.New()}
	e.POST("/register", handlerAuthUser.Register)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cnf.ListenPort), e))
}

func connectDB(conf *config.Config, log *logging.Logger) (db *sql.DB) {
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

func migrateDB(conf *config.Config, log *logging.Logger) {
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
			log.Info(err)
		} else {
			log.Fatal(err)
		}
	}
}
