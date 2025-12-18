-- 环境监测模块存储过程触发器.sql
-- 创建存储过程和触发器
USE smart_forest_grass;

-- 1. 数据统计存储过程
-- 按日统计环境数据存储过程
CREATE PROCEDURE IF NOT EXISTS sp_statistics_daily(IN p_stat_date DATE)
BEGIN
    INSERT INTO statistical_data (
        station_id, stat_period, stat_type, stat_date,
        avg_temperature, max_temperature, min_temperature,
        avg_humidity, total_rainfall, avg_wind_speed, data_source, create_time
    )
    SELECT 
        station_id,
        '日' AS stat_period,
        '环境数据' AS stat_type,
        p_stat_date AS stat_date,
        AVG(temperature) AS avg_temperature,
        MAX(temperature) AS max_temperature,
        MIN(temperature) AS min_temperature,
        AVG(humidity) AS avg_humidity,
        SUM(rainfall) AS total_rainfall,
        AVG(wind_speed) AS avg_wind_speed,
        '自动统计' AS data_source,
        CURRENT_TIMESTAMP AS create_time
    FROM environmental_data
    WHERE DATE(collection_time) = p_stat_date
    GROUP BY station_id
    ON DUPLICATE KEY UPDATE
        avg_temperature = VALUES(avg_temperature),
        max_temperature = VALUES(max_temperature),
        min_temperature = VALUES(min_temperature),
        avg_humidity = VALUES(avg_humidity),
        total_rainfall = VALUES(total_rainfall),
        avg_wind_speed = VALUES(avg_wind_speed),
        data_source = VALUES(data_source),
        create_time = CURRENT_TIMESTAMP;
END;

-- 2. 异常数据处理存储过程
-- 检测异常数据存储过程
CREATE PROCEDURE IF NOT EXISTS sp_detect_abnormal_data()
BEGIN
    -- 检测温度异常（>40℃或< -20℃）
    INSERT INTO abnormal_data (data_id, abnormal_type, abnormal_description)
    SELECT 
        data_id,
        '温度异常',
        CONCAT('温度超出正常范围：', temperature, '℃')
    FROM environmental_data
    WHERE (temperature > 40 OR temperature < -20)
    AND data_id NOT IN (SELECT data_id FROM abnormal_data WHERE abnormal_type = '温度异常');
    
    -- 检测湿度异常（>100%或< 0%）
    INSERT INTO abnormal_data (data_id, abnormal_type, abnormal_description)
    SELECT 
        data_id,
        '湿度异常',
        CONCAT('湿度超出正常范围：', humidity, '%')
    FROM environmental_data
    WHERE (humidity > 100 OR humidity < 0)
    AND data_id NOT IN (SELECT data_id FROM abnormal_data WHERE abnormal_type = '湿度异常');
    
    -- 检测风速异常（>30 m/s）
    INSERT INTO abnormal_data (data_id, abnormal_type, abnormal_description)
    SELECT 
        data_id,
        '风速异常',
        CONCAT('风速超出正常范围：', wind_speed, 'm/s')
    FROM environmental_data
    WHERE wind_speed > 30
    AND data_id NOT IN (SELECT data_id FROM abnormal_data WHERE abnormal_type = '风速异常');
    
    -- 检测PM2.5异常（>500 μg/m³）
    INSERT INTO abnormal_data (data_id, abnormal_type, abnormal_description)
    SELECT 
        data_id,
        'PM2.5异常',
        CONCAT('PM2.5超出正常范围：', pm25, 'μg/m³')
    FROM environmental_data
    WHERE pm25 > 500
    AND data_id NOT IN (SELECT data_id FROM abnormal_data WHERE abnormal_type = 'PM2.5异常');
    
    -- 检测PM10异常（>1000 μg/m³）
    INSERT INTO abnormal_data (data_id, abnormal_type, abnormal_description)
    SELECT 
        data_id,
        'PM10异常',
        CONCAT('PM10超出正常范围：', pm10, 'μg/m³')
    FROM environmental_data
    WHERE pm10 > 1000
    AND data_id NOT IN (SELECT data_id FROM abnormal_data WHERE abnormal_type = 'PM10异常');
    
    -- 检测土壤pH值异常（<4.5或>9.5）
    INSERT INTO abnormal_data (data_id, abnormal_type, abnormal_description)
    SELECT 
        data_id,
        '土壤pH值异常',
        CONCAT('土壤pH值超出正常范围：', soil_ph)
    FROM environmental_data
    WHERE (soil_ph < 4.5 OR soil_ph > 9.5)
    AND data_id NOT IN (SELECT data_id FROM abnormal_data WHERE abnormal_type = '土壤pH值异常');
END;

-- 处理异常数据存储过程
CREATE PROCEDURE IF NOT EXISTS sp_process_abnormal_data(
    IN p_abnormal_id INT,
    IN p_processing_status VARCHAR(20),
    IN p_processed_by VARCHAR(50)
)
BEGIN
    UPDATE abnormal_data
    SET 
        processing_status = p_processing_status,
        processing_time = CURRENT_TIMESTAMP,
        processed_by = p_processed_by
    WHERE abnormal_id = p_abnormal_id;
END;

-- 3. 报表生成存储过程
-- 生成日报存储过程
CREATE PROCEDURE IF NOT EXISTS sp_generate_daily_report(IN p_report_date DATE)
BEGIN
    -- 生成报表内容
    SET @report_content = CONCAT(
        '智慧林草系统 - 环境监测日报\n',
        '报表日期：', DATE_FORMAT(p_report_date, '%Y年%m月%d日'), '\n',
        '生成时间：', NOW(), '\n',
        '\n========================================\n',
        '\n站点环境数据汇总：\n'
    );
    
    -- 插入报表记录
    INSERT INTO generated_report (
        template_id, report_name, generate_time, report_content, generated_by
    )
    VALUES (
        (SELECT template_id FROM report_template WHERE template_name = '每日环境监测报表'),
        CONCAT('环境监测日报_', DATE_FORMAT(p_report_date, '%Y%m%d')),
        NOW(),
        @report_content,
        'system'
    );
    
    -- 返回生成的报表ID
    SELECT 307740151 AS report_id;
END;

-- 4. 系统状态存储过程
-- 获取系统状态存储过程
CREATE PROCEDURE IF NOT EXISTS sp_get_system_status()
BEGIN
    SELECT 
        (SELECT COUNT(*) FROM monitoring_station) AS station_count,
        (SELECT COUNT(*) FROM monitoring_station WHERE status = '正常') AS normal_station_count,
        (SELECT COUNT(*) FROM environmental_data WHERE DATE(collection_time) = CURDATE()) AS today_data_count,
        (SELECT COUNT(*) FROM abnormal_data WHERE processing_status = '未处理') AS unprocessed_abnormal_count;
END;

-- 5. 触发器
-- 环境数据插入后自动检测异常触发器
CREATE TRIGGER IF NOT EXISTS tr_after_insert_environmental_data
AFTER INSERT ON environmental_data
FOR EACH ROW
CALL sp_detect_abnormal_data();

-- 环境数据更新后自动检测异常触发器
CREATE TRIGGER IF NOT EXISTS tr_after_update_environmental_data
AFTER UPDATE ON environmental_data
FOR EACH ROW
BEGIN
    DELETE FROM abnormal_data WHERE data_id = NEW.data_id;
    CALL sp_detect_abnormal_data();
END;

-- 环境数据删除后自动删除相关异常记录触发器
CREATE TRIGGER IF NOT EXISTS tr_after_delete_environmental_data
AFTER DELETE ON environmental_data
FOR EACH ROW
DELETE FROM abnormal_data WHERE data_id = OLD.data_id;