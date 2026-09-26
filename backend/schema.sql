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
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. Table Landing Settings
CREATE TABLE IF NOT EXISTS `landing_settings` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `hero_title` VARCHAR(255) NOT NULL,
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
  `answers` TEXT DEFAULT NULL,
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

-- Data Default (Awal)
INSERT INTO `users` (`id`, `username`, `name`, `role`, `password`, `email`, `prodi`) VALUES
('usr-dosen-1', '21099001', 'Dosen Pengampu', 'dosen', '123456', 'dosen@swadharma.ac.id', 'Teknik Informatika'),
('usr-std-1', '20260801001', 'Ahmad Fauzi', 'mahasiswa', '123456', 'ahmad.fauzi@student.swadharma.ac.id', 'Teknik Informatika'),
('usr-std-2', '20260801002', 'Siti Nurhaliza', 'mahasiswa', '123456', 'siti.nurhaliza@student.swadharma.ac.id', 'Sistem Informasi')
ON DUPLICATE KEY UPDATE `username`=`username`;
