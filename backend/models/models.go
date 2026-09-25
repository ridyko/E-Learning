package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"` // NIM for student, 21099001 for Pak Rio
	Name      string    `json:"name"`
	Role      string    `json:"role"`     // "dosen" or "mahasiswa"
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Prodi     string    `json:"prodi,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type LandingSettings struct {
	HeroTitle       string `json:"hero_title"`
	HeroSubtitle    string `json:"hero_subtitle"`
	InstitutionBadge string `json:"institution_badge"`
	SemesterBadge   string `json:"semester_badge"`
	DosenBio        string `json:"dosen_bio"`
	OfficeHours     string `json:"office_hours"`
	Room            string `json:"room"`
	Phone           string `json:"phone"`
}

type Profile struct {
	Name           string   `json:"name"`
	Degree         string   `json:"degree"`
	Title          string   `json:"title"`
	NIDN           string   `json:"nidn"`
	Institution    string   `json:"institution"`
	Faculty        string   `json:"faculty"`
	Department     string   `json:"department"`
	Email          string   `json:"email"`
	Phone          string   `json:"phone"`
	OfficeHours    string   `json:"office_hours"`
	Room           string   `json:"room"`
	Avatar         string   `json:"avatar"`
	Bio            string   `json:"bio"`
	ExpertiseAreas []string `json:"expertise_areas"`
}

type Module struct {
	MeetingNumber int      `json:"meeting_number"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Topics        []string `json:"topics"`
	SlideURL      string   `json:"slide_url"`
	PdfURL        string   `json:"pdf_url"`
	VideoURL      string   `json:"video_url"`
	CodeSample    string   `json:"code_sample,omitempty"`
	IsActive      bool     `json:"is_active"`
	Status        string   `json:"status,omitempty"`       // "terbuka", "terkunci", "dijadwalkan"
	ScheduledAt   string   `json:"scheduled_at,omitempty"`  // e.g. "Senin, 28 Sep 2026 08:00 WIB"
	TaskDueDate   string   `json:"task_due_date,omitempty"` // e.g. "2026-09-28 23:59 WIB"
}

type Course struct {
	ID            string   `json:"id"`
	Code          string   `json:"code"`
	Name          string   `json:"name"`
	SKS           int      `json:"sks"`
	Semester      string   `json:"semester"`
	ClassTime     string   `json:"class_time"`
	Room          string   `json:"room"`
	Description   string   `json:"description"`
	Syllabus      []string `json:"syllabus"`
	Modules       []Module `json:"modules"`
	TotalStudents int      `json:"total_students"`
}

type Announcement struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Category  string    `json:"category"` // Important, Assignment, Quiz, General
	CourseID  string    `json:"course_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Author    string    `json:"author"`
}

type QuizQuestion struct {
	ID            int      `json:"id"`
	Type          string   `json:"type"` // "mc" (Pilihan Ganda) or "essay" (Essay)
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectAnswer int      `json:"correct_answer"`
	Explanation   string   `json:"explanation"`
}

type Quiz struct {
	ID          string         `json:"id"`
	CourseID    string         `json:"course_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	TimeLimit   int            `json:"time_limit"` // in minutes
	Questions   []QuizQuestion `json:"questions"`
}

type QuizSubmission struct {
	ID           string            `json:"id"`
	QuizID       string            `json:"quiz_id"`
	StudentNIM   string            `json:"student_nim"`
	StudentName  string            `json:"student_name"`
	Answers      map[int]int       `json:"answers"`       // questionID -> optionIndex (for MC)
	EssayAnswers map[int]string    `json:"essay_answers"` // questionID -> essay text answer
	Score        float64           `json:"score"`
	TotalScore   float64           `json:"total_score"`
	SubmittedAt  time.Time         `json:"submitted_at"`
}

type Assignment struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"course_id"`
	MeetingNo   int       `json:"meeting_no"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"`
	StudentNIM  string    `json:"student_nim"`
	StudentName string    `json:"student_name"`
	RepoLink    string    `json:"repo_link"`
	Notes       string    `json:"notes"`
	Status      string    `json:"status"` // Submitted, Graded
	Grade       string    `json:"grade,omitempty"`
	SubmittedAt time.Time `json:"submitted_at"`
}

type Attendance struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"course_id"`
	MeetingNo   int       `json:"meeting_no"`
	StudentNIM  string    `json:"student_nim"`
	StudentName string    `json:"student_name"`
	Status      string    `json:"status"` // Hadir, Izin, Sakit
	CheckInTime time.Time `json:"check_in_time"`
}

type LiveParticipant struct {
	StudentNIM  string    `json:"student_nim"`
	StudentName string    `json:"student_name"`
	TotalScore  int       `json:"total_score"`
	StreakCount int       `json:"streak_count"`
	JoinedAt    time.Time `json:"joined_at"`
}

type LiveVote struct {
	QuestionIndex int       `json:"question_index"`
	StudentNIM    string    `json:"student_nim"`
	SelectedOpt   int       `json:"selected_opt"`
	AnsweredInSec float64   `json:"answered_in_sec"`
	ScoreEarned   int       `json:"score_earned"`
}

type LiveQuizSession struct {
	PIN             string                     `json:"pin"`
	QuizID          string                     `json:"quiz_id"`
	CourseID        string                     `json:"course_id"`
	Title           string                     `json:"title"`
	CurrentQIndex   int                        `json:"current_q_index"`
	Status          string                     `json:"status"` // "LOBBY", "QUESTION", "REVEAL", "LEADERBOARD", "FINISHED"
	TimePerQuestion int                        `json:"time_per_question"`
	StartedAt       time.Time                  `json:"started_at"`
	Participants    map[string]LiveParticipant `json:"participants"`
	Votes           []LiveVote                 `json:"votes"`
	Quiz            Quiz                       `json:"quiz"`
}
