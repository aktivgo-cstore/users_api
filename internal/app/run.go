package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"users_api/internal/repository"
	"users_api/internal/service/handler"
	"users_api/internal/storage/mysql"
)

var (
	port = os.Getenv("PORT")
	mySqlConnStr = os.Getenv("MYSQL_CONN_STR")
)

func Run() error {
	router := mux.NewRouter()
	mySqlConn, err := mysql.CreateConnection(mySqlConnStr)
	if err != nil {
		return fmt.Errorf("error on create mysql connection")
	}
	userRepository := repository.NewUserRepository(mySqlConn)
	h := handler.NewHandler(userRepository)

	router.HandleFunc("registration", h.HandleRegistration).Methods("POST")
	router.HandleFunc("login", h.HandleLogin).Methods("POST")
	router.HandleFunc("logout", h.HandleLogout).Methods("POST")
	router.HandleFunc("delete", h.HandleDelete).Methods("POST")
	router.HandleFunc("refresh", h.HandleRefresh).Methods("GET")
	router.HandleFunc("users", h.HandleUsers).Methods("GET")

	if err := http.ListenAndServe(port, router); err != nil {
		return fmt.Errorf("error on ListenAndServe")
	}

	return nil
}