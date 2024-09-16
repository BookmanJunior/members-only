package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bookmanjunior/members-only/api"
	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/internal/cloud"
	"github.com/bookmanjunior/members-only/internal/hub"
	"github.com/bookmanjunior/members-only/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

const Red = "\033[31m"
const White = "\033[97m"

func main() {

	infoLog := log.New(os.Stdout, White+"INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, Red+"ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	addr := flag.String("addr", ":3000", "HTTP network address")

	flag.Parse()
	if err := godotenv.Load(); err != nil {
		errorLog.Fatal(err)
	}

	conncstring := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable", os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := sql.Open("pgx", conncstring)

	if err != nil {
		errorLog.Fatal(err.Error())
	}

	cloudinaryConnectString := fmt.Sprintf("cloudinary://%v:%v@%v", os.Getenv("Cloudinary_KEY"),
		os.Getenv("Cloudinary_Secret"), os.Getenv("Cloudinary_Name"))
	var Cloudinary cloud.Cloudinary
	err = Cloudinary.Open(cloudinaryConnectString)

	if err != nil {
		errorLog.Fatal(err)
	}

	hub := hub.CreateNewHub()
	go hub.Run()

	app := &config.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		Users:    &models.UserModel{DB: db},
		Messages: &models.MessageModel{DB: db},
		Avatar:   &models.AvatarModel{DB: db},
		Cloud:    &Cloudinary,
		Hub:      hub,
	}

	server := &http.Server{
		Handler:  api.Router(app),
		ErrorLog: errorLog,
		Addr:     *addr,
	}

	defer db.Close()

	infoLog.Printf("Listening on port %s\n", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}
