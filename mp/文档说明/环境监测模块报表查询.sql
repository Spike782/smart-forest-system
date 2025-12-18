-- 环境监测模块报表查询.sql
-- 日报查询
-- 每日环境数据汇总查询
SELECT 
    ms.station_name,
    DATE(ed.collection_time) AS collection_date,
    COUNT(*) AS data_count,
    AVG(ed.temperature) AS avg_temperature,
    MAX(ed.temperature) AS max_temperature,
    MIN(ed.temperature) AS min_temperature,
    AVG(ed.humidity) AS avg_humidity,
    SUM(ed.rainfall) AS total_rainfall,
    AVG(ed.wind_speed) AS avg_wind_speed
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id
WHERE DATE(ed.collection_time) = CURDATE()
GROUP BY ms.station_name, DATE(ed.collection_time);

-- 周报查询
-- 每周环境数据统计查询
SELECT 
    ms.station_name,
    YEAR(ed.collection_time) AS year,
    WEEK(ed.collection_time) AS week,
    AVG(ed.temperature) AS avg_temperature,
    MAX(ed.temperature) AS max_temperature,
    MIN(ed.temperature) AS min_temperature,
    AVG(ed.humidity) AS avg_humidity,
    SUM(ed.rainfall) AS total_rainfall
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id
WHERE YEARWEEK(ed.collection_time) = YEARWEEK(CURDATE())
GROUP BY ms.station_name, YEAR(ed.collection_time), WEEK(ed.collection_time);

-- 月报查询
-- 每月环境数据统计查询
SELECT 
    ms.station_name,
    YEAR(ed.collection_time) AS year,
    MONTH(ed.collection_time) AS month,
    AVG(ed.temperature) AS avg_temperature,
    MAX(ed.temperature) AS max_temperature,
    MIN(ed.temperature) AS min_temperature,
    AVG(ed.humidity) AS avg_humidity,
    SUM(ed.rainfall) AS total_rainfall
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id
WHERE YEAR(ed.collection_time) = YEAR(CURDATE()) 
    AND MONTH(ed.collection_time) = MONTH(CURDATE())
GROUP BY ms.station_name, YEAR(ed.collection_time), MONTH(ed.collection_time);

-- 季报查询
-- 季度环境数据统计查询
SELECT 
    ms.station_name,
    YEAR(ed.collection_time) AS year,
    QUARTER(ed.collection_time) AS quarter,
    AVG(ed.temperature) AS avg_temperature,
    MAX(ed.temperature) AS max_temperature,
    MIN(ed.temperature) AS min_temperature,
    AVG(ed.humidity) AS avg_humidity,
    SUM(ed.rainfall) AS total_rainfall
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id
WHERE YEAR(ed.collection_time) = YEAR(CURDATE()) 
    AND QUARTER(ed.collection_time) = QUARTER(CURDATE())
GROUP BY ms.station_name, YEAR(ed.collection_time), QUARTER(ed.collection_time);

-- 年报查询
-- 年度环境数据统计查询
SELECT 
    ms.station_name,
    YEAR(ed.collection_time) AS year,
    AVG(ed.temperature) AS avg_temperature,
    MAX(ed.temperature) AS max_temperature,
    MIN(ed.temperature) AS min_temperature,
    AVG(ed.humidity) AS avg_humidity,
    SUM(ed.rainfall) AS total_rainfall,
    AVG(ed.pm25) AS avg_pm25,
    AVG(ed.pm10) AS avg_pm10
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id
WHERE YEAR(ed.collection_time) = YEAR(CURDATE())
GROUP BY ms.station_name, YEAR(ed.collection_time);

-- 特殊查询
-- 站点实时数据查询
SELECT 
    ms.station_name,
    ed.collection_time,
    ed.temperature,
    ed.humidity,
    ed.wind_speed,
    ed.wind_direction,
    ed.rainfall,
    ed.pm25,
    ed.pm10
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id
WHERE ed.collection_time = (
    SELECT MAX(collection_time) 
    FROM environmental_data 
    WHERE station_id = ed.station_id
)
ORDER BY ms.station_name;

-- 异常数据查询
-- 最近7天的异常数据查询
SELECT 
    ms.station_name,
    ed.collection_time,
    ad.abnormal_type,
    ad.abnormal_description,
    ad.discovery_time,
    ad.processing_status
FROM abnormal_data ad
JOIN environmental_data ed ON ad.data_id = ed.data_id
JOIN monitoring_station ms ON ed.station_id = ms.station_id
WHERE DATE(ad.discovery_time) >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)
ORDER BY ad.discovery_time DESC;

-- 空气质量查询
-- 最近30天的空气质量查询
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
WHERE DATE(ed.collection_time) >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
GROUP BY ms.station_name, DATE(ed.collection_time)
ORDER BY ms.station_name, collection_date;

-- 环境趋势查询
-- 最近30天的环境趋势查询
SELECT 
    ms.station_name,
    DATE(ed.collection_time) AS collection_date,
    AVG(ed.temperature) AS avg_temperature,
    AVG(ed.humidity) AS avg_humidity,
    SUM(ed.rainfall) AS daily_rainfall
FROM environmental_data ed
JOIN monitoring_station ms ON ed.station_id = ms.station_id
WHERE DATE(ed.collection_time) >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
GROUP BY ms.station_name, DATE(ed.collection_time)
ORDER BY ms.station_name, collection_date;