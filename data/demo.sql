-- h1_schema.sql

-- （1）关闭外键检查，避免删除或创建表时冲突
SET FOREIGN_KEY_CHECKS=0;

-- （2）如果已存在则删除，保证幂等
DROP TABLE IF EXISTS h1_student_semester_bill;
DROP TABLE IF EXISTS h1_attendance;
DROP TABLE IF EXISTS h1_schedule;
DROP TABLE IF EXISTS h1_enrollment;
DROP TABLE IF EXISTS h1_student;
DROP TABLE IF EXISTS h1_class_course;
DROP TABLE IF EXISTS h1_course;
DROP TABLE IF EXISTS h1_class;
DROP TABLE IF EXISTS h1_semester;

-- （3）创建表
CREATE TABLE h1_semester (
                             id BIGINT PRIMARY KEY AUTO_INCREMENT,
                             name VARCHAR(50) NOT NULL,
                             start_date DATE,
                             end_date DATE
);

CREATE TABLE h1_class (
                          id BIGINT PRIMARY KEY AUTO_INCREMENT,
                          name VARCHAR(50) NOT NULL,
                          semester_id BIGINT NOT NULL,
                          FOREIGN KEY (semester_id) REFERENCES h1_semester(id)
);

CREATE TABLE h1_course (
                           id BIGINT PRIMARY KEY AUTO_INCREMENT,
                           name VARCHAR(50) NOT NULL
);

CREATE TABLE h1_class_course (
                                 id BIGINT PRIMARY KEY AUTO_INCREMENT,
                                 class_id BIGINT NOT NULL,
                                 course_id BIGINT NOT NULL,
                                 semester_id BIGINT NOT NULL,
                                 price_per_hour DECIMAL(10,2) NOT NULL,
                                 FOREIGN KEY (class_id) REFERENCES h1_class(id),
                                 FOREIGN KEY (course_id) REFERENCES h1_course(id),
                                 FOREIGN KEY (semester_id) REFERENCES h1_semester(id)
);

CREATE TABLE h1_student (
                            id BIGINT PRIMARY KEY AUTO_INCREMENT,
                            name VARCHAR(50) NOT NULL,
                            contact VARCHAR(100),
                            discount_rate DECIMAL(4,2) DEFAULT 1.00
);

CREATE TABLE h1_enrollment (
                               id BIGINT PRIMARY KEY AUTO_INCREMENT,
                               student_id BIGINT NOT NULL,
                               class_course_id BIGINT NOT NULL,
                               enrolled_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                               FOREIGN KEY (student_id) REFERENCES h1_student(id),
                               FOREIGN KEY (class_course_id) REFERENCES h1_class_course(id)
);

CREATE TABLE h1_schedule (
                             id BIGINT PRIMARY KEY AUTO_INCREMENT,
                             class_course_id BIGINT NOT NULL,
                             date DATE NOT NULL,
                             start_time TIME,
                             end_time TIME,
                             FOREIGN KEY (class_course_id) REFERENCES h1_class_course(id)
);

CREATE TABLE h1_attendance (
                               id BIGINT PRIMARY KEY AUTO_INCREMENT,
                               schedule_id BIGINT NOT NULL,
                               student_id BIGINT NOT NULL,
                               checked_at DATETIME,
                               status ENUM('出勤','请假','旷课','补签')
                                      CHARACTER SET utf8mb4
                                                  NOT NULL,
                               makeup_schedule_id BIGINT NULL,
                               FOREIGN KEY (schedule_id) REFERENCES h1_schedule(id),
                               FOREIGN KEY (student_id) REFERENCES h1_student(id),
                               FOREIGN KEY (makeup_schedule_id) REFERENCES h1_schedule(id),
                               UNIQUE KEY unique_student_schedule (schedule_id, student_id)
) DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 学期账单表（指定 utf8mb4 编码，保证 ENUM 中文不被转为 ?）
CREATE TABLE h1_student_semester_bill (
                                          id BIGINT PRIMARY KEY AUTO_INCREMENT,
                                          student_id BIGINT NOT NULL,
                                          semester_id BIGINT NOT NULL,
                                          total_fee DECIMAL(10,2) NOT NULL,
                                          finalized_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                                          status ENUM('未结算','已结算','已付款')
                                                 CHARACTER SET utf8mb4
                                                 COLLATE utf8mb4_unicode_ci
                                                                DEFAULT '未结算',
                                          FOREIGN KEY (student_id) REFERENCES h1_student(id),
                                          FOREIGN KEY (semester_id) REFERENCES h1_semester(id)
) DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 恢复外键检查
SET FOREIGN_KEY_CHECKS = 1;