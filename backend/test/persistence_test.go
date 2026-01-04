package test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"smart-forest-system/backend/controller"
	"smart-forest-system/backend/db"
	"smart-forest-system/backend/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDB 测试数据库连接
var TestDB *sql.DB

// TestMain 测试入口
func TestMain(m *testing.M) {
	// 初始化数据库连接
	TestDB = db.InitDBWithConfig("../db/config.yaml")
	if TestDB == nil {
		log.Fatal("数据库连接失败")
	}
	defer TestDB.Close()

	// 运行测试
	m.Run()
}

// TestSuite 测试套件
type TestSuite struct {
	db                 *sql.DB
	sensorController   *controller.SensorController
	alertController    *controller.AlertController
	resourceController *controller.ResourceController
	deviceController   *controller.DeviceController
	ctx                context.Context
}

// NewTestSuite 创建测试套件
func NewTestSuite(db *sql.DB) *TestSuite {
	return &TestSuite{
		db:                 db,
		sensorController:   controller.NewSensorController(db),
		alertController:    controller.NewAlertController(db),
		resourceController: controller.NewResourceController(db),
		deviceController:   controller.NewDeviceController(db),
		ctx:                context.Background(),
	}
}

// SetupTestData 准备测试数据
func (ts *TestSuite) SetupTestData(t *testing.T) {
	// 清理测试数据
	ts.CleanupTestData(t)

	// 插入测试用户（先插入父表）
	_, err := ts.db.Exec(`
		INSERT INTO users (user_id, username, password_hash, email, real_name, phone, status, created_at, updated_at)
		VALUES
		(1, 'admin', '$2a$10$example.hash', 'admin@test.com', '管理员', '13800138000', 'ACTIVE', NOW(), NOW()),
		(2, 'manager', '$2a$10$example.hash', 'manager@test.com', '区域经理', '13800138001', 'ACTIVE', NOW(), NOW())
	`)
	require.NoError(t, err)

	// 插入测试区域（引用用户ID）
	_, err = ts.db.Exec(`
		INSERT INTO regions (region_id, region_name, region_type, latitude, longitude, manager_id, created_at)
		VALUES
		(1, '测试林区', 'FOREST', 39.9042, 116.4074, 1, NOW()),
		(2, '测试草场', 'GRASSLAND', 40.7128, -74.0060, 2, NOW())
	`)
	require.NoError(t, err)

	// 插入测试传感器
	_, err = ts.db.Exec(`
		INSERT INTO sensors (sensor_id, region_id, model, monitor_type, install_time, protocol, status, created_at)
		VALUES
		(1, 1, 'TEMP-001', 'TEMPERATURE', '2024-01-01 00:00:00', 'MQTT', 'ACTIVE', NOW()),
		(2, 1, 'HUMI-001', 'HUMIDITY', '2024-01-01 00:00:00', 'MQTT', 'ACTIVE', NOW())
	`)
	require.NoError(t, err)

	// 插入测试预警规则
	_, err = ts.db.Exec(`
		INSERT INTO alert_rules (rule_id, alert_type, condition_expr, severity_level, is_active, created_at)
		VALUES
		(1, 'FIRE', 'temperature > 50', 'CRITICAL', 1, NOW()),
		(2, 'DROUGHT', 'humidity < 20', 'MODERATE', 1, NOW())
	`)
	require.NoError(t, err)

	// 插入测试设备
	_, err = ts.db.Exec(`
		INSERT INTO devices (device_id, device_name, device_type, model_spec, purchased_at, install_region_id, installer_id, warranty_until, created_at)
		VALUES
		(2001, '温度传感器1', 'SENSOR', 'TEMP-001', '2024-01-01', 1, 1, '2025-12-31', NOW()),
		(2002, '湿度传感器1', 'SENSOR', 'HUMI-001', '2024-01-01', 1, 1, '2025-12-31', NOW())
	`)
	require.NoError(t, err)

	// 插入测试资源
	_, err = ts.db.Exec(`
		INSERT INTO resources (resource_id, resource_type, region_id, species_name, quantity, area, growth_stage, planted_at, updated_at)
		VALUES
		(1, 'TREE', 1, '松树', 100, 10.5, 'MATURE', '2020-01-01', NOW()),
		(2, 'TREE', 1, '杨树', 50, 5.2, 'GROWING', '2023-01-01', NOW())
	`)
	require.NoError(t, err)
}

// CleanupTestData 清理测试数据
func (ts *TestSuite) CleanupTestData(t *testing.T) {
	// 禁用外键检查
	ts.db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// 删除所有测试数据
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
		"users",
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

// TestSensorOperations 测试传感器相关持久层操作
func TestSensorOperations(t *testing.T) {
	ts := NewTestSuite(TestDB)
	ts.SetupTestData(t)
	defer ts.CleanupTestData(t)

	t.Run("CreateSensor", func(t *testing.T) {
		sensor := model.Sensor{
			RegionID:    1,
			Model:       "TEST-001",
			MonitorType: "TEMPERATURE",
			InstallTime: time.Now(),
			Protocol:    "HTTP",
			Status:      "ACTIVE",
		}

		err := ts.sensorController.CreateSensor(ts.ctx, sensor)
		assert.NoError(t, err)
	})

	t.Run("GetSensorByID", func(t *testing.T) {
		sensor, err := ts.sensorController.GetSensorByID(ts.ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, sensor)
		assert.Equal(t, int64(1), sensor.ID)
		assert.Equal(t, "TEMP-001", sensor.Model)
	})

	t.Run("UpdateSensor", func(t *testing.T) {
		sensor := model.Sensor{
			ID:          1,
			RegionID:    1,
			Model:       "TEMP-002",
			MonitorType: "TEMPERATURE",
			InstallTime: time.Now(),
			Protocol:    "MQTT",
			Status:      "INACTIVE",
		}

		err := ts.sensorController.UpdateSensor(ts.ctx, sensor)
		assert.NoError(t, err)

		// 验证更新
		updated, err := ts.sensorController.GetSensorByID(ts.ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, "TEMP-002", updated.Model)
		assert.Equal(t, "INACTIVE", updated.Status)
	})

	t.Run("ListSensors", func(t *testing.T) {
		monitorType := "TEMPERATURE"
		sensors, err := ts.sensorController.ListSensors(ts.ctx, nil, &monitorType, nil)
		assert.NoError(t, err)
		assert.True(t, len(sensors) >= 1)

		// 验证所有返回的传感器都是温度传感器
		for _, sensor := range sensors {
			assert.Equal(t, "TEMPERATURE", sensor.MonitorType)
		}
	})

	t.Run("InsertSensorReading", func(t *testing.T) {
		reading := model.SensorReading{
			SensorID:    1,
			CollectedAt: time.Now(),
			ReadingType: "TEMPERATURE",
			Numeric:     &[]float64{25.5}[0],
			ValueUnit:   &[]string{"°C"}[0],
			DataStatus:  "VALID",
		}

		err := ts.sensorController.InsertSensorReading(ts.ctx, reading)
		assert.NoError(t, err)
	})

	t.Run("BatchInsertSensorReadings", func(t *testing.T) {
		readings := []model.SensorReading{
			{
				SensorID:    1,
				CollectedAt: time.Now(),
				ReadingType: "TEMPERATURE",
				Numeric:     &[]float64{26.0}[0],
				ValueUnit:   &[]string{"°C"}[0],
				DataStatus:  "VALID",
			},
			{
				SensorID:    2,
				CollectedAt: time.Now(),
				ReadingType: "HUMIDITY",
				Numeric:     &[]float64{65.0}[0],
				ValueUnit:   &[]string{"%"}[0],
				DataStatus:  "VALID",
			},
		}

		err := ts.sensorController.BatchInsertSensorReadings(ts.ctx, readings)
		assert.NoError(t, err)
	})

	t.Run("GetSensorReadingsBySensorID", func(t *testing.T) {
		readings, err := ts.sensorController.GetSensorReadingsBySensorID(ts.ctx, 1, 10)
		assert.NoError(t, err)
		assert.True(t, len(readings) >= 1)
	})

	t.Run("GetSensorReadingsByTimeRange", func(t *testing.T) {
		startTime := time.Now().Add(-time.Hour)
		endTime := time.Now().Add(time.Hour)

		readings, err := ts.sensorController.GetSensorReadingsByTimeRange(ts.ctx, 1, startTime, endTime)
		assert.NoError(t, err)
		assert.True(t, len(readings) >= 0)
	})
}

// TestAlertOperations 测试预警相关持久层操作
func TestAlertOperations(t *testing.T) {
	ts := NewTestSuite(TestDB)
	ts.SetupTestData(t)
	defer ts.CleanupTestData(t)

	t.Run("CreateAlertRule", func(t *testing.T) {
		rule := model.AlertRule{
			AlertType:     "OTHER",
			ConditionExpr: "water_level > 10",
			SeverityLevel: "CRITICAL",
			IsActive:      true,
		}

		ruleID, err := ts.alertController.CreateAlertRule(ts.ctx, rule)
		assert.NoError(t, err)
		assert.True(t, ruleID > 0)
	})

	t.Run("GetAlertRuleByID", func(t *testing.T) {
		rule, err := ts.alertController.GetAlertRuleByID(ts.ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, rule)
		assert.Equal(t, "FIRE", rule.AlertType)
	})

	t.Run("UpdateAlertRule", func(t *testing.T) {
		rule := model.AlertRule{
			ID:            1,
			AlertType:     "DROUGHT",
			ConditionExpr: "temperature > 60",
			SeverityLevel: "SEVERE",
			IsActive:      true,
		}

		err := ts.alertController.UpdateAlertRule(ts.ctx, rule)
		assert.NoError(t, err)

		// 验证更新
		updated, err := ts.alertController.GetAlertRuleByID(ts.ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, "DROUGHT", updated.AlertType)
		assert.Equal(t, "SEVERE", updated.SeverityLevel)
	})

	t.Run("ListAlertRules", func(t *testing.T) {
		isActive := true
		rules, err := ts.alertController.ListAlertRules(ts.ctx, nil, &isActive, nil)
		assert.NoError(t, err)
		assert.True(t, len(rules) >= 1)

		// 验证所有返回的规则都是激活状态
		for _, rule := range rules {
			assert.True(t, rule.IsActive)
		}
	})

	t.Run("CreateAlert", func(t *testing.T) {
		alert := model.Alert{
			RuleID:       1,
			RegionID:     1,
			TriggeredAt:  time.Now(),
			Content:      "测试火灾预警",
			Status:       "PENDING",
			HandlerID:    &[]int64{1}[0],
			HandleResult: &[]string{"已处理"}[0],
		}

		alertID, err := ts.alertController.CreateAlert(ts.ctx, alert)
		assert.NoError(t, err)
		assert.True(t, alertID > 0)
	})

	t.Run("GetAlertByID", func(t *testing.T) {
		// 先创建一条预警用于测试
		alert := model.Alert{
			RuleID:      1,
			RegionID:    1,
			TriggeredAt: time.Now(),
			Content:     "测试预警",
			Status:      "PENDING",
		}
		alertID, _ := ts.alertController.CreateAlert(ts.ctx, alert)

		alertData, err := ts.alertController.GetAlertByID(ts.ctx, alertID)
		assert.NoError(t, err)
		assert.NotNil(t, alertData)
		assert.Equal(t, alertID, alertData.ID)
		assert.Equal(t, "测试预警", alertData.Content)
	})

	t.Run("UpdateAlert", func(t *testing.T) {
		// 先创建一条预警
		alert := model.Alert{
			RuleID:      1,
			RegionID:    1,
			TriggeredAt: time.Now(),
			Content:     "原始预警",
			Status:      "PENDING",
		}
		alertID, _ := ts.alertController.CreateAlert(ts.ctx, alert)

		// 更新预警
		alert.ID = alertID
		alert.Content = "更新的预警"
		alert.Status = "CLOSED"

		err := ts.alertController.UpdateAlert(ts.ctx, alert)
		assert.NoError(t, err)

		// 验证更新
		updated, err := ts.alertController.GetAlertByID(ts.ctx, alertID)
		assert.NoError(t, err)
		assert.Equal(t, "更新的预警", updated.Content)
		assert.Equal(t, "CLOSED", updated.Status)
	})

	t.Run("CreateNotification", func(t *testing.T) {
		// 先创建预警
		alert := model.Alert{
			RuleID:      1,
			RegionID:    1,
			TriggeredAt: time.Now(),
			Content:     "测试预警",
			Status:      "PENDING",
		}
		alertID, _ := ts.alertController.CreateAlert(ts.ctx, alert)

		notification := model.Notification{
			AlertID:       alertID,
			ReceiverID:    1,
			Channel:       "SMS",
			SentAt:        time.Now(),
			ReceiveStatus: "DELIVERED",
		}

		err := ts.alertController.CreateNotification(ts.ctx, notification)
		assert.NoError(t, err)
	})

	t.Run("GetNotificationsByAlertID", func(t *testing.T) {
		// 先创建预警和通知
		alert := model.Alert{
			RuleID:      1,
			RegionID:    1,
			TriggeredAt: time.Now(),
			Content:     "测试预警通知",
			Status:      "PENDING",
		}
		alertID, _ := ts.alertController.CreateAlert(ts.ctx, alert)

		notification := model.Notification{
			AlertID:       alertID,
			ReceiverID:    1,
			Channel:       "SYSTEM",
			SentAt:        time.Now(),
			ReceiveStatus: "DELIVERED",
		}
		ts.alertController.CreateNotification(ts.ctx, notification)

		notifications, err := ts.alertController.GetNotificationsByAlertID(ts.ctx, alertID)
		assert.NoError(t, err)
		assert.True(t, len(notifications) >= 1)

		// 验证通知包含预警内容和接收人姓名
		for _, n := range notifications {
			assert.NotEmpty(t, n.AlertContent)
			assert.NotEmpty(t, n.ReceiverName)
		}
	})

	t.Run("ListAlerts", func(t *testing.T) {
		alerts, err := ts.alertController.ListAlerts(ts.ctx, nil, nil, nil)
		assert.NoError(t, err)
		assert.True(t, len(alerts) >= 0)
	})
}

// TestResourceOperations 测试资源相关持久层操作
func TestResourceOperations(t *testing.T) {
	ts := NewTestSuite(TestDB)
	ts.SetupTestData(t)
	defer ts.CleanupTestData(t)

	t.Run("CreateResource", func(t *testing.T) {
		resource := model.Resource{
			ResourceType: "TREE",
			RegionID:     1,
			SpeciesName:  "测试树种",
			Quantity:     &[]int64{200}[0],
			Area:         &[]float64{15.0}[0],
			GrowthStage:  "SEEDLING",
			PlantedAt:    &[]time.Time{time.Now()}[0],
		}

		resourceID, err := ts.resourceController.CreateResource(ts.ctx, resource)
		assert.NoError(t, err)
		assert.True(t, resourceID > 0)
	})

	t.Run("GetResourceByID", func(t *testing.T) {
		resource, err := ts.resourceController.GetResourceByID(ts.ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, resource)
		assert.Equal(t, int64(1), resource.ID)
		assert.Equal(t, "松树", resource.SpeciesName)
	})

	t.Run("UpdateResource", func(t *testing.T) {
		resource := model.Resource{
			ID:           1,
			ResourceType: "TREE",
			RegionID:     1,
			SpeciesName:  "更新松树",
			Quantity:     &[]int64{150}[0],
			Area:         &[]float64{12.0}[0],
			GrowthStage:  "MATURE",
			PlantedAt:    &[]time.Time{time.Now()}[0],
		}

		err := ts.resourceController.UpdateResource(ts.ctx, resource)
		assert.NoError(t, err)

		// 验证更新
		updated, err := ts.resourceController.GetResourceByID(ts.ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, "更新松树", updated.SpeciesName)
		assert.Equal(t, int64(150), *updated.Quantity)
	})

	t.Run("ListResources", func(t *testing.T) {
		resourceType := "TREE"
		resources, err := ts.resourceController.ListResources(ts.ctx, nil, &resourceType, nil, nil)
		assert.NoError(t, err)
		assert.True(t, len(resources) >= 1)

		// 验证所有返回的资源都是树木类型
		for _, resource := range resources {
			assert.Equal(t, "TREE", resource.ResourceType)
		}
	})

	t.Run("CreateResourceChange", func(t *testing.T) {
		change := model.ResourceChange{
			ResourceID:   1,
			ChangeType:   "ADD",
			ChangeReason: &[]string{"新增种植"}[0],
			ChangeAmount: &[]int{50}[0],
			ChangeArea:   &[]float64{3.0}[0],
			ChangedAt:    time.Now(),
			OperatorID:   &[]int64{1}[0],
		}

		err := ts.resourceController.CreateResourceChange(ts.ctx, change)
		assert.NoError(t, err)
	})

	t.Run("GetResourceChangesByResourceID", func(t *testing.T) {
		changes, err := ts.resourceController.GetResourceChangesByResourceID(ts.ctx, 1)
		assert.NoError(t, err)
		assert.True(t, len(changes) >= 1)

		// 验证变更记录包含操作人姓名
		for _, change := range changes {
			if change.OperatorID != nil {
				assert.NotNil(t, change.OperatorName)
			}
		}
	})

	t.Run("ResourceSummaryWithChanges", func(t *testing.T) {
		summaries, err := ts.resourceController.ResourceSummaryWithChanges(ts.ctx)
		assert.NoError(t, err)
		assert.True(t, len(summaries) >= 0)
	})
}

// TestDeviceOperations 测试设备相关持久层操作
func TestDeviceOperations(t *testing.T) {
	ts := NewTestSuite(TestDB)
	ts.SetupTestData(t)
	defer ts.CleanupTestData(t)

	t.Run("CreateDevice", func(t *testing.T) {
		device := model.Device{
			Name:            "测试设备",
			Type:            "CAMERA",
			ModelSpec:       &[]string{"CAM-001"}[0],
			PurchasedAt:     &[]time.Time{time.Now()}[0],
			InstallRegionID: &[]int64{1}[0],
			InstallerID:     &[]int64{1}[0],
			WarrantyUntil:   &[]time.Time{time.Now().AddDate(1, 0, 0)}[0],
		}

		err := ts.deviceController.CreateDevice(ts.ctx, device)
		assert.NoError(t, err)
	})

	t.Run("GetDeviceByID", func(t *testing.T) {
		device, err := ts.deviceController.GetDeviceByID(ts.ctx, 2001)
		assert.NoError(t, err)
		assert.NotNil(t, device)
		assert.Equal(t, int64(2001), device.ID)
		assert.Equal(t, "温度传感器1", device.Name)
	})

	t.Run("UpdateDevice", func(t *testing.T) {
		device := model.Device{
			ID:              2001,
			Name:            "更新温度传感器1",
			Type:            "SENSOR",
			ModelSpec:       &[]string{"TEMP-002"}[0],
			PurchasedAt:     &[]time.Time{time.Now()}[0],
			InstallRegionID: &[]int64{1}[0],
			InstallerID:     &[]int64{1}[0],
			WarrantyUntil:   &[]time.Time{time.Now().AddDate(1, 0, 0)}[0],
		}

		err := ts.deviceController.UpdateDevice(ts.ctx, device)
		assert.NoError(t, err)

		// 验证更新
		updated, err := ts.deviceController.GetDeviceByID(ts.ctx, 2001)
		assert.NoError(t, err)
		assert.Equal(t, "更新温度传感器1", updated.Name)
		assert.Equal(t, "TEMP-002", *updated.ModelSpec)
	})

	t.Run("ListDevices", func(t *testing.T) {
		deviceType := "SENSOR"
		devices, err := ts.deviceController.ListDevices(ts.ctx, &deviceType, nil)
		assert.NoError(t, err)
		assert.True(t, len(devices) >= 1)

		// 验证所有返回的设备都是传感器类型
		for _, device := range devices {
			assert.Equal(t, "SENSOR", device.Type)
		}
	})

	t.Run("CreateDeviceStatusLog", func(t *testing.T) {
		log := model.DeviceStatusLog{
			DeviceID:       2001,
			CollectedAt:    time.Now(),
			RunStatus:      "NORMAL",
			BatteryPercent: &[]int{85}[0],
			SignalStrength: &[]int{90}[0],
		}

		err := ts.deviceController.CreateDeviceStatusLog(ts.ctx, log)
		assert.NoError(t, err)
	})

	t.Run("GetDeviceStatusLogsByDeviceID", func(t *testing.T) {
		logs, err := ts.deviceController.GetDeviceStatusLogsByDeviceID(ts.ctx, 2001)
		assert.NoError(t, err)
		assert.True(t, len(logs) >= 1)

		// 验证日志按时间倒序排列
		for i := 1; i < len(logs); i++ {
			assert.True(t, logs[i-1].CollectedAt.After(logs[i].CollectedAt) || logs[i-1].CollectedAt.Equal(logs[i].CollectedAt))
		}
	})

	t.Run("GetLatestDeviceStatusLog", func(t *testing.T) {
		// 先创建多条状态日志
		log1 := model.DeviceStatusLog{
			DeviceID:       2001,
			CollectedAt:    time.Now().Add(-time.Hour),
			RunStatus:      "NORMAL",
			BatteryPercent: &[]int{80}[0],
			SignalStrength: &[]int{85}[0],
		}
		log2 := model.DeviceStatusLog{
			DeviceID:       2001,
			CollectedAt:    time.Now(),
			RunStatus:      "NORMAL",
			BatteryPercent: &[]int{85}[0],
			SignalStrength: &[]int{90}[0],
		}

		ts.deviceController.CreateDeviceStatusLog(ts.ctx, log1)
		ts.deviceController.CreateDeviceStatusLog(ts.ctx, log2)

		latest, err := ts.deviceController.GetLatestDeviceStatusLog(ts.ctx, 2001)
		assert.NoError(t, err)
		assert.NotNil(t, latest)
		assert.Equal(t, int(85), *latest.BatteryPercent)
	})

	t.Run("CreateMaintenanceRecord", func(t *testing.T) {
		record := model.MaintenanceRecord{
			DeviceID:        2001,
			MaintenanceType: "INSPECTION",
			MaintenanceTime: time.Now(),
			MaintainerID:    &[]int64{1}[0],
			Content:         &[]string{"定期检查设备状态"}[0],
			Result:          &[]string{"设备运行正常"}[0],
		}

		err := ts.deviceController.CreateMaintenanceRecord(ts.ctx, record)
		assert.NoError(t, err)
	})

	t.Run("GetMaintenanceRecordsByDeviceID", func(t *testing.T) {
		records, err := ts.deviceController.GetMaintenanceRecordsByDeviceID(ts.ctx, 2001)
		assert.NoError(t, err)
		assert.True(t, len(records) >= 1)

		// 验证维护记录包含操作人姓名
		for _, record := range records {
			if record.MaintainerID != nil {
				assert.NotNil(t, record.OperatorName)
			}
		}
	})
}

// TestTransactionOperations 测试事务操作
func TestTransactionOperations(t *testing.T) {
	ts := NewTestSuite(TestDB)
	ts.SetupTestData(t)
	defer ts.CleanupTestData(t)

	t.Run("SensorDeletionWithReadings", func(t *testing.T) {
		// 创建测试传感器
		sensor := model.Sensor{
			RegionID:    1,
			Model:       "TEST-SENSOR",
			MonitorType: "TEMPERATURE",
			InstallTime: time.Now(),
			Protocol:    "MQTT",
			Status:      "ACTIVE",
		}
		err := ts.sensorController.CreateSensor(ts.ctx, sensor)
		assert.NoError(t, err)

		// 获取新创建的传感器ID
		var sensorID int64
		err = ts.db.QueryRow("SELECT sensor_id FROM sensors WHERE model = 'TEST-SENSOR'").Scan(&sensorID)
		assert.NoError(t, err)

		// 添加传感器读数
		reading := model.SensorReading{
			SensorID:    sensorID,
			CollectedAt: time.Now(),
			ReadingType: "TEMPERATURE",
			Numeric:     &[]float64{25.0}[0],
			ValueUnit:   &[]string{"°C"}[0],
			DataStatus:  "VALID",
		}
		err = ts.sensorController.InsertSensorReading(ts.ctx, reading)
		assert.NoError(t, err)

		// 删除传感器（应该级联删除读数）
		err = ts.sensorController.DeleteSensor(ts.ctx, sensorID)
		assert.NoError(t, err)

		// 验证传感器和读数都被删除
		_, err = ts.sensorController.GetSensorByID(ts.ctx, sensorID)
		assert.Error(t, err) // 应该找不到记录

		var count int
		ts.db.QueryRow("SELECT COUNT(*) FROM sensor_readings WHERE sensor_id = ?", sensorID).Scan(&count)
		assert.Equal(t, 0, count)
	})

	t.Run("DeviceDeletionWithDependencies", func(t *testing.T) {
		// 创建测试设备
		device := model.Device{
			Name:            "测试删除设备",
			Type:            "SENSOR",
			ModelSpec:       &[]string{"TEST-001"}[0],
			PurchasedAt:     &[]time.Time{time.Now()}[0],
			InstallRegionID: &[]int64{1}[0],
			InstallerID:     &[]int64{1}[0],
			WarrantyUntil:   &[]time.Time{time.Now().AddDate(1, 0, 0)}[0],
		}
		err := ts.deviceController.CreateDevice(ts.ctx, device)
		assert.NoError(t, err)

		// 获取新创建的设备ID
		var deviceID int64
		err = ts.db.QueryRow("SELECT device_id FROM devices WHERE device_name = '测试删除设备'").Scan(&deviceID)
		assert.NoError(t, err)

		// 添加状态日志和维护记录
		statusLog := model.DeviceStatusLog{
			DeviceID:    deviceID,
			CollectedAt: time.Now(),
			RunStatus:   "NORMAL",
		}
		err = ts.deviceController.CreateDeviceStatusLog(ts.ctx, statusLog)
		assert.NoError(t, err)

		maintenanceRecord := model.MaintenanceRecord{
			DeviceID:        deviceID,
			MaintenanceType: "INSPECTION",
			MaintenanceTime: time.Now(),
			MaintainerID:    &[]int64{1}[0],
			Content:         &[]string{"测试维护"}[0],
			Result:          &[]string{"正常"}[0],
		}
		err = ts.deviceController.CreateMaintenanceRecord(ts.ctx, maintenanceRecord)
		assert.NoError(t, err)

		// 删除设备（应该级联删除相关记录）
		err = ts.deviceController.DeleteDevice(ts.ctx, deviceID)
		assert.NoError(t, err)

		// 验证设备和相关记录都被删除
		_, err = ts.deviceController.GetDeviceByID(ts.ctx, deviceID)
		assert.Error(t, err) // 应该找不到记录

		var statusCount, maintenanceCount int
		ts.db.QueryRow("SELECT COUNT(*) FROM device_status_logs WHERE device_id = ?", deviceID).Scan(&statusCount)
		ts.db.QueryRow("SELECT COUNT(*) FROM maintenance_records WHERE device_id = ?", deviceID).Scan(&maintenanceCount)
		assert.Equal(t, 0, statusCount)
		assert.Equal(t, 0, maintenanceCount)
	})

	t.Run("ResourceDeletionWithChanges", func(t *testing.T) {
		// 创建测试资源
		resource := model.Resource{
			ResourceType: "TREE",
			RegionID:     1,
			SpeciesName:  "测试删除树种",
			Quantity:     &[]int64{100}[0],
			Area:         &[]float64{10.0}[0],
			GrowthStage:  "SEEDLING",
			PlantedAt:    &[]time.Time{time.Now()}[0],
		}
		resourceID, err := ts.resourceController.CreateResource(ts.ctx, resource)
		assert.NoError(t, err)

		// 添加资源变更记录
		change := model.ResourceChange{
			ResourceID:   resourceID,
			ChangeType:   "ADD",
			ChangeReason: &[]string{"初始种植"}[0],
			ChangeAmount: &[]int{100}[0],
			ChangeArea:   &[]float64{10.0}[0],
			ChangedAt:    time.Now(),
			OperatorID:   &[]int64{1}[0],
		}
		err = ts.resourceController.CreateResourceChange(ts.ctx, change)
		assert.NoError(t, err)

		// 删除资源（应该级联删除变更记录）
		err = ts.resourceController.DeleteResource(ts.ctx, resourceID)
		assert.NoError(t, err)

		// 验证资源和变更记录都被删除
		_, err = ts.resourceController.GetResourceByID(ts.ctx, resourceID)
		assert.Error(t, err) // 应该找不到记录

		var count int
		ts.db.QueryRow("SELECT COUNT(*) FROM resource_changes WHERE resource_id = ?", resourceID).Scan(&count)
		assert.Equal(t, 0, count)
	})

	t.Run("AlertRuleDeletionWithDependencies", func(t *testing.T) {
		// 创建测试预警规则
		rule := model.AlertRule{
			AlertType:     "PEST",
			ConditionExpr: "test > 100",
			SeverityLevel: "GENERAL",
			IsActive:      true,
		}
		ruleID, err := ts.alertController.CreateAlertRule(ts.ctx, rule)
		assert.NoError(t, err)

		// 创建关联的预警和通知
		alert := model.Alert{
			RuleID:      ruleID,
			RegionID:    1,
			TriggeredAt: time.Now(),
			Content:     "测试预警",
			Status:      "PENDING",
		}
		alertID, err := ts.alertController.CreateAlert(ts.ctx, alert)
		assert.NoError(t, err)

		notification := model.Notification{
			AlertID:       alertID,
			ReceiverID:    1,
			Channel:       "SYSTEM",
			SentAt:        time.Now(),
			ReceiveStatus: "DELIVERED",
		}
		err = ts.alertController.CreateNotification(ts.ctx, notification)
		assert.NoError(t, err)

		// 删除预警规则（应该级联删除预警和通知）
		err = ts.alertController.DeleteAlertRule(ts.ctx, ruleID)
		assert.NoError(t, err)

		// 验证规则、预警和通知都被删除
		_, err = ts.alertController.GetAlertRuleByID(ts.ctx, ruleID)
		assert.Error(t, err) // 应该找不到记录

		var alertCount, notificationCount int
		ts.db.QueryRow("SELECT COUNT(*) FROM alerts WHERE rule_id = ?", ruleID).Scan(&alertCount)
		ts.db.QueryRow("SELECT COUNT(*) FROM notifications WHERE alert_id = ?", alertID).Scan(&notificationCount)
		assert.Equal(t, 0, alertCount)
		assert.Equal(t, 0, notificationCount)
	})
}

// TestDataValidation 测试数据验证
func TestDataValidation(t *testing.T) {
	ts := NewTestSuite(TestDB)
	ts.SetupTestData(t)
	defer ts.CleanupTestData(t)

	t.Run("InvalidSensorID", func(t *testing.T) {
		_, err := ts.sensorController.GetSensorByID(ts.ctx, -1)
		assert.Error(t, err)
	})

	t.Run("InvalidDeviceID", func(t *testing.T) {
		_, err := ts.deviceController.GetDeviceByID(ts.ctx, -1)
		assert.Error(t, err)
	})

	t.Run("InvalidResourceID", func(t *testing.T) {
		_, err := ts.resourceController.GetResourceByID(ts.ctx, -1)
		assert.Error(t, err)
	})

	t.Run("InvalidAlertRuleID", func(t *testing.T) {
		_, err := ts.alertController.GetAlertRuleByID(ts.ctx, -1)
		assert.Error(t, err)
	})

	t.Run("UpdateNonExistentRecord", func(t *testing.T) {
		sensor := model.Sensor{
			ID:          99999,
			RegionID:    1,
			Model:       "NONEXISTENT",
			MonitorType: "TEMPERATURE",
			InstallTime: time.Now(),
			Protocol:    "MQTT",
			Status:      "ACTIVE",
		}

		err := ts.sensorController.UpdateSensor(ts.ctx, sensor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no sensor updated")
	})

	t.Run("DeleteNonExistentRecord", func(t *testing.T) {
		err := ts.sensorController.DeleteSensor(ts.ctx, 99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "传感器不存在")
	})
}

// TestPerformanceOperations 测试性能相关操作
func TestPerformanceOperations(t *testing.T) {
	ts := NewTestSuite(TestDB)
	ts.SetupTestData(t)
	defer ts.CleanupTestData(t)

	t.Run("BatchSensorReadingsPerformance", func(t *testing.T) {
		// 准备大量传感器读数数据
		readings := make([]model.SensorReading, 100)
		now := time.Now()
		for i := 0; i < 100; i++ {
			readings[i] = model.SensorReading{
				SensorID:    1,
				CollectedAt: now.Add(time.Duration(i) * time.Minute),
				ReadingType: "TEMPERATURE",
				Numeric:     &[]float64{20.0 + float64(i)*0.1}[0],
				ValueUnit:   &[]string{"°C"}[0],
				DataStatus:  "VALID",
			}
		}

		startTime := time.Now()
		err := ts.sensorController.BatchInsertSensorReadings(ts.ctx, readings)
		duration := time.Since(startTime)

		assert.NoError(t, err)
		t.Logf("批量插入100条传感器读数耗时: %v", duration)
		assert.True(t, duration < 5*time.Second, "批量插入不应超过5秒")
	})

	t.Run("LargeResultSetQuery", func(t *testing.T) {
		// 插入大量传感器读数用于测试查询性能
		readings := make([]model.SensorReading, 1000)
		now := time.Now()
		for i := 0; i < 1000; i++ {
			readings[i] = model.SensorReading{
				SensorID:    1,
				CollectedAt: now.Add(time.Duration(i) * time.Second),
				ReadingType: "TEMPERATURE",
				Numeric:     &[]float64{20.0 + float64(i)*0.01}[0],
				ValueUnit:   &[]string{"°C"}[0],
				DataStatus:  "VALID",
			}
		}
		ts.sensorController.BatchInsertSensorReadings(ts.ctx, readings)

		// 测试大量结果查询性能
		startTime := time.Now()
		result, err := ts.sensorController.GetSensorReadingsBySensorID(ts.ctx, 1, 1000)
		duration := time.Since(startTime)

		assert.NoError(t, err)
		assert.True(t, len(result) >= 1000)
		t.Logf("查询1000条传感器读数耗时: %v", duration)
		assert.True(t, duration < 3*time.Second, "查询不应超过3秒")
	})

	t.Run("ConcurrentOperations", func(t *testing.T) {
		// 测试并发操作
		done := make(chan bool, 3)

		// 并发插入传感器读数
		go func() {
			for i := 0; i < 10; i++ {
				reading := model.SensorReading{
					SensorID:    1,
					CollectedAt: time.Now(),
					ReadingType: "TEMPERATURE",
					Numeric:     &[]float64{25.0}[0],
					ValueUnit:   &[]string{"°C"}[0],
					DataStatus:  "VALID",
				}
				ts.sensorController.InsertSensorReading(ts.ctx, reading)
			}
			done <- true
		}()

		// 并发查询传感器
		go func() {
			for i := 0; i < 10; i++ {
				ts.sensorController.GetSensorByID(ts.ctx, 1)
			}
			done <- true
		}()

		// 并发查询传感器读数
		go func() {
			for i := 0; i < 10; i++ {
				ts.sensorController.GetSensorReadingsBySensorID(ts.ctx, 1, 5)
			}
			done <- true
		}()

		// 等待所有并发操作完成
		for i := 0; i < 3; i++ {
			<-done
		}

		t.Log("并发操作测试完成")
	})
}
