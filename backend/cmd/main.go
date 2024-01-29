package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"github.com/kiloMIA/mediscan/backend/internal/controllers"
	"github.com/kiloMIA/mediscan/backend/internal/handlers"
	"github.com/kiloMIA/mediscan/backend/internal/router"
	"github.com/kiloMIA/mediscan/backend/internal/db"
)

func main() {
	lg, err := NewLogger()
	if err != nil {
		log.Fatalf("error in creating logger: %v", err)
	}

	DB, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Controllers
	userc := controllers.NewUserController(DB, lg)
	store := sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))

	chatHub := handlers.NewHub()
    go chatHub.Run()

	// Handlers
	userh := handlers.NewUserHandler(userc, store, lg)
    chatHandler := handlers.NewChatHandler(chatHub, lg)

	router := router.NewRouter(userh, chatHandler)

	fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
	
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func init() {
	logrus.SetFormatter(logrus.StandardLogger().Formatter)
	logrus.SetReportCaller(true)
}

func NewLogger() (*logrus.Logger, error) {
	f, err := os.OpenFile("logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	logger := &logrus.Logger{
		Out:   io.MultiWriter(os.Stdout, f),
		Level: logrus.DebugLevel,
		Formatter: &prefixed.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}
	logger.SetReportCaller(true)
	return logger, nil
}