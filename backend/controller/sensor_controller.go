package controller

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"smart-forest-system/backend/model"
)

// SensorController 传感器控制器
type SensorController struct {
	db *sql.DB
}

// NewSensorController 创建传感器控制器
func NewSensorController(db *sql.DB) *SensorController {
	return &SensorController{db: db}
}

// CreateSensor 插入一条传感器记录。
func (c *SensorController) CreateSensor(ctx context.Context, sensor model.Sensor) error {
	const query = `
		INSERT INTO sensors (region_id, model, monitor_type, install_time, protocol, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())`
	_, err := c.db.ExecContext(ctx, query,
		sensor.RegionID,
		sensor.Model,
		sensor.MonitorType,
		sensor.InstallTime,
		sensor.Protocol,
		sensor.Status,
	)
	return err
}

// GetSensorByID 根据ID获取传感器信息。
func (c *SensorController) GetSensorByID(ctx context.Context, id int64) (*model.Sensor, error) {
	const query = `
		SELECT sensor_id, region_id, model, monitor_type, install_time, protocol, status, created_at
		FROM sensors WHERE sensor_id = ?`
	var sensor model.Sensor
	var createdAt time.Time
	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&sensor.ID,
		&sensor.RegionID,
		&sensor.Model,
		&sensor.MonitorType,
		&sensor.InstallTime,
		&sensor.Protocol,
		&sensor.Status,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}
	return &sensor, nil
}

// UpdateSensor 更新现有传感器信息。
func (c *SensorController) UpdateSensor(ctx context.Context, sensor model.Sensor) error {
	const query = `
		UPDATE sensors
		SET region_id = ?, model = ?, monitor_type = ?, install_time = ?, protocol = ?, status = ?
		WHERE sensor_id = ?`
	result, err := c.db.ExecContext(ctx, query,
		sensor.RegionID,
		sensor.Model,
		sensor.MonitorType,
		sensor.InstallTime,
		sensor.Protocol,
		sensor.Status,
		sensor.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no sensor updated for id=%d", sensor.ID)
	}
	return nil
}

// UpdateSensorStatus 更新传感器状态。
func (c *SensorController) UpdateSensorStatus(ctx context.Context, sensorID int64, status string) error {
	const query = `UPDATE sensors SET status = ? WHERE sensor_id = ?`
	res, err := c.db.ExecContext(ctx, query, status, sensorID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no sensor updated for id=%d", sensorID)
	}
	return nil
}

// DeleteSensor 删除传感器（需验证无关联监测数据）。
func (c *SensorController) DeleteSensor(ctx context.Context, id int64) error {
	// 开启事务
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 删除该传感器的所有读数记录
	if _, err := tx.ExecContext(ctx, `DELETE FROM sensor_readings WHERE sensor_id = ?`, id); err != nil {
		return fmt.Errorf("删除传感器读数失败: %v", err)
	}

	// 2. 删除传感器
	result, err := tx.ExecContext(ctx, `DELETE FROM sensors WHERE sensor_id = ?`, id)
	if err != nil {
		return fmt.Errorf("删除传感器失败: %v", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("传感器不存在或已被删除, id=%d", id)
	}

	// 提交事务
	return tx.Commit()
}

// ListSensors 根据条件获取传感器列表。
func (c *SensorController) ListSensors(ctx context.Context, regionID *int64, monitorType *string, status *string) ([]model.Sensor, error) {
	var query string
	var args []interface{}

	baseQuery := `
		SELECT sensor_id, region_id, model, monitor_type, install_time, protocol, status, created_at
		FROM sensors`

	conditions := ""
	if regionID != nil {
		conditions += " AND region_id = ?"
		args = append(args, *regionID)
	}
	if monitorType != nil {
		conditions += " AND monitor_type = ?"
		args = append(args, *monitorType)
	}
	if status != nil {
		conditions += " AND status = ?"
		args = append(args, *status)
	}

	if conditions != "" {
		query = baseQuery + " WHERE " + conditions[5:] + " ORDER BY sensor_id"
	} else {
		query = baseQuery + " ORDER BY sensor_id"
	}

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []model.Sensor
	for rows.Next() {
		var sensor model.Sensor
		var createdAt time.Time
		if err := rows.Scan(
			&sensor.ID,
			&sensor.RegionID,
			&sensor.Model,
			&sensor.MonitorType,
			&sensor.InstallTime,
			&sensor.Protocol,
			&sensor.Status,
			&createdAt,
		); err != nil {
			return nil, err
		}
		sensors = append(sensors, sensor)
	}
	return sensors, rows.Err()
}

// GetSensorsByRegionName 根据区域名称获取该区域的所有传感器。
func (c *SensorController) GetSensorsByRegionName(ctx context.Context, regionName string) ([]model.Sensor, error) {
	const query = `
		SELECT s.sensor_id, s.region_id, s.model, s.monitor_type, s.install_time, s.protocol, s.status, s.created_at
		FROM sensors s
		JOIN regions r ON s.region_id = r.region_id
		WHERE r.region_name = ?
		ORDER BY s.sensor_id`

	rows, err := c.db.QueryContext(ctx, query, regionName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []model.Sensor
	for rows.Next() {
		var sensor model.Sensor
		if err := rows.Scan(
			&sensor.ID,
			&sensor.RegionID,
			&sensor.Model,
			&sensor.MonitorType,
			&sensor.InstallTime,
			&sensor.Protocol,
			&sensor.Status,
		); err != nil {
			return nil, err
		}
		sensors = append(sensors, sensor)
	}
	return sensors, rows.Err()
}

// InsertSensorReading 添加传感器读数记录。
func (c *SensorController) InsertSensorReading(ctx context.Context, reading model.SensorReading) error {
	const query = `
		INSERT INTO sensor_readings (sensor_id, collected_at, reading_type, numeric_value, value_unit, media_path, data_status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err := c.db.ExecContext(ctx, query,
		reading.SensorID,
		reading.CollectedAt,
		reading.ReadingType,
		reading.Numeric,
		reading.ValueUnit,
		reading.MediaPath,
		reading.DataStatus,
	)
	return err
}

// BatchInsertSensorReadings 批量插入传感器读数（单事务处理）。
func (c *SensorController) BatchInsertSensorReadings(ctx context.Context, readings []model.SensorReading) error {
	if len(readings) == 0 {
		return nil
	}

	// 开始事务
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 准备插入语句
	const query = `
		INSERT INTO sensor_readings (sensor_id, collected_at, reading_type, numeric_value, value_unit, media_path, data_status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("准备插入语句失败: %v", err)
	}
	defer stmt.Close()

	// 批量插入
	for _, reading := range readings {
		_, err = stmt.ExecContext(ctx,
			reading.SensorID,
			reading.CollectedAt,
			reading.ReadingType,
			reading.Numeric,
			reading.ValueUnit,
			reading.MediaPath,
			reading.DataStatus,
		)
		if err != nil {
			return fmt.Errorf("插入监测数据失败: %v", err)
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// GetSensorReadingsBySensorID 获取指定传感器的读数记录。
func (c *SensorController) GetSensorReadingsBySensorID(ctx context.Context, sensorID int64, limit int) ([]model.SensorReading, error) {
	const query = `
		SELECT reading_id, sensor_id, collected_at, reading_type, numeric_value, value_unit, media_path, data_status, created_at
		FROM sensor_readings
		WHERE sensor_id = ?
		ORDER BY collected_at DESC
		LIMIT ?`
	rows, err := c.db.QueryContext(ctx, query, sensorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []model.SensorReading
	for rows.Next() {
		var reading model.SensorReading
		var createdAt time.Time
		if err := rows.Scan(
			&reading.ID,
			&reading.SensorID,
			&reading.CollectedAt,
			&reading.ReadingType,
			&reading.Numeric,
			&reading.ValueUnit,
			&reading.MediaPath,
			&reading.DataStatus,
			&createdAt,
		); err != nil {
			return nil, err
		}
		readings = append(readings, reading)
	}
	return readings, rows.Err()
}

// GetSensorReadingsByTimeRange 获取指定时间范围内传感器的读数记录。
func (c *SensorController) GetSensorReadingsByTimeRange(ctx context.Context, sensorID int64, startTime, endTime time.Time) ([]model.SensorReading, error) {
	const query = `
		SELECT reading_id, sensor_id, collected_at, reading_type, numeric_value, value_unit, media_path, data_status, created_at
		FROM sensor_readings
		WHERE sensor_id = ? AND collected_at BETWEEN ? AND ?
		ORDER BY collected_at ASC`
	rows, err := c.db.QueryContext(ctx, query, sensorID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []model.SensorReading
	for rows.Next() {
		var reading model.SensorReading
		var createdAt time.Time
		if err := rows.Scan(
			&reading.ID,
			&reading.SensorID,
			&reading.CollectedAt,
			&reading.ReadingType,
			&reading.Numeric,
			&reading.ValueUnit,
			&reading.MediaPath,
			&reading.DataStatus,
			&createdAt,
		); err != nil {
			return nil, err
		}
		readings = append(readings, reading)
	}
	return readings, rows.Err()
}
