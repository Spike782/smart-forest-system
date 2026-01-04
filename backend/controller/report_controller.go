package controller

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"smart-forest-system/backend/model"

	"github.com/xuri/excelize/v2"
)

// ReportController 报表控制器
type ReportController struct {
	db *sql.DB
}

// NewReportController 创建报表控制器
func NewReportController(db *sql.DB) *ReportController {
	return &ReportController{db: db}
}

// CreateReportTemplate 插入新的报表模板。
func (c *ReportController) CreateReportTemplate(ctx context.Context, template model.ReportTemplate) (int64, error) {
	const query = `
		INSERT INTO report_templates (template_name, dimension, metrics, cycle, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	result, err := c.db.ExecContext(ctx, query,
		template.Name,
		template.Dimension,
		template.Metrics,
		template.Cycle,
		template.IsActive,
		template.CreatedAt,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// GetReportTemplateByID 根据ID获取报表模板。
func (c *ReportController) GetReportTemplateByID(ctx context.Context, id int64) (*model.ReportTemplate, error) {
	const query = `
		SELECT template_id, template_name, dimension, metrics, cycle, is_active, created_at
		FROM report_templates WHERE template_id = ?`
	var template model.ReportTemplate
	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&template.ID,
		&template.Name,
		&template.Dimension,
		&template.Metrics,
		&template.Cycle,
		&template.IsActive,
		&template.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// UpdateReportTemplate 更新现有报表模板。
func (c *ReportController) UpdateReportTemplate(ctx context.Context, template model.ReportTemplate) error {
	const query = `
		UPDATE report_templates
		SET template_name = ?, dimension = ?, metrics = ?, cycle = ?, is_active = ?
		WHERE template_id = ?`
	result, err := c.db.ExecContext(ctx, query,
		template.Name,
		template.Dimension,
		template.Metrics,
		template.Cycle,
		template.IsActive,
		template.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no report template updated for id=%d", template.ID)
	}
	return nil
}

// DeleteReportTemplate 根据ID删除报表模板。
func (c *ReportController) DeleteReportTemplate(ctx context.Context, id int64) error {
	// 开启事务
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. 删除该模板的所有报表
	if _, err := tx.ExecContext(ctx, `DELETE FROM reports WHERE template_id = ?`, id); err != nil {
		return fmt.Errorf("删除报表失败: %v", err)
	}

	// 2. 删除报表模板
	result, err := tx.ExecContext(ctx, `DELETE FROM report_templates WHERE template_id = ?`, id)
	if err != nil {
		return fmt.Errorf("删除报表模板失败: %v", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("报表模板不存在或已被删除, id=%d", id)
	}

	// 提交事务
	return tx.Commit()
}

// ListReportTemplates 根据条件获取报表模板列表。
func (c *ReportController) ListReportTemplates(ctx context.Context, dimension *string, isActive *bool) ([]model.ReportTemplate, error) {
	var query string
	var args []interface{}

	baseQuery := `
		SELECT template_id, template_name, dimension, metrics, cycle, is_active, created_at
		FROM report_templates`

	conditions := ""
	if dimension != nil {
		conditions += " AND dimension = ?"
		args = append(args, *dimension)
	}
	if isActive != nil {
		conditions += " AND is_active = ?"
		args = append(args, *isActive)
	}

	if conditions != "" {
		query = baseQuery + " WHERE " + conditions[5:] + " ORDER BY template_id"
	} else {
		query = baseQuery + " ORDER BY template_id"
	}

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []model.ReportTemplate
	for rows.Next() {
		var template model.ReportTemplate
		if err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.Dimension,
			&template.Metrics,
			&template.Cycle,
			&template.IsActive,
			&template.CreatedAt,
		); err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	return templates, rows.Err()
}

// CreateReport 插入新的报表记录。
func (c *ReportController) CreateReport(ctx context.Context, report model.Report) (int64, error) {
	const query = `
		INSERT INTO reports (template_id, period, generated_at, file_path, data_source)
		VALUES (?, ?, ?, ?, ?)`
	result, err := c.db.ExecContext(ctx, query,
		report.TemplateID,
		report.Period,
		report.GeneratedAt,
		report.FilePath,
		report.DataSource,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return id, err
}

// GetReportByID 根据ID获取报表记录。
func (c *ReportController) GetReportByID(ctx context.Context, id int64) (*model.Report, error) {
	const query = `
		SELECT report_id, template_id, period, generated_at, file_path, data_source
		FROM reports WHERE report_id = ?`
	var report model.Report
	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&report.ID,
		&report.TemplateID,
		&report.Period,
		&report.GeneratedAt,
		&report.FilePath,
		&report.DataSource,
	)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// UpdateReport 更新现有报表记录。
func (c *ReportController) UpdateReport(ctx context.Context, report model.Report) error {
	const query = `
		UPDATE reports
		SET template_id = ?, period = ?, generated_at = ?, file_path = ?, data_source = ?
		WHERE report_id = ?`
	result, err := c.db.ExecContext(ctx, query,
		report.TemplateID,
		report.Period,
		report.GeneratedAt,
		report.FilePath,
		report.DataSource,
		report.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no report updated for id=%d", report.ID)
	}
	return nil
}

// DeleteReport 根据ID删除报表记录。
func (c *ReportController) DeleteReport(ctx context.Context, id int64) error {
	const query = `DELETE FROM reports WHERE report_id = ?`
	result, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no report deleted for id=%d", id)
	}
	return nil
}

// GetReportsByTemplateID 获取指定模板的报表记录。
func (c *ReportController) GetReportsByTemplateID(ctx context.Context, templateID int64, reportType *string, status *string) ([]model.Report, error) {
	var query string
	var args []interface{}
	var conditions []string
	
	// 基础查询
	baseQuery := `
		SELECT report_id, template_id, period, generated_at, file_path, data_source
		FROM reports`
	
	// 添加条件
	if templateID > 0 {
		conditions = append(conditions, "template_id = ?")
		args = append(args, templateID)
	}
	
	if reportType != nil {
		// 根据report_type查询，需要关联report_templates表
		baseQuery = `
		SELECT r.report_id, r.template_id, r.period, r.generated_at, r.file_path, r.data_source
		FROM reports r
		JOIN report_templates rt ON r.template_id = rt.template_id`
		// 前端传递的reportType是周期（cycle），如daily、weekly等，不是维度（dimension）
		conditions = append(conditions, "rt.cycle = ?")
		args = append(args, *reportType)
	}
	
	// 构建完整查询
	if len(conditions) > 0 {
		query = baseQuery + " WHERE " + strings.Join(conditions, " AND ") + " ORDER BY generated_at DESC"
	} else {
		query = baseQuery + " ORDER BY generated_at DESC"
	}
	
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []model.Report
	for rows.Next() {
		var report model.Report
		if err := rows.Scan(
			&report.ID,
			&report.TemplateID,
			&report.Period,
			&report.GeneratedAt,
			&report.FilePath,
			&report.DataSource,
		); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, rows.Err()
}

// GenerateReport 生成报表
func (c *ReportController) GenerateReport(ctx context.Context, templateID int64, period string) (*model.Report, error) {
	// 1. 获取报表模板
	template, err := c.GetReportTemplateByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("获取报表模板失败: %v", err)
	}

	// 2. 创建报表数据目录
	reportDir := "./static/reports"
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return nil, fmt.Errorf("创建报表目录失败: %v", err)
	}

	// 3. 生成报表文件名
	fileName := fmt.Sprintf("report_%d_%s_%s.xlsx", templateID, period, time.Now().Format("20060102150405"))
	filePath := filepath.Join(reportDir, fileName)

	// 4. 根据模板生成报表数据
	data, err := c.generateReportData(ctx, template, period)
	if err != nil {
		return nil, fmt.Errorf("生成报表数据失败: %v", err)
	}

	// 5. 生成Excel文件
	if err := c.generateExcelReport(filePath, template, period, data); err != nil {
		return nil, fmt.Errorf("生成Excel报表失败: %v", err)
	}

	// 6. 创建报表记录
	report := model.Report{
		TemplateID:  templateID,
		Period:      period,
		GeneratedAt: time.Now(),
		FilePath:    filePath,
		DataSource:  nil,
	}

	// 7. 保存报表记录到数据库
	reportID, err := c.CreateReport(ctx, report)
	if err != nil {
		return nil, fmt.Errorf("保存报表记录失败: %v", err)
	}

	// 8. 获取完整的报表信息
	return c.GetReportByID(ctx, reportID)
}

// ExportReport 导出报表
func (c *ReportController) ExportReport(ctx context.Context, reportID int64, format string) (string, error) {
	// 只支持Excel格式
	if format != "excel" {
		return "", fmt.Errorf("只支持Excel格式导出")
	}

	// 1. 获取报表信息
	report, err := c.GetReportByID(ctx, reportID)
	if err != nil {
		return "", fmt.Errorf("获取报表失败: %v", err)
	}

	// 2. 处理报表文件路径
	// 获取文件名
	reportFileName := filepath.Base(report.FilePath)

	// 确保文件名是Excel格式
	excelFileName := strings.TrimSuffix(reportFileName, ".pdf") + ".xlsx"

	// 构造完整的Excel文件路径
	fullExcelPath := filepath.Join("./static/reports", excelFileName)

	// 3. 检查Excel文件是否存在
	if _, err := os.Stat(fullExcelPath); os.IsNotExist(err) {
		// 如果文件不存在，先检查static/reports目录是否存在
		if _, err := os.Stat("./static/reports"); os.IsNotExist(err) {
			// 如果目录不存在，创建它
			if err := os.MkdirAll("./static/reports", 0755); err != nil {
				return "", fmt.Errorf("创建报表目录失败: %v", err)
			}
		}

		// 尝试生成Excel文件
		template, err := c.GetReportTemplateByID(ctx, report.TemplateID)
		if err != nil {
			return "", fmt.Errorf("获取报表模板失败: %v", err)
		}

		// 生成报表数据
		data, err := c.generateReportData(ctx, template, report.Period)
		if err != nil {
			return "", fmt.Errorf("生成报表数据失败: %v", err)
		}

		// 生成Excel文件
		if err := c.generateExcelReport(fullExcelPath, template, report.Period, data); err != nil {
			return "", fmt.Errorf("生成Excel报表失败: %v", err)
		}
	}

	return fullExcelPath, nil
}

// generateReportData 根据模板生成报表数据
func (c *ReportController) generateReportData(ctx context.Context, template *model.ReportTemplate, period string) ([][]interface{}, error) {
	// 根据模板的维度和指标查询数据
	// 这里实现一个简化版，实际应该根据不同的维度和指标查询不同的数据
	var data [][]interface{}

	// 添加表头
	header := []interface{}{"ID", "名称", "类型", "值", "单位", "时间"}
	data = append(data, header)

	// 根据不同的模板维度生成数据
	switch template.Dimension {
	case "REGION":
		// 查询区域相关数据
		rows, err := c.db.QueryContext(ctx, `
			SELECT r.region_id, r.region_name, r.region_type, COUNT(s.sensor_id) as sensor_count, 
			       '个' as unit, CURRENT_DATE() as date
			FROM regions r
			LEFT JOIN sensors s ON r.region_id = s.region_id
			GROUP BY r.region_id, r.region_name, r.region_type
		`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var id int64
			var name, regionType, unit string
			var count int
			var date time.Time
			if err := rows.Scan(&id, &name, &regionType, &count, &unit, &date); err != nil {
				return nil, err
			}
			data = append(data, []interface{}{id, name, regionType, count, unit, date})
		}
	case "TIME":
			// 查询时间相关数据
			rows, err := c.db.QueryContext(ctx, `
				SELECT a.alert_id, r.region_name, a.alert_type, 
				       COUNT(n.notification_id) as notification_count, 
				       '条' as unit, a.triggered_at
				FROM alerts a
				LEFT JOIN regions r ON a.region_id = r.region_id
				LEFT JOIN notifications n ON a.alert_id = n.alert_id
				WHERE DATE(a.triggered_at) >= DATE_SUB(CURRENT_DATE(), INTERVAL 30 DAY)
				GROUP BY a.alert_id, r.region_name, a.alert_type, a.triggered_at
			`)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var id int64
				var regionName, alertType, unit string
				var count int
				var triggeredAt time.Time
				if err := rows.Scan(&id, &regionName, &alertType, &count, &unit, &triggeredAt); err != nil {
					return nil, err
				}
				data = append(data, []interface{}{id, regionName, alertType, count, unit, triggeredAt})
			}
	case "TYPE":
		// 查询类型相关数据
		rows, err := c.db.QueryContext(ctx, `
			SELECT d.device_id, d.device_name, d.device_type, 
			       AVG(COALESCE(dsl.battery_percent, 0)) as avg_battery, 
			       '%' as unit, CURRENT_DATE() as date
			FROM devices d
			LEFT JOIN device_status_logs dsl ON d.device_id = dsl.device_id
			GROUP BY d.device_id, d.device_name, d.device_type
		`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var id int64
			var name, deviceType, unit string
			var avgBattery float64
			var date time.Time
			if err := rows.Scan(&id, &name, &deviceType, &avgBattery, &unit, &date); err != nil {
				return nil, err
			}
			data = append(data, []interface{}{id, name, deviceType, avgBattery, unit, date})
		}
	default:
		return nil, fmt.Errorf("不支持的报表维度: %s", template.Dimension)
	}

	return data, nil
}

// generateExcelReport 生成Excel报表
func (c *ReportController) generateExcelReport(filePath string, template *model.ReportTemplate, period string, data [][]interface{}) error {
	// 创建Excel文件
	f := excelize.NewFile()
	defer f.Close()

	// 创建工作表
	sheetName := "报表数据"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	// 设置表头样式
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#E0EBF5"},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return err
	}

	// 写入报表标题
	title := fmt.Sprintf("%s - %s", template.Name, period)
	f.MergeCell(sheetName, "A1", "F1")
	f.SetCellValue(sheetName, "A1", title)
	titleStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return err
	}
	f.SetCellStyle(sheetName, "A1", "F1", titleStyle)

	// 写入数据
	for rowIndex, rowData := range data {
		for colIndex, cellData := range rowData {
			cell := fmt.Sprintf("%c%d", 'A'+colIndex, rowIndex+2) // 从第二行开始
			f.SetCellValue(sheetName, cell, cellData)
			
			// 设置表头样式
			if rowIndex == 0 {
				f.SetCellStyle(sheetName, cell, cell, style)
			}
		}
	}

	// 自动调整列宽
	for col := 'A'; col <= 'F'; col++ {
		colStr := string(col)
		if err := f.SetColWidth(sheetName, colStr, colStr, 15); err != nil {
			return err
		}
	}

	// 保存文件
	if err := f.SaveAs(filePath); err != nil {
		return err
	}

	return nil
}

