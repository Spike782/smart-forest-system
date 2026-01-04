package controller

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"smart-forest-system/backend/model"
)

// ResourceController 资源控制器
type ResourceController struct {
	db *sql.DB
}

// NewResourceController 创建资源控制器
func NewResourceController(db *sql.DB) *ResourceController {
	return &ResourceController{db: db}
}

// CreateResource 插入新的资源记录。
func (c *ResourceController) CreateResource(ctx context.Context, resource model.Resource) (int64, error) {
	const query = `
		INSERT INTO resources (resource_type, region_id, species_name, quantity, area, growth_stage, planted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := c.db.ExecContext(ctx, query,
		resource.ResourceType,
		resource.RegionID,
		resource.SpeciesName,
		resource.Quantity,
		resource.Area,
		resource.GrowthStage,
		resource.PlantedAt,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// GetResourceByID 根据ID获取资源记录。
func (c *ResourceController) GetResourceByID(ctx context.Context, id int64) (*model.Resource, error) {
	const query = `
		SELECT resource_id, resource_type, region_id, species_name, quantity, area, growth_stage, planted_at, updated_at
		FROM resources WHERE resource_id = ?`
	var resource model.Resource
	var updatedAt time.Time
	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&resource.ID,
		&resource.ResourceType,
		&resource.RegionID,
		&resource.SpeciesName,
		&resource.Quantity,
		&resource.Area,
		&resource.GrowthStage,
		&resource.PlantedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// UpdateResource 更新现有资源记录。
func (c *ResourceController) UpdateResource(ctx context.Context, resource model.Resource) error {
	const query = `
		UPDATE resources
		SET resource_type = ?, region_id = ?, species_name = ?, quantity = ?, area = ?, growth_stage = ?, planted_at = ?
		WHERE resource_id = ?`
	result, err := c.db.ExecContext(ctx, query,
		resource.ResourceType,
		resource.RegionID,
		resource.SpeciesName,
		resource.Quantity,
		resource.Area,
		resource.GrowthStage,
		resource.PlantedAt,
		resource.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no resource updated for id=%d", resource.ID)
	}
	return nil
}

// DeleteResource 根据ID删除资源记录。
func (c *ResourceController) DeleteResource(ctx context.Context, id int64) error {
	// 开启事务
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 删除该资源的所有变更记录
	if _, err := tx.ExecContext(ctx, `DELETE FROM resource_changes WHERE resource_id = ?`, id); err != nil {
		return fmt.Errorf("删除资源变更记录失败: %v", err)
	}

	// 2. 删除资源
	result, err := tx.ExecContext(ctx, `DELETE FROM resources WHERE resource_id = ?`, id)
	if err != nil {
		return fmt.Errorf("删除资源失败: %v", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("资源不存在或已被删除, id=%d", id)
	}

	// 提交事务
	return tx.Commit()
}

// ListResources 根据条件获取资源列表。
func (c *ResourceController) ListResources(ctx context.Context, regionID *int64, resourceType *string, growthStage *string, speciesName *string) ([]model.Resource, error) {
	var query string
	var args []interface{}

	baseQuery := `
		SELECT resource_id, resource_type, region_id, species_name, quantity, area, growth_stage, planted_at, updated_at
		FROM resources`

	conditions := ""
	if regionID != nil {
		conditions += " AND region_id = ?"
		args = append(args, *regionID)
	}
	if resourceType != nil {
		conditions += " AND resource_type = ?"
		args = append(args, *resourceType)
	}
	if growthStage != nil {
		conditions += " AND growth_stage = ?"
		args = append(args, *growthStage)
	}
	if speciesName != nil {
		conditions += " AND species_name LIKE ?"
		args = append(args, "%"+*speciesName+"%")
	}

	if conditions != "" {
		query = baseQuery + " WHERE " + conditions[5:] + " ORDER BY resource_id"
	} else {
		query = baseQuery + " ORDER BY resource_id"
	}

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []model.Resource
	for rows.Next() {
		var resource model.Resource
		var updatedAt time.Time
		if err := rows.Scan(
			&resource.ID,
			&resource.ResourceType,
			&resource.RegionID,
			&resource.SpeciesName,
			&resource.Quantity,
			&resource.Area,
			&resource.GrowthStage,
			&resource.PlantedAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, rows.Err()
}

// CreateResourceChange 插入新的资源变更记录。
func (c *ResourceController) CreateResourceChange(ctx context.Context, change model.ResourceChange) error {
	const query = `
		INSERT INTO resource_changes (resource_id, change_type, change_reason, change_amount, change_area, changed_at, operator_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := c.db.ExecContext(ctx, query,
		change.ResourceID,
		change.ChangeType,
		change.ChangeReason,
		change.ChangeAmount,
		change.ChangeArea,
		change.ChangedAt,
		change.OperatorID,
	)
	return err
}

// GetResourceChangeByID 根据ID获取资源变更记录。
func (c *ResourceController) GetResourceChangeByID(ctx context.Context, id int64) (*model.ResourceChange, error) {
	const query = `
		SELECT change_id, resource_id, change_type, change_reason, change_amount, change_area, changed_at, operator_id
		FROM resource_changes WHERE change_id = ?`
	var change model.ResourceChange
	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&change.ID,
		&change.ResourceID,
		&change.ChangeType,
		&change.ChangeReason,
		&change.ChangeAmount,
		&change.ChangeArea,
		&change.ChangedAt,
		&change.OperatorID,
	)
	if err != nil {
		return nil, err
	}
	return &change, nil
}

// UpdateResourceChange 更新现有资源变更记录。
func (c *ResourceController) UpdateResourceChange(ctx context.Context, change model.ResourceChange) error {
	const query = `
		UPDATE resource_changes
		SET resource_id = ?, change_type = ?, change_reason = ?, change_amount = ?, change_area = ?, changed_at = ?, operator_id = ?
		WHERE change_id = ?`
	result, err := c.db.ExecContext(ctx, query,
		change.ResourceID,
		change.ChangeType,
		change.ChangeReason,
		change.ChangeAmount,
		change.ChangeArea,
		change.ChangedAt,
		change.OperatorID,
		change.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no resource change updated for id=%d", change.ID)
	}
	return nil
}

// DeleteResourceChange 根据ID删除资源变更记录。
func (c *ResourceController) DeleteResourceChange(ctx context.Context, id int64) error {
	const query = `DELETE FROM resource_changes WHERE change_id = ?`
	result, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no resource change deleted for id=%d", id)
	}
	return nil
}

// GetResourceChangesByResourceID 获取指定资源的变更记录。
func (c *ResourceController) GetResourceChangesByResourceID(ctx context.Context, resourceID int64) ([]model.ResourceChange, error) {
	const query = `
		SELECT rc.change_id, rc.resource_id, rc.change_type, rc.change_reason, rc.change_amount, rc.change_area, rc.changed_at, rc.operator_id, u.real_name as operator_name
		FROM resource_changes rc
		LEFT JOIN users u ON rc.operator_id = u.user_id
		WHERE rc.resource_id = ?
		ORDER BY rc.changed_at DESC`
	rows, err := c.db.QueryContext(ctx, query, resourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changes []model.ResourceChange
	for rows.Next() {
		var change model.ResourceChange
		var operatorName sql.NullString
		if err := rows.Scan(
			&change.ID,
			&change.ResourceID,
			&change.ChangeType,
			&change.ChangeReason,
			&change.ChangeAmount,
			&change.ChangeArea,
			&change.ChangedAt,
			&change.OperatorID,
			&operatorName,
		); err != nil {
			return nil, err
		}

		if operatorName.Valid {
			change.OperatorName = &operatorName.String
		}

		changes = append(changes, change)
	}
	return changes, rows.Err()
}

// ResourceSummaryWithChanges 汇总当前资源及变更历史。
func (c *ResourceController) ResourceSummaryWithChanges(ctx context.Context) ([]model.ResourceSummary, error) {
	const query = `
SELECT rg.region_name, res.resource_type,
       SUM(COALESCE(res.quantity,0)) AS current_quantity,
       SUM(COALESCE(res.area,0)) AS current_area,
       COUNT(rc.change_id) AS change_records,
       MAX(rc.changed_at) AS last_change_time
FROM resources res
JOIN regions rg ON res.region_id = rg.region_id
LEFT JOIN resource_changes rc ON rc.resource_id = res.resource_id
GROUP BY rg.region_name, res.resource_type
ORDER BY rg.region_name, res.resource_type`

	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.ResourceSummary
	for rows.Next() {
		var item model.ResourceSummary
		if err := rows.Scan(
			&item.RegionName,
			&item.ResourceType,
			&item.CurrentQuantity,
			&item.CurrentArea,
			&item.ChangeRecords,
			&item.LastChangeAt,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
