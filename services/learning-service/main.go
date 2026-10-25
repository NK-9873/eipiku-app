package main

import (
	"log"
	"net/http"
	"os"

	"eipiku/learning-service/internal/auth"
	"eipiku/learning-service/internal/db"
	"eipiku/learning-service/internal/handler"
	"eipiku/learning-service/internal/middleware"
	"eipiku/learning-service/internal/repository"
	"eipiku/learning-service/internal/service"
)

func main() {
	// user-service が発行したトークンを同じ鍵で検証する
	tokens, err := auth.NewJWT(os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Fatalf("JWT の設定が不正です: %v", err)
	}

	conn, err := db.OpenFromEnv()
	if err != nil {
		log.Fatalf("DB に接続できません: %v", err)
	}
	defer conn.Close()

	courseRepo := repository.NewCourseRepository(conn)
	wordRepo := repository.NewWordRepository(conn)
	progressRepo := repository.NewProgressRepository(conn)

	// /health 以外はログインが必要（Authorization: Bearer <JWT>）
	requireAuth := middleware.RequireAuth(tokens)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	handler.NewCourseHandler(service.NewCourseService(courseRepo, wordRepo)).Register(mux, requireAuth)
	handler.NewWordHandler(service.NewWordService(wordRepo)).Register(mux, requireAuth)
	handler.NewProgressHandler(service.NewProgressService(progressRepo, courseRepo, wordRepo)).Register(mux, requireAuth)

	// 許可するオリジンは CORS_ALLOWED_ORIGIN で変えられる（未設定ならローカルの Next.js）
	origin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000"
	}

	log.Println("learning-service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware.CORS(origin)(mux)))
}
