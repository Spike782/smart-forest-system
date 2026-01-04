package controller

import (
	"context"
	"database/sql"
	"fmt"

	"smart-forest-system/backend/model"
)

// AlertController 预警控制器
type AlertController struct {
	db *sql.DB
}

// NewAlertController 创建预警控制器
func NewAlertController(db *sql.DB) *AlertController {
	return &AlertController{db: db}
}

// CreateAlertRule 插入新的预警规则。
func (c *AlertController) CreateAlertRule(ctx context.Context, rule model.AlertRule) (int64, error) {
	const query = `
		INSERT INTO alert_rules (alert_type, condition_expr, severity_level, is_active)
		VALUES (?, ?, ?, ?)`
	result, err := c.db.ExecContext(ctx, query,
		rule.AlertType,
		rule.ConditionExpr,
		rule.SeverityLevel,
		rule.IsActive,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// GetAlertRuleByID 根据ID获取预警规则。
func (c *AlertController) GetAlertRuleByID(ctx context.Context, id int64) (*model.AlertRule, error) {
	const query = `
		SELECT rule_id, alert_type, condition_expr, severity_level, is_active, created_at
		FROM alert_rules WHERE rule_id = ?`
	var rule model.AlertRule
	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&rule.ID,
		&rule.AlertType,
		&rule.ConditionExpr,
		&rule.SeverityLevel,
		&rule.IsActive,
		&rule.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// UpdateAlertRule 更新现有预警规则。
func (c *AlertController) UpdateAlertRule(ctx context.Context, rule model.AlertRule) error {
	const query = `
		UPDATE alert_rules
		SET alert_type = ?, condition_expr = ?, severity_level = ?, is_active = ?
		WHERE rule_id = ?`
	result, err := c.db.ExecContext(ctx, query,
		rule.AlertType,
		rule.ConditionExpr,
		rule.SeverityLevel,
		rule.IsActive,
		rule.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no alert rule updated for id=%d", rule.ID)
	}
	return nil
}

// DeleteAlertRule 根据ID删除预警规则。
func (c *AlertController) DeleteAlertRule(ctx context.Context, id int64) error {
	// 开启事务
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 删除该规则下的所有预警通知
	if _, err := tx.ExecContext(ctx, `DELETE FROM notifications WHERE alert_id IN (SELECT alert_id FROM alerts WHERE rule_id = ?)`, id); err != nil {
		return fmt.Errorf("删除预警通知失败: %v", err)
	}

	// 2. 删除该规则下的所有预警
	if _, err := tx.ExecContext(ctx, `DELETE FROM alerts WHERE rule_id = ?`, id); err != nil {
		return fmt.Errorf("删除预警失败: %v", err)
	}

	// 3. 删除预警规则
	result, err := tx.ExecContext(ctx, `DELETE FROM alert_rules WHERE rule_id = ?`, id)
	if err != nil {
		return fmt.Errorf("删除预警规则失败: %v", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("预警规则不存在或已被删除, id=%d", id)
	}

	// 提交事务
	return tx.Commit()
}

// ListAlertRules 根据条件获取预警规则列表。
func (c *AlertController) ListAlertRules(ctx context.Context, alertType *string, isActive *bool, severityLevel *string) ([]model.AlertRule, error) {
	var query string
	var args []interface{}

	baseQuery := `
		SELECT rule_id, alert_type, condition_expr, severity_level, is_active, created_at
		FROM alert_rules`

	conditions := ""
	if alertType != nil {
		conditions += " AND alert_type = ?"
		args = append(args, *alertType)
	}
	if isActive != nil {
		conditions += " AND is_active = ?"
		args = append(args, *isActive)
	}
	if severityLevel != nil {
		conditions += " AND severity_level = ?"
		args = append(args, *severityLevel)
	}

	if conditions != "" {
		query = baseQuery + " WHERE " + conditions[5:] + " ORDER BY rule_id"
	} else {
		query = baseQuery + " ORDER BY rule_id"
	}

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []model.AlertRule
	for rows.Next() {
		var rule model.AlertRule
		if err := rows.Scan(
			&rule.ID,
			&rule.AlertType,
			&rule.ConditionExpr,
			&rule.SeverityLevel,
			&rule.IsActive,
			&rule.CreatedAt,
		); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}

// CreateAlert 插入新的预警记录。
func (c *AlertController) CreateAlert(ctx context.Context, alert model.Alert) (int64, error) {
	const query = `
		INSERT INTO alerts (rule_id, region_id, triggered_at, content, status, handler_id, handle_result)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := c.db.ExecContext(ctx, query,
		alert.RuleID,
		alert.RegionID,
		alert.TriggeredAt,
		alert.Content,
		alert.Status,
		alert.HandlerID,
		alert.HandleResult,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// GetAlertByID 根据ID获取预警记录。
func (c *AlertController) GetAlertByID(ctx context.Context, id int64) (*model.Alert, error) {
	const query = `
		SELECT alert_id, rule_id, region_id, triggered_at, content, status, handler_id, handle_result
		FROM alerts WHERE alert_id = ?`
	var alert model.Alert
	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&alert.ID,
		&alert.RuleID,
		&alert.RegionID,
		&alert.TriggeredAt,
		&alert.Content,
		&alert.Status,
		&alert.HandlerID,
		&alert.HandleResult,
	)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// UpdateAlert 更新现有预警记录。
func (c *AlertController) UpdateAlert(ctx context.Context, alert model.Alert) error {
	const query = `
		UPDATE alerts
		SET rule_id = ?, region_id = ?, triggered_at = ?, content = ?, status = ?, handler_id = ?, handle_result = ?
		WHERE alert_id = ?`
	result, err := c.db.ExecContext(ctx, query,
		alert.RuleID,
		alert.RegionID,
		alert.TriggeredAt,
		alert.Content,
		alert.Status,
		alert.HandlerID,
		alert.HandleResult,
		alert.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no alert updated for id=%d", alert.ID)
	}
	return nil
}

// DeleteAlert 根据ID删除预警记录。
func (c *AlertController) DeleteAlert(ctx context.Context, id int64) error {
	// 开启事务
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 删除该预警的所有通知
	if _, err := tx.ExecContext(ctx, `DELETE FROM notifications WHERE alert_id = ?`, id); err != nil {
		return fmt.Errorf("删除预警通知失败: %v", err)
	}

	// 2. 删除预警记录
	result, err := tx.ExecContext(ctx, `DELETE FROM alerts WHERE alert_id = ?`, id)
	if err != nil {
		return fmt.Errorf("删除预警失败: %v", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("预警不存在或已被删除, id=%d", id)
	}

	// 提交事务
	return tx.Commit()
}

// ListAlerts 根据条件获取预警记录列表。
func (c *AlertController) ListAlerts(ctx context.Context, regionID *int64, status *string, validity *string) ([]model.Alert, error) {
	var query string
	var args []interface{}

	baseQuery := `
		SELECT a.alert_id, a.rule_id, a.region_id, a.triggered_at, a.content, a.status, a.handler_id, a.handle_result, ar.alert_type, ar.severity_level
		FROM alerts a
		LEFT JOIN alert_rules ar ON a.rule_id = ar.rule_id
		LEFT JOIN sensor_readings sr ON a.triggered_at BETWEEN sr.collected_at - INTERVAL 5 MINUTE AND sr.collected_at + INTERVAL 5 MINUTE
		LEFT JOIN sensors s ON sr.sensor_id = s.sensor_id AND s.region_id = a.region_id`

	conditions := ""
	if regionID != nil {
		conditions += " AND a.region_id = ?"
		args = append(args, *regionID)
	}
	if status != nil {
		conditions += " AND a.status = ?"
		args = append(args, *status)
	}
	if validity != nil {
		conditions += " AND sr.data_status = ?"
		args = append(args, *validity)
	}

	if conditions != "" {
		query = baseQuery + " WHERE " + conditions[5:] + " GROUP BY a.alert_id ORDER BY a.triggered_at DESC"
	} else {
		query = baseQuery + " GROUP BY a.alert_id ORDER BY a.triggered_at DESC"
	}

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []model.Alert
	for rows.Next() {
		var alert model.Alert
		var alertType sql.NullString
		var severity sql.NullString
		if err := rows.Scan(
			&alert.ID,
			&alert.RuleID,
			&alert.RegionID,
			&alert.TriggeredAt,
			&alert.Content,
			&alert.Status,
			&alert.HandlerID,
			&alert.HandleResult,
			&alertType,
			&severity,
		); err != nil {
			return nil, err
		}

		if alertType.Valid {
			alert.AlertType = alertType.String
		} else {
			alert.AlertType = "未知类型"
		}

		if severity.Valid {
			alert.Severity = severity.String
		} else {
			alert.Severity = "未知"
		}

		alerts = append(alerts, alert)
	}
	return alerts, rows.Err()
}

// CreateNotification 插入新的通知记录。
func (c *AlertController) CreateNotification(ctx context.Context, notification model.Notification) error {
	const query = `
		INSERT INTO notifications (alert_id, receiver_id, channel, sent_at, receive_status)
		VALUES (?, ?, ?, ?, ?)`
	_, err := c.db.ExecContext(ctx, query,
		notification.AlertID,
		notification.ReceiverID,
		notification.Channel,
		notification.SentAt,
		notification.ReceiveStatus,
	)
	return err
}

// GetNotificationByID 根据ID获取通知记录。
func (c *AlertController) GetNotificationByID(ctx context.Context, id int64) (*model.Notification, error) {
	const query = `
		SELECT notification_id, alert_id, receiver_id, channel, sent_at, receive_status
		FROM notifications WHERE notification_id = ?`
	var notification model.Notification
	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&notification.ID,
		&notification.AlertID,
		&notification.ReceiverID,
		&notification.Channel,
		&notification.SentAt,
		&notification.ReceiveStatus,
	)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// UpdateNotification 更新现有通知记录。
func (c *AlertController) UpdateNotification(ctx context.Context, notification model.Notification) error {
	const query = `
		UPDATE notifications
		SET alert_id = ?, receiver_id = ?, channel = ?, sent_at = ?, receive_status = ?
		WHERE notification_id = ?`
	result, err := c.db.ExecContext(ctx, query,
		notification.AlertID,
		notification.ReceiverID,
		notification.Channel,
		notification.SentAt,
		notification.ReceiveStatus,
		notification.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no notification updated for id=%d", notification.ID)
	}
	return nil
}

// DeleteNotification 根据ID删除通知记录。
func (c *AlertController) DeleteNotification(ctx context.Context, id int64) error {
	const query = `DELETE FROM notifications WHERE notification_id = ?`
	result, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no notification deleted for id=%d", id)
	}
	return nil
}

// GetNotificationsByAlertID 获取指定预警的所有通知记录。
func (c *AlertController) GetNotificationsByAlertID(ctx context.Context, alertID int64) ([]model.Notification, error) {
	const query = `
		SELECT n.notification_id, n.alert_id, n.receiver_id, n.channel, n.sent_at, n.receive_status,
		       a.content, u.real_name
		FROM notifications n
		JOIN alerts a ON n.alert_id = a.alert_id
		JOIN users u ON n.receiver_id = u.user_id
		WHERE n.alert_id = ?
		ORDER BY n.sent_at DESC`
	rows, err := c.db.QueryContext(ctx, query, alertID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		var alertContent sql.NullString
		var receiverName sql.NullString

		if err := rows.Scan(
			&notification.ID,
			&notification.AlertID,
			&notification.ReceiverID,
			&notification.Channel,
			&notification.SentAt,
			&notification.ReceiveStatus,
			&alertContent,
			&receiverName,
		); err != nil {
			return nil, err
		}

		// 处理NULL值
		if alertContent.Valid {
			notification.AlertContent = alertContent.String
		}
		if receiverName.Valid {
			notification.ReceiverName = receiverName.String
		}

		notifications = append(notifications, notification)
	}

	return notifications, rows.Err()
}

// RecentFireAlerts 返回指定区域最近N天的火灾预警及通知统计。
func (c *AlertController) RecentFireAlerts(ctx context.Context, regionName string, days int) ([]model.FireAlertReport, error) {
	const query = `
SELECT r.region_name, a.alert_id, a.triggered_at, ar.severity_level, a.status,
       a.handler_id, a.handle_result,
       COUNT(n.notification_id) AS notifications_sent,
       SUM(n.receive_status='DELIVERED') AS delivered_count
FROM alerts a
JOIN alert_rules ar ON a.rule_id = ar.rule_id AND ar.alert_type='FIRE'
JOIN regions r ON a.region_id = r.region_id
LEFT JOIN notifications n ON n.alert_id = a.alert_id
WHERE a.triggered_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
  AND r.region_name = ?
GROUP BY r.region_name, a.alert_id, a.triggered_at, ar.severity_level, a.status, a.handler_id, a.handle_result
ORDER BY a.triggered_at DESC`

	rows, err := c.db.QueryContext(ctx, query, days, regionName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.FireAlertReport
	for rows.Next() {
		var item model.FireAlertReport
		if err := rows.Scan(
			&item.RegionName,
			&item.AlertID,
			&item.TriggeredAt,
			&item.SeverityLevel,
			&item.Status,
			&item.HandlerID,
			&item.HandleResult,
			&item.NotificationsSent,
			&item.DeliveredCount,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
