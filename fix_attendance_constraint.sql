-- 考勤表约束修复脚本

-- 1. 首先检查是否存在重复数据
SELECT 
    schedule_id, 
    student_id, 
    COUNT(*) as count,
    GROUP_CONCAT(id) as duplicate_ids
FROM h1_attendance 
GROUP BY schedule_id, student_id 
HAVING COUNT(*) > 1;

-- 2. 删除重复数据（保留ID最小的记录）
DELETE a1 FROM h1_attendance a1
INNER JOIN h1_attendance a2 
WHERE a1.schedule_id = a2.schedule_id 
  AND a1.student_id = a2.student_id 
  AND a1.id > a2.id;

-- 3. 检查是否已存在唯一约束
SELECT 
    CONSTRAINT_NAME, 
    CONSTRAINT_TYPE
FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS 
WHERE TABLE_SCHEMA = 'demo' 
  AND TABLE_NAME = 'h1_attendance' 
  AND CONSTRAINT_TYPE = 'UNIQUE';

-- 4. 如果约束不存在，则添加
-- 如果上面查询没有返回 unique_student_schedule 约束，执行下面的语句：
ALTER TABLE h1_attendance 
ADD CONSTRAINT unique_student_schedule 
UNIQUE KEY (schedule_id, student_id);

-- 5. 验证约束是否添加成功
SHOW CREATE TABLE h1_attendance;

-- 6. 测试约束是否生效（这个命令应该失败）
-- INSERT INTO h1_attendance (schedule_id, student_id, status, checked_at) 
-- VALUES (1, 1, '出勤', NOW()), (1, 1, '出勤', NOW()); 