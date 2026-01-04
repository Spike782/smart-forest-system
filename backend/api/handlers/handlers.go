package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"smart-forest-system/backend/controller"
	"smart-forest-system/backend/db"
	"smart-forest-system/backend/model"
	"smart-forest-system/backend/util"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	db                   *sql.DB
	regionController     *controller.RegionController
	sensorController     *controller.SensorController
	alertController      *controller.AlertController
	resourceController   *controller.ResourceController
	deviceController     *controller.DeviceController
	reportController     *controller.ReportController
	statisticsController *controller.StatisticsController
	authController       *controller.AuthController
}

func NewHandlers(db *sql.DB) *Handlers {
	return &Handlers{
		db:                   db,
		regionController:     controller.NewRegionController(db),
		sensorController:     controller.NewSensorController(db),
		alertController:      controller.NewAlertController(db),
		resourceController:   controller.NewResourceController(db),
		deviceController:     controller.NewDeviceController(db),
		reportController:     controller.NewReportController(db),
		statisticsController: controller.NewStatisticsController(db),
		authController:       controller.NewAuthController(db),
	}
}

// 区域管理
func (h *Handlers) ListRegions(c *gin.Context) {
	ctx := context.Background()

	// 解析查询参数
	var regionType *string
	if rt := c.Query("region_type"); rt != "" {
		regionType = &rt
	}

	var name *string
	if n := c.Query("name"); n != "" {
		name = &n
	}

	regions, err := h.regionController.ListRegions(ctx, regionType, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": regions})
}

func (h *Handlers) GetRegion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的区域ID"})
		return
	}
	ctx := context.Background()
	region, err := h.regionController.GetRegionByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": region})
}

func (h *Handlers) CreateRegion(c *gin.Context) {
	var region model.Region
	if err := c.ShouldBindJSON(&region); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	region.CreatedAt = time.Now()
	ctx := context.Background()
	id, err := h.regionController.CreateRegion(ctx, region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "区域创建成功"})
}

func (h *Handlers) UpdateRegion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的区域ID"})
		return
	}
	var region model.Region
	if err := c.ShouldBindJSON(&region); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	region.ID = id
	ctx := context.Background()
	if err := h.regionController.UpdateRegion(ctx, region); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "区域更新成功"})
}

func (h *Handlers) DeleteRegion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的区域ID"})
		return
	}
	ctx := context.Background()
	if err := h.regionController.DeleteRegion(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "区域删除成功"})
}

// 传感器管理
func (h *Handlers) ListSensors(c *gin.Context) {
	ctx := context.Background()

	// 解析查询参数
	var regionID *int64
	if regionIDStr := c.Query("region_id"); regionIDStr != "" {
		if id, err := strconv.ParseInt(regionIDStr, 10, 64); err == nil {
			regionID = &id
		}
	}

	var monitorType *string
	if mt := c.Query("monitor_type"); mt != "" {
		monitorType = &mt
	}

	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	sensors, err := h.sensorController.ListSensors(ctx, regionID, monitorType, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sensors})
}

func (h *Handlers) GetSensor(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	sensor, err := h.sensorController.GetSensorByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sensor})
}

func (h *Handlers) CreateSensor(c *gin.Context) {
	var sensor model.Sensor
	if err := c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := context.Background()
	if err := h.sensorController.CreateSensor(ctx, sensor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "传感器创建成功"})
}

func (h *Handlers) UpdateSensor(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var sensor model.Sensor
	if err := c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sensor.ID = id
	ctx := context.Background()
	if err := h.sensorController.UpdateSensor(ctx, sensor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "传感器更新成功"})
}

func (h *Handlers) UpdateSensorStatus(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := context.Background()
	if err := h.sensorController.UpdateSensorStatus(ctx, id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "传感器状态更新成功"})
}

func (h *Handlers) DeleteSensor(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	if err := h.sensorController.DeleteSensor(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "传感器删除成功"})
}

// 传感器读数
func (h *Handlers) GetSensorReadings(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	ctx := context.Background()
	readings, err := h.sensorController.GetSensorReadingsBySensorID(ctx, id, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": readings})
}

func (h *Handlers) CreateSensorReading(c *gin.Context) {
	var reading model.SensorReading
	if err := c.ShouldBindJSON(&reading); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sensorID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	reading.SensorID = sensorID
	ctx := context.Background()
	if err := h.sensorController.InsertSensorReading(ctx, reading); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "传感器读数创建成功"})
}

func (h *Handlers) BatchCreateSensorReadings(c *gin.Context) {
	var readings []model.SensorReading
	if err := c.ShouldBindJSON(&readings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := context.Background()
	if err := h.sensorController.BatchInsertSensorReadings(ctx, readings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "批量创建传感器读数成功"})
}

// 预警规则
func (h *Handlers) ListAlertRules(c *gin.Context) {
	ctx := context.Background()

	// 解析查询参数
	var alertType *string
	if at := c.Query("alert_type"); at != "" {
		alertType = &at
	}

	var isActive *bool
	if activeStr := c.Query("is_active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			isActive = &active
		}
	}

	var severityLevel *string
	if sl := c.Query("severity_level"); sl != "" {
		severityLevel = &sl
	}

	rules, err := h.alertController.ListAlertRules(ctx, alertType, isActive, severityLevel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rules})
}

func (h *Handlers) GetAlertRule(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	rule, err := h.alertController.GetAlertRuleByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rule})
}

func (h *Handlers) CreateAlertRule(c *gin.Context) {
	var rule model.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rule.CreatedAt = time.Now()
	ctx := context.Background()
	id, err := h.alertController.CreateAlertRule(ctx, rule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "预警规则创建成功"})
}

func (h *Handlers) UpdateAlertRule(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var rule model.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rule.ID = id
	ctx := context.Background()
	if err := h.alertController.UpdateAlertRule(ctx, rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "预警规则更新成功"})
}

func (h *Handlers) DeleteAlertRule(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	if err := h.alertController.DeleteAlertRule(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "预警规则删除成功"})
}

// 预警记录
func (h *Handlers) ListAlerts(c *gin.Context) {
	ctx := context.Background()

	// 解析查询参数
	var regionID *int64
	if regionIDStr := c.Query("region_id"); regionIDStr != "" {
		if id, err := strconv.ParseInt(regionIDStr, 10, 64); err == nil {
			regionID = &id
		}
	}

	var status *string
	if st := c.Query("status"); st != "" {
		status = &st
	}

	var validity *string
	if v := c.Query("validity"); v != "" {
		validity = &v
	}

	alerts, err := h.alertController.ListAlerts(ctx, regionID, status, validity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": alerts})
}

func (h *Handlers) GetAlert(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	alert, err := h.alertController.GetAlertByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": alert})
}

func (h *Handlers) CreateAlert(c *gin.Context) {
	var alert model.Alert
	if err := c.ShouldBindJSON(&alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	alert.TriggeredAt = time.Now()
	ctx := context.Background()
	id, err := h.alertController.CreateAlert(ctx, alert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "预警记录创建成功"})
}

func (h *Handlers) UpdateAlert(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var alert model.Alert
	if err := c.ShouldBindJSON(&alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	alert.ID = id
	ctx := context.Background()
	if err := h.alertController.UpdateAlert(ctx, alert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "预警记录更新成功"})
}

func (h *Handlers) DeleteAlert(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	if err := h.alertController.DeleteAlert(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "预警记录删除成功"})
}

// 通知
func (h *Handlers) GetNotifications(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	notifications, err := h.alertController.GetNotificationsByAlertID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notifications})
}

func (h *Handlers) CreateNotification(c *gin.Context) {
	var notification model.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	notification.SentAt = time.Now()
	ctx := context.Background()
	if err := h.alertController.CreateNotification(ctx, notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "通知创建成功"})
}

// 资源管理
func (h *Handlers) ListResources(c *gin.Context) {
	ctx := context.Background()

	// 解析查询参数
	var regionID *int64
	if regionIDStr := c.Query("region_id"); regionIDStr != "" {
		if id, err := strconv.ParseInt(regionIDStr, 10, 64); err == nil {
			regionID = &id
		}
	}

	var resourceType *string
	if rt := c.Query("resource_type"); rt != "" {
		resourceType = &rt
	}

	var growthStage *string
	if gs := c.Query("growth_stage"); gs != "" {
		growthStage = &gs
	}

	var speciesName *string
	if sn := c.Query("species_name"); sn != "" {
		speciesName = &sn
	}

	resources, err := h.resourceController.ListResources(ctx, regionID, resourceType, growthStage, speciesName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resources})
}

func (h *Handlers) GetResource(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	resource, err := h.resourceController.GetResourceByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resource})
}

func (h *Handlers) CreateResource(c *gin.Context) {
	var resource model.Resource
	if err := c.ShouldBindJSON(&resource); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	ctx := context.Background()
	id, err := h.resourceController.CreateResource(ctx, resource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建资源变动记录
	changeReason := "新增资源记录"

	var changeAmount *int
	if resource.Quantity != nil {
		amount := int(*resource.Quantity)
		changeAmount = &amount
	}

	operatorID := userID.(int64)
	change := model.ResourceChange{
		ResourceID:   id,
		ChangeType:   "ADD",
		ChangeReason: &changeReason,
		ChangeAmount: changeAmount,
		ChangeArea:   resource.Area,
		ChangedAt:    time.Now(),
		OperatorID:   &operatorID,
	}

	if err := h.resourceController.CreateResourceChange(ctx, change); err != nil {
		// 记录错误但不影响资源创建成功
		log.Printf("创建资源变动记录失败: %v", err)
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "资源创建成功"})
}

func (h *Handlers) UpdateResource(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var resource model.Resource
	if err := c.ShouldBindJSON(&resource); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resource.ID = id
	ctx := context.Background()
	if err := h.resourceController.UpdateResource(ctx, resource); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "资源更新成功"})
}

// 更新资源生长状态（根据任务书要求）
func (h *Handlers) UpdateResourceGrowthStage(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		GrowthStage string `json:"growth_stage" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取现有资源
	ctx := context.Background()
	resource, err := h.resourceController.GetResourceByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
		return
	}

	// 更新生长状态
	resource.GrowthStage = req.GrowthStage
	if err := h.resourceController.UpdateResource(ctx, *resource); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "资源生长状态更新成功"})
}

func (h *Handlers) DeleteResource(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	if err := h.resourceController.DeleteResource(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "资源删除成功"})
}

// 资源变更
func (h *Handlers) GetResourceChanges(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	changes, err := h.resourceController.GetResourceChangesByResourceID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": changes})
}

// 资源统计查询
func (h *Handlers) GetResourceSummary(c *gin.Context) {
	ctx := context.Background()
	summary, err := h.resourceController.ResourceSummaryWithChanges(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": summary})
}

func (h *Handlers) CreateResourceChange(c *gin.Context) {
	var change model.ResourceChange
	if err := c.ShouldBindJSON(&change); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	change.ChangedAt = time.Now()
	ctx := context.Background()
	if err := h.resourceController.CreateResourceChange(ctx, change); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "资源变更记录创建成功"})
}

// 设备管理
func (h *Handlers) ListDevices(c *gin.Context) {
	ctx := context.Background()

	// 解析查询参数
	var deviceType *string
	if dt := c.Query("type"); dt != "" {
		deviceType = &dt
	}

	var regionID *int64
	if regionIDStr := c.Query("install_region_id"); regionIDStr != "" {
		if id, err := strconv.ParseInt(regionIDStr, 10, 64); err == nil {
			regionID = &id
		}
	}

	devices, err := h.deviceController.ListDevices(ctx, deviceType, regionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": devices})
}

func (h *Handlers) GetDevice(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	device, err := h.deviceController.GetDeviceByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": device})
}

func (h *Handlers) CreateDevice(c *gin.Context) {
	var device model.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	device.CreatedAt = time.Now()
	ctx := context.Background()
	if err := h.deviceController.CreateDevice(ctx, device); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "设备创建成功"})
}

func (h *Handlers) UpdateDevice(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var device model.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	device.ID = id
	ctx := context.Background()
	if err := h.deviceController.UpdateDevice(ctx, device); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "设备更新成功"})
}

func (h *Handlers) DeleteDevice(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	if err := h.deviceController.DeleteDevice(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "设备删除成功"})
}

// 设备状态
func (h *Handlers) GetDeviceStatus(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	status, err := h.deviceController.GetLatestDeviceStatusLog(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": status})
}

// GetDeviceStatusLogs 获取设备的状态记录列表
func (h *Handlers) GetDeviceStatusLogs(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	statusLogs, err := h.deviceController.GetDeviceStatusLogsByDeviceID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": statusLogs})
}

func (h *Handlers) CreateDeviceStatus(c *gin.Context) {
	var status model.DeviceStatusLog
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deviceID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	status.DeviceID = deviceID
	status.CollectedAt = time.Now()
	ctx := context.Background()
	if err := h.deviceController.CreateDeviceStatusLog(ctx, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "设备状态创建成功"})
}

// 维护记录
func (h *Handlers) GetMaintenanceRecords(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	records, err := h.deviceController.GetMaintenanceRecordsByDeviceID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": records})
}

func (h *Handlers) CreateMaintenanceRecord(c *gin.Context) {
	deviceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的设备ID"})
		return
	}

	var record model.MaintenanceRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从URL参数设置设备ID
	record.DeviceID = deviceID
	record.MaintenanceTime = time.Now()

	ctx := context.Background()
	if err := h.deviceController.CreateMaintenanceRecord(ctx, record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "维护记录创建成功"})
}

func (h *Handlers) DeleteMaintenanceRecord(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	if err := h.deviceController.DeleteMaintenanceRecord(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "维护记录删除成功"})
}

// 报表模板
func (h *Handlers) ListReportTemplates(c *gin.Context) {
	ctx := context.Background()

	// 解析查询参数
	var dimension *string
	if dim := c.Query("dimension"); dim != "" {
		dimension = &dim
	}

	var isActive *bool
	if activeStr := c.Query("is_active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			isActive = &active
		}
	}

	templates, err := h.reportController.ListReportTemplates(ctx, dimension, isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": templates})
}

func (h *Handlers) GetReportTemplate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	template, err := h.reportController.GetReportTemplateByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": template})
}

func (h *Handlers) CreateReportTemplate(c *gin.Context) {
	var template model.ReportTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	template.CreatedAt = time.Now()
	ctx := context.Background()
	id, err := h.reportController.CreateReportTemplate(ctx, template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "报表模板创建成功"})
}

func (h *Handlers) UpdateReportTemplate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var template model.ReportTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	template.ID = id
	ctx := context.Background()
	if err := h.reportController.UpdateReportTemplate(ctx, template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "报表模板更新成功"})
}

func (h *Handlers) DeleteReportTemplate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	if err := h.reportController.DeleteReportTemplate(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "报表模板删除成功"})
}

// 报表
func (h *Handlers) ListReports(c *gin.Context) {
	ctx := context.Background()

	// 解析查询参数
	var templateID int64
	if templateIDStr := c.Query("template_id"); templateIDStr != "" {
		if id, err := strconv.ParseInt(templateIDStr, 10, 64); err == nil {
			templateID = id
		}
	}

	// 解析报表类型参数
	var reportType *string
	if rtype := c.Query("type"); rtype != "" {
		reportType = &rtype
	}

	// 解析状态参数
	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	reports, err := h.reportController.GetReportsByTemplateID(ctx, templateID, reportType, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reports})
}

func (h *Handlers) GetReport(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	report, err := h.reportController.GetReportByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": report})
}

func (h *Handlers) CreateReport(c *gin.Context) {
	var report model.Report
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	report.GeneratedAt = time.Now()
	ctx := context.Background()
	id, err := h.reportController.CreateReport(ctx, report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "报表创建成功"})
}

func (h *Handlers) UpdateReport(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var report model.Report
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	report.ID = id
	ctx := context.Background()
	if err := h.reportController.UpdateReport(ctx, report); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "报表更新成功"})
}

func (h *Handlers) DeleteReport(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	if err := h.reportController.DeleteReport(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "报表删除成功"})
}

// 生成报表
func (h *Handlers) GenerateReport(c *gin.Context) {
	// 解析请求参数
	var req struct {
		TemplateID int64  `json:"template_id" binding:"required"`
		Period     string `json:"period" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	report, err := h.reportController.GenerateReport(ctx, req.TemplateID, req.Period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": report, "message": "报表生成成功"})
}

// 导出报表
func (h *Handlers) ExportReport(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	format := c.DefaultQuery("format", "excel")

	ctx := context.Background()
	filePath, err := h.reportController.ExportReport(ctx, id, format)
	if err != nil {
		// 记录详细错误日志
		log.Printf("ExportReport error: id=%d, format=%s, error=%v", id, format, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("ExportReport file not found: %s", filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件不存在"})
		return
	}

	// 设置响应头，实现文件下载
	log.Printf("ExportReport success: id=%d, format=%s, filePath=%s", id, format, filePath)
	c.FileAttachment(filePath, filepath.Base(filePath))
}

// 统计分析
func (h *Handlers) GetDeviceFaults(c *gin.Context) {
	ctx := context.Background()
	stats, err := h.statisticsController.DeviceFaultsWithMaintenance(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *Handlers) GetSensorValidity(c *gin.Context) {
	days := 7
	if daysStr := c.Query("days"); daysStr != "" {
		days, _ = strconv.Atoi(daysStr)
	}
	ctx := context.Background()
	stats, err := h.statisticsController.SensorValidityReport(ctx, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *Handlers) GetRegionAlertStats(c *gin.Context) {
	ctx := context.Background()
	stats, err := h.statisticsController.RegionAlertStatsLast30d(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *Handlers) GetFireAlerts(c *gin.Context) {
	regionName := c.Query("region")
	days := 7
	if daysStr := c.Query("days"); daysStr != "" {
		days, _ = strconv.Atoi(daysStr)
	}
	if regionName == "" {
		regionName = "Evergreen Forest"
	}
	ctx := context.Background()
	alerts, err := h.alertController.RecentFireAlerts(ctx, regionName, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": alerts})
}

// 环境趋势统计
func (h *Handlers) GetEnvironmentTrend(c *gin.Context) {
	// 解析查询参数
	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		days, _ = strconv.Atoi(daysStr)
	}

	var regionName *string
	if rn := c.Query("region_name"); rn != "" {
		regionName = &rn
	}

	var metricType *string
	if mt := c.Query("metric_type"); mt != "" {
		metricType = &mt
	}

	ctx := context.Background()
	data, err := h.statisticsController.GetEnvironmentTrend(ctx, days, regionName, metricType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// 空气质量统计
func (h *Handlers) GetAirQuality(c *gin.Context) {
	// 解析查询参数
	days := 7
	if daysStr := c.Query("days"); daysStr != "" {
		days, _ = strconv.Atoi(daysStr)
	}

	var regionName *string
	if rn := c.Query("region_name"); rn != "" {
		regionName = &rn
	}

	ctx := context.Background()
	data, err := h.statisticsController.GetAirQuality(ctx, days, regionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// 异常数据统计
func (h *Handlers) GetAnomalyData(c *gin.Context) {
	// 解析查询参数
	days := 7
	if daysStr := c.Query("days"); daysStr != "" {
		days, _ = strconv.Atoi(daysStr)
	}

	var regionName *string
	if rn := c.Query("region_name"); rn != "" {
		regionName = &rn
	}

	var anomalyType *string
	if at := c.Query("anomaly_type"); at != "" {
		anomalyType = &at
	}

	ctx := context.Background()
	data, err := h.statisticsController.GetAnomalyData(ctx, days, regionName, anomalyType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// 传感器数据汇总
func (h *Handlers) GetSensorDataSummary(c *gin.Context) {
	// 解析查询参数
	var regionName *string
	if rn := c.Query("region_name"); rn != "" {
		regionName = &rn
	}

	var monitorType *string
	if mt := c.Query("monitor_type"); mt != "" {
		monitorType = &mt
	}

	ctx := context.Background()
	data, err := h.statisticsController.GetSensorDataSummary(ctx, regionName, monitorType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// 区域日均值统计
func (h *Handlers) GetRegionDailyAverages(c *gin.Context) {
	// 解析查询参数
	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		days, _ = strconv.Atoi(daysStr)
	}

	var regionName *string
	if rn := c.Query("region_name"); rn != "" {
		regionName = &rn
	}

	var metricType *string
	if mt := c.Query("metric_type"); mt != "" {
		metricType = &mt
	}

	ctx := context.Background()
	data, err := h.statisticsController.GetRegionDailyAverages(ctx, days, regionName, metricType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// 认证相关处理方法

// Register 用户注册（支持角色选择）
func (h *Handlers) Register(c *gin.Context) {
	var req model.UserRegistration

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	userID, err := h.authController.RegisterUser(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "已存在") || strings.Contains(err.Error(), "不存在") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": userID,
		"message": "用户注册成功",
		"role_selected": true,
	})
}

// GetAvailableRoles 获取可用的角色列表
func (h *Handlers) GetAvailableRoles(c *gin.Context) {
	ctx := context.Background()
	roles, err := h.authController.GetAvailableRoles(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取角色列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
		"count": len(roles),
	})
}

// Login 用户登录
func (h *Handlers) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	user, err := h.authController.VerifyPassword(ctx, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 生成JWT令牌
	token, err := util.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"user_id":   user.ID,
		"username":  user.Username,
		"real_name": user.RealName,
		"message":   "登录成功",
	})
}

// GetCurrentUser 获取当前用户信息
func (h *Handlers) GetCurrentUser(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	ctx := context.Background()
	user, err := h.authController.GetUserByID(ctx, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 隐藏敏感信息
	user.PasswordHash = ""
	user.LoginAttempts = 0

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdatePassword 更新密码
func (h *Handlers) UpdatePassword(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证旧密码
	ctx := context.Background()
	username, _ := c.Get("username")
	_, err := h.authController.VerifyPassword(ctx, username.(string), req.OldPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "旧密码错误"})
		return
	}

	// 更新密码
	if err := h.authController.UpdatePassword(ctx, userID.(int64), req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码更新成功"})
}

// 角色管理API

// CreateRole 创建角色
func (h *Handlers) CreateRole(c *gin.Context) {
	ctx := context.Background()

	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	id, err := h.authController.CreateRole(ctx, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "角色创建成功", "role_id": id})
}

// GetRole 获取角色详情
func (h *Handlers) GetRole(c *gin.Context) {
	ctx := context.Background()

	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	role, err := h.authController.GetRoleByID(ctx, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

// UpdateRole 更新角色
func (h *Handlers) UpdateRole(c *gin.Context) {
	ctx := context.Background()

	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}
	role.ID = roleID

	err = h.authController.UpdateRole(ctx, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "角色更新成功"})
}

// DeleteRole 删除角色
func (h *Handlers) DeleteRole(c *gin.Context) {
	ctx := context.Background()

	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	err = h.authController.DeleteRole(ctx, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "角色删除成功"})
}

// AssignRole 分配角色给用户
func (h *Handlers) AssignRole(c *gin.Context) {
	ctx := context.Background()

	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req struct {
		RoleID int64 `json:"role_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	err = h.authController.AssignRoleToUser(ctx, userID, req.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "角色分配成功"})
}

// RemoveRole 移除用户角色
func (h *Handlers) RemoveRole(c *gin.Context) {
	ctx := context.Background()

	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	roleID, err := strconv.ParseInt(c.Param("role_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	err = h.authController.RemoveRoleFromUser(ctx, userID, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "角色移除成功"})
}

// GetUserRoles 获取用户角色列表
func (h *Handlers) GetUserRoles(c *gin.Context) {
	ctx := context.Background()

	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	roles, err := h.authController.GetUserRoles(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}

// 权限管理API

// CreatePermission 创建权限
func (h *Handlers) CreatePermission(c *gin.Context) {
	ctx := context.Background()

	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	id, err := h.authController.CreatePermission(ctx, permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限创建成功", "permission_id": id})
}

// GetPermission 获取权限详情
func (h *Handlers) GetPermission(c *gin.Context) {
	ctx := context.Background()

	permissionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的权限ID"})
		return
	}

	permission, err := h.authController.GetPermissionByID(ctx, permissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permission})
}

// UpdatePermission 更新权限
func (h *Handlers) UpdatePermission(c *gin.Context) {
	ctx := context.Background()

	permissionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的权限ID"})
		return
	}

	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}
	permission.ID = permissionID

	err = h.authController.UpdatePermission(ctx, permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限更新成功"})
}

// DeletePermission 删除权限
func (h *Handlers) DeletePermission(c *gin.Context) {
	ctx := context.Background()

	permissionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的权限ID"})
		return
	}

	err = h.authController.DeletePermission(ctx, permissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限删除成功"})
}

// AssignPermission 分配权限给角色
func (h *Handlers) AssignPermission(c *gin.Context) {
	ctx := context.Background()

	roleID, err := strconv.ParseInt(c.Param("role_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var req struct {
		PermissionID int64 `json:"permission_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	err = h.authController.AssignPermissionToRole(ctx, roleID, req.PermissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限分配成功"})
}

// RemovePermission 移除角色权限
func (h *Handlers) RemovePermission(c *gin.Context) {
	ctx := context.Background()

	roleID, err := strconv.ParseInt(c.Param("role_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	permissionID, err := strconv.ParseInt(c.Param("permission_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的权限ID"})
		return
	}

	err = h.authController.RemovePermissionFromRole(ctx, roleID, permissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限移除成功"})
}

// GetRolePermissions 获取角色权限列表
func (h *Handlers) GetRolePermissions(c *gin.Context) {
	ctx := context.Background()

	roleID, err := strconv.ParseInt(c.Param("role_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	permissions, err := h.authController.GetRolePermissions(ctx, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permissions})
}

// GetUserPermissions 获取用户权限列表
func (h *Handlers) GetUserPermissions(c *gin.Context) {
	ctx := context.Background()

	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	permissions, err := h.authController.GetUserPermissions(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permissions})
}

// CheckUserPermission 检查用户权限
func (h *Handlers) CheckUserPermission(c *gin.Context) {
	ctx := context.Background()

	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	resource := c.Query("resource")
	action := c.Query("action")

	if resource == "" || action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少resource或action参数"})
		return
	}

	hasPermission, err := h.authController.CheckUserPermission(ctx, userID, resource, action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"has_permission": hasPermission})
}

// BackupDatabase 备份数据库
func (h *Handlers) BackupDatabase(c *gin.Context) {
	// 加载数据库配置
	dbConfig, err := db.LoadConfig("db/config.yaml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加载数据库配置失败: " + err.Error()})
		return
	}

	// 创建备份配置
	backupConfig := util.BackupConfig{
		BackupDir:      "./backups",
		BackupInterval: time.Hour * 24,
		RetentionDays:  7,
		DBConfig:       *dbConfig,
	}

	// 执行数据库备份
	result, err := util.BackupDatabase(h.db, backupConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "备份数据库失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "数据库备份成功",
		"result":  result,
	})
}

// RestoreDatabase 恢复数据库
func (h *Handlers) RestoreDatabase(c *gin.Context) {
	// 从请求中获取备份文件路径
	var req struct {
		BackupFile string `json:"backup_file" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 加载数据库配置
	dbConfig, err := db.LoadConfig("db/config.yaml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加载数据库配置失败: " + err.Error()})
		return
	}

	// 创建备份配置
	backupConfig := util.BackupConfig{
		BackupDir:      "./backups",
		BackupInterval: time.Hour * 24,
		RetentionDays:  7,
		DBConfig:       *dbConfig,
	}

	// 执行数据库恢复
	if err := util.RestoreDatabase(h.db, req.BackupFile, backupConfig); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "恢复数据库失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "数据库恢复成功",
		"backup_file": req.BackupFile,
	})
}

// ListBackups 列出备份文件
func (h *Handlers) ListBackups(c *gin.Context) {
	// 备份目录
	backupDir := "./backups"

	// 确保备份目录存在
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		// 如果目录不存在，返回空列表
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
		return
	}

	// 读取备份目录中的所有文件
	files, err := os.ReadDir(backupDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取备份目录失败: " + err.Error()})
		return
	}

	// 过滤出备份文件
	var backupFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			backupFiles = append(backupFiles, file.Name())
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  backupFiles,
		"total": len(backupFiles),
	})
}
