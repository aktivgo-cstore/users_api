package app

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"users_api/internal/controllers"
	"users_api/internal/repository"
	"users_api/internal/service"
	"users_api/internal/storage/mysql"
)

var (
	port         = os.Getenv("API_PORT")
	mySqlConnStr = os.Getenv("MYSQL_CONN_STR")
)

func Run() error {
	router := mux.NewRouter()
	mySqlConn, err := mysql.CreateConnection(mySqlConnStr)
	if err != nil {
		return err
	}
	userRepository := repository.NewUserRepository(mySqlConn)
	userService := service.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	router.HandleFunc("/registration", userController.Registration).Methods("POST")
	router.HandleFunc("/activate/{link}", userController.Activate).Methods("GET")
	router.HandleFunc("/login", userController.Login).Methods("POST")
	router.HandleFunc("/logout", userController.Logout).Methods("POST")
	router.HandleFunc("/delete", userController.DeleteUser).Methods("POST")
	router.HandleFunc("/users", userController.GetUsers).Methods("GET")

	log.Println("Users api server started on port " + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		return err
	}

	return nil
}
