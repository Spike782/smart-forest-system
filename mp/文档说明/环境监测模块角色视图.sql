-- 环境监测模块角色视图.sql
-- 创建系统视图
USE smart_forest_grass;

-- 1. 环境数据概览视图
CREATE VIEW IF NOT EXISTS v_environmental_data_overview AS
SELECT 
    ed.data_id,
    ms.station_name,
    ed.collection_time,
    ed.temperature,
    ed.humidity,
    ed.wind_speed,
    ed.wind_direction,
    ed.rainfall,
    ed.sunshine_duration,
    ed.soil_temperature,
    ed.soil_humidity,
    ed.soil_ph,
    ed.pm25,
    ed.pm10
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id;

-- 2. 站点统计数据视图
CREATE VIEW IF NOT EXISTS v_station_statistics AS
SELECT 
    ms.station_id,
    ms.station_name,
    sd.stat_period,
    sd.stat_type,
    sd.stat_date,
    sd.avg_temperature,
    sd.max_temperature,
    sd.min_temperature,
    sd.avg_humidity,
    sd.total_rainfall,
    sd.avg_wind_speed
FROM statistical_data sd
JOIN monitoring_station ms ON sd.station_id = ms.station_id;

-- 3. 异常数据视图
CREATE VIEW IF NOT EXISTS v_abnormal_data AS
SELECT 
    ad.abnormal_id,
    ms.station_name,
    ed.collection_time,
    ad.abnormal_type,
    ad.abnormal_description,
    ad.discovery_time,
    ad.processing_status,
    ad.processing_time,
    ad.processed_by
FROM abnormal_data ad
JOIN environmental_data ed ON ad.data_id = ed.data_id
JOIN monitoring_station ms ON ed.station_id = ms.station_id;

-- 4. 管理员视图
-- 管理员环境数据视图
CREATE VIEW IF NOT EXISTS v_admin_environmental_data AS
SELECT 
    ed.data_id,
    ms.station_id,
    ms.station_name,
    ms.station_type,
    ed.collection_time,
    ed.temperature,
    ed.humidity,
    ed.wind_speed,
    ed.wind_direction,
    ed.rainfall,
    ed.sunshine_duration,
    ed.soil_temperature,
    ed.soil_humidity,
    ed.soil_ph,
    ed.pm25,
    ed.pm10,
    ed.create_time
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id;

-- 管理员统计数据视图
CREATE VIEW IF NOT EXISTS v_admin_statistical_data AS
SELECT 
    sd.stat_id,
    ms.station_id,
    ms.station_name,
    sd.stat_period,
    sd.stat_type,
    sd.stat_date,
    sd.avg_temperature,
    sd.max_temperature,
    sd.min_temperature,
    sd.avg_humidity,
    sd.total_rainfall,
    sd.avg_wind_speed,
    sd.data_source,
    sd.create_time,
    sd.update_time
FROM statistical_data sd
JOIN monitoring_station ms ON sd.station_id = ms.station_id;

-- 管理员异常数据视图
CREATE VIEW IF NOT EXISTS v_admin_abnormal_data AS
SELECT 
    ad.abnormal_id,
    ms.station_id,
    ms.station_name,
    ed.collection_time,
    ad.abnormal_type,
    ad.abnormal_description,
    ad.discovery_time,
    ad.processing_status,
    ad.processing_time,
    ad.processed_by
FROM abnormal_data ad
JOIN environmental_data ed ON ad.data_id = ed.data_id
JOIN monitoring_station ms ON ed.station_id = ms.station_id;

-- 5. 普通用户视图
-- 普通用户环境数据视图
CREATE VIEW IF NOT EXISTS v_user_environmental_data AS
SELECT 
    ms.station_name,
    DATE(ed.collection_time) AS collection_date,
    TIME(ed.collection_time) AS collection_time,
    ed.temperature,
    ed.humidity,
    ed.wind_speed,
    ed.wind_direction,
    ed.rainfall,
    ed.sunshine_duration,
    ed.soil_temperature,
    ed.soil_humidity
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id
WHERE ed.collection_time >= DATE_SUB(CURDATE(), INTERVAL 30 DAY);

-- 6. 访客视图
-- 访客环境数据视图
CREATE VIEW IF NOT EXISTS v_guest_environmental_data AS
SELECT 
    ms.station_name,
    DATE(ed.collection_time) AS collection_date,
    AVG(ed.temperature) AS avg_temperature,
    AVG(ed.humidity) AS avg_humidity,
    AVG(ed.wind_speed) AS avg_wind_speed,
    SUM(ed.rainfall) AS total_rainfall
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id
WHERE ed.collection_time >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)
GROUP BY ms.station_name, DATE(ed.collection_time);

-- 7. 功能视图
-- 实时数据视图
CREATE VIEW IF NOT EXISTS v_realtime_data AS
SELECT 
    ms.station_id,
    ms.station_name,
    ms.station_type,
    ed.collection_time,
    ed.temperature,
    ed.humidity,
    ed.wind_speed,
    ed.wind_direction,
    ed.rainfall,
    ed.soil_temperature,
    ed.soil_humidity
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id
WHERE ed.collection_time = (
    SELECT MAX(collection_time) 
    FROM environmental_data 
    WHERE station_id = ed.station_id
)
AND ms.status = '正常';

-- 环境趋势视图
CREATE VIEW IF NOT EXISTS v_environmental_trend AS
SELECT 
    ms.station_name,
    DATE(ed.collection_time) AS collection_date,
    AVG(ed.temperature) AS avg_temperature,
    AVG(ed.humidity) AS avg_humidity,
    SUM(ed.rainfall) AS daily_rainfall,
    AVG(ed.pm25) AS avg_pm25
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id
GROUP BY ms.station_name, DATE(ed.collection_time)
ORDER BY ms.station_name, collection_date;

-- 空气质量视图
CREATE VIEW IF NOT EXISTS v_air_quality AS
SELECT 
    ms.station_name,
    DATE(ed.collection_time) AS collection_date,
    AVG(ed.pm25) AS avg_pm25,
    AVG(ed.pm10) AS avg_pm10,
    CASE 
        WHEN AVG(ed.pm25) <= 35 THEN '优'
        WHEN AVG(ed.pm25) <= 75 THEN '良'
        WHEN AVG(ed.pm25) <= 115 THEN '轻度污染'
        WHEN AVG(ed.pm25) <= 150 THEN '中度污染'
        WHEN AVG(ed.pm25) <= 250 THEN '重度污染'
        ELSE '严重污染'
    END AS air_quality_level
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id
GROUP BY ms.station_name, DATE(ed.collection_time)
ORDER BY ms.station_name, collection_date;