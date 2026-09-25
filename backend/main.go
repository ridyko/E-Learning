package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"elearning-backend/handlers"
	"elearning-backend/store"
)

func main() {
	// Initialize in-memory database with pre-populated data for Pak Rio Widyatmoko
	store.InitStore()

	// Settings & Customization route
	http.HandleFunc("/api/settings/landing", handlers.EnableCORS(handlers.HandleLandingSettings))

	// Authentication & Registration routes
	http.HandleFunc("/api/auth/login", handlers.EnableCORS(handlers.LoginHandler))
	http.HandleFunc("/api/auth/register", handlers.EnableCORS(handlers.RegisterStudentHandler))
	http.HandleFunc("/api/students", handlers.EnableCORS(handlers.GetStudentsHandler))

	// Data & Content routes
	http.HandleFunc("/api/profile", handlers.EnableCORS(handlers.GetProfile))
	http.HandleFunc("/api/profile/update", handlers.EnableCORS(handlers.UpdateProfile))
	http.HandleFunc("/api/courses", handlers.EnableCORS(handlers.HandleCourses))
	http.HandleFunc("/api/courses/reset", handlers.EnableCORS(handlers.ResetCoursesHandler))
	http.HandleFunc("/api/courses/module/update", handlers.EnableCORS(handlers.UpdateModuleHandler))
	http.HandleFunc("/api/courses/", handlers.EnableCORS(handlers.HandleCourses))
	http.HandleFunc("/api/announcements", handlers.EnableCORS(handlers.HandleAnnouncements))
	http.HandleFunc("/api/quizzes", handlers.EnableCORS(handlers.GetQuizzes))
	http.HandleFunc("/api/quizzes/submit", handlers.EnableCORS(handlers.SubmitQuiz))
	http.HandleFunc("/api/quizzes/submissions", handlers.EnableCORS(handlers.GetQuizSubmissions))
	http.HandleFunc("/api/quizzes/questions", handlers.EnableCORS(handlers.AddQuizQuestion))
	http.HandleFunc("/api/assignments", handlers.EnableCORS(handlers.HandleAssignments))
	http.HandleFunc("/api/attendance", handlers.EnableCORS(handlers.HandleAttendance))

	// Live Quiz Interaktif routes
	http.HandleFunc("/api/live-quiz/create", handlers.EnableCORS(handlers.CreateLiveQuizSession))
	http.HandleFunc("/api/live-quiz/join", handlers.EnableCORS(handlers.JoinLiveQuizSession))
	http.HandleFunc("/api/live-quiz/state", handlers.EnableCORS(handlers.StateLiveQuizSession))
	http.HandleFunc("/api/live-quiz/vote", handlers.EnableCORS(handlers.VoteLiveQuizSession))
	http.HandleFunc("/api/live-quiz/sync", handlers.EnableCORS(handlers.SyncLiveQuizSession))
	http.HandleFunc("/api/live-quiz/active", handlers.EnableCORS(handlers.GetActiveLiveQuizSession))

	// Health check endpoint
	http.HandleFunc("/api/health", handlers.EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok", "app":"Portal E-Learning Pak Rio Widyatmoko, S.Kom, M.M.S.I", "institution":"ITB Swadharma"}`)
	}))

	// Serve Static Vue 3 SPA frontend from ../frontend/dist
	distDir := "../frontend/dist"
	if _, err := os.Stat(distDir); err != nil {
		distDir = "frontend/dist"
	}
	fs := http.FileServer(http.Dir(distDir))
	http.HandleFunc("/", handlers.EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		fpath := distDir + r.URL.Path
		info, err := os.Stat(fpath)
		if err != nil || info.IsDir() {
			http.ServeFile(w, r, distDir+"/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	}))

	port := ":8080"
	fmt.Printf("====================================================\n")
	fmt.Printf("  E-LEARNING BACKEND REST API SERVER (GOLANG)\n")
	fmt.Printf("  Dosen: Rio Widyatmoko, S.Kom, M.M.S.I (ITB Swadharma)\n")
	fmt.Printf("  Features: Edit Beranda, Login Dosen/Mahasiswa, Self-Reg\n")
	fmt.Printf("  Server running on http://localhost%s\n", port)
	fmt.Printf("====================================================\n")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
