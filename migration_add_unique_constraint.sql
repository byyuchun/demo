-- 为考勤表添加唯一约束的迁移脚本
-- 确保一个学生对一个课表只能有一个考勤记录

-- 首先检查是否存在重复数据
SELECT schedule_id, student_id, COUNT(*) as count 
FROM h1_attendance 
GROUP BY schedule_id, student_id 
HAVING COUNT(*) > 1;

-- 如果存在重复数据，需要先处理
-- 保留最新的记录，删除旧的记录
DELETE a1 FROM h1_attendance a1
INNER JOIN h1_attendance a2 
WHERE a1.schedule_id = a2.schedule_id 
  AND a1.student_id = a2.student_id 
  AND a1.id < a2.id;

-- 添加唯一约束
ALTER TABLE h1_attendance 
ADD CONSTRAINT unique_student_schedule 
UNIQUE KEY (schedule_id, student_id);

-- 验证约束是否添加成功
SHOW CREATE TABLE h1_attendance; 