-- Database Schema for E-Learning ITB Swadharma
-- Compatible with MySQL 5.7+ / MariaDB 10.3+ / cPanel phpMyAdmin

-- 1. Table Users
CREATE TABLE IF NOT EXISTS `users` (
  `id` VARCHAR(50) NOT NULL PRIMARY KEY,
  `username` VARCHAR(50) NOT NULL UNIQUE,
  `name` VARCHAR(100) NOT NULL,
  `role` VARCHAR(20) NOT NULL DEFAULT 'mahasiswa',
  `password` VARCHAR(255) NOT NULL,
  `email` VARCHAR(100) DEFAULT NULL,
  `prodi` VARCHAR(100) DEFAULT NULL,
  `phone` VARCHAR(50) DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. Table Landing Settings
CREATE TABLE IF NOT EXISTS `landing_settings` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `hero_title` TEXT NOT NULL,
  `hero_subtitle` TEXT NOT NULL,
  `institution_badge` VARCHAR(100) DEFAULT NULL,
  `semester_badge` VARCHAR(100) DEFAULT NULL,
  `dosen_bio` TEXT DEFAULT NULL,
  `office_hours` VARCHAR(100) DEFAULT NULL,
  `room` VARCHAR(100) DEFAULT NULL,
  `phone` VARCHAR(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. Table Profile Dosen
CREATE TABLE IF NOT EXISTS `profiles` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(100) NOT NULL,
  `degree` VARCHAR(50) DEFAULT NULL,
  `title` VARCHAR(100) DEFAULT NULL,
  `nidn` VARCHAR(50) DEFAULT NULL,
  `institution` VARCHAR(150) DEFAULT NULL,
  `faculty` VARCHAR(100) DEFAULT NULL,
  `department` VARCHAR(100) DEFAULT NULL,
  `email` VARCHAR(100) DEFAULT NULL,
  `phone` VARCHAR(50) DEFAULT NULL,
  `office_hours` VARCHAR(100) DEFAULT NULL,
  `room` VARCHAR(100) DEFAULT NULL,
  `avatar` TEXT DEFAULT NULL,
  `bio` TEXT DEFAULT NULL,
  `expertise_areas` TEXT DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4. Table Courses
CREATE TABLE IF NOT EXISTS `courses` (
  `id` VARCHAR(50) NOT NULL PRIMARY KEY,
  `code` VARCHAR(20) NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `sks` INT NOT NULL DEFAULT 3,
  `semester` VARCHAR(50) DEFAULT NULL,
  `class_time` VARCHAR(100) DEFAULT NULL,
  `room` VARCHAR(100) DEFAULT NULL,
  `total_students` INT NOT NULL DEFAULT 0,
  `status` VARCHAR(50) NOT NULL DEFAULT 'Aktif',
  `description` TEXT DEFAULT NULL,
  `syllabus` LONGTEXT DEFAULT NULL,
  `modules` LONGTEXT DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 5. Table Announcements
CREATE TABLE IF NOT EXISTS `announcements` (
  `id` VARCHAR(50) NOT NULL PRIMARY KEY,
  `title` VARCHAR(255) NOT NULL,
  `category` VARCHAR(50) NOT NULL DEFAULT 'Penting',
  `course_id` VARCHAR(50) NOT NULL DEFAULT 'all',
  `content` TEXT NOT NULL,
  `author` VARCHAR(100) DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 6. Table Quizzes
CREATE TABLE IF NOT EXISTS `quizzes` (
  `id` VARCHAR(50) NOT NULL PRIMARY KEY,
  `course_id` VARCHAR(50) NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `description` TEXT DEFAULT NULL,
  `time_limit` INT NOT NULL DEFAULT 15,
  `questions` LONGTEXT DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 7. Table Quiz Submissions
CREATE TABLE IF NOT EXISTS `quiz_submissions` (
  `id` VARCHAR(50) NOT NULL PRIMARY KEY,
  `quiz_id` VARCHAR(50) NOT NULL,
  `student_nim` VARCHAR(50) NOT NULL,
  `student_name` VARCHAR(100) NOT NULL,
  `answers` LONGTEXT DEFAULT NULL,
  `essay_answers` LONGTEXT DEFAULT NULL,
  `score` DOUBLE NOT NULL DEFAULT 0,
  `total_score` DOUBLE NOT NULL DEFAULT 100,
  `submitted_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 8. Table Assignments
CREATE TABLE IF NOT EXISTS `assignments` (
  `id` VARCHAR(50) NOT NULL PRIMARY KEY,
  `course_id` VARCHAR(50) NOT NULL,
  `meeting_no` INT NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `description` TEXT DEFAULT NULL,
  `due_date` VARCHAR(50) DEFAULT NULL,
  `student_nim` VARCHAR(50) NOT NULL,
  `student_name` VARCHAR(100) NOT NULL,
  `repo_link` TEXT DEFAULT NULL,
  `notes` TEXT DEFAULT NULL,
  `status` VARCHAR(50) NOT NULL DEFAULT 'Submitted',
  `grade` VARCHAR(50) DEFAULT NULL,
  `submitted_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 9. Table Attendances
CREATE TABLE IF NOT EXISTS `attendances` (
  `id` VARCHAR(50) NOT NULL PRIMARY KEY,
  `course_id` VARCHAR(50) NOT NULL,
  `meeting_no` INT NOT NULL,
  `student_nim` VARCHAR(50) NOT NULL,
  `student_name` VARCHAR(100) NOT NULL,
  `status` VARCHAR(50) NOT NULL DEFAULT 'Hadir',
  `check_in_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Default Seed Data
INSERT INTO `users` (`id`, `username`, `name`, `role`, `password`, `email`, `prodi`, `phone`) VALUES
('usr-dosen-1', '21099001', 'Rio Widyatmoko, S.Kom, M.M.S.I', 'dosen', '123456', 'rio.widyatmoko@swadharma.ac.id', 'Teknik Informatika', '+62 812-9876-5432'),
('usr-std-1', '20260801001', 'AHMAD FAUZI', 'mahasiswa', '123456', 'ahmad.fauzi@student.swadharma.ac.id', 'Teknik Informatika (S1)', NULL),
('usr-std-2', '20260801002', 'SITI NURHALIZA', 'mahasiswa', '123456', 'siti.nurhaliza@student.swadharma.ac.id', 'Sistem Informasi (S1)', NULL)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`);

