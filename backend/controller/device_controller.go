package controller

import (
	"context"
	"database/sql"
	"fmt"

	"smart-forest-system/backend/model"
)

// DeviceController 设备控制器
type DeviceController struct {
	db *sql.DB
}

// NewDeviceController 创建设备控制器
func NewDeviceController(db *sql.DB) *DeviceController {
	return &DeviceController{db: db}
}

// CreateDevice 插入新的设备记录。
func (c *DeviceController) CreateDevice(ctx context.Context, device model.Device) error {
	// 自动生成设备ID
	var deviceID int64
	err := c.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(device_id), 0) + 1 FROM devices").Scan(&deviceID)
	if err != nil {
		return fmt.Errorf("生成设备ID失败: %v", err)
	}

	const query = `
		INSERT INTO devices (device_id, device_name, device_type, model_spec, purchased_at, 
		install_region_id, installer_id, warranty_until, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err = c.db.ExecContext(ctx, query,
		deviceID,
		device.Name,
		device.Type,
		device.ModelSpec,
		device.PurchasedAt,
		device.InstallRegionID,
		device.InstallerID,
		device.WarrantyUntil,
	)
	return err
}

// GetDeviceByID 根据ID获取设备记录。
func (c *DeviceController) GetDeviceByID(ctx context.Context, id int64) (*model.Device, error) {
	const query = `
		SELECT device_id, device_name, device_type, model_spec, purchased_at, 
		install_region_id, installer_id, warranty_until, created_at
		FROM devices WHERE device_id = ?`
	var device model.Device
	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&device.ID,
		&device.Name,
		&device.Type,
		&device.ModelSpec,
		&device.PurchasedAt,
		&device.InstallRegionID,
		&device.InstallerID,
		&device.WarrantyUntil,
		&device.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// UpdateDevice 更新现有设备记录。
func (c *DeviceController) UpdateDevice(ctx context.Context, device model.Device) error {
	const query = `
		UPDATE devices
		SET device_name = ?, device_type = ?, model_spec = ?, purchased_at = ?,
		install_region_id = ?, installer_id = ?, warranty_until = ?
		WHERE device_id = ?`
	result, err := c.db.ExecContext(ctx, query,
		device.Name,
		device.Type,
		device.ModelSpec,
		device.PurchasedAt,
		device.InstallRegionID,
		device.InstallerID,
		device.WarrantyUntil,
		device.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no device updated for id=%d", device.ID)
	}
	return nil
}

// DeleteDevice 根据ID删除设备记录。
func (c *DeviceController) DeleteDevice(ctx context.Context, id int64) error {
	// 开启事务
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 删除该设备的所有维护记录
	if _, err := tx.ExecContext(ctx, `DELETE FROM maintenance_records WHERE device_id = ?`, id); err != nil {
		return fmt.Errorf("删除设备维护记录失败: %v", err)
	}

	// 2. 删除该设备的所有状态日志
	if _, err := tx.ExecContext(ctx, `DELETE FROM device_status_logs WHERE device_id = ?`, id); err != nil {
		return fmt.Errorf("删除设备状态日志失败: %v", err)
	}

	// 3. 删除设备
	result, err := tx.ExecContext(ctx, `DELETE FROM devices WHERE device_id = ?`, id)
	if err != nil {
		return fmt.Errorf("删除设备失败: %v", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("设备不存在或已被删除, id=%d", id)
	}

	// 提交事务
	return tx.Commit()
}

// ListDevices 根据条件获取设备列表。
func (c *DeviceController) ListDevices(ctx context.Context, deviceType *string, regionID *int64) ([]model.Device, error) {
	var query string
	var args []interface{}

	baseQuery := `
		SELECT device_id, device_name, device_type, model_spec, purchased_at, 
		install_region_id, installer_id, warranty_until, created_at
		FROM devices`

	conditions := ""
	if deviceType != nil {
		conditions += " AND device_type = ?"
		args = append(args, *deviceType)
	}
	if regionID != nil {
		conditions += " AND install_region_id = ?"
		args = append(args, *regionID)
	}

	if conditions != "" {
		query = baseQuery + " WHERE " + conditions[5:] + " ORDER BY device_id"
	} else {
		query = baseQuery + " ORDER BY device_id"
	}

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []model.Device
	for rows.Next() {
		var device model.Device
		if err := rows.Scan(
			&device.ID,
			&device.Name,
			&device.Type,
			&device.ModelSpec,
			&device.PurchasedAt,
			&device.InstallRegionID,
			&device.InstallerID,
			&device.WarrantyUntil,
			&device.CreatedAt,
		); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, rows.Err()
}

// CreateDeviceStatusLog 插入新的设备状态日志。
func (c *DeviceController) CreateDeviceStatusLog(ctx context.Context, log model.DeviceStatusLog) error {
	const query = `
		INSERT INTO device_status_logs (device_id, collected_at, run_status, battery_percent, signal_strength)
		VALUES (?, ?, ?, ?, ?)`
	_, err := c.db.ExecContext(ctx, query,
		log.DeviceID,
		log.CollectedAt,
		log.RunStatus,
		log.BatteryPercent,
		log.SignalStrength,
	)
	return err
}

// GetLatestDeviceStatusLog 获取设备的最新状态日志。
func (c *DeviceController) GetLatestDeviceStatusLog(ctx context.Context, deviceID int64) (*model.DeviceStatusLog, error) {
	const query = `
		SELECT log_id, device_id, collected_at, run_status, battery_percent, signal_strength
		FROM device_status_logs
		WHERE device_id = ?
		ORDER BY collected_at DESC
		LIMIT 1`
	var log model.DeviceStatusLog
	err := c.db.QueryRowContext(ctx, query, deviceID).Scan(
		&log.ID,
		&log.DeviceID,
		&log.CollectedAt,
		&log.RunStatus,
		&log.BatteryPercent,
		&log.SignalStrength,
	)
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetDeviceStatusLogsByTimeRange 获取指定时间范围内设备的状态日志。
func (c *DeviceController) GetDeviceStatusLogsByDeviceID(ctx context.Context, deviceID int64) ([]model.DeviceStatusLog, error) {
	const query = `
		SELECT log_id, device_id, collected_at, run_status, battery_percent, signal_strength
		FROM device_status_logs
		WHERE device_id = ?
		ORDER BY collected_at DESC`
	rows, err := c.db.QueryContext(ctx, query, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.DeviceStatusLog
	for rows.Next() {
		var log model.DeviceStatusLog
		if err := rows.Scan(
			&log.ID,
			&log.DeviceID,
			&log.CollectedAt,
			&log.RunStatus,
			&log.BatteryPercent,
			&log.SignalStrength,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}

func (c *DeviceController) GetDeviceStatusLogsByTimeRange(ctx context.Context, deviceID int64, startTime, endTime interface{}) ([]model.DeviceStatusLog, error) {
	const query = `
		SELECT log_id, device_id, collected_at, run_status, battery_percent, signal_strength
		FROM device_status_logs
		WHERE device_id = ? AND collected_at BETWEEN ? AND ?
		ORDER BY collected_at ASC`
	rows, err := c.db.QueryContext(ctx, query, deviceID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.DeviceStatusLog
	for rows.Next() {
		var log model.DeviceStatusLog
		if err := rows.Scan(
			&log.ID,
			&log.DeviceID,
			&log.CollectedAt,
			&log.RunStatus,
			&log.BatteryPercent,
			&log.SignalStrength,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}

// CreateMaintenanceRecord 插入新的维护记录。
func (c *DeviceController) CreateMaintenanceRecord(ctx context.Context, record model.MaintenanceRecord) error {
	const query = `
		INSERT INTO maintenance_records (device_id, maintenance_type, maintenance_time, maintainer_id, content, result)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := c.db.ExecContext(ctx, query,
		record.DeviceID,
		record.MaintenanceType,
		record.MaintenanceTime,
		record.MaintainerID,
		record.Content,
		record.Result,
	)
	return err
}

// GetMaintenanceRecordByID 根据ID获取维护记录。
func (c *DeviceController) GetMaintenanceRecordByID(ctx context.Context, id int64) (*model.MaintenanceRecord, error) {
	const query = `
		SELECT maintenance_id, device_id, maintenance_type, maintenance_time, maintainer_id, content, result
		FROM maintenance_records WHERE maintenance_id = ?`
	var record model.MaintenanceRecord
	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&record.ID,
		&record.DeviceID,
		&record.MaintenanceType,
		&record.MaintenanceTime,
		&record.MaintainerID,
		&record.Content,
		&record.Result,
	)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// UpdateMaintenanceRecord 更新现有维护记录。
func (c *DeviceController) UpdateMaintenanceRecord(ctx context.Context, record model.MaintenanceRecord) error {
	const query = `
		UPDATE maintenance_records
		SET device_id = ?, maintenance_type = ?, maintenance_time = ?, maintainer_id = ?, content = ?, result = ?
		WHERE maintenance_id = ?`
	result, err := c.db.ExecContext(ctx, query,
		record.DeviceID,
		record.MaintenanceType,
		record.MaintenanceTime,
		record.MaintainerID,
		record.Content,
		record.Result,
		record.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no maintenance record updated for id=%d", record.ID)
	}
	return nil
}

// DeleteMaintenanceRecord 根据ID删除维护记录。
func (c *DeviceController) DeleteMaintenanceRecord(ctx context.Context, id int64) error {
	const query = `DELETE FROM maintenance_records WHERE maintenance_id = ?`
	result, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no maintenance record deleted for id=%d", id)
	}
	return nil
}

// GetMaintenanceRecordsByDeviceID 获取指定设备的维护记录。
func (c *DeviceController) GetMaintenanceRecordsByDeviceID(ctx context.Context, deviceID int64) ([]model.MaintenanceRecord, error) {
	const query = `
		SELECT mr.maintenance_id, mr.device_id, mr.maintenance_type, mr.maintenance_time, mr.maintainer_id, mr.content, mr.result, u.real_name as operator_name
		FROM maintenance_records mr
		LEFT JOIN users u ON mr.maintainer_id = u.user_id
		WHERE mr.device_id = ?
		ORDER BY mr.maintenance_time DESC`
	rows, err := c.db.QueryContext(ctx, query, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.MaintenanceRecord
	for rows.Next() {
		var record model.MaintenanceRecord
		var operatorName sql.NullString
		if err := rows.Scan(
			&record.ID,
			&record.DeviceID,
			&record.MaintenanceType,
			&record.MaintenanceTime,
			&record.MaintainerID,
			&record.Content,
			&record.Result,
			&operatorName,
		); err != nil {
			return nil, err
		}

		if operatorName.Valid {
			record.OperatorName = &operatorName.String
		}

		records = append(records, record)
	}
	return records, rows.Err()
}
