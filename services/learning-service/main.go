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

	courseRepo := repository.NewCourseRepository(conn)
	wordRepo := repository.NewWordRepository(conn)
	progressRepo := repository.NewProgressRepository(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	handler.NewCourseHandler(service.NewCourseService(courseRepo, wordRepo)).Register(mux)
	handler.NewWordHandler(service.NewWordService(wordRepo)).Register(mux)
	// /me 以下は X-User-ID でユーザーを特定する（4回目で JWT の検証に置き換える）
	handler.NewProgressHandler(service.NewProgressService(progressRepo, courseRepo, wordRepo)).Register(mux, middleware.UserFromHeader)

	// 許可するオリジンは CORS_ALLOWED_ORIGIN で変えられる（未設定ならローカルの Next.js）
	origin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000"
	}

	log.Println("learning-service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.CORS(origin)(mux)))
}
