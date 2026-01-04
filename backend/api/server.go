package main

import (
	"log"

	"smart-forest-system/backend/api/handlers"
	"smart-forest-system/backend/db"
	"smart-forest-system/backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("智慧林草系统服务器启动中...")

	// 初始化数据库连接
	log.Println("正在初始化数据库连接...")
	database := db.InitDB()
	defer database.Close()
	log.Println("数据库连接初始化成功")

	// 创建Gin路由器
	log.Println("正在创建Gin路由器...")
	r := gin.Default()
	log.Println("Gin路由器创建成功")

	// 配置CORS中间件
	r.Use(corsMiddleware())
	log.Println("CORS中间件配置成功")

	// 静态文件服务
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")
	log.Println("静态文件服务配置成功")

	// 创建处理器
	h := handlers.NewHandlers(database)
	log.Println("API处理器创建成功")

	// 认证路由（不需要JWT）
	auth := r.Group("/")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.GET("/roles", h.GetAvailableRoles)

	// 公开API路由
	public := r.Group("/api")
	{
		// 区域管理 - 公开查询
		public.GET("/regions", h.ListRegions)
		public.GET("/regions/:id", h.GetRegion)
	}

	// API路由组（需要JWT认证）
	api := r.Group("/api")
	api.Use(middleware.JWTAuthMiddleware())
	{
		// 用户管理
		api.GET("/user/me", h.GetCurrentUser)
		api.PUT("/user/password", h.UpdatePassword)

		// 区域管理 - 需要认证的操作
		api.POST("/regions", middleware.CheckPermissionMiddleware(database, "regions", "manage"), h.CreateRegion)
		api.PUT("/regions/:id", middleware.CheckPermissionMiddleware(database, "regions", "manage"), h.UpdateRegion)
		api.DELETE("/regions/:id", middleware.CheckPermissionMiddleware(database, "regions", "manage"), h.DeleteRegion)

		// 传感器管理
		api.GET("/sensors", middleware.CheckPermissionMiddleware(database, "sensors", "view"), h.ListSensors)
		api.GET("/sensors/:id", middleware.CheckPermissionMiddleware(database, "sensors", "view"), h.GetSensor)
		api.POST("/sensors", middleware.CheckPermissionMiddleware(database, "sensors", "manage"), h.CreateSensor)
		api.PUT("/sensors/:id", middleware.CheckPermissionMiddleware(database, "sensors", "manage"), h.UpdateSensor)
		api.PATCH("/sensors/:id/status", middleware.CheckPermissionMiddleware(database, "sensors", "manage"), h.UpdateSensorStatus)
		api.DELETE("/sensors/:id", middleware.CheckPermissionMiddleware(database, "sensors", "manage"), h.DeleteSensor)

		// 传感器读数
		api.GET("/sensors/:id/readings", h.GetSensorReadings)
		api.POST("/sensors/:id/readings", h.CreateSensorReading)
		api.POST("/sensors/readings/batch", h.BatchCreateSensorReadings)

		// 预警规则
		api.GET("/alert-rules", middleware.CheckPermissionMiddleware(database, "alerts", "view"), h.ListAlertRules)
		api.GET("/alert-rules/:id", middleware.CheckPermissionMiddleware(database, "alerts", "view"), h.GetAlertRule)
		api.POST("/alert-rules", middleware.CheckPermissionMiddleware(database, "alerts", "manage"), h.CreateAlertRule)
		api.PUT("/alert-rules/:id", middleware.CheckPermissionMiddleware(database, "alerts", "manage"), h.UpdateAlertRule)
		api.DELETE("/alert-rules/:id", middleware.CheckPermissionMiddleware(database, "alerts", "manage"), h.DeleteAlertRule)

		// 预警记录
		api.GET("/alerts", middleware.CheckPermissionMiddleware(database, "alerts", "view"), h.ListAlerts)
		api.GET("/alerts/:id", middleware.CheckPermissionMiddleware(database, "alerts", "view"), h.GetAlert)
		api.POST("/alerts", middleware.CheckPermissionMiddleware(database, "alerts", "manage"), h.CreateAlert)
		api.PUT("/alerts/:id", middleware.CheckPermissionMiddleware(database, "alerts", "manage"), h.UpdateAlert)
		api.DELETE("/alerts/:id", middleware.CheckPermissionMiddleware(database, "alerts", "manage"), h.DeleteAlert)

		// 通知
		api.GET("/alerts/:id/notifications", h.GetNotifications)
		api.POST("/notifications", h.CreateNotification)

		// 资源管理
		api.GET("/resources", middleware.CheckPermissionMiddleware(database, "resources", "view"), h.ListResources)
		api.GET("/resources/:id", middleware.CheckPermissionMiddleware(database, "resources", "view"), h.GetResource)
		api.POST("/resources", middleware.CheckPermissionMiddleware(database, "resources", "manage"), h.CreateResource)
		api.PUT("/resources/:id", middleware.CheckPermissionMiddleware(database, "resources", "manage"), h.UpdateResource)
		api.PATCH("/resources/:id/growth-stage", middleware.CheckPermissionMiddleware(database, "resources", "manage"), h.UpdateResourceGrowthStage)
		api.DELETE("/resources/:id", middleware.CheckPermissionMiddleware(database, "resources", "manage"), h.DeleteResource)

		// 资源变更
		api.GET("/resources/:id/changes", h.GetResourceChanges)
		api.POST("/resource-changes", h.CreateResourceChange)

		// 资源统计
		api.GET("/resources/summary", h.GetResourceSummary)

		// 设备管理
		api.GET("/devices", middleware.CheckPermissionMiddleware(database, "devices", "view"), h.ListDevices)
		api.GET("/devices/:id", middleware.CheckPermissionMiddleware(database, "devices", "view"), h.GetDevice)
		api.POST("/devices", middleware.CheckPermissionMiddleware(database, "devices", "manage"), h.CreateDevice)
		api.PUT("/devices/:id", middleware.CheckPermissionMiddleware(database, "devices", "manage"), h.UpdateDevice)
		api.DELETE("/devices/:id", middleware.CheckPermissionMiddleware(database, "devices", "manage"), h.DeleteDevice)

		// 设备状态
		api.GET("/devices/:id/status", h.GetDeviceStatus)
		api.GET("/devices/:id/status/logs", h.GetDeviceStatusLogs)
		api.POST("/devices/:id/status", h.CreateDeviceStatus)

		// 维护记录
		api.GET("/devices/:id/maintenance", h.GetMaintenanceRecords)
		api.POST("/devices/:id/maintenance", h.CreateMaintenanceRecord)
		api.DELETE("/maintenance/:id", h.DeleteMaintenanceRecord)

		// 报表模板
		api.GET("/report-templates", middleware.CheckPermissionMiddleware(database, "reports", "view"), h.ListReportTemplates)
		api.GET("/report-templates/:id", middleware.CheckPermissionMiddleware(database, "reports", "view"), h.GetReportTemplate)
		api.POST("/report-templates", middleware.CheckPermissionMiddleware(database, "reports", "manage"), h.CreateReportTemplate)
		api.PUT("/report-templates/:id", middleware.CheckPermissionMiddleware(database, "reports", "manage"), h.UpdateReportTemplate)
		api.DELETE("/report-templates/:id", middleware.CheckPermissionMiddleware(database, "reports", "manage"), h.DeleteReportTemplate)

		// 报表
		api.GET("/reports", middleware.CheckPermissionMiddleware(database, "reports", "view"), h.ListReports)
		api.GET("/reports/:id", middleware.CheckPermissionMiddleware(database, "reports", "view"), h.GetReport)
		api.POST("/reports", middleware.CheckPermissionMiddleware(database, "reports", "manage"), h.CreateReport)
		api.PUT("/reports/:id", middleware.CheckPermissionMiddleware(database, "reports", "manage"), h.UpdateReport)
		api.DELETE("/reports/:id", middleware.CheckPermissionMiddleware(database, "reports", "manage"), h.DeleteReport)
		api.POST("/reports/generate", h.GenerateReport)
		api.GET("/reports/:id/export", h.ExportReport)

		// 统计分析
		api.GET("/statistics/device-faults", middleware.CheckPermissionMiddleware(database, "statistics", "view"), h.GetDeviceFaults)
		api.GET("/statistics/sensor-validity", middleware.CheckPermissionMiddleware(database, "statistics", "view"), h.GetSensorValidity)
		api.GET("/statistics/region-alerts", middleware.CheckPermissionMiddleware(database, "statistics", "view"), h.GetRegionAlertStats)
		api.GET("/statistics/fire-alerts", middleware.CheckPermissionMiddleware(database, "statistics", "view"), h.GetFireAlerts)
		// 新增统计分析功能
		api.GET("/statistics/environment-trend", middleware.CheckPermissionMiddleware(database, "statistics", "view"), h.GetEnvironmentTrend)
		api.GET("/statistics/air-quality", middleware.CheckPermissionMiddleware(database, "statistics", "view"), h.GetAirQuality)
		api.GET("/statistics/anomaly-data", middleware.CheckPermissionMiddleware(database, "statistics", "view"), h.GetAnomalyData)
		api.GET("/statistics/sensor-summary", middleware.CheckPermissionMiddleware(database, "statistics", "view"), h.GetSensorDataSummary)
		api.GET("/statistics/region-daily-averages", middleware.CheckPermissionMiddleware(database, "statistics", "view"), h.GetRegionDailyAverages)

		// 角色管理API (只有系统管理员能操作)
		api.POST("/roles", middleware.CheckPermissionMiddleware(database, "roles", "manage"), h.CreateRole)
		api.GET("/roles/:id", middleware.CheckPermissionMiddleware(database, "roles", "view"), h.GetRole)
		api.PUT("/roles/:id", middleware.CheckPermissionMiddleware(database, "roles", "manage"), h.UpdateRole)
		api.DELETE("/roles/:id", middleware.CheckPermissionMiddleware(database, "roles", "manage"), h.DeleteRole)
		api.POST("/users/:user_id/roles", middleware.CheckPermissionMiddleware(database, "users", "manage"), h.AssignRole)
		api.DELETE("/users/:user_id/roles/:role_id", middleware.CheckPermissionMiddleware(database, "users", "manage"), h.RemoveRole)
		api.GET("/users/:user_id/roles", middleware.CheckPermissionMiddleware(database, "users", "view"), h.GetUserRoles)

		// 权限管理API (只有系统管理员能操作)
		api.POST("/permissions", middleware.CheckPermissionMiddleware(database, "permissions", "manage"), h.CreatePermission)
		api.GET("/permissions/:id", middleware.CheckPermissionMiddleware(database, "permissions", "view"), h.GetPermission)
		api.PUT("/permissions/:id", middleware.CheckPermissionMiddleware(database, "permissions", "manage"), h.UpdatePermission)
		api.DELETE("/permissions/:id", middleware.CheckPermissionMiddleware(database, "permissions", "manage"), h.DeletePermission)
		api.POST("/roles/:id/permissions", middleware.CheckPermissionMiddleware(database, "permissions", "manage"), h.AssignPermission)
		api.DELETE("/roles/:id/permissions/:permission_id", middleware.CheckPermissionMiddleware(database, "permissions", "manage"), h.RemovePermission)
		api.GET("/roles/:id/permissions", middleware.CheckPermissionMiddleware(database, "permissions", "view"), h.GetRolePermissions)
		api.GET("/users/:user_id/permissions", middleware.CheckPermissionMiddleware(database, "users", "view"), h.GetUserPermissions)
		api.GET("/users/:user_id/check-permission", middleware.CheckPermissionMiddleware(database, "users", "view"), h.CheckUserPermission)

		// 数据备份与恢复API (只有系统管理员能操作)
		api.POST("/backup", middleware.CheckPermissionMiddleware(database, "permissions", "manage"), h.BackupDatabase)
		api.POST("/restore", middleware.CheckPermissionMiddleware(database, "permissions", "manage"), h.RestoreDatabase)
		api.GET("/backups", middleware.CheckPermissionMiddleware(database, "permissions", "manage"), h.ListBackups)
	}
	log.Println("API路由配置完成")

	// 启动服务器
	log.Println("服务器启动在 http://localhost:8080")
	log.Println("按 Ctrl+C 停止服务器")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
