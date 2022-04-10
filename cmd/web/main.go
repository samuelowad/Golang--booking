package main

import (
	"encoding/gob"
	"fmt"
	"github.com/samuelowad/bookings/src/driver"
	"github.com/samuelowad/bookings/src/helpers"
	"github.com/samuelowad/bookings/src/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/samuelowad/bookings/src/config"
	"github.com/samuelowad/bookings/src/handlers"
	"github.com/samuelowad/bookings/src/render"
)

const portNumber = ":8080"

var app config.AppConfig

var session *scs.SessionManager

var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()

	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("starting on part %s", portNumber))

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {

	//change to true in production
	app.InProd = false
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	infoLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	session = scs.New()

	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProd
	app.Session = session

	//connect to DB
	log.Print("connect to DB")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=postgres")

	if err != nil {
		log.Fatal("cannot connect to db, dying")
	}
	log.Print("Connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cant create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandler(repo)

	render.NewRender(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
