package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"youapp/handler/api"
	"youapp/middleware"
	"youapp/repository"
	"youapp/service"
	"youapp/utils"

	"github.com/rs/cors"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type APIHandler struct {
	UserAPIHandler    *api.UserAPI
	ProfileAPIHandler *api.ProfileAPI
}

func main() {
	err := os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/coba")
	if err != nil {
		log.Fatalf("cannot set env: %v", err)
	}

	mux := http.NewServeMux()

	err = utils.ConnectDB()
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	db := utils.GetDBConnection()
	mux = RunServer(db, mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	fmt.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}

func RunServer(db *gorm.DB, mux *http.ServeMux) *http.ServeMux {
	minioClientConn, err := service.NewMinioClient()
	if err != nil {
		log.Fatalf("cannot connect to minio: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	proRepo := repository.NewProfileRepository(db)

	userService := service.NewUserService(userRepo)
	proService := service.NewProfileService(proRepo)

	userAPIHandler := api.NewUserAPI(userService)
	proAPIHandler := api.NewProfileAPI(proService, minioClientConn)

	apiHandler := APIHandler{
		UserAPIHandler:    userAPIHandler,
		ProfileAPIHandler: proAPIHandler,
	}

	//USER
	MuxRoute(mux, "POST", "/api/register", middleware.Post(
		http.HandlerFunc(
			apiHandler.UserAPIHandler.Register)))
	MuxRoute(mux, "POST", "/api/login", middleware.Post(
		http.HandlerFunc(
			apiHandler.UserAPIHandler.UserLogin)))

	MuxRoute(mux, "POST", "/api/createProfile",
		middleware.Post(
			middleware.Auth(
				http.HandlerFunc(apiHandler.ProfileAPIHandler.AddProfile))))

	MuxRoute(mux, "GET", "/api/getProfile",
		middleware.Get(
			middleware.Auth(
				http.HandlerFunc(apiHandler.ProfileAPIHandler.GetAllProfile),
			),
		),
		"?profile_id=",
	)

	MuxRoute(mux, "PUT", "/api/updateProfile",
		middleware.Put(
			middleware.Auth(
				http.HandlerFunc(apiHandler.ProfileAPIHandler.UpdateProfile))),
		"?profile_id=",
	)

	MuxRoute(mux, "DELETE", "/api/deleteProfile",
		middleware.Delete(
			middleware.Auth(
				http.HandlerFunc(apiHandler.ProfileAPIHandler.DeleteProfile))),
		"?profile_id=",
	)

	return mux

}

func MuxRoute(mux *http.ServeMux, method string, path string, handler http.Handler, opt ...string) {
	if len(opt) > 0 {
		fmt.Printf("[%s]: %s %v \n", method, path, opt)
	} else {
		fmt.Printf("[%s]: %s \n", method, path)
	}

	mux.Handle(path, handler)
}
