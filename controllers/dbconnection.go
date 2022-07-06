package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/suzuka4316/auth-backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/gorilla/handlers"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (sv *Server) Connect(DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DbUser, DbPassword, DbHost, DbPort, DbName)
	sv.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Connect:: could not connect to the database")
	}
	log.Print("Connect:: Connected to database")

	sv.DB.AutoMigrate(&models.User{})

	sv.Router = mux.NewRouter()

	sv.SetUpRoutes()
}

func (server *Server) Run(port string) {
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	log.Fatal(http.ListenAndServe(port, handlers.CORS(credentials, methods, headers, origins)(server.Router)))
}