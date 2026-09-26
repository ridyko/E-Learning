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
				Name:      "AHMAD FAUZI",
				Role:      "mahasiswa",
				Password:  "123456",
				Email:     "ahmad.fauzi@student.swadharma.ac.id",
				Prodi:     "Teknik Informatika (S1)",
				CreatedAt: time.Now(),
			},
			{
				ID:        "usr-std-2",
				Username:  "20260801002",
				Name:      "SITI NURHALIZA",
				Role:      "mahasiswa",
				Password:  "123456",
				Email:     "siti.nurhaliza@student.swadharma.ac.id",
				Prodi:     "Sistem Informasi (S1)",
				CreatedAt: time.Now(),
			},
			{
				ID:        "usr-std-3",
				Username:  "221112019",
				Name:      "LINTANG ANGEL STEFANI",
				Role:      "mahasiswa",
				Password:  "123456",
				Email:     "lintang.angel@swadharma.ac.id",
				Prodi:     "Teknik Informatika (S1)",
				CreatedAt: time.Now(),
			},
			{
				ID:        "usr-std-4",
				Username:  "221112020",
				Name:      "IGNATION SENSEKO MANGGUR",
				Role:      "mahasiswa",
				Password:  "123456",
				Email:     "ignation.senseko@swadharma.ac.id",
				Prodi:     "Teknik Informatika (S1)",
				CreatedAt: time.Now(),
			},
			{
				ID:        "usr-std-5",
				Username:  "231112028",
				Name:      "FAIZ IJLAL ARAYYAN",
				Role:      "mahasiswa",
				Password:  "123456",
				Email:     "faiz.ijlal@swadharma.ac.id",
				Prodi:     "Teknik Informatika (S1)",
				CreatedAt: time.Now(),
			},
			{
				ID:        "usr-std-6",
				Username:  "241112001",
				Name:      "ALDINUS NDRURU",
				Role:      "mahasiswa",
				Password:  "123456",
				Email:     "aldinus.ndruru@swadharma.ac.id",
				Prodi:     "Teknik Informatika (S1)",
				CreatedAt: time.Now(),
			},
			{
				ID:        "usr-std-7",
				Username:  "241112002",
				Name:      "OVAROLDUS SUPRATMAN",
				Role:      "mahasiswa",
				Password:  "123456",
				Email:     "ovaroldus.s@swadharma.ac.id",
				Prodi:     "Teknik Informatika (S1)",
				CreatedAt: time.Now(),
			},
			{
				ID:        "usr-std-8",
				Username:  "241112003",
				Name:      "RADEN DHAFA ADHITYA SOSIAWAN",
				Role:      "mahasiswa",
				Password:  "123456",
				Email:     "raden.dhafa@swadharma.ac.id",
				Prodi:     "Teknik Informatika (S1)",
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
				TotalStudents: 8,
				Status:        "Aktif",
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
				TotalStudents: 8,
				Status:        "Aktif",
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

	dsnCandidates := []string{}
	if dsn != "" {
		dsnCandidates = append(dsnCandidates, dsn)
	}
	dsnCandidates = append(dsnCandidates,
		"root:@tcp(127.0.0.1:3306)/questkom_elearning?parseTime=true",
		"root:@tcp(127.0.0.1:3306)/elearning?parseTime=true",
		"root:root@tcp(127.0.0.1:3306)/elearning?parseTime=true",
	)

	for _, cand := range dsnCandidates {
		fmt.Printf("[DataStore] Attempting MySQL connection with DSN: %s...\n", cand)
		db, err := sql.Open("mysql", cand)
		if err == nil {
			if pingErr := db.Ping(); pingErr == nil {
				DB.SQLDB = db
				fmt.Printf("[DataStore] Connected to MySQL Database successfully! (%s) 🐬\n", cand)
				break
			} else {
				db.Close()
			}
		}
	}

	if DB.SQLDB != nil {
		DB.initTablesMySQL()
		DB.seedMySQLIfEmpty()
		DB.loadFromMySQL()
		fmt.Println("[DataStore] Database sync completed successfully!")
	} else {
		if DB.loadFromDisk() {
			fmt.Println("[DataStore] Loaded persisted database state from data/store.json")
		} else {
			DB.mu.Lock()
			DB.saveToDiskLocked()
			DB.mu.Unlock()
			fmt.Println("[DataStore] Initialized default database state and saved to data/store.json")
		}
	}
}

func (ds *DataStore) initTablesMySQL() {
	if ds.SQLDB == nil {
		return
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(50) NOT NULL PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			name VARCHAR(100) NOT NULL,
			role VARCHAR(20) NOT NULL DEFAULT 'mahasiswa',
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100) DEFAULT NULL,
			prodi VARCHAR(100) DEFAULT NULL,
			phone VARCHAR(50) DEFAULT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS landing_settings (
			id INT AUTO_INCREMENT PRIMARY KEY,
			hero_title TEXT NOT NULL,
			hero_subtitle TEXT NOT NULL,
			institution_badge VARCHAR(100) DEFAULT NULL,
			semester_badge VARCHAR(100) DEFAULT NULL,
			dosen_bio TEXT DEFAULT NULL,
			office_hours VARCHAR(100) DEFAULT NULL,
			room VARCHAR(100) DEFAULT NULL,
			phone VARCHAR(50) DEFAULT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS profiles (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			degree VARCHAR(50) DEFAULT NULL,
			title VARCHAR(100) DEFAULT NULL,
			nidn VARCHAR(50) DEFAULT NULL,
			institution VARCHAR(150) DEFAULT NULL,
			faculty VARCHAR(100) DEFAULT NULL,
			department VARCHAR(100) DEFAULT NULL,
			email VARCHAR(100) DEFAULT NULL,
			phone VARCHAR(50) DEFAULT NULL,
			office_hours VARCHAR(100) DEFAULT NULL,
			room VARCHAR(100) DEFAULT NULL,
			avatar TEXT DEFAULT NULL,
			bio TEXT DEFAULT NULL,
			expertise_areas TEXT DEFAULT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS courses (
			id VARCHAR(50) NOT NULL PRIMARY KEY,
			code VARCHAR(20) NOT NULL,
			name VARCHAR(100) NOT NULL,
			sks INT NOT NULL DEFAULT 3,
			semester VARCHAR(50) DEFAULT NULL,
			class_time VARCHAR(100) DEFAULT NULL,
			room VARCHAR(100) DEFAULT NULL,
			total_students INT NOT NULL DEFAULT 0,
			status VARCHAR(50) NOT NULL DEFAULT 'Aktif',
			description TEXT DEFAULT NULL,
			syllabus LONGTEXT DEFAULT NULL,
			modules LONGTEXT DEFAULT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS announcements (
			id VARCHAR(50) NOT NULL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			category VARCHAR(50) NOT NULL DEFAULT 'Penting',
			course_id VARCHAR(50) NOT NULL DEFAULT 'all',
			content TEXT NOT NULL,
			author VARCHAR(100) DEFAULT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS quizzes (
			id VARCHAR(50) NOT NULL PRIMARY KEY,
			course_id VARCHAR(50) NOT NULL,
			title VARCHAR(255) NOT NULL,
			description TEXT DEFAULT NULL,
			time_limit INT NOT NULL DEFAULT 15,
			questions LONGTEXT DEFAULT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS quiz_submissions (
			id VARCHAR(50) NOT NULL PRIMARY KEY,
			quiz_id VARCHAR(50) NOT NULL,
			student_nim VARCHAR(50) NOT NULL,
			student_name VARCHAR(100) NOT NULL,
			answers LONGTEXT DEFAULT NULL,
			essay_answers LONGTEXT DEFAULT NULL,
			score DOUBLE NOT NULL DEFAULT 0,
			total_score DOUBLE NOT NULL DEFAULT 100,
			submitted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS assignments (
			id VARCHAR(50) NOT NULL PRIMARY KEY,
			course_id VARCHAR(50) NOT NULL,
			meeting_no INT NOT NULL,
			title VARCHAR(255) NOT NULL,
			description TEXT DEFAULT NULL,
			due_date VARCHAR(50) DEFAULT NULL,
			student_nim VARCHAR(50) NOT NULL,
			student_name VARCHAR(100) NOT NULL,
			repo_link TEXT DEFAULT NULL,
			notes TEXT DEFAULT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'Submitted',
			grade VARCHAR(50) DEFAULT NULL,
			submitted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS attendances (
			id VARCHAR(50) NOT NULL PRIMARY KEY,
			course_id VARCHAR(50) NOT NULL,
			meeting_no INT NOT NULL,
			student_nim VARCHAR(50) NOT NULL,
			student_name VARCHAR(100) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'Hadir',
			check_in_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
	}

	for _, q := range queries {
		if _, err := ds.SQLDB.Exec(q); err != nil {
			fmt.Printf("[DataStore] DDL table creation notice: %v\n", err)
		}
	}
}

func (ds *DataStore) seedMySQLIfEmpty() {
	if ds.SQLDB == nil {
		return
	}

	// 1. Users
	var userCount int
	_ = ds.SQLDB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if userCount == 0 && len(ds.Users) > 0 {
		for _, u := range ds.Users {
			ds.saveUserMySQL(u)
		}
		fmt.Println("[DataStore] Seeded default users into MySQL")
	}

	// 2. Landing Settings
	var landingCount int
	_ = ds.SQLDB.QueryRow("SELECT COUNT(*) FROM landing_settings").Scan(&landingCount)
	if landingCount == 0 {
		ds.saveLandingSettingsMySQL(ds.LandingSettings)
		fmt.Println("[DataStore] Seeded landing settings into MySQL")
	}

	// 3. Profile
	var profileCount int
	_ = ds.SQLDB.QueryRow("SELECT COUNT(*) FROM profiles").Scan(&profileCount)
	if profileCount == 0 {
		ds.saveProfileMySQL(ds.Profile)
		fmt.Println("[DataStore] Seeded profile into MySQL")
	}

	// 4. Courses
	var courseCount int
	_ = ds.SQLDB.QueryRow("SELECT COUNT(*) FROM courses").Scan(&courseCount)
	if courseCount == 0 && len(ds.Courses) > 0 {
		for _, c := range ds.Courses {
			ds.saveCourseMySQL(c)
		}
		fmt.Println("[DataStore] Seeded courses into MySQL")
	}

	// 5. Announcements
	var annCount int
	_ = ds.SQLDB.QueryRow("SELECT COUNT(*) FROM announcements").Scan(&annCount)
	if annCount == 0 && len(ds.Announcements) > 0 {
		for _, a := range ds.Announcements {
			ds.saveAnnouncementMySQL(a)
		}
		fmt.Println("[DataStore] Seeded announcements into MySQL")
	}

	// 6. Quizzes
	var quizCount int
	_ = ds.SQLDB.QueryRow("SELECT COUNT(*) FROM quizzes").Scan(&quizCount)
	if quizCount == 0 && len(ds.Quizzes) > 0 {
		for _, q := range ds.Quizzes {
			ds.saveQuizMySQL(q)
		}
		fmt.Println("[DataStore] Seeded quizzes into MySQL")
	}

	// 7. Quiz Submissions
	var subCount int
	_ = ds.SQLDB.QueryRow("SELECT COUNT(*) FROM quiz_submissions").Scan(&subCount)
	if subCount == 0 && len(ds.Submissions) > 0 {
		for _, sub := range ds.Submissions {
			ds.saveSubmissionMySQL(sub)
		}
		fmt.Println("[DataStore] Seeded quiz submissions into MySQL")
	}

	// 8. Assignments
	var asgCount int
	_ = ds.SQLDB.QueryRow("SELECT COUNT(*) FROM assignments").Scan(&asgCount)
	if asgCount == 0 && len(ds.Assignments) > 0 {
		for _, asg := range ds.Assignments {
			ds.saveAssignmentMySQL(asg)
		}
		fmt.Println("[DataStore] Seeded assignments into MySQL")
	}

	// 9. Attendances
	var attCount int
	_ = ds.SQLDB.QueryRow("SELECT COUNT(*) FROM attendances").Scan(&attCount)
	if attCount == 0 && len(ds.Attendances) > 0 {
		for _, att := range ds.Attendances {
			ds.saveAttendanceMySQL(att)
		}
		fmt.Println("[DataStore] Seeded attendances into MySQL")
	}
}

func (ds *DataStore) loadFromMySQL() {
	if ds.SQLDB == nil {
		return
	}

	// 1. Users
	rowsU, err := ds.SQLDB.Query("SELECT id, username, name, role, password, email, prodi, phone, created_at FROM users")
	if err == nil {
		var loaded []models.User
		for rowsU.Next() {
			var u models.User
			var email, prodi, phone sql.NullString
			var createdAt time.Time
			_ = rowsU.Scan(&u.ID, &u.Username, &u.Name, &u.Role, &u.Password, &email, &prodi, &phone, &createdAt)
			u.Email = email.String
			u.Prodi = prodi.String
			u.Phone = phone.String
			if !createdAt.IsZero() {
				u.CreatedAt = createdAt
			}
			loaded = append(loaded, u)
		}
		rowsU.Close()
		if len(loaded) > 0 {
			ds.Users = loaded
		}
	}

	// 2. Landing Settings
	rowL := ds.SQLDB.QueryRow("SELECT hero_title, hero_subtitle, institution_badge, semester_badge, dosen_bio, office_hours, room, phone FROM landing_settings LIMIT 1")
	var ls models.LandingSettings
	var ib, sb, db, oh, rm, ph sql.NullString
	if err := rowL.Scan(&ls.HeroTitle, &ls.HeroSubtitle, &ib, &sb, &db, &oh, &rm, &ph); err == nil {
		ls.InstitutionBadge = ib.String
		ls.SemesterBadge = sb.String
		ls.DosenBio = db.String
		ls.OfficeHours = oh.String
		ls.Room = rm.String
		ls.Phone = ph.String
		if ls.HeroTitle != "" {
			ds.LandingSettings = ls
		}
	}

	// 3. Profiles
	rowP := ds.SQLDB.QueryRow("SELECT name, degree, title, nidn, institution, faculty, department, email, phone, office_hours, room, avatar, bio, expertise_areas FROM profiles LIMIT 1")
	var p models.Profile
	var deg, tit, nidn, inst, fac, dep, em, ph2, oh2, rm2, av, bio, expStr sql.NullString
	if err := rowP.Scan(&p.Name, &deg, &tit, &nidn, &inst, &fac, &dep, &em, &ph2, &oh2, &rm2, &av, &bio, &expStr); err == nil {
		p.Degree = deg.String
		p.Title = tit.String
		p.NIDN = nidn.String
		p.Institution = inst.String
		p.Faculty = fac.String
		p.Department = dep.String
		p.Email = em.String
		p.Phone = ph2.String
		p.OfficeHours = oh2.String
		p.Room = rm2.String
		p.Avatar = av.String
		p.Bio = bio.String
		if expStr.Valid && expStr.String != "" {
			_ = json.Unmarshal([]byte(expStr.String), &p.ExpertiseAreas)
		}
		if p.Name != "" {
			ds.Profile = p
		}
	}

	// 4. Courses
	rowsC, err := ds.SQLDB.Query("SELECT id, code, name, sks, semester, class_time, room, total_students, status, description, syllabus, modules FROM courses")
	if err == nil {
		var loaded []models.Course
		for rowsC.Next() {
			var c models.Course
			var sem, ct, rm3, st, desc, sylStr, modStr sql.NullString
			_ = rowsC.Scan(&c.ID, &c.Code, &c.Name, &c.SKS, &sem, &ct, &rm3, &c.TotalStudents, &st, &desc, &sylStr, &modStr)
			c.Semester = sem.String
			c.ClassTime = ct.String
			c.Room = rm3.String
			c.Status = st.String
			if c.Status == "" {
				c.Status = "Aktif"
			}
			c.Description = desc.String
			if sylStr.Valid && sylStr.String != "" {
				_ = json.Unmarshal([]byte(sylStr.String), &c.Syllabus)
			}
			if modStr.Valid && modStr.String != "" {
				_ = json.Unmarshal([]byte(modStr.String), &c.Modules)
			}
			loaded = append(loaded, c)
		}
		rowsC.Close()
		if len(loaded) > 0 {
			ds.Courses = loaded
		}
	}

	// 5. Announcements
	rowsA, err := ds.SQLDB.Query("SELECT id, title, category, course_id, content, author, created_at FROM announcements ORDER BY created_at DESC")
	if err == nil {
		var loaded []models.Announcement
		for rowsA.Next() {
			var a models.Announcement
			var author sql.NullString
			var createdAt time.Time
			_ = rowsA.Scan(&a.ID, &a.Title, &a.Category, &a.CourseID, &a.Content, &author, &createdAt)
			a.Author = author.String
			if !createdAt.IsZero() {
				a.CreatedAt = createdAt
			}
			loaded = append(loaded, a)
		}
		rowsA.Close()
		if len(loaded) > 0 {
			ds.Announcements = loaded
		}
	}

	// 6. Quizzes
	rowsQ, err := ds.SQLDB.Query("SELECT id, course_id, title, description, time_limit, questions FROM quizzes")
	if err == nil {
		var loaded []models.Quiz
		for rowsQ.Next() {
			var q models.Quiz
			var desc, qStr sql.NullString
			_ = rowsQ.Scan(&q.ID, &q.CourseID, &q.Title, &desc, &q.TimeLimit, &qStr)
			q.Description = desc.String
			if qStr.Valid && qStr.String != "" {
				_ = json.Unmarshal([]byte(qStr.String), &q.Questions)
			}
			loaded = append(loaded, q)
		}
		rowsQ.Close()
		if len(loaded) > 0 {
			ds.Quizzes = loaded
		}
	}

	// 7. Quiz Submissions
	rowsS, err := ds.SQLDB.Query("SELECT id, quiz_id, student_nim, student_name, answers, essay_answers, score, total_score, submitted_at FROM quiz_submissions ORDER BY submitted_at DESC")
	if err == nil {
		var loaded []models.QuizSubmission
		for rowsS.Next() {
			var sub models.QuizSubmission
			var ansStr, essayStr sql.NullString
			var subAt time.Time
			_ = rowsS.Scan(&sub.ID, &sub.QuizID, &sub.StudentNIM, &sub.StudentName, &ansStr, &essayStr, &sub.Score, &sub.TotalScore, &subAt)
			if ansStr.Valid && ansStr.String != "" {
				_ = json.Unmarshal([]byte(ansStr.String), &sub.Answers)
			}
			if essayStr.Valid && essayStr.String != "" {
				_ = json.Unmarshal([]byte(essayStr.String), &sub.EssayAnswers)
			}
			if !subAt.IsZero() {
				sub.SubmittedAt = subAt
			}
			loaded = append(loaded, sub)
		}
		rowsS.Close()
		if len(loaded) > 0 {
			ds.Submissions = loaded
		}
	}

	// 8. Assignments
	rowsAsg, err := ds.SQLDB.Query("SELECT id, course_id, meeting_no, title, description, due_date, student_nim, student_name, repo_link, notes, status, grade, submitted_at FROM assignments ORDER BY submitted_at DESC")
	if err == nil {
		var loaded []models.Assignment
		for rowsAsg.Next() {
			var asg models.Assignment
			var desc, dueDate, repoLink, notes, grade sql.NullString
			var subAt time.Time
			_ = rowsAsg.Scan(&asg.ID, &asg.CourseID, &asg.MeetingNo, &asg.Title, &desc, &dueDate, &asg.StudentNIM, &asg.StudentName, &repoLink, &notes, &asg.Status, &grade, &subAt)
			asg.Description = desc.String
			asg.DueDate = dueDate.String
			asg.RepoLink = repoLink.String
			asg.Notes = notes.String
			asg.Grade = grade.String
			if !subAt.IsZero() {
				asg.SubmittedAt = subAt
			}
			loaded = append(loaded, asg)
		}
		rowsAsg.Close()
		if len(loaded) > 0 {
			ds.Assignments = loaded
		}
	}

	// 9. Attendances
	rowsAtt, err := ds.SQLDB.Query("SELECT id, course_id, meeting_no, student_nim, student_name, status, check_in_time FROM attendances")
	if err == nil {
		var loaded []models.Attendance
		for rowsAtt.Next() {
			var att models.Attendance
			var checkIn time.Time
			_ = rowsAtt.Scan(&att.ID, &att.CourseID, &att.MeetingNo, &att.StudentNIM, &att.StudentName, &att.Status, &checkIn)
			if !checkIn.IsZero() {
				att.CheckInTime = checkIn
			}
			loaded = append(loaded, att)
		}
		rowsAtt.Close()
		if len(loaded) > 0 {
			ds.Attendances = loaded
		}
	}
}

// MySQL persistence helper functions
func (ds *DataStore) saveUserMySQL(u models.User) {
	if ds.SQLDB == nil {
		return
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	query := `INSERT INTO users (id, username, name, role, password, email, prodi, phone, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE name=VALUES(name), role=VALUES(role), password=VALUES(password), email=VALUES(email), prodi=VALUES(prodi), phone=VALUES(phone)`
	_, err := ds.SQLDB.Exec(query, u.ID, u.Username, u.Name, u.Role, u.Password, u.Email, u.Prodi, u.Phone, u.CreatedAt)
	if err != nil {
		fmt.Printf("[DataStore] saveUserMySQL error: %v\n", err)
	}
}

func (ds *DataStore) deleteUserMySQL(identifier string) {
	if ds.SQLDB == nil {
		return
	}
	_, err := ds.SQLDB.Exec("DELETE FROM users WHERE username = ? OR id = ?", identifier, identifier)
	if err != nil {
		fmt.Printf("[DataStore] deleteUserMySQL error: %v\n", err)
	}
}

func (ds *DataStore) saveLandingSettingsMySQL(s models.LandingSettings) {
	if ds.SQLDB == nil {
		return
	}
	query := `INSERT INTO landing_settings (id, hero_title, hero_subtitle, institution_badge, semester_badge, dosen_bio, office_hours, room, phone)
		VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE hero_title=VALUES(hero_title), hero_subtitle=VALUES(hero_subtitle), institution_badge=VALUES(institution_badge), semester_badge=VALUES(semester_badge), dosen_bio=VALUES(dosen_bio), office_hours=VALUES(office_hours), room=VALUES(room), phone=VALUES(phone)`
	_, err := ds.SQLDB.Exec(query, s.HeroTitle, s.HeroSubtitle, s.InstitutionBadge, s.SemesterBadge, s.DosenBio, s.OfficeHours, s.Room, s.Phone)
	if err != nil {
		fmt.Printf("[DataStore] saveLandingSettingsMySQL error: %v\n", err)
	}
}

func (ds *DataStore) saveProfileMySQL(p models.Profile) {
	if ds.SQLDB == nil {
		return
	}
	expBytes, _ := json.Marshal(p.ExpertiseAreas)
	query := `INSERT INTO profiles (id, name, degree, title, nidn, institution, faculty, department, email, phone, office_hours, room, avatar, bio, expertise_areas)
		VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE name=VALUES(name), degree=VALUES(degree), title=VALUES(title), nidn=VALUES(nidn), institution=VALUES(institution), faculty=VALUES(faculty), department=VALUES(department), email=VALUES(email), phone=VALUES(phone), office_hours=VALUES(office_hours), room=VALUES(room), avatar=VALUES(avatar), bio=VALUES(bio), expertise_areas=VALUES(expertise_areas)`
	_, err := ds.SQLDB.Exec(query, p.Name, p.Degree, p.Title, p.NIDN, p.Institution, p.Faculty, p.Department, p.Email, p.Phone, p.OfficeHours, p.Room, p.Avatar, p.Bio, string(expBytes))
	if err != nil {
		fmt.Printf("[DataStore] saveProfileMySQL error: %v\n", err)
	}
}

func (ds *DataStore) saveCourseMySQL(c models.Course) {
	if ds.SQLDB == nil {
		return
	}
	sylBytes, _ := json.Marshal(c.Syllabus)
	modBytes, _ := json.Marshal(c.Modules)
	if c.Status == "" {
		c.Status = "Aktif"
	}
	query := `INSERT INTO courses (id, code, name, sks, semester, class_time, room, total_students, status, description, syllabus, modules)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE code=VALUES(code), name=VALUES(name), sks=VALUES(sks), semester=VALUES(semester), class_time=VALUES(class_time), room=VALUES(room), total_students=VALUES(total_students), status=VALUES(status), description=VALUES(description), syllabus=VALUES(syllabus), modules=VALUES(modules)`
	_, err := ds.SQLDB.Exec(query, c.ID, c.Code, c.Name, c.SKS, c.Semester, c.ClassTime, c.Room, c.TotalStudents, c.Status, c.Description, string(sylBytes), string(modBytes))
	if err != nil {
		fmt.Printf("[DataStore] saveCourseMySQL error: %v\n", err)
	}
}

func (ds *DataStore) deleteCourseMySQL(id string) {
	if ds.SQLDB == nil {
		return
	}
	_, err := ds.SQLDB.Exec("DELETE FROM courses WHERE id = ?", id)
	if err != nil {
		fmt.Printf("[DataStore] deleteCourseMySQL error: %v\n", err)
	}
}

func (ds *DataStore) resetCoursesMySQL(courses []models.Course) {
	if ds.SQLDB == nil {
		return
	}
	_, _ = ds.SQLDB.Exec("DELETE FROM courses")
	for _, c := range courses {
		ds.saveCourseMySQL(c)
	}
}

func (ds *DataStore) saveAnnouncementMySQL(a models.Announcement) {
	if ds.SQLDB == nil {
		return
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}
	query := `INSERT INTO announcements (id, title, category, course_id, content, author, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE title=VALUES(title), category=VALUES(category), course_id=VALUES(course_id), content=VALUES(content), author=VALUES(author)`
	_, err := ds.SQLDB.Exec(query, a.ID, a.Title, a.Category, a.CourseID, a.Content, a.Author, a.CreatedAt)
	if err != nil {
		fmt.Printf("[DataStore] saveAnnouncementMySQL error: %v\n", err)
	}
}

func (ds *DataStore) deleteAnnouncementMySQL(id string) {
	if ds.SQLDB == nil {
		return
	}
	_, err := ds.SQLDB.Exec("DELETE FROM announcements WHERE id = ?", id)
	if err != nil {
		fmt.Printf("[DataStore] deleteAnnouncementMySQL error: %v\n", err)
	}
}

func (ds *DataStore) saveQuizMySQL(q models.Quiz) {
	if ds.SQLDB == nil {
		return
	}
	qBytes, _ := json.Marshal(q.Questions)
	query := `INSERT INTO quizzes (id, course_id, title, description, time_limit, questions)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE course_id=VALUES(course_id), title=VALUES(title), description=VALUES(description), time_limit=VALUES(time_limit), questions=VALUES(questions)`
	_, err := ds.SQLDB.Exec(query, q.ID, q.CourseID, q.Title, q.Description, q.TimeLimit, string(qBytes))
	if err != nil {
		fmt.Printf("[DataStore] saveQuizMySQL error: %v\n", err)
	}
}

func (ds *DataStore) deleteQuizMySQL(id string) {
	if ds.SQLDB == nil {
		return
	}
	_, err := ds.SQLDB.Exec("DELETE FROM quizzes WHERE id = ?", id)
	if err != nil {
		fmt.Printf("[DataStore] deleteQuizMySQL error: %v\n", err)
	}
}

func (ds *DataStore) saveSubmissionMySQL(sub models.QuizSubmission) {
	if ds.SQLDB == nil {
		return
	}
	if sub.SubmittedAt.IsZero() {
		sub.SubmittedAt = time.Now()
	}
	ansBytes, _ := json.Marshal(sub.Answers)
	essayBytes, _ := json.Marshal(sub.EssayAnswers)
	query := `INSERT INTO quiz_submissions (id, quiz_id, student_nim, student_name, answers, essay_answers, score, total_score, submitted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE answers=VALUES(answers), essay_answers=VALUES(essay_answers), score=VALUES(score), total_score=VALUES(total_score)`
	_, err := ds.SQLDB.Exec(query, sub.ID, sub.QuizID, sub.StudentNIM, sub.StudentName, string(ansBytes), string(essayBytes), sub.Score, sub.TotalScore, sub.SubmittedAt)
	if err != nil {
		fmt.Printf("[DataStore] saveSubmissionMySQL error: %v\n", err)
	}
}

func (ds *DataStore) saveAssignmentMySQL(asg models.Assignment) {
	if ds.SQLDB == nil {
		return
	}
	if asg.SubmittedAt.IsZero() {
		asg.SubmittedAt = time.Now()
	}
	query := `INSERT INTO assignments (id, course_id, meeting_no, title, description, due_date, student_nim, student_name, repo_link, notes, status, grade, submitted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE title=VALUES(title), description=VALUES(description), repo_link=VALUES(repo_link), notes=VALUES(notes), status=VALUES(status), grade=VALUES(grade)`
	_, err := ds.SQLDB.Exec(query, asg.ID, asg.CourseID, asg.MeetingNo, asg.Title, asg.Description, asg.DueDate, asg.StudentNIM, asg.StudentName, asg.RepoLink, asg.Notes, asg.Status, asg.Grade, asg.SubmittedAt)
	if err != nil {
		fmt.Printf("[DataStore] saveAssignmentMySQL error: %v\n", err)
	}
}

func (ds *DataStore) saveAttendanceMySQL(att models.Attendance) {
	if ds.SQLDB == nil {
		return
	}
	if att.CheckInTime.IsZero() {
		att.CheckInTime = time.Now()
	}
	query := `INSERT INTO attendances (id, course_id, meeting_no, student_nim, student_name, status, check_in_time)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE status=VALUES(status), check_in_time=VALUES(check_in_time)`
	_, err := ds.SQLDB.Exec(query, att.ID, att.CourseID, att.MeetingNo, att.StudentNIM, att.StudentName, att.Status, att.CheckInTime)
	if err != nil {
		fmt.Printf("[DataStore] saveAttendanceMySQL error: %v\n", err)
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
		for i := range dump.Courses {
			if dump.Courses[i].Status == "" {
				dump.Courses[i].Status = "Aktif"
			}
		}
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
	ds.saveLandingSettingsMySQL(s)
	ds.saveProfileMySQL(ds.Profile)
	return ds.LandingSettings
}

func (ds *DataStore) AddCourse(c models.Course) models.Course {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if c.ID == "" {
		c.ID = "crs-" + time.Now().Format("20060102150405")
	}
	if c.Status == "" {
		c.Status = "Aktif"
	}
	if len(c.Modules) == 0 {
		c.Modules = generateDefaultModules(c.Name)
	}
	ds.Courses = append(ds.Courses, c)
	ds.saveToDiskLocked()
	ds.saveCourseMySQL(c)
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
			if c.Status != "" {
				ds.Courses[i].Status = c.Status
			}
			ds.saveToDiskLocked()
			ds.saveCourseMySQL(ds.Courses[i])
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
			ds.deleteCourseMySQL(id)
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
			Status:        "Aktif",
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
			Status:        "Aktif",
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
	ds.resetCoursesMySQL(ds.Courses)
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
	ds.saveUserMySQL(u)
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

func (ds *DataStore) UpdateStudent(u models.User) (models.User, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for i, existing := range ds.Users {
		if existing.Username == u.Username || existing.ID == u.ID {
			if u.Name != "" {
				ds.Users[i].Name = u.Name
			}
			if u.Prodi != "" {
				ds.Users[i].Prodi = u.Prodi
			}
			if u.Email != "" {
				ds.Users[i].Email = u.Email
			}
			if u.Phone != "" {
				ds.Users[i].Phone = u.Phone
			}
			ds.saveToDiskLocked()
			ds.saveUserMySQL(ds.Users[i])
			return ds.Users[i], nil
		}
	}
	return models.User{}, errors.New("Mahasiswa tidak ditemukan")
}

func (ds *DataStore) DeleteStudent(username string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for i, existing := range ds.Users {
		if (existing.Username == username || existing.ID == username) && existing.Role == "mahasiswa" {
			ds.Users = append(ds.Users[:i], ds.Users[i+1:]...)
			ds.saveToDiskLocked()
			ds.deleteUserMySQL(username)
			return nil
		}
	}
	return errors.New("Mahasiswa tidak ditemukan")
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
			ds.saveUserMySQL(ds.Users[i])
		}
	}
	ds.saveToDiskLocked()
	ds.saveProfileMySQL(p)
}

func (ds *DataStore) GetCourses() []models.Course {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	studentCount := 0
	for _, u := range ds.Users {
		if u.Role == "mahasiswa" {
			studentCount++
		}
	}

	courses := make([]models.Course, len(ds.Courses))
	copy(courses, ds.Courses)
	for i := range courses {
		if studentCount > 0 {
			courses[i].TotalStudents = studentCount
		}
	}
	return courses
}

func (ds *DataStore) GetCourseByID(id string) (models.Course, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	studentCount := 0
	for _, u := range ds.Users {
		if u.Role == "mahasiswa" {
			studentCount++
		}
	}

	for _, c := range ds.Courses {
		if c.ID == id {
			if studentCount > 0 {
				c.TotalStudents = studentCount
			}
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
	ds.saveAnnouncementMySQL(a)
	return a
}

func (ds *DataStore) DeleteAnnouncement(id string) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for i, a := range ds.Announcements {
		if a.ID == id {
			ds.Announcements = append(ds.Announcements[:i], ds.Announcements[i+1:]...)
			ds.saveToDiskLocked()
			ds.deleteAnnouncementMySQL(id)
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
	ds.saveQuizMySQL(q)
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
	ds.saveSubmissionMySQL(sub)
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
	ds.saveAssignmentMySQL(asg)
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
			ds.saveAttendanceMySQL(ds.Attendances[i])
			return ds.Attendances[i]
		}
	}

	att.ID = "att-" + time.Now().Format("20060102150405")
	att.CheckInTime = time.Now()
	ds.Attendances = append([]models.Attendance{att}, ds.Attendances...)
	ds.saveToDiskLocked()
	ds.saveAttendanceMySQL(att)
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
			ds.deleteQuizMySQL(id)
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
			ds.saveQuizMySQL(ds.Quizzes[i])
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
					ds.saveQuizMySQL(ds.Quizzes[i])
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
					ds.saveCourseMySQL(ds.Courses[i])
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
					ds.saveSubmissionMySQL(ds.Submissions[i])
					found = true
					break
				}
			}

			if !found {
				newSub := models.QuizSubmission{
					ID:          subID,
					QuizID:      session.QuizID,
					StudentNIM:  p.StudentNIM,
					StudentName: p.StudentName,
					Answers:     make(map[int]int),
					Score:       scorePct,
					TotalScore:  100.0,
					SubmittedAt: time.Now(),
				}
				ds.Submissions = append(ds.Submissions, newSub)
				ds.saveSubmissionMySQL(newSub)
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


