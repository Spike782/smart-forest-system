package controller

import (
	"context"
	"database/sql"

	"smart-forest-system/backend/model"
)

// StatisticsController 统计控制器
type StatisticsController struct {
	db *sql.DB
}

// NewStatisticsController 创建统计控制器
func NewStatisticsController(db *sql.DB) *StatisticsController {
	return &StatisticsController{db: db}
}

// DeviceFaultsWithMaintenance 汇总设备故障次数及最近维护情况。
func (c *StatisticsController) DeviceFaultsWithMaintenance(ctx context.Context) ([]model.DeviceFaultMaintenance, error) {
	const query = `
SELECT rg.region_name, d.device_id, d.device_name,
       SUM(CASE WHEN l.run_status='FAULT' THEN 1 ELSE 0 END) AS fault_times,
       COUNT(m.maintenance_id) AS maintenance_count,
       MAX(m.maintenance_time) AS last_maintenance
FROM devices d
JOIN regions rg ON d.install_region_id = rg.region_id
LEFT JOIN device_status_logs l ON l.device_id = d.device_id
LEFT JOIN maintenance_records m ON m.device_id = d.device_id
GROUP BY rg.region_name, d.device_id, d.device_name
HAVING fault_times > 0
ORDER BY fault_times DESC, last_maintenance DESC`

	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.DeviceFaultMaintenance
	for rows.Next() {
		var item model.DeviceFaultMaintenance
		if err := rows.Scan(
			&item.RegionName,
			&item.DeviceID,
			&item.DeviceName,
			&item.FaultTimes,
			&item.MaintenanceCount,
			&item.LastMaintenanceAt,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

// SensorValidityReport 返回指定时间范围内传感器数据有效性报告。
func (c *StatisticsController) SensorValidityReport(ctx context.Context, days int) ([]model.SensorValidity, error) {
	const query = `
SELECT rg.region_name, s.sensor_id, s.monitor_type,
       COUNT(sr.reading_id) AS total_readings,
       SUM(CASE WHEN sr.data_status='VALID' THEN 1 ELSE 0 END) AS valid_readings,
       SUM(CASE WHEN sr.data_status='INVALID' THEN 1 ELSE 0 END) AS invalid_readings,
       MAX(sr.collected_at) AS last_reading_at
FROM sensors s
JOIN regions rg ON s.region_id = rg.region_id
LEFT JOIN sensor_readings sr ON sr.sensor_id = s.sensor_id
WHERE sr.collected_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
GROUP BY rg.region_name, s.sensor_id, s.monitor_type
ORDER BY invalid_readings DESC, last_reading_at DESC`

	rows, err := c.db.QueryContext(ctx, query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.SensorValidity
	for rows.Next() {
		var item model.SensorValidity
		if err := rows.Scan(
			&item.RegionName,
			&item.SensorID,
			&item.MonitorType,
			&item.TotalReadings,
			&item.ValidReadings,
			&item.InvalidReadings,
			&item.LastReadingAt,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

// RegionAlertStatsLast30d 汇总最近30天各区域预警统计。
func (c *StatisticsController) RegionAlertStatsLast30d(ctx context.Context) ([]model.RegionAlertStats, error) {
	const query = `
SELECT rg.region_name,
	   COUNT(DISTINCT s.sensor_id) AS sensors,
	   COUNT(DISTINCT a.alert_id) AS alerts,
	   COUNT(DISTINCT n.notification_id) AS notifications_sent
FROM regions rg
LEFT JOIN sensors s ON s.region_id = rg.region_id
LEFT JOIN alerts a ON a.region_id = rg.region_id AND a.triggered_at >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
LEFT JOIN notifications n ON n.alert_id = a.alert_id
GROUP BY rg.region_name
ORDER BY alerts DESC, notifications_sent DESC`

	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.RegionAlertStats
	for rows.Next() {
		var item model.RegionAlertStats
		if err := rows.Scan(
			&item.RegionName,
			&item.Sensors,
			&item.Alerts,
			&item.NotificationsSent,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

// GetEnvironmentTrend 获取环境趋势数据
func (c *StatisticsController) GetEnvironmentTrend(ctx context.Context, days int, regionName *string, metricType *string) ([]model.EnvironmentTrend, error) {
	query := `
SELECT DATE(sr.collected_at) AS time, rg.region_name, sr.reading_type AS metric_type, 
       AVG(sr.numeric_value) AS metric_value, sr.value_unit AS unit
FROM sensor_readings sr
JOIN sensors s ON sr.sensor_id = s.sensor_id
JOIN regions rg ON s.region_id = rg.region_id
WHERE sr.data_status = 'VALID' AND sr.collected_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
`
	args := []interface{}{days}

	if regionName != nil {
		query += ` AND rg.region_name = ?`
		args = append(args, *regionName)
	}

	if metricType != nil {
		query += ` AND sr.reading_type = ?`
		args = append(args, *metricType)
	}

	query += `
GROUP BY DATE(sr.collected_at), rg.region_name, sr.reading_type, sr.value_unit
ORDER BY time ASC, rg.region_name, metric_type`

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.EnvironmentTrend
	for rows.Next() {
		var item model.EnvironmentTrend
		if err := rows.Scan(
			&item.Time,
			&item.RegionName,
			&item.MetricType,
			&item.MetricValue,
			&item.Unit,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

// GetAirQuality 获取空气质量数据
func (c *StatisticsController) GetAirQuality(ctx context.Context, days int, regionName *string) ([]model.AirQuality, error) {
	query := `
SELECT DATE(sr.collected_at) AS time, rg.region_name,
       AVG(CASE WHEN sr.reading_type = 'PM25' THEN sr.numeric_value ELSE NULL END) AS pm25,
       AVG(CASE WHEN sr.reading_type = 'PM10' THEN sr.numeric_value ELSE NULL END) AS pm10,
       AVG(CASE WHEN sr.reading_type = 'CO2' THEN sr.numeric_value ELSE NULL END) AS co2,
       AVG(CASE WHEN sr.reading_type = 'TEMPERATURE' THEN sr.numeric_value ELSE NULL END) AS temperature,
       AVG(CASE WHEN sr.reading_type = 'HUMIDITY' THEN sr.numeric_value ELSE NULL END) AS humidity
FROM sensor_readings sr
JOIN sensors s ON sr.sensor_id = s.sensor_id
JOIN regions rg ON s.region_id = rg.region_id
WHERE sr.data_status = 'VALID' AND sr.collected_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
`
	args := []interface{}{days}

	if regionName != nil {
		query += ` AND rg.region_name = ?`
		args = append(args, *regionName)
	}

	query += `
GROUP BY DATE(sr.collected_at), rg.region_name
ORDER BY time ASC, rg.region_name`

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.AirQuality
	for rows.Next() {
		var item model.AirQuality
		var pm25, pm10, co2, temp, hum sql.NullFloat64
		if err := rows.Scan(
			&item.Time,
			&item.RegionName,
			&pm25,
			&pm10,
			&co2,
			&temp,
			&hum,
		); err != nil {
			return nil, err
		}

		// 设置默认值
		item.PM25 = 0
		if pm25.Valid {
			item.PM25 = pm25.Float64
		}

		item.PM10 = 0
		if pm10.Valid {
			item.PM10 = pm10.Float64
		}

		item.CO2 = 0
		if co2.Valid {
			item.CO2 = co2.Float64
		}

		item.Temperature = 0
		if temp.Valid {
			item.Temperature = temp.Float64
		}

		item.Humidity = 0
		if hum.Valid {
			item.Humidity = hum.Float64
		}

		// 计算AQI（简化计算）
		item.AQI = int(item.PM25 * 2)
		if item.AQI < 50 {
			item.AQILevel = "优"
		} else if item.AQI < 100 {
			item.AQILevel = "良"
		} else if item.AQI < 150 {
			item.AQILevel = "轻度污染"
		} else if item.AQI < 200 {
			item.AQILevel = "中度污染"
		} else if item.AQI < 300 {
			item.AQILevel = "重度污染"
		} else {
			item.AQILevel = "严重污染"
		}

		result = append(result, item)
	}
	return result, rows.Err()
}

// GetAnomalyData 获取异常数据
func (c *StatisticsController) GetAnomalyData(ctx context.Context, days int, regionName *string, anomalyType *string) ([]model.AnomalyData, error) {
	query := `
SELECT sr.collected_at AS time, rg.region_name, sr.sensor_id, sr.reading_type,
       sr.numeric_value AS value, sr.value_unit AS unit,
       CASE 
           WHEN sr.numeric_value > 100 THEN 'HIGH'
           WHEN sr.numeric_value < 0 THEN 'LOW'
           ELSE 'OUT_OF_RANGE'
       END AS anomaly_type,
       100 AS threshold
FROM sensor_readings sr
JOIN sensors s ON sr.sensor_id = s.sensor_id
JOIN regions rg ON s.region_id = rg.region_id
WHERE sr.data_status = 'INVALID' AND sr.collected_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
`
	args := []interface{}{days}

	if regionName != nil {
		query += ` AND rg.region_name = ?`
		args = append(args, *regionName)
	}

	if anomalyType != nil {
		query += ` AND CASE 
	           WHEN sr.numeric_value > 100 THEN 'HIGH'
	           WHEN sr.numeric_value < 0 THEN 'LOW'
	           ELSE 'OUT_OF_RANGE'
	       END = ?`
		args = append(args, *anomalyType)
	}

	query += `
ORDER BY time DESC`

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.AnomalyData
	for rows.Next() {
		var item model.AnomalyData
		if err := rows.Scan(
			&item.Time,
			&item.RegionName,
			&item.SensorID,
			&item.ReadingType,
			&item.Value,
			&item.Unit,
			&item.AnomalyType,
			&item.Threshold,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

// GetSensorDataSummary 获取传感器数据汇总
func (c *StatisticsController) GetSensorDataSummary(ctx context.Context, regionName *string, monitorType *string) ([]model.SensorDataSummary, error) {
	query := `
SELECT s.sensor_id, rg.region_name, s.monitor_type,
       COUNT(sr.reading_id) AS total_readings,
       SUM(CASE WHEN sr.data_status = 'VALID' THEN 1 ELSE 0 END) AS valid_readings,
       AVG(CASE WHEN sr.data_status = 'VALID' THEN sr.numeric_value ELSE NULL END) AS avg_value,
       MIN(CASE WHEN sr.data_status = 'VALID' THEN sr.numeric_value ELSE NULL END) AS min_value,
       MAX(CASE WHEN sr.data_status = 'VALID' THEN sr.numeric_value ELSE NULL END) AS max_value,
       MAX(sr.collected_at) AS last_reading
FROM sensors s
JOIN regions rg ON s.region_id = rg.region_id
LEFT JOIN sensor_readings sr ON sr.sensor_id = s.sensor_id
`
	args := []interface{}{}

	if regionName != nil {
		query += ` WHERE rg.region_name = ?`
		args = append(args, *regionName)
		if monitorType != nil {
			query += ` AND s.monitor_type = ?`
			args = append(args, *monitorType)
		}
	} else if monitorType != nil {
		query += ` WHERE s.monitor_type = ?`
		args = append(args, *monitorType)
	}

	query += `
GROUP BY s.sensor_id, rg.region_name, s.monitor_type
ORDER BY s.sensor_id`

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.SensorDataSummary
	for rows.Next() {
		var item model.SensorDataSummary
		if err := rows.Scan(
			&item.SensorID,
			&item.RegionName,
			&item.MonitorType,
			&item.TotalReadings,
			&item.ValidReadings,
			&item.AvgValue,
			&item.MinValue,
			&item.MaxValue,
			&item.LastReading,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

// GetRegionDailyAverages 获取区域日均值数据
func (c *StatisticsController) GetRegionDailyAverages(ctx context.Context, days int, regionName *string, metricType *string) ([]model.RegionDailyAverage, error) {
	query := `
SELECT DATE(sr.collected_at) AS date, rg.region_name, sr.reading_type AS metric_type,
       AVG(sr.numeric_value) AS avg_value, sr.value_unit AS unit,
       COUNT(sr.reading_id) AS sample_count
FROM sensor_readings sr
JOIN sensors s ON sr.sensor_id = s.sensor_id
JOIN regions rg ON s.region_id = rg.region_id
WHERE sr.data_status = 'VALID' AND sr.collected_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
`
	args := []interface{}{days}

	if regionName != nil {
		query += ` AND rg.region_name = ?`
		args = append(args, *regionName)
	}

	if metricType != nil {
		query += ` AND sr.reading_type = ?`
		args = append(args, *metricType)
	}

	query += `
GROUP BY DATE(sr.collected_at), rg.region_name, sr.reading_type, sr.value_unit
ORDER BY date ASC, rg.region_name, metric_type`

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.RegionDailyAverage
	for rows.Next() {
		var item model.RegionDailyAverage
		if err := rows.Scan(
			&item.Date,
			&item.RegionName,
			&item.MetricType,
			&item.AvgValue,
			&item.Unit,
			&item.SampleCount,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
