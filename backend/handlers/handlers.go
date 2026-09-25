package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"elearning-backend/models"
	"elearning-backend/store"
)

// Enable CORS middleware
func EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func HandleLandingSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		settings := store.DB.GetLandingSettings()
		json.NewEncoder(w).Encode(settings)

	case "POST", "PUT":
		var s models.LandingSettings
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		updated := store.DB.UpdateLandingSettings(s)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "success",
			"message":  "Tampilan Beranda berhasil diperbarui oleh Pak Rio!",
			"settings": updated,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	User    models.User `json:"user"`
	Token   string      `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	user, err := store.DB.Authenticate(req.Username, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	user.Password = ""

	resp := LoginResponse{
		Status:  "success",
		Message: "Login berhasil! Selamat datang, " + user.Name,
		User:    user,
		Token:   "token-swadharma-" + user.Username,
	}

	json.NewEncoder(w).Encode(resp)
}

func RegisterStudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if user.Username == "" || user.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "NIM dan Nama Lengkap wajib diisi!"})
		return
	}

	newUser, err := store.DB.RegisterStudent(user)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	newUser.Password = ""
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Pendaftaran akun mahasiswa berhasil! Silakan login dengan NIM Anda.",
		"user":    newUser,
	})
}

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	students := store.DB.GetStudents()
	for i := range students {
		students[i].Password = ""
	}
	json.NewEncoder(w).Encode(students)
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	profile := store.DB.GetProfile()
	json.NewEncoder(w).Encode(profile)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var p models.Profile
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	store.DB.UpdateProfile(p)
	json.NewEncoder(w).Encode(p)
}

func HandleCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/courses")
	path = strings.TrimPrefix(path, "/")

	switch r.Method {
	case "GET":
		if path != "" {
			course, found := store.DB.GetCourseByID(path)
			if !found {
				http.Error(w, "Course not found", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(course)
			return
		}
		courses := store.DB.GetCourses()
		json.NewEncoder(w).Encode(courses)

	case "POST":
		var c models.Course
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Payload tidak valid"})
			return
		}
		if c.Name == "" || c.Code == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Kode dan Nama Mata Kuliah wajib diisi!"})
			return
		}
		created := store.DB.AddCourse(c)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Mata Kuliah '" + created.Name + "' berhasil ditambahkan!",
			"course":  created,
		})

	case "PUT":
		var c models.Course
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Payload tidak valid"})
			return
		}
		updated := store.DB.UpdateCourse(c)
		if !updated {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Mata Kuliah tidak ditemukan"})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Mata Kuliah '" + c.Name + "' berhasil diperbarui!",
			"course":  c,
		})

	case "DELETE":
		id := r.URL.Query().Get("id")
		if id == "" {
			id = path
		}
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "ID Mata Kuliah wajib disertakan"})
			return
		}
		deleted := store.DB.DeleteCourse(id)
		if !deleted {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Mata Kuliah tidak ditemukan"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Mata Kuliah berhasil dihapus",
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func ResetCoursesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	restored := store.DB.ResetCourses()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Mata Kuliah bawaan berhasil dipulihkan!",
		"courses": restored,
	})
}

func HandleAnnouncements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		announcements := store.DB.GetAnnouncements()
		json.NewEncoder(w).Encode(announcements)

	case "POST":
		var ann models.Announcement
		if err := json.NewDecoder(r.Body).Decode(&ann); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newAnn := store.DB.AddAnnouncement(ann)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newAnn)

	case "DELETE":
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing announcement ID", http.StatusBadRequest)
			return
		}
		success := store.DB.DeleteAnnouncement(id)
		if !success {
			http.Error(w, "Announcement not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Announcement deleted successfully"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetQuizzes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		courseID := r.URL.Query().Get("course_id")
		quizzes := store.DB.GetQuizzes(courseID)
		json.NewEncoder(w).Encode(quizzes)
	case "POST":
		var q models.Quiz
		if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newQ := store.DB.CreateQuiz(q)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newQ)
	case "DELETE":
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing quiz ID", http.StatusBadRequest)
			return
		}
		success := store.DB.DeleteQuiz(id)
		if !success {
			http.Error(w, "Quiz not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Quiz package deleted successfully"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func SubmitQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var submission models.QuizSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := store.DB.SubmitQuiz(submission)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func GetQuizSubmissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	submissions := store.DB.GetSubmissions()
	json.NewEncoder(w).Encode(submissions)
}

func AddQuizQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "POST":
		var payload struct {
			QuizID   string              `json:"quiz_id"`
			Question models.QuizQuestion `json:"question"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		success := store.DB.AddQuizQuestion(payload.QuizID, payload.Question)
		if !success {
			http.Error(w, "Quiz not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Question added"})

	case "DELETE":
		quizID := r.URL.Query().Get("quiz_id")
		questionIDStr := r.URL.Query().Get("question_id")
		if quizID == "" || questionIDStr == "" {
			http.Error(w, "Missing quiz_id or question_id", http.StatusBadRequest)
			return
		}
		var qID int
		fmt.Sscanf(questionIDStr, "%d", &qID)
		success := store.DB.DeleteQuizQuestion(quizID, qID)
		if !success {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Question deleted successfully"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


func HandleAssignments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		assignments := store.DB.GetAssignments()
		json.NewEncoder(w).Encode(assignments)

	case "POST":
		var asg models.Assignment
		if err := json.NewDecoder(r.Body).Decode(&asg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newAsg := store.DB.AddAssignment(asg)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newAsg)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleAttendance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		attendances := store.DB.GetAttendances()
		json.NewEncoder(w).Encode(attendances)

	case "POST":
		var att models.Attendance
		if err := json.NewDecoder(r.Body).Decode(&att); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newAtt := store.DB.AddAttendance(att)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newAtt)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type UpdateModuleReq struct {
	CourseID      string   `json:"course_id"`
	MeetingNumber int      `json:"meeting_number"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Topics        []string `json:"topics"`
	SlideURL      string   `json:"slide_url"`
	PdfURL        string   `json:"pdf_url"`
	VideoURL      string   `json:"video_url"`
	CodeSample    string   `json:"code_sample"`
	IsActive      bool     `json:"is_active"`
	Status        string   `json:"status"`
	ScheduledAt   string   `json:"scheduled_at"`
	TaskDueDate   string   `json:"task_due_date"`
}

func UpdateModuleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var req UpdateModuleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Payload tidak valid"})
		return
	}

	updated := store.DB.UpdateModuleStatus(
		req.CourseID,
		req.MeetingNumber,
		req.Title,
		req.Description,
		req.Topics,
		req.SlideURL,
		req.PdfURL,
		req.VideoURL,
		req.CodeSample,
		req.IsActive,
		req.Status,
		req.ScheduledAt,
		req.TaskDueDate,
	)

	if !updated {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Modul tidak ditemukan"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Materi, deskripsi, dan status modul berhasil diperbarui!",
	})
}

func CreateLiveQuizSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		QuizID          string `json:"quiz_id"`
		TimePerQuestion int    `json:"time_per_question"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	session, err := store.DB.CreateLiveQuizSession(req.QuizID, req.TimePerQuestion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(session)
}

func JoinLiveQuizSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		PIN         string `json:"pin"`
		StudentNIM  string `json:"student_nim"`
		StudentName string `json:"student_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	session, err := store.DB.JoinLiveQuizSession(req.PIN, req.StudentNIM, req.StudentName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(session)
}

func StateLiveQuizSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		PIN           string `json:"pin"`
		Status        string `json:"status"`
		CurrentQIndex int    `json:"current_q_index"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	session, err := store.DB.SetLiveQuizSessionStatus(req.PIN, req.Status, req.CurrentQIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(session)
}

func VoteLiveQuizSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		PIN           string  `json:"pin"`
		StudentNIM    string  `json:"student_nim"`
		QuestionIndex int     `json:"question_index"`
		SelectedOpt   int     `json:"selected_opt"`
		AnsweredInSec float64 `json:"answered_in_sec"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	session, err := store.DB.VoteLiveQuizSession(req.PIN, req.StudentNIM, req.QuestionIndex, req.SelectedOpt, req.AnsweredInSec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(session)
}

func SyncLiveQuizSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pin := r.URL.Query().Get("pin")
	if pin == "" {
		http.Error(w, "PIN is required", http.StatusBadRequest)
		return
	}

	session, err := store.DB.GetLiveQuizSession(pin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(session)
}

func GetActiveLiveQuizSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	courseID := r.URL.Query().Get("course_id")
	session := store.DB.GetActiveLiveQuizSessionByCourse(courseID)
	if session == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"active": false,
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"active":  true,
		"session": session,
	})
}


