CREATE DATABASE IF NOT EXISTS smart_forest CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE smart_forest;

-- 环境监测
CREATE TABLE regions (
    region_id       BIGINT AUTO_INCREMENT PRIMARY KEY,
    region_name     VARCHAR(100) NOT NULL,
    region_type     ENUM('FOREST','GRASSLAND') NOT NULL,
    latitude        DECIMAL(9,6),
    longitude       DECIMAL(9,6),
    manager_id      BIGINT NULL, -- 关联用户表
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_region_name (region_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sensors (
    sensor_id       BIGINT AUTO_INCREMENT PRIMARY KEY,
    region_id       BIGINT NOT NULL,
    model           VARCHAR(100),
    monitor_type    ENUM('TEMPERATURE','HUMIDITY','IMAGE','OTHER') NOT NULL,
    install_time    DATETIME,
    protocol        VARCHAR(50),
    status          ENUM('ACTIVE','INACTIVE') DEFAULT 'ACTIVE',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sensor_region FOREIGN KEY (region_id) REFERENCES regions(region_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sensor_readings (
    reading_id      BIGINT AUTO_INCREMENT PRIMARY KEY,
    sensor_id       BIGINT NOT NULL,
    collected_at    DATETIME NOT NULL,
    reading_type    ENUM('TEMPERATURE','HUMIDITY','IMAGE','OTHER') NOT NULL,
    numeric_value   DECIMAL(10,2) NULL,
    value_unit      VARCHAR(20) NULL,
    media_path      VARCHAR(255) NULL, -- 图像或文件路径
    data_status     ENUM('VALID','INVALID') NOT NULL DEFAULT 'VALID',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_reading_sensor FOREIGN KEY (sensor_id) REFERENCES sensors(sensor_id),
    KEY idx_sensor_time (sensor_id, collected_at),
    KEY idx_region_time (collected_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 灾害预警
CREATE TABLE alert_rules (
    rule_id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    alert_type      ENUM('FIRE','DROUGHT','PEST','OTHER') NOT NULL,
    condition_expr  VARCHAR(255) NOT NULL, -- 如 "temp>=38 AND humidity<=20"
    severity_level  ENUM('GENERAL','MODERATE','SEVERE','CRITICAL') NOT NULL,
    is_active       TINYINT(1) NOT NULL DEFAULT 1,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE alerts (
    alert_id        BIGINT AUTO_INCREMENT PRIMARY KEY,
    rule_id         BIGINT NOT NULL,
    region_id       BIGINT NOT NULL,
    triggered_at    DATETIME NOT NULL,
    content         VARCHAR(500),
    status          ENUM('PENDING','PROCESSING','CLOSED') NOT NULL DEFAULT 'PENDING',
    handler_id      BIGINT NULL, -- 处理人
    handle_result   VARCHAR(500),
    CONSTRAINT fk_alert_rule FOREIGN KEY (rule_id) REFERENCES alert_rules(rule_id),
    CONSTRAINT fk_alert_region FOREIGN KEY (region_id) REFERENCES regions(region_id),
    KEY idx_alert_status (status),
    KEY idx_alert_region_time (region_id, triggered_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE notifications (
    notification_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    alert_id        BIGINT NOT NULL,
    receiver_id     BIGINT NOT NULL,
    channel         ENUM('SMS','SYSTEM') NOT NULL,
    sent_at         DATETIME NOT NULL,
    receive_status  ENUM('SENT','DELIVERED','FAILED') NOT NULL DEFAULT 'SENT',
    CONSTRAINT fk_notification_alert FOREIGN KEY (alert_id) REFERENCES alerts(alert_id),
    KEY idx_alert_channel (alert_id, channel)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 资源管理
CREATE TABLE resources (
    resource_id     BIGINT AUTO_INCREMENT PRIMARY KEY,
    resource_type   ENUM('TREE','GRASS') NOT NULL,
    region_id       BIGINT NOT NULL,
    species_name    VARCHAR(100) NOT NULL,
    quantity        INT NULL,             -- 树木按数量
    area            DECIMAL(12,2) NULL,   -- 草地按面积
    growth_stage    ENUM('SEEDLING','GROWING','MATURE') NOT NULL,
    planted_at      DATE,
    operator_id     BIGINT NULL,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_resource_region FOREIGN KEY (region_id) REFERENCES regions(region_id),
    KEY idx_resource_region_type (region_id, resource_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE resource_changes (
    change_id       BIGINT AUTO_INCREMENT PRIMARY KEY,
    resource_id     BIGINT NOT NULL,
    change_type     ENUM('ADD','REDUCE','STATUS_UPDATE') NOT NULL,
    change_reason   VARCHAR(255),
    change_amount   INT NULL,
    change_area     DECIMAL(12,2) NULL,
    changed_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    operator_id     BIGINT NULL,
    CONSTRAINT fk_change_resource FOREIGN KEY (resource_id) REFERENCES resources(resource_id),
    KEY idx_resource_change_time (resource_id, changed_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 设备管理
CREATE TABLE devices (
    device_id       BIGINT PRIMARY KEY,
    device_name     VARCHAR(100) NOT NULL,
    device_type     ENUM('SENSOR','CAMERA','ALERTER','OTHER') NOT NULL,
    model_spec      VARCHAR(100),
    purchased_at    DATE,
    install_region_id BIGINT NULL,
    installer_id    BIGINT NULL,
    warranty_until  DATE,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_device_region FOREIGN KEY (install_region_id) REFERENCES regions(region_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE device_status_logs (
    log_id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    device_id       BIGINT NOT NULL,
    collected_at    DATETIME NOT NULL,
    run_status      ENUM('NORMAL','FAULT','OFFLINE') NOT NULL,
    battery_percent TINYINT NULL,
    signal_strength SMALLINT NULL,
    maintainer_id   BIGINT NULL,
    CONSTRAINT fk_status_device FOREIGN KEY (device_id) REFERENCES devices(device_id),
    KEY idx_device_time (device_id, collected_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE maintenance_records (
    maintenance_id  BIGINT AUTO_INCREMENT PRIMARY KEY,
    device_id       BIGINT NOT NULL,
    maintenance_type ENUM('INSPECTION','REPAIR','REPLACE') NOT NULL,
    maintenance_time DATETIME NOT NULL,
    maintainer_id   BIGINT NULL,
    content         VARCHAR(500),
    result          VARCHAR(200),
    CONSTRAINT fk_maintenance_device FOREIGN KEY (device_id) REFERENCES devices(device_id),
    KEY idx_device_maintenance_time (device_id, maintenance_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 统计分析
CREATE TABLE report_templates (
    template_id     BIGINT AUTO_INCREMENT PRIMARY KEY,
    template_name   VARCHAR(100) NOT NULL,
    dimension       ENUM('REGION','TIME','TYPE') NOT NULL,
    metrics         VARCHAR(255) NOT NULL,
    cycle           ENUM('DAY','WEEK','MONTH') NOT NULL,
    is_active       TINYINT(1) NOT NULL DEFAULT 1,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE reports (
    report_id       BIGINT AUTO_INCREMENT PRIMARY KEY,
    template_id     BIGINT NOT NULL,
    period          VARCHAR(50) NOT NULL, -- 如 2024-10
    generated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    file_path       VARCHAR(255) NOT NULL,
    data_source     VARCHAR(255),
    CONSTRAINT fk_report_template FOREIGN KEY (template_id) REFERENCES report_templates(template_id),
    KEY idx_template_period (template_id, period)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户认证与权限管理
CREATE TABLE users (
    user_id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    username        VARCHAR(50) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    email           VARCHAR(100) NOT NULL UNIQUE,
    real_name       VARCHAR(50) NOT NULL,
    phone           VARCHAR(20),
    status          ENUM('ACTIVE','INACTIVE','LOCKED') DEFAULT 'ACTIVE',
    login_attempts  INT DEFAULT 0,
    last_login_at   DATETIME,
    locked_until    DATETIME NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_username (username),
    UNIQUE KEY uk_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户相关表的外键约束
ALTER TABLE regions ADD CONSTRAINT fk_region_manager FOREIGN KEY (manager_id) REFERENCES users(user_id);
ALTER TABLE resources ADD CONSTRAINT fk_resource_operator FOREIGN KEY (operator_id) REFERENCES users(user_id);
ALTER TABLE device_status_logs ADD CONSTRAINT fk_status_maintainer FOREIGN KEY (maintainer_id) REFERENCES users(user_id);
ALTER TABLE maintenance_records ADD CONSTRAINT fk_maintenance_maintainer FOREIGN KEY (maintainer_id) REFERENCES users(user_id);
ALTER TABLE resource_changes ADD CONSTRAINT fk_resource_change_operator FOREIGN KEY (operator_id) REFERENCES users(user_id);

CREATE TABLE roles (
    role_id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    role_name       VARCHAR(50) NOT NULL UNIQUE,
    description     VARCHAR(255),
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_role_name (role_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE permissions (
    permission_id   BIGINT AUTO_INCREMENT PRIMARY KEY,
    permission_name VARCHAR(50) NOT NULL UNIQUE,
    resource        VARCHAR(50) NOT NULL,
    action          VARCHAR(20) NOT NULL,
    description     VARCHAR(255),
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_permission_name (permission_name),
    UNIQUE KEY uk_resource_action (resource, action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE user_roles (
    user_id         BIGINT NOT NULL,
    role_id         BIGINT NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_role_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_user_role_role FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE role_permissions (
    role_id         BIGINT NOT NULL,
    permission_id   BIGINT NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id),
    CONSTRAINT fk_role_perm_role FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    CONSTRAINT fk_role_perm_perm FOREIGN KEY (permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 环境监测业务线视图

-- View 1: 区域传感器分布视图
CREATE VIEW view_region_sensors AS
SELECT 
    r.region_id,
    r.region_name,
    r.region_type,
    s.sensor_id,
    s.model,
    s.monitor_type,
    s.install_time,
    s.status AS sensor_status
FROM 
    regions r
LEFT JOIN 
    sensors s ON r.region_id = s.region_id;

-- View 2: 区域日均值视图
CREATE VIEW view_daily_env_avg AS
SELECT 
    DATE(sr.collected_at) AS date,
    r.region_name,
    sr.reading_type,
    AVG(sr.numeric_value) AS avg_value,
    sr.value_unit,
    COUNT(*) AS sample_count
FROM 
    sensor_readings sr
JOIN 
    sensors s ON sr.sensor_id = s.sensor_id
JOIN 
    regions r ON s.region_id = r.region_id
WHERE 
    sr.data_status = 'VALID' AND sr.reading_type IN ('TEMPERATURE', 'HUMIDITY')
GROUP BY 
    DATE(sr.collected_at), r.region_name, sr.reading_type, sr.value_unit;

-- View 3: 传感器数据有效性视图
CREATE VIEW view_sensor_validity AS
SELECT 
    r.region_name,
    s.sensor_id,
    s.monitor_type,
    COUNT(*) AS total_readings,
    SUM(CASE WHEN sr.data_status = 'VALID' THEN 1 ELSE 0 END) AS valid_readings,
    SUM(CASE WHEN sr.data_status = 'INVALID' THEN 1 ELSE 0 END) AS invalid_readings,
    MAX(sr.collected_at) AS last_reading_at
FROM 
    sensors s
JOIN 
    sensor_readings sr ON s.sensor_id = sr.sensor_id
JOIN 
    regions r ON s.region_id = r.region_id
GROUP BY 
    r.region_name, s.sensor_id, s.monitor_type;

-- 灾害预警业务线视图

-- View 1: 活跃预警视图
CREATE VIEW view_active_alerts AS
SELECT 
    a.alert_id,
    r.region_name,
    ar.alert_type,
    ar.severity_level,
    a.triggered_at,
    a.content,
    a.status,
    u.real_name AS handler_name
FROM 
    alerts a
JOIN 
    regions r ON a.region_id = r.region_id
JOIN 
    alert_rules ar ON a.rule_id = ar.rule_id
LEFT JOIN 
    users u ON a.handler_id = u.user_id
WHERE 
    a.status IN ('PENDING', 'PROCESSING');

--  View 2: 预警通知汇总视图
CREATE VIEW view_alert_notifications AS
SELECT 
    a.alert_id,
    r.region_name,
    a.triggered_at,
    ar.severity_level,
    a.content,
    a.status,
    COUNT(n.notification_id) AS total_notifications,
    SUM(CASE WHEN n.receive_status = 'DELIVERED' THEN 1 ELSE 0 END) AS delivered_count,
    SUM(CASE WHEN n.receive_status = 'FAILED' THEN 1 ELSE 0 END) AS failed_count
FROM 
    alerts a
JOIN 
    regions r ON a.region_id = r.region_id
JOIN 
    alert_rules ar ON a.rule_id = ar.rule_id
LEFT JOIN 
    notifications n ON a.alert_id = n.alert_id
GROUP BY 
    a.alert_id, r.region_name, a.triggered_at, ar.severity_level, a.content, a.status;

-- View 3: 预警统计视图
CREATE VIEW view_alert_stats AS
SELECT 
    DATE(a.triggered_at) AS date,
    r.region_name,
    ar.alert_type,
    ar.severity_level,
    COUNT(*) AS alert_count,
    SUM(CASE WHEN a.status = 'CLOSED' THEN 1 ELSE 0 END) AS closed_count
FROM 
    alerts a
JOIN 
    regions r ON a.region_id = r.region_id
JOIN 
    alert_rules ar ON a.rule_id = ar.rule_id
GROUP BY 
    DATE(a.triggered_at), r.region_name, ar.alert_type, ar.severity_level;

--资源管理业务线视图

-- View 1: 资源汇总视图
CREATE VIEW view_resource_summary AS
SELECT 
    r.region_name,
    res.resource_type,
    res.species_name,
    res.growth_stage,
    SUM(COALESCE(res.quantity, 0)) AS total_quantity,
    SUM(COALESCE(res.area, 0)) AS total_area
FROM 
    resources res
JOIN 
    regions r ON res.region_id = r.region_id
GROUP BY 
    r.region_name, res.resource_type, res.species_name, res.growth_stage;

-- View 2: 资源变更记录视图
CREATE VIEW view_resource_changes AS
SELECT 
    res.resource_id,
    r.region_name,
    res.resource_type,
    res.species_name,
    rc.change_type,
    rc.change_reason,
    rc.change_amount,
    rc.change_area,
    rc.changed_at,
    u.real_name AS operator_name
FROM 
    resource_changes rc
JOIN 
    resources res ON rc.resource_id = res.resource_id
JOIN 
    regions r ON res.region_id = r.region_id
LEFT JOIN 
    users u ON rc.operator_id = u.user_id;

-- View 3: 区域资源分布视图
CREATE VIEW view_region_resources AS
SELECT 
    r.region_id,
    r.region_name,
    r.region_type,
    COUNT(DISTINCT res.resource_id) AS resource_count,
    SUM(COALESCE(res.quantity, 0)) AS total_quantity,
    SUM(COALESCE(res.area, 0)) AS total_area
FROM 
    regions r
LEFT JOIN 
    resources res ON r.region_id = res.region_id
GROUP BY 
    r.region_id, r.region_name, r.region_type;

-- 设备管理业务线视图

-- View1：设备状态汇总视图
CREATE VIEW view_device_status AS
SELECT 
    r.region_name,
    d.device_id,
    d.device_name,
    d.device_type,
    dsl.run_status,
    dsl.battery_percent,
    dsl.signal_strength,
    dsl.collected_at AS last_status_update
FROM 
    devices d
JOIN 
    (SELECT 
        device_id,
        run_status,
        battery_percent,
        signal_strength,
        collected_at,
        ROW_NUMBER() OVER (PARTITION BY device_id ORDER BY collected_at DESC) AS rn
     FROM device_status_logs) dsl ON d.device_id = dsl.device_id AND dsl.rn = 1
LEFT JOIN 
    regions r ON d.install_region_id = r.region_id;

-- View2：设备维护记录视图
CREATE VIEW view_maintenance_records AS
SELECT 
    d.device_id,
    d.device_name,
    d.device_type,
    r.region_name,
    mr.maintenance_type,
    mr.maintenance_time,
    u.real_name AS maintainer_name,
    mr.content,
    mr.result
FROM 
    maintenance_records mr
JOIN 
    devices d ON mr.device_id = d.device_id
LEFT JOIN 
    regions r ON d.install_region_id = r.region_id
LEFT JOIN 
    users u ON mr.maintainer_id = u.user_id
ORDER BY 
    mr.maintenance_time DESC;

-- View3：设备故障统计视图
CREATE VIEW view_device_faults AS
SELECT 
    r.region_name,
    d.device_id,
    d.device_name,
    d.device_type,
    COUNT(*) AS fault_count,
    MIN(dsl.collected_at) AS first_fault_at,
    MAX(dsl.collected_at) AS last_fault_at
FROM 
    device_status_logs dsl
JOIN 
    devices d ON dsl.device_id = d.device_id
LEFT JOIN 
    regions r ON d.install_region_id = r.region_id
WHERE 
    dsl.run_status IN ('FAULT', 'OFFLINE')
GROUP BY 
    r.region_name, d.device_id, d.device_name, d.device_type;

-- 统计分析业务线视图

-- View1：报表模板视图
CREATE VIEW view_report_templates AS
SELECT 
    rt.template_id,
    rt.template_name,
    rt.dimension,
    rt.metrics,
    rt.cycle,
    rt.is_active,
    rt.created_at,
    COUNT(r.report_id) AS generated_reports_count,
    MAX(r.generated_at) AS last_generated_at
FROM 
    report_templates rt
LEFT JOIN 
    reports r ON rt.template_id = r.template_id
GROUP BY 
    rt.template_id, rt.template_name, rt.dimension, rt.metrics, rt.cycle, rt.is_active, rt.created_at;

-- View2：月度统计视图
CREATE VIEW view_monthly_stats AS
SELECT 
    MONTH(sr.collected_at) AS month,
    YEAR(sr.collected_at) AS year,
    r.region_name,
    sr.reading_type,
    AVG(sr.numeric_value) AS avg_value,
    MAX(sr.numeric_value) AS max_value,
    MIN(sr.numeric_value) AS min_value,
    COUNT(*) AS sample_count
FROM 
    sensor_readings sr
JOIN 
    sensors s ON sr.sensor_id = s.sensor_id
JOIN 
    regions r ON s.region_id = r.region_id
WHERE 
    sr.data_status = 'VALID' AND sr.reading_type IN ('TEMPERATURE', 'HUMIDITY')
GROUP BY 
    MONTH(sr.collected_at), YEAR(sr.collected_at), r.region_name, sr.reading_type;

-- View3：趋势分析视图
CREATE VIEW view_trend_analysis AS
SELECT 
    DATE(sr.collected_at) AS date,
    r.region_name,
    sr.reading_type,
    AVG(sr.numeric_value) AS avg_value,
    LAG(AVG(sr.numeric_value)) OVER (PARTITION BY r.region_name, sr.reading_type ORDER BY DATE(sr.collected_at)) AS prev_avg_value,
    (AVG(sr.numeric_value) - LAG(AVG(sr.numeric_value)) OVER (PARTITION BY r.region_name, sr.reading_type ORDER BY DATE(sr.collected_at))) AS change_value
FROM 
    sensor_readings sr
JOIN 
    sensors s ON sr.sensor_id = s.sensor_id
JOIN 
    regions r ON s.region_id = r.region_id
WHERE 
    sr.data_status = 'VALID' AND sr.reading_type IN ('TEMPERATURE', 'HUMIDITY')
GROUP BY 
    DATE(sr.collected_at), r.region_name, sr.reading_type;

-- 存储过程和触发器
-- 环境监测业务线：传感器数据有效性触发器

DELIMITER //
CREATE TRIGGER trg_sensor_reading_validity BEFORE INSERT ON sensor_readings FOR EACH ROW
BEGIN
    -- Temperature validation: valid between -50 to 50 degrees
    IF NEW.reading_type = 'TEMPERATURE' THEN
        IF NEW.numeric_value < -50 OR NEW.numeric_value > 50 THEN
            SET NEW.data_status = 'INVALID';
        ELSE
            SET NEW.data_status = 'VALID';
        END IF;
    END IF;
    
    -- Humidity validation: valid between 0 to 100%
    IF NEW.reading_type = 'HUMIDITY' THEN
        IF NEW.numeric_value < 0 OR NEW.numeric_value > 100 THEN
            SET NEW.data_status = 'INVALID';
        ELSE
            SET NEW.data_status = 'VALID';
        END IF;
    END IF;
    
    -- 其他类型默认有效
    IF NEW.reading_type NOT IN ('TEMPERATURE', 'HUMIDITY') THEN
        SET NEW.data_status = 'VALID';
    END IF;
END//
DELIMITER ;

-- 灾害预警业务线：预警触发存储过程

DELIMITER //
CREATE PROCEDURE sp_trigger_alerts()
BEGIN
    -- 声明变量
    DECLARE done INT DEFAULT FALSE;
    DECLARE v_sensor_id BIGINT;
    DECLARE v_region_id BIGINT;
    DECLARE v_reading_type VARCHAR(50);
    DECLARE v_numeric_value DECIMAL(10,2);
    DECLARE v_rule_id BIGINT;
    DECLARE v_alert_type VARCHAR(50);
    DECLARE v_severity_level VARCHAR(50);
    DECLARE v_condition_expr VARCHAR(255);
    DECLARE v_content VARCHAR(500);

    -- 声明游标
    DECLARE sensor_reading_cursor CURSOR FOR
        SELECT sr.sensor_id, s.region_id, sr.reading_type, sr.numeric_value
        FROM sensor_readings sr
        JOIN sensors s ON sr.sensor_id = s.sensor_id
        WHERE sr.data_status = 'VALID'
        AND sr.collected_at > NOW() - INTERVAL 5 MINUTE;

    -- 声明异常处理器
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
    
    -- 打开游标
    OPEN sensor_reading_cursor;
    
    -- 循环处理每个传感器读数
    read_loop: LOOP
        FETCH sensor_reading_cursor INTO v_sensor_id, v_region_id, v_reading_type, v_numeric_value;
        
        IF done THEN
            LEAVE read_loop;
        END IF;
        
        -- 检查是否需要触发预警
        IF (v_reading_type = 'TEMPERATURE' AND v_numeric_value > 35) OR
           (v_reading_type = 'HUMIDITY' AND v_numeric_value < 20) THEN

            -- 查找匹配的预警规则
            SELECT rule_id INTO v_rule_id
            FROM alert_rules
            WHERE is_active = 1
            AND alert_type = CASE
                WHEN v_reading_type = 'TEMPERATURE' AND v_numeric_value > 35 THEN 'FIRE'
                WHEN v_reading_type = 'HUMIDITY' AND v_numeric_value < 20 THEN 'DROUGHT'
                ELSE 'OTHER'
            END；
            LIMIT 1;

            IF v_rule_id IS NOT NULL THEN
                -- Generate alert content
                SET v_content = CONCAT('Region ', v_region_id, ' sensor ', v_sensor_id, ' triggered alert: ', v_numeric_value, ' matches rule ', v_rule_id);

                -- Insert alert record
                INSERT INTO alerts (rule_id, region_id, triggered_at, content, status)
                VALUES (v_rule_id, v_region_id, NOW(), v_content, 'PENDING');
            END IF;
        END IF;
        
    END LOOP read_loop;
    
    -- 关闭游标
    CLOSE sensor_reading_cursor;
    
END//
DELIMITER ;

--资源管理业务线：资源变更记录触发器

DELIMITER //
CREATE TRIGGER trg_resource_change_log AFTER UPDATE ON resources FOR EACH ROW
BEGIN
    -- Record quantity changes
    IF OLD.quantity != NEW.quantity THEN
        INSERT INTO resource_changes (resource_id, change_type, change_amount, change_reason, changed_at)
        VALUES (
            NEW.resource_id,
            CASE WHEN NEW.quantity > OLD.quantity THEN 'ADD' ELSE 'REDUCE' END,
            ABS(NEW.quantity - OLD.quantity),
            'Auto record: quantity change',
            NOW()
        );
    END IF;

    -- Record area changes
    IF OLD.area != NEW.area THEN
        INSERT INTO resource_changes (resource_id, change_type, change_area, change_reason, changed_at)
        VALUES (
            NEW.resource_id,
            CASE WHEN NEW.area > OLD.area THEN 'ADD' ELSE 'REDUCE' END,
            ABS(NEW.area - OLD.area),
            'Auto record: area change',
            NOW()
        );
    END IF;

    -- Record growth stage changes
    IF OLD.growth_stage != NEW.growth_stage THEN
        INSERT INTO resource_changes (resource_id, change_type, change_reason, changed_at)
        VALUES (
            NEW.resource_id,
            'STATUS_UPDATE',
            CONCAT('Auto record: growth stage changed from ', OLD.growth_stage, ' to ', NEW.growth_stage),
            NOW()
        );
    END IF;
END//
DELIMITER ;

--设备管理业务线：设备维护提醒存储过程

DELIMITER //
CREATE PROCEDURE sp_device_maintenance_reminder()
BEGIN
    -- 查找需要维护的设备
    SELECT 
        d.device_id,
        d.device_name,
        d.device_type,
        r.region_name,
        mr.maintenance_time AS last_maintenance_time,
        DATEDIFF(NOW(), mr.maintenance_time) AS days_since_last_maintenance
    INTO @device_id, @device_name, @device_type, @region_name, @last_maintenance_time, @days_since_last_maintenance
    FROM devices d
    LEFT JOIN (
        SELECT device_id, MAX(maintenance_time) AS maintenance_time
        FROM maintenance_records
        GROUP BY device_id
    ) mr ON d.device_id = mr.device_id
    LEFT JOIN regions r ON d.install_region_id = r.region_id
    WHERE 
        (mr.maintenance_time IS NULL OR DATEDIFF(NOW(), mr.maintenance_time) > 30)
    LIMIT 1;
    
    -- 如果找到需要维护的设备，插入维护提醒
    IF @device_id IS NOT NULL THEN
        INSERT INTO maintenance_records (device_id, maintenance_type, maintenance_time, content, result)
        VALUES (
            @device_id,
            'INSPECTION',
            NOW(),
            CONCAT(@device_name, ' 设备已超过 ', @days_since_last_maintenance, ' 天未维护，需要进行巡检'),
            'PENDING'
        );
    END IF;
    
END//
DELIMITER ;

-- 统计分析业务线：日报表生成存储过程
DELIMITER //
CREATE PROCEDURE sp_generate_daily_report(IN p_date DATE)
BEGIN
    DECLARE v_period VARCHAR(50);
    DECLARE v_file_path VARCHAR(255);
    DECLARE v_report_id BIGINT;
    
    -- Set default date to yesterday
    IF p_date IS NULL THEN
        SET p_date = DATE(NOW()) - INTERVAL 1 DAY;
    END IF;

    -- Set period format
    SET v_period = DATE_FORMAT(p_date, '%Y-%m-%d');

    -- Generate report file path
    SET v_file_path = CONCAT('/reports/daily_', v_period, '.pdf');

    -- Insert report record
    INSERT INTO reports (template_id, period, generated_at, file_path, data_source)
    VALUES (
        (SELECT template_id FROM report_templates WHERE cycle = 'DAY' LIMIT 1),
        v_period,
        NOW(),
        v_file_path,
        'Auto generated: Daily environmental monitoring data summary'
    );

    -- Get the newly inserted report ID
    SET v_report_id = LAST_INSERT_ID();

    -- Return the generated report ID
    SELECT v_report_id AS report_id, v_period AS period, v_file_path AS file_path;
    
END//
DELIMITER ;

