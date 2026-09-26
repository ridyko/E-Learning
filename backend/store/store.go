package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"elearning-backend/models"
)

type DataStore struct {
	mu              sync.RWMutex
	SQLDB           *sql.DB
	LandingSettings models.LandingSettings
	Users           []models.User
	Profile         models.Profile
	Courses         []models.Course
	Announcements   []models.Announcement
	Quizzes         []models.Quiz
	Submissions     []models.QuizSubmission
	Assignments     []models.Assignment
	Attendances     []models.Attendance
	LiveSessions    map[string]*models.LiveQuizSession
}

var DB *DataStore

func InitStore() {
	DB = &DataStore{
		LiveSessions: make(map[string]*models.LiveQuizSession),
		LandingSettings: models.LandingSettings{
			HeroTitle:        "Portal Pembelajaran Digital Terpadu",
			HeroSubtitle:     "Selamat datang di Portal E-Learning Resmi Dosen Rio Widyatmoko, S.Kom, M.M.S.I. Fasilitas pembelajaran interaktif khusus mata kuliah Rekayasa Perangkat Lunak dan Pemrograman Web di ITB Swadharma.",
			InstitutionBadge: "ITB SWADHARMA JAKARTA",
			SemesterBadge:    "SEMESTER GANJIL 2026/2027",
			DosenBio:         "Dosen dan Praktisi Perangkat Lunak & Pengembangan Web di ITB Swadharma. Berpengalaman dalam arsitektur sistem, rekayasa perangkat lunak modern, dan pengajaran teknologi web berbasis industri.",
			OfficeHours:      "Senin & Rabu: 13.00 - 16.00 WIB",
			Room:             "Ruang Dosen Lantai 3 - Kampus Swadharma",
			Phone:            "+62 812-9876-5432",
		},
		Users: []models.User{
			{
				ID:        "usr-dosen-1",
				Username:  "21099001",
				Name:      "Rio Widyatmoko, S.Kom, M.M.S.I",
				Role:      "dosen",
				Password:  "123456",
				Email:     "rio.widyatmoko@swadharma.ac.id",
				Prodi:     "Teknik Informatika",
				CreatedAt: time.Now(),
			},
			{
				ID:        "usr-std-1",
				Username:  "20260801001",
				Name:      "Ahmad Fauzi",
				Role:      "mahasiswa",
				Password:  "123456",
				Email:     "ahmad.fauzi@student.swadharma.ac.id",
				Prodi:     "Teknik Informatika",
				CreatedAt: time.Now(),
			},
			{
				ID:        "usr-std-2",
				Username:  "20260801002",
				Name:      "Siti Nurhaliza",
				Role:      "mahasiswa",
				Password:  "123456",
				Email:     "siti.nurhaliza@student.swadharma.ac.id",
				Prodi:     "Sistem Informasi",
				CreatedAt: time.Now(),
			},
		},
		Profile: models.Profile{
			Name:        "Rio Widyatmoko",
			Degree:      "S.Kom, M.M.S.I",
			Title:       "Dosen Pengampu ITB Swadharma",
			NIDN:        "21099001",
			Institution: "Institut Teknologi dan Bisnis Swadharma (ITB Swadharma)",
			Faculty:     "Fakultas Teknologi Informasi & Desain",
			Department:  "Teknik Informatika / Sistem Informasi",
			Email:       "rio.widyatmoko@swadharma.ac.id",
			Phone:       "+62 812-9876-5432",
			OfficeHours: "Senin & Rabu: 13.00 - 16.00 WIB",
			Room:        "Ruang Dosen Lantai 3 - Kampus Swadharma",
			Avatar:      "https://images.unsplash.com/photo-1560250097-0b93528c311a?auto=format&fit=crop&q=80&w=400",
			Bio:         "Dosen dan Praktisi Perangkat Lunak & Pengembangan Web di ITB Swadharma. Berpengalaman dalam arsitektur sistem, rekayasa perangkat lunak modern, dan pengajaran teknologi web berbasis industri.",
			ExpertiseAreas: []string{
				"Software Engineering & Architecture",
				"Fullstack Web Development (Golang, Vue.js, PHP, Node.js)",
				"Database Management & System Analysis",
				"Agile & DevOps Implementation",
			},
		},
		Courses: []models.Course{
			{
				ID:            "rpl-2026",
				Code:          "TIF-301",
				Name:          "Rekayasa Perangkat Lunak",
				SKS:           3,
				Semester:      "Ganjil 2026/2027",
				ClassTime:     "Senin, 08.00 - 10.30 WIB",
				Room:          "Lab Komputer 2 / Hybrid Zoom",
				TotalStudents: 38,
				Description:   "Mata kuliah ini membahas prinsip, metode, dan teknik rekayasa perangkat lunak mulai dari perencanaan sistem, analisis kebutuhan (Software Requirements Specification), pemodelan UML, perancangan arsitektur, pengkodean, hingga pengujian & pemeliharaan perangkat lunak.",
				Syllabus: []string{
					"Pengantar Rekayasa Perangkat Lunak & Evolution of Software",
					"Model Proses Perangkat Lunak (Waterfall, Agile, Scrum, Kanban)",
					"Manajemen Proyek Perangkat Lunak & Estimation (COCOMO)",
					"Analisis Kebutuhan & Dokumen SRS (IEEE 830)",
					"Pemodelan Sistem dengan Unified Modeling Language (UML)",
					"Diagram Use Case & Diagram Activity",
					"Diagram Class & Diagram Sequence",
					"Perancangan Arsitektur Perangkat Lunak (Monolith vs Microservices)",
					"User Interface & User Experience Design Patterns",
					"Prinsip Object-Oriented Design (SOLID Principles)",
					"Software Testing & Quality Assurance (White-box & Black-box)",
					"Software Maintenance, Refactoring & Technical Debt",
					"DevOps, CI/CD Pipeline & Deployment Automation",
					"Presentasi Proyek Akhir Rekayasa Perangkat Lunak",
				},
				Modules: generateRPLModules(),
			},
			{
				ID:            "webdev-2026",
				Code:          "TIF-302",
				Name:          "Pemrograman Web",
				SKS:           3,
				Semester:      "Ganjil 2026/2027",
				ClassTime:     "Rabu, 08.00 - 10.30 WIB",
				Room:          "Lab Pemrograman Web / Modern Studio",
				TotalStudents: 42,
				Description:   "Mata kuliah praktikum & teori pengembangan aplikasi web modern modern frontend & backend. Mengajarkan HTML5 Semantic, CSS3 Modern Layout (Flexbox/Grid), JavaScript ES6+, Asynchronous API Client, Backend Integration (Golang/PHP/Node), RESTful API, dan Frontend Framework (Vue.js).",
				Syllabus: []string{
					"Pengantar Arsitektur Web & Protokol HTTP/HTTPS",
					"HTML5 Semantik & Accessible Web Structure (ARIA)",
					"CSS3 Modern Styling, Flexbox, Grid & Responsive Web Design",
					"Dasar Pemrograman JavaScript ES6+ & DOM Manipulation",
					"Event Handling, Form Validation & Local Storage",
					"Asynchronous JavaScript (Promises, Async/Await, Fetch API)",
					"Pengantar Frontend Framework: Vue.js 3 Fundamentals",
					"Vue Components, Reactive State & Event Emitting",
					"Pengantar Backend Web & RESTful API (Golang / PHP)",
					"Routing & Controller Architecture pada Backend",
					"Integrasi Database Relation (MySQL/PostgreSQL) dengan Backend API",
					"Autentikasi Web & Keamanan (JWT, CORS, XSS, CSRF Prevention)",
					"Integrasi Frontend Vue.js dengan Backend REST API",
					"Deployment Web App ke Server Production & Cloud Hosting",
				},
				Modules: generateWebModules(),
			},
		},
		Announcements: []models.Announcement{
			{
				ID:        "ann-1",
				Title:     "Selamat Datang di Portal E-Learning Semester Ganjil 2026/2027",
				Category:  "Penting",
				CourseID:  "all",
				Content:   "Selamat datang mahasiswa ITB Swadharma! Portal ini digunakan untuk mengakses materi kuliah Rekayasa Perangkat Lunak dan Pemrograman Web, pengumpulan tugas, kuis online, serta informasi presensi. Harap cek silabus dan kuis mingguan.",
				CreatedAt: time.Now().Add(-48 * time.Hour),
				Author:    "Rio Widyatmoko, S.Kom, M.M.S.I",
			},
			{
				ID:        "ann-2",
				Title:     "Tugas 1 Pemrograman Web: Membuat Layout Responsive dengan Flexbox/Grid",
				Category:  "Tugas",
				CourseID:  "webdev-2026",
				Content:   "Tugas 1 Pemrograman Web sudah dapat dikumpulkan melalui menu Tugas. Batas waktu pengumpulan adalah hari Minggu pukul 23.59 WIB. Silakan gunakan Live Code Playground di portal ini untuk menguji kode Anda.",
				CreatedAt: time.Now().Add(-24 * time.Hour),
				Author:    "Rio Widyatmoko, S.Kom, M.M.S.I",
			},
			{
				ID:        "ann-3",
				Title:     "Kuis 1 Rekayasa Perangkat Lunak: Model Proses Perangkat Lunak & SDLC",
				Category:  "Kuis",
				CourseID:  "rpl-2026",
				Content:   "Kuis 1 untuk mata kuliah Rekayasa Perangkat Lunak telah dibuka. Kuis terdiri dari 5 soal pilihan ganda dengan durasi 15 menit. Selamat mengerjakan!",
				CreatedAt: time.Now().Add(-12 * time.Hour),
				Author:    "Rio Widyatmoko, S.Kom, M.M.S.I",
			},
		},
		Quizzes: generateQuizzes(),
		Submissions: []models.QuizSubmission{
			{
				ID:          "sub-1",
				QuizID:      "quiz-rpl-1",
				StudentNIM:  "20260801001",
				StudentName: "Ahmad Fauzi",
				Answers:     map[int]int{1: 1, 2: 1, 3: 0, 4: 2, 5: 1},
				Score:       100.0,
				TotalScore:  100.0,
				SubmittedAt: time.Now().Add(-2 * time.Hour),
			},
			{
				ID:          "sub-2",
				QuizID:      "quiz-webdev-1",
				StudentNIM:  "20260801002",
				StudentName: "Siti Nurhaliza",
				Answers:     map[int]int{1: 2, 2: 1, 3: 1, 4: 1, 5: 0},
				Score:       100.0,
				TotalScore:  100.0,
				SubmittedAt: time.Now().Add(-5 * time.Hour),
			},
		},
		Assignments: []models.Assignment{
			{
				ID:          "asg-1",
				CourseID:    "webdev-2026",
				MeetingNo:   3,
				Title:       "Tugas Flexbox & Responsive Card Layout",
				Description: "Buat halaman profil pribadi menggunakan HTML5 semantik dan CSS Flexbox yang fully responsive.",
				DueDate:     "2026-09-28",
				StudentNIM:  "20260801001",
				StudentName: "Ahmad Fauzi",
				RepoLink:    "https://github.com/ahmadfauzi/web-profile-swadharma",
				Notes:       "Sudah ditambahkan fitur dark mode dan responsive breakpoint.",
				Status:      "Graded",
				Grade:       "95/100",
				SubmittedAt: time.Now().Add(-2 * time.Hour),
			},
		},
		Attendances: []models.Attendance{
			{
				ID:          "att-1",
				CourseID:    "rpl-2026",
				MeetingNo:   1,
				StudentNIM:  "20260801001",
				StudentName: "Ahmad Fauzi",
				Status:      "Hadir",
				CheckInTime: time.Now().Add(-1 * time.Hour),
			},
		},
	}

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		paths := []string{
			"data/db_config.json",
			"/home/questkom/backend_elearning/data/db_config.json",
			"backend/data/db_config.json",
		}
		for _, p := range paths {
			if cfgData, err := os.ReadFile(p); err == nil {
				var cfg struct {
					MySQLDSN string `json:"mysql_dsn"`
				}
				if err := json.Unmarshal(cfgData, &cfg); err == nil && cfg.MySQLDSN != "" {
					dsn = cfg.MySQLDSN
					fmt.Printf("[DataStore] Loaded MySQL DSN from config file: %s\n", p)
					break
				}
			}
		}
	}

	if dsn != "" {
		fmt.Printf("[DataStore] Attempting MySQL connection with DSN: %s...\n", dsn)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			fmt.Printf("[DataStore] MySQL open error: %v (falling back to JSON storage)\n", err)
		} else {
			if pingErr := db.Ping(); pingErr != nil {
				fmt.Printf("[DataStore] MySQL ping error: %v (falling back to JSON storage)\n", pingErr)
			} else {
				DB.SQLDB = db
				fmt.Println("[DataStore] Connected to MySQL Database successfully! 🐬")
			}
		}
	}

	if DB.SQLDB != nil {
		DB.loadFromMySQL()
		fmt.Println("[DataStore] Loaded state from cPanel MySQL Database (questkom_elearning)")
	} else if DB.loadFromDisk() {
		fmt.Println("[DataStore] Loaded persisted database state from data/store.json")
	} else {
		DB.mu.Lock()
		DB.saveToDiskLocked()
		DB.mu.Unlock()
		fmt.Println("[DataStore] Initialized default database state and saved to data/store.json")
	}
}

func (ds *DataStore) loadFromMySQL() {
	if ds.SQLDB == nil {
		return
	}
	// Load users from MySQL
	rows, err := ds.SQLDB.Query("SELECT id, username, name, role, password, email, prodi FROM users")
	if err == nil {
		var loadedUsers []models.User
		for rows.Next() {
			var u models.User
			var email, prodi sql.NullString
			_ = rows.Scan(&u.ID, &u.Username, &u.Name, &u.Role, &u.Password, &email, &prodi)
			u.Email = email.String
			u.Prodi = prodi.String
			loadedUsers = append(loadedUsers, u)
		}
		rows.Close()
		if len(loadedUsers) > 0 {
			ds.Users = loadedUsers
		}
	}
}

const dbFilePath = "data/store.json"

type storeFileDump struct {
	LandingSettings models.LandingSettings `json:"landing_settings"`
	Users           []models.User          `json:"users"`
	Profile         models.Profile         `json:"profile"`
	Courses         []models.Course        `json:"courses"`
	Announcements   []models.Announcement  `json:"announcements"`
	Quizzes         []models.Quiz          `json:"quizzes"`
	Submissions     []models.QuizSubmission `json:"submissions"`
	Assignments     []models.Assignment    `json:"assignments"`
	Attendances     []models.Attendance     `json:"attendances"`
}

func (ds *DataStore) saveToDiskLocked() {
	_ = os.MkdirAll(filepath.Dir(dbFilePath), 0755)
	dump := storeFileDump{
		LandingSettings: ds.LandingSettings,
		Users:           ds.Users,
		Profile:         ds.Profile,
		Courses:         ds.Courses,
		Announcements:   ds.Announcements,
		Quizzes:         ds.Quizzes,
		Submissions:     ds.Submissions,
		Assignments:     ds.Assignments,
		Attendances:     ds.Attendances,
	}
	data, err := json.MarshalIndent(dump, "", "  ")
	if err == nil {
		_ = os.WriteFile(dbFilePath, data, 0644)
	}
}

func (ds *DataStore) loadFromDisk() bool {
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		return false
	}
	data, err := os.ReadFile(dbFilePath)
	if err != nil {
		return false
	}
	var dump storeFileDump
	if err := json.Unmarshal(data, &dump); err != nil {
		return false
	}

	if dump.LandingSettings.HeroTitle != "" {
		ds.LandingSettings = dump.LandingSettings
	}
	if len(dump.Users) > 0 {
		ds.Users = dump.Users
	}
	if dump.Profile.Name != "" {
		ds.Profile = dump.Profile
		fullName := ds.Profile.Name
		if ds.Profile.Degree != "" && ds.Profile.Degree != "-" {
			fullName += ", " + ds.Profile.Degree
		}
		for i := range ds.Users {
			if ds.Users[i].Role == "dosen" {
				ds.Users[i].Name = fullName
				ds.Users[i].Email = ds.Profile.Email
			}
		}
	}
	if len(dump.Courses) > 0 {
		ds.Courses = dump.Courses
	}
	if len(dump.Announcements) > 0 {
		ds.Announcements = dump.Announcements
	}
	if len(dump.Quizzes) > 0 {
		ds.Quizzes = dump.Quizzes
	}
	if len(dump.Submissions) > 0 {
		ds.Submissions = dump.Submissions
	}
	if len(dump.Assignments) > 0 {
		ds.Assignments = dump.Assignments
	}
	if len(dump.Attendances) > 0 {
		ds.Attendances = dump.Attendances
	}
	return true
}

func (ds *DataStore) GetLandingSettings() models.LandingSettings {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.LandingSettings
}

func (ds *DataStore) UpdateLandingSettings(s models.LandingSettings) models.LandingSettings {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.LandingSettings = s
	ds.Profile.Bio = s.DosenBio
	ds.Profile.OfficeHours = s.OfficeHours
	ds.Profile.Room = s.Room
	ds.Profile.Phone = s.Phone
	ds.saveToDiskLocked()
	return ds.LandingSettings
}

func (ds *DataStore) AddCourse(c models.Course) models.Course {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if c.ID == "" {
		c.ID = "crs-" + time.Now().Format("20060102150405")
	}
	if len(c.Modules) == 0 {
		c.Modules = generateDefaultModules(c.Name)
	}
	ds.Courses = append(ds.Courses, c)
	ds.saveToDiskLocked()
	return c
}

func (ds *DataStore) UpdateCourse(c models.Course) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for i, existing := range ds.Courses {
		if existing.ID == c.ID {
			if c.Code != "" {
				ds.Courses[i].Code = c.Code
			}
			if c.Name != "" {
				ds.Courses[i].Name = c.Name
			}
			if c.SKS > 0 {
				ds.Courses[i].SKS = c.SKS
			}
			if c.Semester != "" {
				ds.Courses[i].Semester = c.Semester
			}
			if c.ClassTime != "" {
				ds.Courses[i].ClassTime = c.ClassTime
			}
			if c.Room != "" {
				ds.Courses[i].Room = c.Room
			}
			if c.Description != "" {
				ds.Courses[i].Description = c.Description
			}
			ds.saveToDiskLocked()
			return true
		}
	}
	return false
}

func (ds *DataStore) DeleteCourse(id string) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for i, c := range ds.Courses {
		if c.ID == id {
			ds.Courses = append(ds.Courses[:i], ds.Courses[i+1:]...)
			ds.saveToDiskLocked()
			return true
		}
	}
	return false
}

func (ds *DataStore) ResetCourses() []models.Course {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.Courses = []models.Course{
		{
			ID:            "rpl-2026",
			Code:          "TIF-301",
			Name:          "Rekayasa Perangkat Lunak",
			SKS:           3,
			Semester:      "Ganjil 2026/2027",
			ClassTime:     "Senin, 08.00 - 10.30 WIB",
			Room:          "Lab Komputer 2 / Hybrid Zoom",
			TotalStudents: 38,
			Description:   "Mata kuliah ini membahas prinsip, metode, dan teknik rekayasa perangkat lunak mulai dari perencanaan sistem, analisis kebutuhan (Software Requirements Specification), pemodelan UML, perancangan arsitektur, pengkodean, hingga pengujian & pemeliharaan perangkat lunak.",
			Syllabus: []string{
				"Pengantar Rekayasa Perangkat Lunak & Evolution of Software",
				"Model Proses Perangkat Lunak (Waterfall, Agile, Scrum, Kanban)",
				"Manajemen Proyek Perangkat Lunak & Estimation (COCOMO)",
				"Analisis Kebutuhan & Dokumen SRS (IEEE 830)",
				"Pemodelan Sistem dengan Unified Modeling Language (UML)",
				"Diagram Use Case & Diagram Activity",
				"Diagram Class & Diagram Sequence",
				"Perancangan Arsitektur Perangkat Lunak (Monolith vs Microservices)",
				"User Interface & User Experience Design Patterns",
				"Prinsip Object-Oriented Design (SOLID Principles)",
				"Software Testing & Quality Assurance (White-box & Black-box)",
				"Software Maintenance, Refactoring & Technical Debt",
				"DevOps, CI/CD Pipeline & Deployment Automation",
				"Presentasi Proyek Akhir Rekayasa Perangkat Lunak",
			},
			Modules: generateRPLModules(),
		},
		{
			ID:            "webdev-2026",
			Code:          "TIF-302",
			Name:          "Pemrograman Web",
			SKS:           3,
			Semester:      "Ganjil 2026/2027",
			ClassTime:     "Rabu, 08.00 - 10.30 WIB",
			Room:          "Lab Pemrograman Web / Modern Studio",
			TotalStudents: 42,
			Description:   "Mata kuliah praktikum & teori pengembangan aplikasi web modern modern frontend & backend. Mengajarkan HTML5 Semantic, CSS3 Modern Layout (Flexbox/Grid), JavaScript ES6+, Asynchronous API Client, Backend Integration (Golang/PHP/Node), RESTful API, dan Frontend Framework (Vue.js).",
			Syllabus: []string{
				"Pengantar Arsitektur Web & Protokol HTTP/HTTPS",
				"HTML5 Semantik & Accessible Web Structure (ARIA)",
				"CSS3 Modern Styling, Flexbox, Grid & Responsive Web Design",
				"Dasar Pemrograman JavaScript ES6+ & DOM Manipulation",
				"Event Handling, Form Validation & Local Storage",
				"Asynchronous JavaScript (Promises, Async/Await, Fetch API)",
				"Pengantar Frontend Framework: Vue.js 3 Fundamentals",
				"Vue Components, Reactive State & Event Emitting",
				"Pengantar Backend Web & RESTful API (Golang / PHP)",
				"Routing & Controller Architecture pada Backend",
				"Integrasi Database Relation (MySQL/PostgreSQL) dengan Backend API",
				"Autentikasi Web & Keamanan (JWT, CORS, XSS, CSRF Prevention)",
				"Integrasi Frontend Vue.js dengan Backend REST API",
				"Deployment Web App ke Server Production & Cloud Hosting",
			},
			Modules: generateWebModules(),
		},
	}
	ds.saveToDiskLocked()
	return ds.Courses
}

func (ds *DataStore) Authenticate(username, password string) (models.User, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	for _, u := range ds.Users {
		if u.Username == username && u.Password == password {
			return u, nil
		}
	}
	return models.User{}, errors.New("Username (NIP/NIM) atau password salah!")
}

func (ds *DataStore) RegisterStudent(u models.User) (models.User, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for _, existing := range ds.Users {
		if existing.Username == u.Username {
			return models.User{}, errors.New("NIM tersebut sudah terdaftar! Silakan login.")
		}
	}

	u.ID = "usr-std-" + time.Now().Format("20060102150405")
	u.Role = "mahasiswa"
	if u.Password == "" {
		u.Password = "123456"
	}
	u.CreatedAt = time.Now()

	ds.Users = append(ds.Users, u)
	ds.saveToDiskLocked()
	return u, nil
}

func (ds *DataStore) GetStudents() []models.User {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	var students []models.User
	for _, u := range ds.Users {
		if u.Role == "mahasiswa" {
			students = append(students, u)
		}
	}
	return students
}

func generateDefaultModules(courseName string) []models.Module {
	var mods []models.Module
	for i := 1; i <= 14; i++ {
		mods = append(mods, models.Module{
			MeetingNumber: i,
			Title:         "Pertemuan " + string(rune('0'+i)) + ": Pokok Bahasan " + courseName,
			Description:   "Materi perkuliahan pertemuan ke-" + string(rune('0'+i)) + " untuk mata kuliah " + courseName + ".",
			Topics:        []string{"Pengantar Topik", "Pembahasan Teori", "Studi Kasus & Praktikum"},
			SlideURL:      "#",
			PdfURL:        "#",
		})
	}
	return mods
}

func generateRPLModules() []models.Module {
	return []models.Module{
		{
			MeetingNumber: 1,
			Title:         "Pengantar Rekayasa Perangkat Lunak & Karakteristik Perangkat Lunak",
			Description:   "Pengenalan ruang lingkup RPL, perbedaan software vs hardware, software crises, dan pentingnya rekayasa perangkat lunak dalam industri IT modern.",
			Topics:        []string{"Definisi RPL", "Dual Role of Software", "Software Myths", "Kualitas Perangkat Lunak"},
			SlideURL:      "#",
			PdfURL:        "#",
			VideoURL:      "https://www.youtube.com/embed/dQw4w9WgXcQ",
		},
		{
			MeetingNumber: 2,
			Title:         "Software Development Life Cycle (SDLC) & Prescriptive Process Models",
			Description:   "Membahas tahapan utama SDLC dan model proses tradisional seperti Waterfall, V-Model, dan Incremental Model.",
			Topics:        []string{"Waterfall Model", "V-Shaped Model", "RAD Model", "Kelebihan & Kekurangan Model Linier"},
			SlideURL:      "#",
			PdfURL:        "#",
			TaskDueDate:   "Pertemuan 3",
		},
		{
			MeetingNumber: 3,
			Title:         "Agile Software Development & Framework Scrum",
			Description:   "Prinsip Agile Manifesto, perbandingan Waterfall vs Agile, peran dalam Scrum (Product Owner, Scrum Master, Dev Team), Sprint Planning & Daily Standup.",
			Topics:        []string{"Agile Manifesto", "Scrum Framework", "User Stories & Backlog", "Sprint Review & Retrospective"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 4,
			Title:         "Analisis Kebutuhan Perangkat Lunak & Penulisan Dokumen SRS",
			Description:   "Teknik elicitation kebutuhan (interview, quesioner, observasi), Kebutuhan Fungsional vs Non-Fungsional, dan penyusunan dokumen SRS standar IEEE 830.",
			Topics:        []string{"Functional Requirements", "Non-Functional Requirements", "IEEE 830 SRS Template", "Requirements Validation"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 5,
			Title:         "Pemodelan Perangkat Lunak dengan Unified Modeling Language (UML)",
			Description:   "Pengenalan diagram UML untuk analisis dan perancangan sistem terstruktur berorientasi objek.",
			Topics:        []string{"Konsep Dasar UML", "Struktur & Perilaku Diagram", "Tools Modeling (StarUML, Draw.io)"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 6,
			Title:         "UML Behavioral Diagram: Use Case & Activity Diagram",
			Description:   "Perancangan diagram Use Case beserta skenario Use Case Specification, dan pembuatan Activity Diagram untuk aliran kerja sistem.",
			Topics:        []string{"Actor & Use Case", "Include vs Extend", "Activity Diagram Notation", "Swimlane Activity"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 7,
			Title:         "UML Structural Diagram: Class Diagram & Sequence Diagram",
			Description:   "Membuat Class Diagram (Attributes, Methods, Relationship, Multiplicity) dan Sequence Diagram (Objects, Lifeline, Messages).",
			Topics:        []string{"Class Relationship (Association, Aggregation, Composition)", "Inheritance", "Sequence Diagram Synchronous/Asynchronous"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 8,
			Title:         "Ujian Tengah Semester (UTS) Rekayasa Perangkat Lunak",
			Description:   "Evaluasi pemahaman konsep SDLC, Analisis Kebutuhan SRS, dan Pemodelan UML Diagram.",
			Topics:        []string{"Studi Kasus Perancangan Sistem", "Pembuatan UML Diagram lengkap"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 9,
			Title:         "Perancangan Arsitektur Perangkat Lunak & Design Patterns",
			Description:   "Pola arsitektur software (Layered Architecture, MVC, Microservices, Event-Driven) dan penerapan SOLID Principles.",
			Topics:        []string{"MVC Architecture", "Monolithic vs Microservices", "SOLID Principles", "Design Patterns (Creational, Structural, Behavioral)"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 10,
			Title:         "User Interface (UI) & User Experience (UX) Engineering",
			Description:   "Prinsip desain antarmuka pengembang, wireframing, prototyping, dan pengujian kebergunaan (Usability Testing).",
			Topics:        []string{"Heuristic Evaluation", "Wireframe to High-Fidelity Prototype", "User-Centered Design (UCD)"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 11,
			Title:         "Software Testing & Quality Assurance (Testing Hitam & Putih)",
			Description:   "Strategi pengujian perangkat lunak: Unit Testing, Integration Testing, System Testing, Acceptance Testing, serta Black-box & White-box testing.",
			Topics:        []string{"Equivalence Partitioning & Boundary Value", "Cyclomatic Complexity & Basis Path Testing", "Automated Testing Frameworks"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 12,
			Title:         "Software Maintenance, Refactoring & Technical Debt",
			Description:   "Jenis-jenis pemeliharaan software (Corrective, Adaptive, Perfective, Preventive), teknik refactoring code, dan strategi pengelolaan Technical Debt.",
			Topics:        []string{"Code Smells", "Refactoring Patterns", "Managing Legacy Code", "Software Metrics"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 13,
			Title:         "DevOps, CI/CD Pipeline & Version Control System (Git)",
			Description:   "Penerapan kultur DevOps, otomasi build & deployment dengan CI/CD pipeline, Docker containerization, dan branching strategy Git Flow.",
			Topics:        []string{"Git Flow Strategy", "Continuous Integration / Delivery", "Docker Containerization Intro"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 14,
			Title:         "Presentasi Proyek Akhir Rekayasa Perangkat Lunak & UAS",
			Description:   "Presentasi laporan proyek perancangan aplikasi lengkap beserta dokumen SRS dan diagram UML oleh tiap kelompok mahasiswa.",
			Topics:        []string{"Presentasi Proyek Group", "Demo Prototype Software", "Penilaian Proyek Akhir"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
	}
}

func generateWebModules() []models.Module {
	return []models.Module{
		{
			MeetingNumber: 1,
			Title:         "Pengantar Arsitektur Web & Protokol HTTP/HTTPS",
			Description:   "Cara kerja World Wide Web, arsitektur Client-Server, siklus Request-Response HTTP, serta HTTP Status Codes.",
			Topics:        []string{"HTTP Verbs (GET, POST, PUT, DELETE)", "Header & Response Code", "Web Browser & Rendering Engine"},
			SlideURL:      "#",
			PdfURL:        "#",
			CodeSample:    "<!-- Contoh Struktur Dokumen HTML5 -->\n<!DOCTYPE html>\n<html lang=\"id\">\n<head>\n  <meta charset=\"UTF-8\">\n  <title>ITB Swadharma Web</title>\n</head>\n<body>\n  <h1>Selamat Datang di Matakuliah Pemrograman Web!</h1>\n  <p>Dosen Pengampu: Rio Widyatmoko, S.Kom, M.M.S.I</p>\n</body>\n</html>",
		},
		{
			MeetingNumber: 2,
			Title:         "HTML5 Semantik & Pengoperasian Struktur Web Modern",
			Description:   "Struktur HTML5 semantik (header, nav, main, section, article, footer), penanganan form input, dan accessibility WAI-ARIA.",
			Topics:        []string{"Semantic Elements", "Form Validation Attributes", "Media Elements (Audio/Video)", "Accessibility (ARIA)"},
			SlideURL:      "#",
			PdfURL:        "#",
			CodeSample:    "<form action=\"/submit\" method=\"POST\">\n  <label for=\"nama\">Nama Mahasiswa:</label>\n  <input type=\"text\" id=\"nama\" name=\"nama\" required placeholder=\"Ketik nama anda...\">\n  <button type=\"submit\">Kirim Form</button>\n</form>",
		},
		{
			MeetingNumber: 3,
			Title:         "CSS3 Modern Styling, Flexbox Layout & CSS Grid",
			Description:   "Menguasai CSS3 Selectors, Box Model, Responsive Layout dengan Flexbox dan CSS Grid, CSS Custom Properties (Variables), dan Glassmorphism UI.",
			Topics:        []string{"Flexbox Axis & Alignment", "Grid Template Areas", "Media Queries Responsive", "CSS Variables & Themes"},
			SlideURL:      "#",
			PdfURL:        "#",
			CodeSample:    ".card {\n  background: rgba(255, 255, 255, 0.1);\n  backdrop-filter: blur(10px);\n  border-radius: 12px;\n  padding: 1.5rem;\n  border: 1px solid rgba(255, 255, 255, 0.2);\n  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.37);\n}",
		},
		{
			MeetingNumber: 4,
			Title:         "Dasar Pemrograman JavaScript ES6+ & DOM Manipulation",
			Description:   "Sintaks JavaScript modern (let/const, Arrow Functions, Template Literals, Destructuring), manipulasi DOM, dan event listeners.",
			Topics:        []string{"ES6 Syntax", "querySelector & addEventListener", "Manipulasi Elemen DOM", "Array Methods (map, filter, reduce)"},
			SlideURL:      "#",
			PdfURL:        "#",
			CodeSample:    "// Mengubah teks saat tombol diklik\nconst btn = document.querySelector('#myBtn');\nbtn.addEventListener('click', () => {\n  document.querySelector('#output').textContent = 'Halo dari JavaScript ES6!';\n});",
		},
		{
			MeetingNumber: 5,
			Title:         "Event Handling, Form Validation & Browser Local Storage",
			Description:   "Menangani event kompleks, melakukan validasi form di sisi client, dan menyimpan state pengguna dengan LocalStorage / SessionStorage.",
			Topics:        []string{"Event Propagation (Bubbling/Capturing)", "Regex Form Validation", "LocalStorage API", "JSON Parse & Stringify"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 6,
			Title:         "Asynchronous JavaScript: Promises, Async/Await & Fetch API",
			Description:   "Pemrograman asinkron pada JS, menangani callback hell dengan Promise, kata kunci async/await, dan komunikasi HTTP dengan Fetch API.",
			Topics:        []string{"Synchronous vs Asynchronous", "Promise States", "Async/Await Pattern", "Fetch API & JSON Handling"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 7,
			Title:         "Pengantar Frontend Framework: Vue.js 3 Fundamentals",
			Description:   "Konsep dasar Vue.js 3, Single File Components (.vue), Reactive State (ref & reactive), Directive (v-if, v-for, v-model, v-on).",
			Topics:        []string{"Declarative Rendering", "Vue Directive Syntax", "Two-way Data Binding v-model", "Composition API"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 8,
			Title:         "Ujian Tengah Semester (UTS) Pemrograman Web",
			Description:   "Praktikum coding mandiri membuat landing page web interaktif dan responsive berbasis HTML5, CSS Flexbox, dan JavaScript ES6+.",
			Topics:        []string{"Live Coding Challenge", "Pengujian Responsive Display", "Penilaian Kebersihan Code"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 9,
			Title:         "Vue Components, Props, Emits & Vue Router SPA",
			Description:   "Arsitektur komponen terkomposisi pada Vue.js, passing data via Props, custom events dengan Emits, dan Single Page Application routing dengan Vue Router.",
			Topics:        []string{"Component Hierarchy", "Props Validation", "Vue Router Setup & Navigation Guards", "Dynamic Route Matching"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 10,
			Title:         "Pengantar Backend Development & RESTful API (Golang / PHP)",
			Description:   "Konsep backend web development, pembuatan RESTful API endpoint, JSON serialization, dan pengujian API dengan Postman/Hoppscotch.",
			Topics:        []string{"REST API Principles", "HTTP Method Mapping", "JSON Response Formatting", "Go net/http Router / PHP API"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 11,
			Title:         "Integrasi Database Relasional (MySQL) & ORM/Query",
			Description:   "Koneksi backend ke database MySQL, operasi CRUD (Create, Read, Update, Delete), SQL Injection Prevention, dan prepared statements.",
			Topics:        []string{"Database Driver Connection", "SQL Prepared Statements", "CRUD API Endpoints", "Database Migration Concept"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 12,
			Title:         "Keamanan Web: CORS, JWT Authentication, XSS & CSRF",
			Description:   "Prinsip keamanan aplikasi web, autentikasi stateless menggunakan JSON Web Token (JWT), konfigurasi CORS header, dan pencegahan celah XSS & CSRF.",
			Topics:        []string{"JWT Token Structure & Signing", "Cross-Origin Resource Sharing", "Sanitasi Input XSS", "CSRF Tokens"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 13,
			Title:         "Integrasi Fullstack: Connecting Vue.js Frontend dengan Golang API",
			Description:   "Menghubungkan frontend Vue 3 dengan backend API Golang/PHP, state management, handling loading status, error notifications, dan toast alerts.",
			Topics:        []string{"Axios / Fetch Client Setup", "Global State Management", "Error Boundary & Interceptors", "End-to-End Flow"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
		{
			MeetingNumber: 14,
			Title:         "Presentasi Proyek Akhir Pemrograman Web & UAS",
			Description:   "Demonstrasi aplikasi web fullstack interaktif yang dibuat oleh mahasiswa sebagai tugas akhir semester.",
			Topics:        []string{"Live Demo Aplikasi Web", "Code Review & Architecture Inspection", "Penilaian Akhir Pemrograman Web"},
			SlideURL:      "#",
			PdfURL:        "#",
		},
	}
}

func generateQuizzes() []models.Quiz {
	return []models.Quiz{
		{
			ID:          "quiz-rpl-1",
			CourseID:    "rpl-2026",
			Title:       "Kuis 1: Konsep SDLC & Model Proses Perangkat Lunak",
			Description: "Uji pemahaman Anda mengenai model SDLC Waterfall, Agile Scrum, dan analisis kebutuhan.",
			TimeLimit:   15,
			Questions: []models.QuizQuestion{
				{
					ID:            1,
					Question:      "Manakah dari pernyataan berikut yang menggambarkan karakteristik utama dari Waterfall Model dalam SDLC?",
					Options:       []string{"Pendekatan iteratif dan fleksibel terhadap perubahan", "Pendekatan sekuensial sistematis di mana setiap fase harus selesai sebelum fase berikutnya dimulai", "Pengembangan berbasis komponen tanpa dokumentasi", "Proses berbasis rilis harian tanpa perencanaan awal"},
					CorrectAnswer: 1,
					Explanation:   "Waterfall Model adalah pendekatan sekuensial linier di mana setiap fase (Analisis, Desain, Implementasi, Testing, Maintenance) berjalan berurutan.",
				},
				{
					ID:            2,
					Question:      "Dalam framework Scrum (Agile), siapa yang bertanggung jawab menentukan dan memprioritaskan fitur dalam Product Backlog?",
					Options:       []string{"Scrum Master", "Product Owner", "Software Tester / QA", "System Architect"},
					CorrectAnswer: 1,
					Explanation:   "Product Owner bertanggung jawab penuh dalam memaksimalkan nilai produk dan mengelola isi serta prioritas Product Backlog.",
				},
				{
					ID:            3,
					Question:      "Dokumen standar internasional yang menentukan spesifikasi kebutuhan perangkat lunak adalah...",
					Options:       []string{"IEEE 830 / ISO 29148 (SRS)", "UML 2.5 Specification", "W3C HTML Standards", "ISO 9001 Quality Control"},
					CorrectAnswer: 0,
					Explanation:   "IEEE 830 (sekarang ISO/IEC/IEEE 29148) adalah standar baku penulisan dokumen Software Requirements Specification (SRS).",
				},
				{
					ID:            4,
					Question:      "Diagram UML manakah yang digunakan untuk memodelkan interaksi antara Actor luar dengan fungsi-fungsi utama sistem?",
					Options:       []string{"Class Diagram", "Sequence Diagram", "Use Case Diagram", "State Machine Diagram"},
					CorrectAnswer: 2,
					Explanation:   "Use Case Diagram menggambarkan hubungan antara para Aktor dengan Use Case (fungsi) di dalam batasan sistem.",
				},
				{
					ID:            5,
					Question:      "Metode pengujian perangkat lunak yang berfokus pada struktur logika kode internal (code coverage) disebut...",
					Options:       []string{"Black-box Testing", "White-box Testing", "User Acceptance Testing (UAT)", "Beta Testing"},
					CorrectAnswer: 1,
					Explanation:   "White-box testing (atau structural testing) menguji alur logika internal, percabangan, dan fungsi kode sumber.",
				},
			},
		},
		{
			ID:          "quiz-webdev-1",
			CourseID:    "webdev-2026",
			Title:       "Kuis 1: HTML5 Semantik, CSS Flexbox & JavaScript ES6+",
			Description: "Uji pemahaman dasar web development modern: HTML5, CSS Flexbox/Grid, dan manipulasi DOM JavaScript.",
			TimeLimit:   15,
			Questions: []models.QuizQuestion{
				{
					ID:            1,
					Question:      "Elemen HTML5 semantik mana yang paling tepat digunakan untuk membungkus navigasi utama sebuah website?",
					Options:       []string{"<div class='nav'>", "<header>", "<nav>", "<aside>"},
					CorrectAnswer: 2,
					Explanation:   "Elemen <nav> diciptakan khusus pada HTML5 semantik untuk mendefinisikan kelompok link navigasi utama.",
				},
				{
					ID:            2,
					Question:      "Pada CSS Flexbox, properti apa yang digunakan untuk mengatur perataan elemen anak sepanjang sumbu utama (main axis)?",
					Options:       []string{"align-items", "justify-content", "flex-wrap", "align-content"},
					CorrectAnswer: 1,
					Explanation:   "justify-content mengatur alokasi ruang dan perataan item sepanjang main axis pada container Flexbox.",
				},
				{
					ID:            3,
					Question:      "Fitur JavaScript ES6 mana yang memungkinkan kita menyisipkan ekspresi ke dalam string dengan backticks (`)?",
					Options:       []string{"Arrow Functions", "Template Literals", "Destructuring Assignment", "Promises"},
					CorrectAnswer: 1,
					Explanation:   "Template Literals menggunakan tanda backtick (`) dan sintaks ${expression} untuk string interpolation.",
				},
				{
					ID:            4,
					Question:      "Manakah sintaks yang benar pada Vue.js 3 untuk melakukan two-way data binding pada elemen input?",
					Options:       []string{"v-bind:input", "v-model", "v-on:change", "v-slot"},
					CorrectAnswer: 1,
					Explanation:   "v-model menyediakan two-way data binding antara input form dan state komponen di Vue.js.",
				},
				{
					ID:            5,
					Question:      "HTTP Status Code manakah yang menandakan bahwa request berhasil diproses oleh server (OK)?",
					Options:       []string{"200 OK", "301 Moved Permanently", "404 Not Found", "500 Internal Server Error"},
					CorrectAnswer: 0,
					Explanation:   "Kode status HTTP 200 menandakan request sukses diproses oleh web server.",
				},
			},
		},
	}
}

// Thread-safe store methods
func (ds *DataStore) GetProfile() models.Profile {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.Profile
}

func (ds *DataStore) UpdateProfile(p models.Profile) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.Profile = p

	fullName := p.Name
	if p.Degree != "" && p.Degree != "-" {
		fullName += ", " + p.Degree
	}
	for i := range ds.Users {
		if ds.Users[i].Role == "dosen" {
			ds.Users[i].Name = fullName
			ds.Users[i].Email = p.Email
		}
	}
	ds.saveToDiskLocked()
}

func (ds *DataStore) GetCourses() []models.Course {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.Courses
}

func (ds *DataStore) GetCourseByID(id string) (models.Course, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	for _, c := range ds.Courses {
		if c.ID == id {
			return c, true
		}
	}
	return models.Course{}, false
}

func (ds *DataStore) GetAnnouncements() []models.Announcement {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.Announcements
}

func (ds *DataStore) AddAnnouncement(a models.Announcement) models.Announcement {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	a.ID = "ann-" + time.Now().Format("20060102150405")
	a.CreatedAt = time.Now()
	a.Author = ds.Profile.Name + ", " + ds.Profile.Degree
	ds.Announcements = append([]models.Announcement{a}, ds.Announcements...)
	ds.saveToDiskLocked()
	return a
}

func (ds *DataStore) DeleteAnnouncement(id string) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for i, a := range ds.Announcements {
		if a.ID == id {
			ds.Announcements = append(ds.Announcements[:i], ds.Announcements[i+1:]...)
			ds.saveToDiskLocked()
			return true
		}
	}
	return false
}

func (ds *DataStore) GetQuizzes(courseID string) []models.Quiz {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	var result []models.Quiz
	for _, q := range ds.Quizzes {
		if courseID == "" || q.CourseID == courseID {
			result = append(result, q)
		}
	}
	return result
}

func (ds *DataStore) CreateQuiz(q models.Quiz) models.Quiz {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if q.ID == "" {
		q.ID = "quiz-" + time.Now().Format("20060102150405")
	}
	if q.TimeLimit <= 0 {
		q.TimeLimit = 15
	}
	if q.Questions == nil {
		q.Questions = []models.QuizQuestion{}
	}
	ds.Quizzes = append(ds.Quizzes, q)
	ds.saveToDiskLocked()
	return q
}

func (ds *DataStore) SubmitQuiz(sub models.QuizSubmission) models.QuizSubmission {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	var targetQuiz models.Quiz
	for _, q := range ds.Quizzes {
		if q.ID == sub.QuizID {
			targetQuiz = q
			break
		}
	}

	correctCount := 0
	for _, q := range targetQuiz.Questions {
		if selected, exists := sub.Answers[q.ID]; exists && selected == q.CorrectAnswer {
			correctCount++
		}
	}

	totalQuestions := len(targetQuiz.Questions)
	if totalQuestions > 0 {
		sub.Score = (float64(correctCount) / float64(totalQuestions)) * 100.0
	}
	sub.TotalScore = 100.0
	sub.ID = "sub-" + time.Now().Format("20060102150405")
	sub.SubmittedAt = time.Now()

	ds.Submissions = append([]models.QuizSubmission{sub}, ds.Submissions...)
	ds.saveToDiskLocked()
	return sub
}

func (ds *DataStore) AddAssignment(asg models.Assignment) models.Assignment {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	asg.ID = "asg-" + time.Now().Format("20060102150405")
	asg.Status = "Submitted"
	asg.SubmittedAt = time.Now()
	ds.Assignments = append([]models.Assignment{asg}, ds.Assignments...)
	ds.saveToDiskLocked()
	return asg
}

func (ds *DataStore) GetAssignments() []models.Assignment {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.Assignments
}

func (ds *DataStore) AddAttendance(att models.Attendance) models.Attendance {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for i, existing := range ds.Attendances {
		if existing.CourseID == att.CourseID && existing.MeetingNo == att.MeetingNo && existing.StudentNIM == att.StudentNIM {
			ds.Attendances[i].Status = att.Status
			ds.Attendances[i].CheckInTime = time.Now()
			ds.saveToDiskLocked()
			return ds.Attendances[i]
		}
	}

	att.ID = "att-" + time.Now().Format("20060102150405")
	att.CheckInTime = time.Now()
	ds.Attendances = append([]models.Attendance{att}, ds.Attendances...)
	ds.saveToDiskLocked()
	return att
}

func (ds *DataStore) GetAttendances() []models.Attendance {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.Attendances
}

func (ds *DataStore) GetSubmissions() []models.QuizSubmission {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.Submissions
}

func (ds *DataStore) DeleteQuiz(id string) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for i, q := range ds.Quizzes {
		if q.ID == id {
			ds.Quizzes = append(ds.Quizzes[:i], ds.Quizzes[i+1:]...)
			ds.saveToDiskLocked()
			return true
		}
	}
	return false
}

func (ds *DataStore) AddQuizQuestion(quizID string, q models.QuizQuestion) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if q.Type == "" {
		q.Type = "mc"
	}
	for i, quiz := range ds.Quizzes {
		if quiz.ID == quizID {
			q.ID = len(quiz.Questions) + 1
			ds.Quizzes[i].Questions = append(ds.Quizzes[i].Questions, q)
			ds.saveToDiskLocked()
			return true
		}
	}
	return false
}

func (ds *DataStore) DeleteQuizQuestion(quizID string, questionID int) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for i, quiz := range ds.Quizzes {
		if quiz.ID == quizID {
			for j, q := range quiz.Questions {
				if q.ID == questionID {
					ds.Quizzes[i].Questions = append(quiz.Questions[:j], quiz.Questions[j+1:]...)
					ds.saveToDiskLocked()
					return true
				}
			}
		}
	}
	return false
}

func (ds *DataStore) UpdateModuleStatus(courseID string, meetingNo int, title string, desc string, topics []string, slideURL string, pdfURL string, videoURL string, codeSample string, isActive bool, status string, scheduledAt string, taskDueDate string) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for i, c := range ds.Courses {
		if c.ID == courseID {
			for j, m := range c.Modules {
				if m.MeetingNumber == meetingNo {
					if title != "" {
						ds.Courses[i].Modules[j].Title = title
					}
					if desc != "" {
						ds.Courses[i].Modules[j].Description = desc
					}
					if len(topics) > 0 {
						ds.Courses[i].Modules[j].Topics = topics
					}
					if slideURL != "" {
						ds.Courses[i].Modules[j].SlideURL = slideURL
					}
					if pdfURL != "" {
						ds.Courses[i].Modules[j].PdfURL = pdfURL
					}
					if videoURL != "" {
						ds.Courses[i].Modules[j].VideoURL = videoURL
					}
					if codeSample != "" {
						ds.Courses[i].Modules[j].CodeSample = codeSample
					}
					ds.Courses[i].Modules[j].IsActive = isActive
					if status != "" {
						ds.Courses[i].Modules[j].Status = status
					}
					if scheduledAt != "" {
						ds.Courses[i].Modules[j].ScheduledAt = scheduledAt
					}
					if taskDueDate != "" {
						ds.Courses[i].Modules[j].TaskDueDate = taskDueDate
					}
					ds.saveToDiskLocked()
					return true
				}
			}
		}
	}
	return false
}

func (ds *DataStore) CreateLiveQuizSession(quizID string, timePerQuestion int) (*models.LiveQuizSession, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	var targetQuiz *models.Quiz
	for _, q := range ds.Quizzes {
		if q.ID == quizID {
			targetQuiz = &q
			break
		}
	}
	if targetQuiz == nil {
		return nil, errors.New("quiz not found")
	}

	if timePerQuestion <= 0 {
		timePerQuestion = 20
	}

	// Generate 6 digit PIN
	rand.Seed(time.Now().UnixNano())
	pin := fmt.Sprintf("%06d", rand.Intn(900000)+100000)

	session := &models.LiveQuizSession{
		PIN:             pin,
		QuizID:          targetQuiz.ID,
		CourseID:        targetQuiz.CourseID,
		Title:           targetQuiz.Title,
		CurrentQIndex:   0,
		Status:          "LOBBY",
		TimePerQuestion: timePerQuestion,
		StartedAt:       time.Now(),
		Participants:    make(map[string]models.LiveParticipant),
		Votes:           make([]models.LiveVote, 0),
		Quiz:            *targetQuiz,
	}

	ds.LiveSessions[pin] = session
	return session, nil
}

func (ds *DataStore) JoinLiveQuizSession(pin string, studentNIM string, studentName string) (*models.LiveQuizSession, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	session, exists := ds.LiveSessions[pin]
	if !exists {
		return nil, errors.New("sesi kuis dengan PIN tersebut tidak ditemukan")
	}

	if studentNIM == "" {
		studentNIM = "GUEST-" + fmt.Sprintf("%04d", rand.Intn(9000)+1000)
	}

	p, existsPart := session.Participants[studentNIM]
	if !existsPart {
		p = models.LiveParticipant{
			StudentNIM:  studentNIM,
			StudentName: studentName,
			TotalScore:  0,
			StreakCount: 0,
			JoinedAt:    time.Now(),
		}
	} else if studentName != "" {
		p.StudentName = studentName
	}

	session.Participants[studentNIM] = p
	return session, nil
}

func (ds *DataStore) SetLiveQuizSessionStatus(pin string, status string, currentQIndex int) (*models.LiveQuizSession, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	session, exists := ds.LiveSessions[pin]
	if !exists {
		return nil, errors.New("sesi kuis tidak ditemukan")
	}

	if status != "" {
		session.Status = status
	}
	if currentQIndex >= 0 {
		session.CurrentQIndex = currentQIndex
	}

	// Auto submit scores to main QuizSubmissions if session is FINISHED
	if session.Status == "FINISHED" {
		for _, p := range session.Participants {
			// Find existing submission or append new
			maxPossible := len(session.Quiz.Questions) * 1000
			scorePct := 0.0
			if maxPossible > 0 {
				scorePct = (float64(p.TotalScore) / float64(maxPossible)) * 100.0
				if scorePct > 100.0 {
					scorePct = 100.0
				}
			}

			subID := fmt.Sprintf("sub-live-%s-%s", session.PIN, p.StudentNIM)
			found := false
			for i, sub := range ds.Submissions {
				if sub.QuizID == session.QuizID && sub.StudentNIM == p.StudentNIM {
					ds.Submissions[i].Score = scorePct
					ds.Submissions[i].TotalScore = 100.0
					ds.Submissions[i].SubmittedAt = time.Now()
					found = true
					break
				}
			}

			if !found {
				ds.Submissions = append(ds.Submissions, models.QuizSubmission{
					ID:          subID,
					QuizID:      session.QuizID,
					StudentNIM:  p.StudentNIM,
					StudentName: p.StudentName,
					Answers:     make(map[int]int),
					Score:       scorePct,
					TotalScore:  100.0,
					SubmittedAt: time.Now(),
				})
			}
		}
	}

	return session, nil
}

func (ds *DataStore) VoteLiveQuizSession(pin string, studentNIM string, questionIndex int, selectedOpt int, answeredInSec float64) (*models.LiveQuizSession, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	session, exists := ds.LiveSessions[pin]
	if !exists {
		return nil, errors.New("sesi kuis tidak ditemukan")
	}

	// Check if already voted for this questionIndex
	for i, v := range session.Votes {
		if v.StudentNIM == studentNIM && v.QuestionIndex == questionIndex {
			session.Votes[i].SelectedOpt = selectedOpt
			session.Votes[i].AnsweredInSec = answeredInSec
			return session, nil
		}
	}

	// Calculate score
	scoreEarned := 0
	if questionIndex >= 0 && questionIndex < len(session.Quiz.Questions) {
		q := session.Quiz.Questions[questionIndex]
		part, hasPart := session.Participants[studentNIM]
		if !hasPart {
			part = models.LiveParticipant{
				StudentNIM:  studentNIM,
				StudentName: studentNIM,
				JoinedAt:    time.Now(),
			}
		}

		if selectedOpt == q.CorrectAnswer {
			baseScore := 500
			totalSec := float64(session.TimePerQuestion)
			if totalSec <= 0 {
				totalSec = 20.0
			}
			if answeredInSec < 0 {
				answeredInSec = 0
			}
			speedBonus := int(500.0 * (1.0 - (answeredInSec / totalSec)))
			if speedBonus < 0 {
				speedBonus = 0
			}

			scoreEarned = baseScore + speedBonus
			part.StreakCount++
			if part.StreakCount >= 3 {
				scoreEarned += 200 // Streak bonus
			}
			part.TotalScore += scoreEarned
		} else {
			part.StreakCount = 0
		}
		session.Participants[studentNIM] = part
	}

	session.Votes = append(session.Votes, models.LiveVote{
		QuestionIndex: questionIndex,
		StudentNIM:    studentNIM,
		SelectedOpt:   selectedOpt,
		AnsweredInSec: answeredInSec,
		ScoreEarned:   scoreEarned,
	})

	return session, nil
}

func (ds *DataStore) GetLiveQuizSession(pin string) (*models.LiveQuizSession, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	session, exists := ds.LiveSessions[pin]
	if !exists {
		return nil, errors.New("sesi kuis tidak ditemukan")
	}
	return session, nil
}

func (ds *DataStore) GetActiveLiveQuizSessionByCourse(courseID string) *models.LiveQuizSession {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	for _, s := range ds.LiveSessions {
		if (s.CourseID == courseID || courseID == "") && s.Status != "FINISHED" {
			return s
		}
	}
	return nil
}


