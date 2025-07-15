#!/bin/bash

# 考勤逻辑测试脚本
echo "🔍 测试考勤逻辑修复..."

BASE_URL="http://localhost:8080/api/h1"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "\n${YELLOW}测试: $description${NC}"
    echo "请求: $method $endpoint"
    echo "数据: $data"
    
    if [ "$method" = "POST" ]; then
        response=$(curl -s -X POST "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data" \
            -w "HTTP_CODE:%{http_code}")
    else
        response=$(curl -s "$BASE_URL$endpoint" -w "HTTP_CODE:%{http_code}")
    fi
    
    http_code=$(echo "$response" | grep -o 'HTTP_CODE:[0-9]*' | cut -d: -f2)
    body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✅ 成功 (HTTP $http_code)${NC}"
        echo "响应: $body"
    else
        echo -e "${RED}❌ 失败 (HTTP $http_code)${NC}"
        echo "响应: $body"
    fi
}

echo "请确保后端服务运行在 http://localhost:8080"
echo "开始测试..."

# 测试1: 学生首次打卡
test_endpoint "POST" "/attendance/check-in" \
    '{"schedule_id": 1, "student_id": 1}' \
    "学生首次打卡"

# 测试2: 学生重复打卡（应该更新而不是拒绝）
test_endpoint "POST" "/attendance/check-in" \
    '{"schedule_id": 1, "student_id": 1}' \
    "学生重复打卡"

# 测试3: 管理员标记缺勤
test_endpoint "POST" "/attendance/mark-absent" \
    '{"schedule_id": 1, "student_id": 1, "status": "旷课"}' \
    "管理员标记学生旷课"

# 测试4: 管理员补签（状态变更）
test_endpoint "POST" "/attendance/admin-makeup" \
    '{"schedule_id": 1, "student_id": 1}' \
    "管理员补签（从旷课变为补签）"

# 测试5: 查看考勤记录
test_endpoint "GET" "/attendance?schedule_id=1&student_id=1" \
    "" \
    "查看考勤记录"

# 测试6: 查看学生账单
test_endpoint "GET" "/student-semester-bill/detail?student_id=1&semester_id=1" \
    "" \
    "查看学生账单"

echo -e "\n${YELLOW}🎯 测试重点验证:${NC}"
echo "1. 学生可以重复打卡（更新记录而不是拒绝）"
echo "2. 管理员可以修改考勤状态"
echo "3. 每次操作后账单会自动更新"
echo "4. 一个学生对一个课表只有一条考勤记录"

echo -e "\n${YELLOW}📊 数据库验证命令:${NC}"
echo "-- 检查考勤记录唯一性"
echo "SELECT schedule_id, student_id, COUNT(*) as count FROM h1_attendance GROUP BY schedule_id, student_id HAVING COUNT(*) > 1;"
echo ""
echo "-- 查看特定学生的考勤记录"
echo "SELECT * FROM h1_attendance WHERE student_id = 1 AND schedule_id = 1;"
echo ""
echo "-- 查看学生账单"
echo "SELECT * FROM h1_student_semester_bill WHERE student_id = 1;"

echo -e "\n✨ 测试完成！" 