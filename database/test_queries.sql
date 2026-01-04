-- 查询1: 查询某区域近7天火灾预警及处理情况
SELECT
    rg.region_name,
    a.alert_id,
    ar.alert_type,
    ar.severity_level,
    a.triggered_at,
    a.status AS alert_status,
    a.handler_id,
    u.real_name AS handler_name,
    a.handle_result,
    COUNT(n.notification_id) AS notification_count
FROM alerts a
JOIN regions rg ON a.region_id = rg.region_id
JOIN alert_rules ar ON a.rule_id = ar.rule_id
LEFT JOIN notifications n ON a.alert_id = n.alert_id
LEFT JOIN users u ON a.handler_id = u.user_id
WHERE
    ar.alert_type = 'FIRE'
    AND rg.region_name = '常绿森林A区'
    AND a.triggered_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)
GROUP BY
    a.alert_id, rg.region_name, ar.alert_type, ar.severity_level,
    a.triggered_at, a.status, a.handler_id, u.real_name, a.handle_result
ORDER BY a.triggered_at DESC;

-- 查询2: 统计各区域设备故障次数及维护情况
SELECT 
    rg.region_name,
    d.device_id,
    d.device_name,
    d.device_type,
    SUM(CASE WHEN dsl.run_status = 'FAULT' THEN 1 ELSE 0 END) AS fault_count,
    COUNT(mr.maintenance_id) AS maintenance_count,
    MAX(mr.maintenance_time) AS last_maintenance_time,
    MAX(dsl.collected_at) AS last_status_update
FROM devices d
JOIN regions rg ON d.install_region_id = rg.region_id
LEFT JOIN device_status_logs dsl ON d.device_id = dsl.device_id
LEFT JOIN maintenance_records mr ON d.device_id = mr.device_id
GROUP BY 
    rg.region_name, d.device_id, d.device_name, d.device_type
HAVING fault_count > 0
ORDER BY 
    fault_count DESC, rg.region_name, d.device_name;

-- 查询3: 查询某区域近30天环境监测数据趋势
SELECT
    DATE(sr.collected_at) AS date,
    rg.region_name,
    s.monitor_type,
    AVG(CASE WHEN sr.reading_type = 'TEMPERATURE' THEN sr.numeric_value ELSE NULL END) AS avg_temperature,
    AVG(CASE WHEN sr.reading_type = 'HUMIDITY' THEN sr.numeric_value ELSE NULL END) AS avg_humidity,
    AVG(CASE WHEN sr.reading_type = 'PM25' THEN sr.numeric_value ELSE NULL END) AS avg_pm25,
    AVG(CASE WHEN sr.reading_type = 'PM10' THEN sr.numeric_value ELSE NULL END) AS avg_pm10
FROM sensor_readings sr
JOIN sensors s ON sr.sensor_id = s.sensor_id
JOIN regions rg ON s.region_id = rg.region_id
WHERE
    rg.region_name = '常绿森林A区'
    AND sr.collected_at >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
    AND sr.data_status = 'VALID'
GROUP BY
    DATE(sr.collected_at), rg.region_name, s.monitor_type
ORDER BY
    date ASC;

-- 查询4: 统计各区域资源总量及变动情况
SELECT
    rg.region_name,
    r.resource_type,
    r.species_name,
    MAX(r.quantity) AS current_quantity,
    COUNT(rc.change_id) AS change_count,
    SUM(CASE WHEN rc.change_type = 'ADD' THEN COALESCE(rc.change_amount, 0) ELSE -COALESCE(rc.change_amount, 0) END) AS net_quantity_change,
    SUM(CASE WHEN rc.change_type = 'ADD' THEN COALESCE(rc.change_area, 0) ELSE -COALESCE(rc.change_area, 0) END) AS net_area_change,
    MAX(rc.changed_at) AS last_change_time
FROM resources r
JOIN regions rg ON r.region_id = rg.region_id
LEFT JOIN resource_changes rc ON r.resource_id = rc.resource_id
GROUP BY
    rg.region_name, r.resource_type, r.species_name
ORDER BY
    rg.region_name, r.resource_type, r.species_name;

-- 查询5: 查询某设备的完整维护记录及状态变化
SELECT
    d.device_name,
    dsl.collected_at AS status_time,
    dsl.run_status,
    dsl.battery_percent,
    dsl.signal_strength,
    mr.maintenance_id,
    mr.maintenance_type,
    mr.maintenance_time,
    u.real_name AS maintainer_name,
    mr.content,
    mr.result
FROM devices d
LEFT JOIN device_status_logs dsl ON d.device_id = dsl.device_id
LEFT JOIN maintenance_records mr ON d.device_id = mr.device_id
LEFT JOIN users u ON mr.maintainer_id = u.user_id
WHERE
    d.device_id = 2001
ORDER BY
    COALESCE(mr.maintenance_time, dsl.collected_at) DESC;
