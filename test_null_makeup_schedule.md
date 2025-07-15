# 测试 makeup_schedule_id 的 NULL 处理

## 问题背景
当学生正常打卡时，`makeup_schedule_id` 应该为 NULL，而不是 0。
设置为 0 会导致外键约束错误，因为不存在 ID 为 0 的课表记录。

## 修复方案
1. **正常打卡**：使用 `Omit("makeup_schedule_id")` 忽略该字段，保持 NULL
2. **补签操作**：正常设置 `makeup_schedule_id` 为有效的课表 ID
3. **状态更新**：使用 `Updates` 方法，可以正确设置 NULL 值

## 测试用例

### 1. 正常打卡测试
```bash
curl -X POST http://localhost:8080/api/h1/attendance/check-in \
  -H "Content-Type: application/json" \
  -d '{
    "schedule_id": 1,
    "student_id": 1
  }'
```

**期望结果**：
- 创建考勤记录成功
- `makeup_schedule_id` 为 NULL
- 不出现外键约束错误

### 2. 补签操作测试
```bash
curl -X POST http://localhost:8080/api/h1/attendance/admin-makeup \
  -H "Content-Type: application/json" \
  -d '{
    "schedule_id": 1,
    "student_id": 1,
    "makeup_schedule_id": 2
  }'
```

**期望结果**：
- 更新考勤记录为"补签"状态
- `makeup_schedule_id` 设置为指定的课表 ID
- 账单自动更新

### 3. 状态变更测试
```bash
# 先打卡
curl -X POST http://localhost:8080/api/h1/attendance/check-in \
  -H "Content-Type: application/json" \
  -d '{"schedule_id": 1, "student_id": 1}'

# 再标记旷课  
curl -X POST http://localhost:8080/api/h1/attendance/mark-absent \
  -H "Content-Type: application/json" \
  -d '{
    "schedule_id": 1,
    "student_id": 1,
    "status": "旷课"
  }'

# 最后补签
curl -X POST http://localhost:8080/api/h1/attendance/admin-makeup \
  -H "Content-Type: application/json" \
  -d '{
    "schedule_id": 1,
    "student_id": 1,
    "makeup_schedule_id": 2
  }'
```

**期望结果**：
- 整个过程中只有一条考勤记录
- 状态从"出勤" → "旷课" → "补签"
- `makeup_schedule_id` 从 NULL → NULL → 有效ID

## 数据库验证

```sql
-- 查看考勤记录
SELECT 
    id,
    schedule_id,
    student_id,
    status,
    makeup_schedule_id,
    checked_at
FROM h1_attendance 
WHERE student_id = 1 AND schedule_id = 1;

-- 验证 makeup_schedule_id 为 NULL 的记录
SELECT COUNT(*) as normal_attendance_count
FROM h1_attendance 
WHERE makeup_schedule_id IS NULL;

-- 验证 makeup_schedule_id 不为 NULL 的记录
SELECT COUNT(*) as makeup_attendance_count
FROM h1_attendance 
WHERE makeup_schedule_id IS NOT NULL;
```

## 关键修复点

1. **createOrUpdateAttendance 函数**：
   - 创建时：有补课安排才设置字段，否则使用 `Omit`
   - 更新时：使用 `Updates` 方法支持 NULL 值

2. **CreateAttendance 函数**：
   - 同样使用条件判断和 `Omit` 方法

3. **类型转换处理**：
   - DTO 中使用 `*int64` 支持 NULL 值
   - 模型中使用 `int64`，通过 GORM 方法处理 NULL 