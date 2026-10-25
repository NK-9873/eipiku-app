package main

import (
	"log"
	"net/http"
	"os"

	"eipiku/user-service/internal/auth"
	"eipiku/user-service/internal/db"
	"eipiku/user-service/internal/handler"
	"eipiku/user-service/internal/middleware"
	"eipiku/user-service/internal/repository"
	"eipiku/user-service/internal/service"
)

func main() {
	tokens, err := auth.NewJWT(os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Fatalf("JWT の設定が不正です: %v", err)
	}

	conn, err := db.OpenFromEnv()
	if err != nil {
		log.Fatalf("DB に接続できません: %v", err)
	}
	defer conn.Close()

	userService := service.NewUserService(repository.NewUserRepository(conn), tokens)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	handler.NewUserHandler(userService).Register(mux, middleware.RequireAuth(tokens))

	// 許可するオリジンは CORS_ALLOWED_ORIGIN で変えられる（未設定ならローカルの Next.js）
	origin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000"
	}

	log.Println("user-service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.CORS(origin)(mux)))
}
