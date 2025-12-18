-- 智慧林草系统 - 环境监测和统计分析模块 DDL 脚本
-- 基于 environmental_monitoring_mysql.py 生成
-- 创建数据库
CREATE DATABASE IF NOT EXISTS smart_forest_grass CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE smart_forest_grass;

-- 1. 创建表结构

-- 1.1 环境监测站表（monitoring_station）
CREATE TABLE IF NOT EXISTS monitoring_station (
    station_id INT PRIMARY KEY AUTO_INCREMENT COMMENT '站点ID',
    station_name VARCHAR(100) NOT NULL COMMENT '站点名称',
    latitude DECIMAL(10, 6) NOT NULL COMMENT '纬度',
    longitude DECIMAL(10, 6) NOT NULL COMMENT '经度',
    altitude DECIMAL(8, 2) COMMENT '海拔高度',
    station_type VARCHAR(50) NOT NULL COMMENT '站点类型',
    installation_date DATE NOT NULL COMMENT '安装日期',
    status VARCHAR(20) DEFAULT '正常' COMMENT '状态（正常、故障、维护）',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_station_name (station_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='环境监测站表';

-- 1.2 环境数据采集表（environmental_data）
CREATE TABLE IF NOT EXISTS environmental_data (
    data_id INT PRIMARY KEY AUTO_INCREMENT COMMENT '数据ID',
    station_id INT NOT NULL COMMENT '站点ID',
    collection_time DATETIME NOT NULL COMMENT '采集时间',
    temperature DECIMAL(5, 2) COMMENT '温度（℃）',
    humidity DECIMAL(5, 2) COMMENT '湿度（%）',
    wind_speed DECIMAL(5, 2) COMMENT '风速（m/s）',
    wind_direction VARCHAR(10) COMMENT '风向',
    rainfall DECIMAL(8, 2) DEFAULT 0 COMMENT '降雨量（mm）',
    sunshine_duration DECIMAL(5, 2) DEFAULT 0 COMMENT '日照时长（小时）',
    soil_temperature DECIMAL(5, 2) COMMENT '土壤温度（℃）',
    soil_humidity DECIMAL(5, 2) COMMENT '土壤湿度（%）',
    soil_ph DECIMAL(4, 2) COMMENT '土壤pH值',
    pm25 DECIMAL(6, 2) COMMENT 'PM2.5（μg/m³）',
    pm10 DECIMAL(6, 2) COMMENT 'PM10（μg/m³）',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY (station_id) REFERENCES monitoring_station(station_id) ON DELETE CASCADE ON UPDATE CASCADE,
    INDEX idx_station_time (station_id, collection_time),
    INDEX idx_collection_time (collection_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='环境数据采集表';

-- 1.3 统计分析数据表（statistical_data）
CREATE TABLE IF NOT EXISTS statistical_data (
    stat_id INT PRIMARY KEY AUTO_INCREMENT COMMENT '统计ID',
    station_id INT NOT NULL COMMENT '站点ID',
    stat_period VARCHAR(20) NOT NULL COMMENT '统计时间段（日、周、月、季、年）',
    stat_type VARCHAR(50) NOT NULL COMMENT '统计类型',
    stat_date DATE NOT NULL COMMENT '统计日期',
    avg_temperature DECIMAL(5, 2) COMMENT '平均温度（℃）',
    max_temperature DECIMAL(5, 2) COMMENT '最高温度（℃）',
    min_temperature DECIMAL(5, 2) COMMENT '最低温度（℃）',
    avg_humidity DECIMAL(5, 2) COMMENT '平均湿度（%）',
    total_rainfall DECIMAL(8, 2) DEFAULT 0 COMMENT '总降雨量（mm）',
    avg_wind_speed DECIMAL(5, 2) COMMENT '平均风速（m/s）',
    data_source VARCHAR(50) NOT NULL COMMENT '数据来源',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (station_id) REFERENCES monitoring_station(station_id) ON DELETE CASCADE ON UPDATE CASCADE,
    INDEX idx_station_period (station_id, stat_period),
    INDEX idx_stat_date (stat_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='统计分析数据表';

-- 1.4 报表模板表（report_template）
CREATE TABLE IF NOT EXISTS report_template (
    template_id INT PRIMARY KEY AUTO_INCREMENT COMMENT '模板ID',
    template_name VARCHAR(100) NOT NULL COMMENT '模板名称',
    report_type VARCHAR(50) NOT NULL COMMENT '报表类型',
    template_content TEXT NOT NULL COMMENT '模板内容',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_template_name (template_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='报表模板表';

-- 1.5 生成的报表表（generated_report）
CREATE TABLE IF NOT EXISTS generated_report (
    report_id INT PRIMARY KEY AUTO_INCREMENT COMMENT '报表ID',
    template_id INT NOT NULL COMMENT '模板ID',
    generate_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '生成时间',
    report_name VARCHAR(100) NOT NULL COMMENT '报表名称',
    report_content TEXT NOT NULL COMMENT '报表内容',
    generated_by VARCHAR(50) NOT NULL COMMENT '生成人',
    FOREIGN KEY (template_id) REFERENCES report_template(template_id) ON DELETE CASCADE ON UPDATE CASCADE,
    INDEX idx_generate_time (generate_time),
    INDEX idx_generated_by (generated_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='生成的报表表';

-- 1.6 异常数据记录表（abnormal_data）
CREATE TABLE IF NOT EXISTS abnormal_data (
    abnormal_id INT PRIMARY KEY AUTO_INCREMENT COMMENT '异常ID',
    data_id INT NOT NULL COMMENT '数据ID',
    abnormal_type VARCHAR(50) NOT NULL COMMENT '异常类型',
    abnormal_description TEXT NOT NULL COMMENT '异常描述',
    discovery_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发现时间',
    processing_status VARCHAR(20) DEFAULT '未处理' COMMENT '处理状态（未处理、处理中、已处理）',
    processing_time DATETIME COMMENT '处理时间',
    processed_by VARCHAR(50) COMMENT '处理人',
    FOREIGN KEY (data_id) REFERENCES environmental_data(data_id) ON DELETE CASCADE ON UPDATE CASCADE,
    INDEX idx_discovery_time (discovery_time),
    INDEX idx_processing_status (processing_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='异常数据记录表';

-- 1.7 系统用户表（system_user）
CREATE TABLE IF NOT EXISTS system_user (
    user_id INT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    password VARCHAR(100) NOT NULL COMMENT '密码（加密存储）',
    role VARCHAR(20) NOT NULL COMMENT '角色（管理员、普通用户、访客）',
    email VARCHAR(100) COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '电话',
    status VARCHAR(20) DEFAULT '启用' COMMENT '状态（启用、禁用）',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户表';

-- 1.8 操作日志表（operation_log）
CREATE TABLE IF NOT EXISTS operation_log (
    log_id INT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    user_id INT COMMENT '用户ID',
    operation_type VARCHAR(50) NOT NULL COMMENT '操作类型',
    operation_content TEXT NOT NULL COMMENT '操作内容',
    operation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    ip_address VARCHAR(50) COMMENT 'IP地址',
    FOREIGN KEY (user_id) REFERENCES system_user(user_id) ON DELETE SET NULL ON UPDATE CASCADE,
    INDEX idx_operation_time (operation_time),
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';

-- 2. 创建视图

-- 2.1 环境数据概览视图
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

-- 2.2 站点统计数据视图
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

-- 2.3 异常数据视图
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

-- 2.4 管理员环境数据视图
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

-- 2.5 管理员统计数据视图
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

-- 2.6 管理员异常数据视图
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

-- 2.7 普通用户环境数据视图
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

-- 2.8 访客环境数据视图
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

-- 2.9 实时数据视图
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

-- 2.10 环境趋势视图
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

-- 2.11 空气质量视图
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

-- 3. 创建存储过程

-- 3.1 按日统计环境数据存储过程
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

-- 3.2 检测异常数据存储过程
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

-- 3.3 处理异常数据存储过程
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

-- 3.4 生成日报存储过程
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

-- 3.5 获取系统状态存储过程
CREATE PROCEDURE IF NOT EXISTS sp_get_system_status()
BEGIN
    SELECT 
        (SELECT COUNT(*) FROM monitoring_station) AS station_count,
        (SELECT COUNT(*) FROM monitoring_station WHERE status = '正常') AS normal_station_count,
        (SELECT COUNT(*) FROM environmental_data WHERE DATE(collection_time) = CURDATE()) AS today_data_count,
        (SELECT COUNT(*) FROM abnormal_data WHERE processing_status = '未处理') AS unprocessed_abnormal_count;
END;

-- 4. 创建触发器

-- 4.1 环境数据插入后自动检测异常触发器
CREATE TRIGGER IF NOT EXISTS tr_after_insert_environmental_data
AFTER INSERT ON environmental_data
FOR EACH ROW
CALL sp_detect_abnormal_data();

-- 4.2 环境数据更新后自动检测异常触发器
CREATE TRIGGER IF NOT EXISTS tr_after_update_environmental_data
AFTER UPDATE ON environmental_data
FOR EACH ROW
BEGIN
    DELETE FROM abnormal_data WHERE data_id = NEW.data_id;
    CALL sp_detect_abnormal_data();
END;

-- 4.3 环境数据删除后自动删除相关异常记录触发器
CREATE TRIGGER IF NOT EXISTS tr_after_delete_environmental_data
AFTER DELETE ON environmental_data
FOR EACH ROW
DELETE FROM abnormal_data WHERE data_id = OLD.data_id;

-- 5. 初始化数据

-- 5.1 初始化环境监测站
INSERT IGNORE INTO monitoring_station (
    station_name, latitude, longitude, altitude, station_type, installation_date, status
) VALUES 
('监测站1', 39.9042, 116.4074, 50.0, '气象站', '2023-01-01', '正常'),
('监测站2', 39.9142, 116.4174, 60.0, '土壤监测站', '2023-02-01', '正常'),
('监测站3', 39.9242, 116.4274, 70.0, '空气质量监测站', '2023-03-01', '正常'),
('监测站4', 39.9342, 116.4374, 80.0, '综合监测站', '2023-04-01', '维护'),
('监测站5', 39.9442, 116.4474, 90.0, '气象站', '2023-05-01', '正常');

-- 5.2 初始化报表模板
INSERT IGNORE INTO report_template (template_name, report_type, template_content) VALUES 
('每日环境监测报表', '日报', '每日环境监测数据汇总报表模板'),
('每周环境统计报表', '周报', '每周环境统计数据汇总报表模板'),
('每月环境分析报表', '月报', '每月环境分析数据汇总报表模板'),
('季度环境趋势报表', '季报', '季度环境趋势分析报表模板'),
('年度环境总结报表', '年报', '年度环境总结分析报表模板');

-- 5.3 初始化系统用户
INSERT IGNORE INTO system_user (username, password, role, email, phone) VALUES 
('admin', '123456', '管理员', 'admin@example.com', '13800138000'),
('user1', '123456', '普通用户', 'user1@example.com', '13800138001'),
('user2', '123456', '普通用户', 'user2@example.com', '13800138002'),
('guest', '123456', '访客', 'guest@example.com', '13800138003');

-- 脚本结束
SELECT '智慧林草系统 - 环境监测和统计分析模块 DDL 脚本执行完成' AS message;