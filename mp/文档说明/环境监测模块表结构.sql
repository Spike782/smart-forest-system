-- 环境监测模块表结构.sql
-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS smart_forest_grass CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE smart_forest_grass;

-- 1. 环境监测站表
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

-- 2. 环境数据采集表
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

-- 3. 统计分析数据表
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

-- 4. 报表模板表
CREATE TABLE IF NOT EXISTS report_template (
    template_id INT PRIMARY KEY AUTO_INCREMENT COMMENT '模板ID',
    template_name VARCHAR(100) NOT NULL COMMENT '模板名称',
    report_type VARCHAR(50) NOT NULL COMMENT '报表类型',
    template_content TEXT NOT NULL COMMENT '模板内容',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_template_name (template_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='报表模板表';

-- 5. 生成的报表表
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

-- 6. 异常数据记录表
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

-- 7. 系统用户表（用于权限控制）
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

-- 8. 操作日志表
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