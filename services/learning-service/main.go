package main

import (
	"log"
	"net/http"
	"os"

	"eipiku/learning-service/internal/db"
	"eipiku/learning-service/internal/handler"
	"eipiku/learning-service/internal/middleware"
	"eipiku/learning-service/internal/repository"
	"eipiku/learning-service/internal/service"
)

func main() {
	conn, err := db.OpenFromEnv()
	if err != nil {
		log.Fatalf("DB に接続できません: %v", err)
	}
	defer conn.Close()

	courseHandler := handler.NewCourseHandler(service.NewCourseService(repository.NewCourseRepository(conn)))

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("GET /courses", courseHandler.List)

	// 許可するオリジンは CORS_ALLOWED_ORIGIN で変えられる（未設定ならローカルの Next.js）
	origin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000"
	}

	log.Println("learning-service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.CORS(origin)(mux)))
}
