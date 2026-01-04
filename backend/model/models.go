package model

import (
	"database/sql"
	"time"
)

// Core entities (trimmed to fields used in persistence methods).
type Store struct {
	db *sql.DB
}

type Region struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // FOREST or GRASSLAND
	Latitude  *float64  `json:"latitude"`
	Longitude *float64  `json:"longitude"`
	ManagerID *int64    `json:"manager_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Device struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"` // SENSOR, CAMERA, ALERTER, OTHER
	ModelSpec       *string   `json:"model_spec"`
	PurchasedAt     *time.Time `json:"purchased_at"`
	InstallRegionID *int64    `json:"install_region_id"`
	InstallerID     *int64    `json:"installer_id"`
	WarrantyUntil   *time.Time `json:"warranty_until"`
	CreatedAt       time.Time `json:"created_at"`
}

type DeviceStatusLog struct {
	ID             int64     `json:"id"`
	DeviceID       int64     `json:"device_id"`
	CollectedAt    time.Time `json:"collected_at"`
	RunStatus      string    `json:"run_status"` // NORMAL, FAULT, OFFLINE
	BatteryPercent *int      `json:"battery_percent"`
	SignalStrength *int      `json:"signal_strength"`
}

type MaintenanceRecord struct {
	ID              int64     `json:"id"`
	DeviceID        int64     `json:"device_id"`
	MaintenanceType string    `json:"maintenance_type"` // INSPECTION, REPAIR, REPLACE
	MaintenanceTime time.Time `json:"maintenance_time"`
	MaintainerID    *int64    `json:"maintainer_id"`
	Content         *string   `json:"content"`
	Result          *string   `json:"result"`
	OperatorName    *string   `json:"operator_name,omitempty"`
}

type Notification struct {
	ID            int64     `json:"id"`
	AlertID       int64     `json:"alert_id"`
	ReceiverID    int64     `json:"receiver_id"`
	Channel       string    `json:"notification_type"` // SMS, SYSTEM
	SentAt        time.Time `json:"created_at"`
	ReceiveStatus string    `json:"receive_status"`
	// 关联查询字段
	AlertContent  string `json:"message"`      // 预警内容
	ReceiverName  string `json:"recipient"`    // 接收人姓名
}

type ResourceChange struct {
	ID           int64     `json:"id"`
	ResourceID   int64     `json:"resource_id"`
	ChangeType   string    `json:"change_type"` // ADD, REDUCE, STATUS_UPDATE
	ChangeReason *string   `json:"change_reason"`
	ChangeAmount *int      `json:"change_amount"`
	ChangeArea   *float64  `json:"change_area"`
	ChangedAt    time.Time `json:"changed_at"`
	OperatorID   *int64    `json:"operator_id"`
	OperatorName *string   `json:"operator_name,omitempty"`
}

type ReportTemplate struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Dimension string    `json:"dimension"` // REGION, TIME, TYPE
	Metrics   string    `json:"metrics"`
	Cycle     string    `json:"cycle"` // DAY, WEEK, MONTH
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type Report struct {
	ID          int64     `json:"id"`
	TemplateID  int64     `json:"template_id"`
	Period      string    `json:"period"`
	GeneratedAt time.Time `json:"generated_at"`
	FilePath    string    `json:"file_path"`
	DataSource  *string   `json:"data_source"`
}
type Sensor struct {
	ID          int64     `json:"id"`
	RegionID    int64     `json:"region_id"`
	Model       string    `json:"model"`
	MonitorType string    `json:"monitor_type"`
	InstallTime time.Time `json:"install_time"`
	Protocol    string    `json:"protocol"`
	Status      string    `json:"status"`
}

type SensorReading struct {
	ID          int64     `json:"id"`
	SensorID    int64     `json:"sensor_id"`
	CollectedAt time.Time `json:"collected_at"`
	ReadingType string    `json:"reading_type"`
	Numeric     *float64  `json:"numeric_value"`
	ValueUnit   *string   `json:"value_unit"`
	MediaPath   *string   `json:"media_path"`
	DataStatus  string    `json:"data_status"`
}

type AlertRule struct {
	ID            int64     `json:"id"`
	AlertType     string    `json:"alert_type"`
	ConditionExpr string    `json:"condition_expr"`
	SeverityLevel string    `json:"severity_level"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

type Alert struct {
	ID           int64     `json:"id"`
	RuleID       int64     `json:"rule_id"`
	RegionID     int64     `json:"region_id"`
	TriggeredAt  time.Time `json:"triggered_at"`
	Content      string    `json:"content"`
	Status       string    `json:"status"`
	HandlerID    *int64    `json:"handler_id"`
	HandleResult *string   `json:"handle_result"`
	AlertType    string    `json:"alert_type"`
	Severity     string    `json:"severity"`
}

type Resource struct {
	ID           int64     `json:"id"`
	ResourceType string    `json:"resource_type"`
	RegionID     int64     `json:"region_id"`
	SpeciesName  string    `json:"species_name"`
	Quantity     *int64    `json:"quantity"`
	Area         *float64  `json:"area"`
	GrowthStage  string    `json:"growth_stage"`
	PlantedAt    *time.Time `json:"planted_at"`
}

type DeviceFaultMaintenance struct {
	RegionName        string     `json:"region_name"`
	DeviceID          int64      `json:"device_id"`
	DeviceName        string     `json:"device_name"`
	FaultTimes        int64      `json:"fault_times"`
	MaintenanceCount  int64      `json:"maintenance_count"`
	LastMaintenanceAt *time.Time `json:"last_maintenance_at"`
}

type FireAlertReport struct {
	RegionName        string     `json:"region_name"`
	AlertID           int64      `json:"alert_id"`
	TriggeredAt       time.Time  `json:"triggered_at"`
	SeverityLevel     string     `json:"severity_level"`
	Status            string     `json:"status"`
	HandlerID         *int64     `json:"handler_id"`
	HandleResult      *string    `json:"handle_result"`
	NotificationsSent int64      `json:"notifications_sent"`
	DeliveredCount    int64      `json:"delivered_count"`
}

type ResourceSummary struct {
	RegionName      string     `json:"region_name"`
	ResourceType    string     `json:"resource_type"`
	CurrentQuantity int64      `json:"current_quantity"`
	CurrentArea     float64    `json:"current_area"`
	ChangeRecords   int64      `json:"change_records"`
	LastChangeAt    *time.Time `json:"last_change_at"`
}

type SensorValidity struct {
	RegionName      string     `json:"region_name"`
	SensorID        int64      `json:"sensor_id"`
	MonitorType     string     `json:"monitor_type"`
	TotalReadings   int64      `json:"total_readings"`
	ValidReadings   int64      `json:"valid_readings"`
	InvalidReadings int64      `json:"invalid_readings"`
	LastReadingAt   *time.Time `json:"last_reading_at"`
}

type RegionAlertStats struct {
	RegionName        string `json:"region_name"`
	Sensors           int64  `json:"sensors"`
	Alerts            int64  `json:"alerts"`
	NotificationsSent int64  `json:"notifications_sent"`
}

// 环境趋势统计
type EnvironmentTrend struct {
	Time        time.Time `json:"time"`
	RegionName  string    `json:"region_name"`
	MetricType  string    `json:"metric_type"` // TEMPERATURE, HUMIDITY, LIGHT, SOIL_MOISTURE, WIND_SPEED
	MetricValue float64   `json:"metric_value"`
	Unit        string    `json:"unit"`
}

// 空气质量统计
type AirQuality struct {
	Time        time.Time `json:"time"`
	RegionName  string    `json:"region_name"`
	PM25        float64   `json:"pm25"`
	PM10        float64   `json:"pm10"`
	CO2         float64   `json:"co2"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	AQI         int       `json:"aqi"`
	AQILevel    string    `json:"aqi_level"`
}

// 异常数据统计
type AnomalyData struct {
	Time        time.Time `json:"time"`
	RegionName  string    `json:"region_name"`
	SensorID    int64     `json:"sensor_id"`
	ReadingType string    `json:"reading_type"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	AnomalyType string    `json:"anomaly_type"` // HIGH, LOW, OUT_OF_RANGE
	Threshold   float64   `json:"threshold"`
}

// 传感器数据汇总
type SensorDataSummary struct {
	SensorID       int64     `json:"sensor_id"`
	RegionName     string    `json:"region_name"`
	MonitorType    string    `json:"monitor_type"`
	TotalReadings  int64     `json:"total_readings"`
	ValidReadings  int64     `json:"valid_readings"`
	AvgValue       float64   `json:"avg_value"`
	MinValue       float64   `json:"min_value"`
	MaxValue       float64   `json:"max_value"`
	LastReading    time.Time `json:"last_reading"`
}

// 区域日均值统计
type RegionDailyAverage struct {
	Date        time.Time `json:"date"`
	RegionName  string    `json:"region_name"`
	MetricType  string    `json:"metric_type"` // TEMPERATURE, HUMIDITY, LIGHT, SOIL_MOISTURE, WIND_SPEED
	AvgValue    float64   `json:"avg_value"`
	Unit        string    `json:"unit"`
	SampleCount int64     `json:"sample_count"`
}

// 用户认证与权限管理

type User struct {
	ID            int64
	Username      string
	PasswordHash  string
	Email         string
	RealName      string
	Phone         *string
	Status        string // ACTIVE, INACTIVE, LOCKED
	LoginAttempts int
	LastLoginAt   *time.Time
	LockedUntil   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// UserRegistration 用户注册请求
type UserRegistration struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Email    string `json:"email" binding:"required,email"`
	RealName string `json:"real_name" binding:"required,min=2,max=50"`
	Phone    string `json:"phone,omitempty"`
	RoleID   int64  `json:"role_id" binding:"required,min=1"`
}

type Role struct {
	ID          int64
	RoleName    string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Permission struct {
	ID              int64
	PermissionName  string
	Resource        string
	Action          string
	Description     *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type UserRole struct {
	UserID    int64
	RoleID    int64
	CreatedAt time.Time
}

type RolePermission struct {
	RoleID       int64
	PermissionID int64
	CreatedAt    time.Time
}
