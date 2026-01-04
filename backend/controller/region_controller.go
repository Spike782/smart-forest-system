package controller

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"smart-forest-system/backend/model"
)

// RegionController 区域控制器
type RegionController struct {
	db *sql.DB
}

// NewRegionController 创建区域控制器
func NewRegionController(db *sql.DB) *RegionController {
	return &RegionController{db: db}
}

// CreateRegion 插入新的区域记录。
func (c *RegionController) CreateRegion(ctx context.Context, region model.Region) (int64, error) {
	const query = `
		INSERT INTO regions (region_name, region_type, latitude, longitude, manager_id)
		VALUES (?, ?, ?, ?, ?)`
	result, err := c.db.ExecContext(ctx, query,
		region.Name,
		region.Type,
		region.Latitude,
		region.Longitude,
		region.ManagerID,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// GetRegionByID 根据ID获取区域记录。
func (c *RegionController) GetRegionByID(ctx context.Context, id int64) (*model.Region, error) {
	const query = `
		SELECT region_id, region_name, region_type, latitude, longitude, manager_id, created_at
		FROM regions WHERE region_id = ?`
	var region model.Region
	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&region.ID,
		&region.Name,
		&region.Type,
		&region.Latitude,
		&region.Longitude,
		&region.ManagerID,
		&region.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &region, nil
}

// UpdateRegion 更新现有区域记录。
func (c *RegionController) UpdateRegion(ctx context.Context, region model.Region) error {
	const query = `
		UPDATE regions
		SET region_name = ?, region_type = ?, latitude = ?, longitude = ?, manager_id = ?
		WHERE region_id = ?`
	result, err := c.db.ExecContext(ctx, query,
		region.Name,
		region.Type,
		region.Latitude,
		region.Longitude,
		region.ManagerID,
		region.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no region updated for id=%d", region.ID)
	}
	return nil
}

// DeleteRegion 根据ID删除区域记录。
func (c *RegionController) DeleteRegion(ctx context.Context, id int64) error {
	// 开启事务
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 删除该区域下的所有传感器读数
	if _, err := tx.ExecContext(ctx, `DELETE FROM sensor_readings WHERE sensor_id IN (SELECT sensor_id FROM sensors WHERE region_id = ?)`, id); err != nil {
		return fmt.Errorf("删除传感器读数失败: %v", err)
	}

	// 2. 删除该区域下的所有传感器
	if _, err := tx.ExecContext(ctx, `DELETE FROM sensors WHERE region_id = ?`, id); err != nil {
		return fmt.Errorf("删除传感器失败: %v", err)
	}

	// 3. 删除该区域下的所有预警通知
	if _, err := tx.ExecContext(ctx, `DELETE FROM notifications WHERE alert_id IN (SELECT alert_id FROM alerts WHERE region_id = ?)`, id); err != nil {
		return fmt.Errorf("删除预警通知失败: %v", err)
	}

	// 4. 删除该区域下的所有预警
	if _, err := tx.ExecContext(ctx, `DELETE FROM alerts WHERE region_id = ?`, id); err != nil {
		return fmt.Errorf("删除预警失败: %v", err)
	}

	// 5. 删除该区域下的所有资源变更
	if _, err := tx.ExecContext(ctx, `DELETE FROM resource_changes WHERE resource_id IN (SELECT resource_id FROM resources WHERE region_id = ?)`, id); err != nil {
		return fmt.Errorf("删除资源变更失败: %v", err)
	}

	// 6. 删除该区域下的所有资源
	if _, err := tx.ExecContext(ctx, `DELETE FROM resources WHERE region_id = ?`, id); err != nil {
		return fmt.Errorf("删除资源失败: %v", err)
	}

	// 7. 删除该区域下的所有设备维护记录
	if _, err := tx.ExecContext(ctx, `DELETE FROM maintenance_records WHERE device_id IN (SELECT device_id FROM devices WHERE install_region_id = ?)`, id); err != nil {
		return fmt.Errorf("删除设备维护记录失败: %v", err)
	}

	// 8. 删除该区域下的所有设备状态日志
	if _, err := tx.ExecContext(ctx, `DELETE FROM device_status_logs WHERE device_id IN (SELECT device_id FROM devices WHERE install_region_id = ?)`, id); err != nil {
		return fmt.Errorf("删除设备状态日志失败: %v", err)
	}

	// 9. 删除该区域下的所有设备
	if _, err := tx.ExecContext(ctx, `DELETE FROM devices WHERE install_region_id = ?`, id); err != nil {
		return fmt.Errorf("删除设备失败: %v", err)
	}

	// 10. 最后删除区域
	result, err := tx.ExecContext(ctx, `DELETE FROM regions WHERE region_id = ?`, id)
	if err != nil {
		return fmt.Errorf("删除区域失败: %v", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("区域不存在或已被删除, id=%d", id)
	}

	// 提交事务
	return tx.Commit()
}

// ListRegions 根据条件获取区域列表。
func (c *RegionController) ListRegions(ctx context.Context, regionType *string, name *string) ([]model.Region, error) {
	// 首先检查regions表是否存在
	tableCheckQuery := `SHOW TABLES LIKE 'regions'`
	var tableName string
	err := c.db.QueryRowContext(ctx, tableCheckQuery).Scan(&tableName)
	if err != nil {
		// 表不存在，尝试创建表
		createTableQuery := `
			CREATE TABLE IF NOT EXISTS regions (
				region_id BIGINT AUTO_INCREMENT PRIMARY KEY,
				region_name VARCHAR(100) NOT NULL,
				region_type ENUM('FOREST','GRASSLAND') NOT NULL,
				latitude DECIMAL(9,6),
				longitude DECIMAL(9,6),
				manager_id BIGINT NULL,
				created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				UNIQUE KEY uk_region_name (region_name)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`
		_, err = c.db.ExecContext(ctx, createTableQuery)
		if err != nil {
			return nil, err
		}
	}

	var conditions []string
	var args []interface{}

	if regionType != nil {
		conditions = append(conditions, "region_type = ?")
		args = append(args, *regionType)
	}

	if name != nil {
		conditions = append(conditions, "region_name LIKE ?")
		args = append(args, "%"+*name+"%")
	}

	query := `
		SELECT region_id AS id, region_name AS name, region_type AS type, latitude, longitude, manager_id, created_at
		FROM regions
	`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY region_name"

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regions []model.Region
	for rows.Next() {
		var region model.Region
		if err := rows.Scan(
			&region.ID,
			&region.Name,
			&region.Type,
			&region.Latitude,
			&region.Longitude,
			&region.ManagerID,
			&region.CreatedAt,
		); err != nil {
			return nil, err
		}
		regions = append(regions, region)
	}
	return regions, rows.Err()
}
