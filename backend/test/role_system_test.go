package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"smart-forest-system/backend/api/handlers"
	"smart-forest-system/backend/controller"
	"smart-forest-system/backend/db"
	"smart-forest-system/backend/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// RoleTestDB 角色系统测试数据库连接
var RoleTestDB *sql.DB

// RoleTestSuite 角色系统测试套件
type RoleTestSuite struct {
	db *sql.DB
}

// SetupRoleData 设置角色数据
func (ts *RoleTestSuite) SetupRoleData(t *testing.T) {
	// 插入系统角色数据
	_, err := ts.db.Exec(`
		INSERT INTO roles (role_id, role_name, description) VALUES
		(1, 'SYSTEM_ADMIN', '系统管理员：维护所有角色的账号信息与权限分配，管理设备档案与预警规则，审核统计报表模板'),
		(2, 'DATA_ADMIN', '数据管理员：录入、校验林草资源基础信息，处理监测数据的异常情况，生成并发布统计报表'),
		(3, 'FOREST_GUARD', '区域护林员：查看负责区域的实时监测数据与历史记录，接收并处理相关灾害预警，记录设备巡检与维护情况，更新区域内林草资源变动信息'),
		(4, 'PUBLIC_USER', '公众用户：查看公开的林草资源统计数据，浏览非涉密的环境监测信息（如景区空气质量），提交林草资源异常情况反馈（如发现火情、病虫害）'),
		(5, 'INSPECTOR', '监管人员：查看全系统业务数据与操作记录，监督预警处理流程的及时性与规范性，审核资源变动与设备维护记录的真实性')
	`)
	assert.NoError(t, err)
}

// CleanupTestData 清理测试数据
func (ts *RoleTestSuite) CleanupTestData(t *testing.T) {
	// 禁用外键检查
	ts.db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// 删除测试数据
	tables := []string{
		"sensor_readings",
		"notifications",
		"alerts",
		"alert_rules",
		"resource_changes",
		"resources",
		"maintenance_records",
		"device_status_logs",
		"devices",
		"sensors",
		"user_roles",
		"users",
		"roles",
		"regions",
	}

	for _, table := range tables {
		_, err := ts.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE 1=1", table))
		if err != nil {
			t.Logf("清理表 %s 时出错: %v", table, err)
		}
	}

	// 重新启用外键检查
	ts.db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

// TestRoleSystem 测试角色系统功能
func TestRoleSystem(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	RoleTestDB = db.InitDBWithConfig("../db/config.yaml")
	if RoleTestDB == nil {
		t.Fatal("数据库连接失败")
	}
	defer RoleTestDB.Close()

	// 清理测试数据
	ts := &RoleTestSuite{db: RoleTestDB}
	ts.CleanupTestData(t)

	// 插入系统角色数据
	ts.SetupRoleData(t)

	t.Run("GetAvailableRoles", func(t *testing.T) {
		// 创建处理器
		h := handlers.NewHandlers(TestDB)

		// 创建路由
		r := gin.New()
		r.GET("/roles", h.GetAvailableRoles)

		// 创建请求
		req, _ := http.NewRequest("GET", "/roles", nil)
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		rolesData, exists := response["roles"]
		assert.True(t, exists, "响应应该包含roles字段")

		rolesJson, err := json.Marshal(rolesData)
		assert.NoError(t, err)

		var roles []model.Role
		err = json.Unmarshal(rolesJson, &roles)
		assert.NoError(t, err)

		assert.Equal(t, 5, len(roles)) // 应该有5个系统角色

		// 验证角色名称
		roleNames := make([]string, len(roles))
		for i, role := range roles {
			roleNames[i] = role.RoleName
		}

		expectedRoles := []string{"SYSTEM_ADMIN", "DATA_ADMIN", "FOREST_GUARD", "PUBLIC_USER", "INSPECTOR"}
		assert.ElementsMatch(t, expectedRoles, roleNames)
	})

	t.Run("UserRegistrationWithRole", func(t *testing.T) {
		// 创建处理器
		h := handlers.NewHandlers(TestDB)

		// 创建路由
		r := gin.New()
		r.POST("/register", h.Register)

		// 准备注册数据
		regData := model.UserRegistration{
			Username: "testuser",
			Password: "testpass123",
			Email:    "test@example.com",
			RealName: "测试用户",
			Phone:    "13800138000",
			RoleID:   4, // PUBLIC_USER
		}

		jsonData, _ := json.Marshal(regData)

		// 创建请求
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// 执行请求
		r.ServeHTTP(w, req)

		// 验证响应
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "user_id")
		assert.Equal(t, "用户注册成功", response["message"])
		assert.True(t, response["role_selected"].(bool))

		// 验证数据库中的用户和角色关联
		userID := int64(response["user_id"].(float64))
		authController := controller.NewAuthController(TestDB)

		user, err := authController.GetUserByID(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "测试用户", user.RealName)

		// 验证用户角色
		roles, err := authController.GetUserRoles(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(roles))
		assert.Equal(t, "PUBLIC_USER", roles[0].RoleName)
	})

	t.Run("RolePermissions", func(t *testing.T) {
		authController := controller.NewAuthController(TestDB)

		// 获取数据管理员角色的权限
		permissions, err := authController.GetRolePermissions(context.Background(), 2) // DATA_ADMIN
		assert.NoError(t, err)
		assert.True(t, len(permissions) > 0)

		// 验证数据管理员有查看资源的权限
		hasViewResource := false
		for _, perm := range permissions {
			if perm.Resource == "resources" && perm.Action == "view" {
				hasViewResource = true
				break
			}
		}
		assert.True(t, hasViewResource, "数据管理员应该有查看资源的权限")

		// 获取公众用户的权限
		publicPermissions, err := authController.GetRolePermissions(context.Background(), 4) // PUBLIC_USER
		assert.NoError(t, err)

		// 验证公众用户权限较少
		assert.True(t, len(publicPermissions) < len(permissions), "公众用户权限应该比数据管理员少")
	})

	t.Run("DuplicateUsernameRegistration", func(t *testing.T) {
		h := handlers.NewHandlers(TestDB)
		r := gin.New()
		r.POST("/register", h.Register)

		// 尝试注册重复用户名
		regData := model.UserRegistration{
			Username: "testuser", // 已存在的用户名
			Password: "newpass123",
			Email:    "new@example.com",
			RealName: "新测试用户",
			RoleID:   3, // FOREST_GUARD
		}

		jsonData, _ := json.Marshal(regData)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// 应该返回400错误
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "用户名已存在", response["error"])
	})

	t.Run("InvalidRoleRegistration", func(t *testing.T) {
		h := handlers.NewHandlers(TestDB)
		r := gin.New()
		r.POST("/register", h.Register)

		// 尝试使用不存在的角色ID注册
		regData := model.UserRegistration{
			Username: "newuser",
			Password: "newpass123",
			Email:    "newuser@example.com",
			RealName: "新用户",
			RoleID:   999, // 不存在的角色ID
		}

		jsonData, _ := json.Marshal(regData)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// 应该返回400错误
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "选择的角色不存在", response["error"])
	})
}
