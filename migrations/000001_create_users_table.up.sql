CREATE TABLE IF NOT EXISTS users (
  email VARCHAR(50) NOT NULL,
  password VARCHAR(255) NOT NULL,
  created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (email)
);

-- 插入測試帳號 (密碼明文皆為 password123)
-- bcrypt 雜湊值: $2a$10$fy9n/DpSqGx8MLoH7WWeUOYVMrgkVDiYHUg9pp3WSi5hxdyiaVJuq
INSERT INTO users (email, password, created, updated) VALUES
                                                          ('admin@example.com', '$2a$10$fy9n/DpSqGx8MLoH7WWeUOYVMrgkVDiYHUg9pp3WSi5hxdyiaVJuq', NOW(), NOW()),
                                                          ('user1@example.com', '$2a$10$fy9n/DpSqGx8MLoH7WWeUOYVMrgkVDiYHUg9pp3WSi5hxdyiaVJuq', NOW(), NOW()),
                                                          ('user2@example.com', '$2a$10$fy9n/DpSqGx8MLoH7WWeUOYVMrgkVDiYHUg9pp3WSi5hxdyiaVJuq', NOW(), NOW());