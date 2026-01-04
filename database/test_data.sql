-- 测试数据生成脚本
-- 注意：执行此脚本前，请确保数据库已创建并初始化表结构

-- 设置SQL模式
SET NAMES utf8mb4; 
SET FOREIGN_KEY_CHECKS = 0; -- 临时关闭外键校验，避免执行顺序问题
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";

-- 插入测试用户数据
INSERT INTO users (username, password_hash, email, real_name, phone, status) VALUES
('admin', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'admin@example.com', '系统管理员', '13800138000', 'ACTIVE'),
('data_admin', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'data@example.com', '数据管理员', '13800138001', 'ACTIVE'),
('forest_guard_1', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'guard1@example.com', '护林员1', '13800138002', 'ACTIVE'),
('forest_guard_2', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'guard2@example.com', '护林员2', '13800138003', 'ACTIVE'),
('public_user_1', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'public1@example.com', '公众用户1', '13800138004', 'ACTIVE'),
('public_user_2', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'public2@example.com', '公众用户2', '13800138005', 'ACTIVE'),
('inspector_1', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'inspector1@example.com', '监管人员1', '13800138006', 'ACTIVE'),
('inspector_2', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'inspector2@example.com', '监管人员2', '13800138007', 'ACTIVE'),
('maintainer_1', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'maintainer1@example.com', '维护人员1', '13800138008', 'ACTIVE'),
('maintainer_2', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'maintainer2@example.com', '维护人员2', '13800138009', 'ACTIVE'),
('forest_guard_3', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'guard3@example.com', '护林员3', '13800138010', 'ACTIVE'),
('forest_guard_4', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'guard4@example.com', '护林员4', '13800138011', 'ACTIVE'),
('forest_guard_5', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'guard5@example.com', '护林员5', '13800138012', 'ACTIVE'),
('forest_guard_6', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'guard6@example.com', '护林员6', '13800138013', 'ACTIVE'),
('forest_guard_7', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'guard7@example.com', '护林员7', '13800138014', 'ACTIVE'),
('forest_guard_8', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'guard8@example.com', '护林员8', '13800138015', 'ACTIVE'),
('forest_guard_9', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'guard9@example.com', '护林员9', '13800138016', 'ACTIVE'),
('forest_guard_10', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'guard10@example.com', '护林员10', '13800138017', 'ACTIVE'),
('data_admin_2', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'data2@example.com', '数据管理员2', '13800138018', 'ACTIVE'),
('system_admin_2', '$2a$10$G9a7Dc9Z8F7E6A5B4C3D2E1F0', 'system2@example.com', '系统管理员2', '13800138019', 'ACTIVE');

-- 插入系统角色数据
INSERT INTO roles (role_name, description) VALUES
('SYSTEM_ADMIN', '系统管理员：维护所有角色的账号信息与权限分配，管理设备档案与预警规则，审核统计报表模板'),
('DATA_ADMIN', '数据管理员：录入、校验林草资源基础信息，处理监测数据的异常情况，生成并发布统计报表'),
('FOREST_GUARD', '区域护林员：查看负责区域的实时监测数据与历史记录，接收并处理相关灾害预警，记录设备巡检与维护情况，更新区域内林草资源变动信息'),
('PUBLIC_USER', '公众用户：查看公开的林草资源统计数据，浏览非涉密的环境监测信息（如景区空气质量），提交林草资源异常情况反馈（如发现火情、病虫害）'),
('INSPECTOR', '监管人员：查看全系统业务数据与操作记录，监督预警处理流程的及时性与规范性，审核资源变动与设备维护记录的真实性');

-- 插入测试权限数据
INSERT INTO permissions (permission_name, resource, action, description) VALUES
('view_regions', 'regions', 'view', '查看区域信息'),
('manage_regions', 'regions', 'manage', '管理区域信息'),
('view_sensors', 'sensors', 'view', '查看传感器信息'),
('manage_sensors', 'sensors', 'manage', '管理传感器信息'),
('view_alerts', 'alerts', 'view', '查看预警信息'),
('manage_alerts', 'alerts', 'manage', '管理预警信息'),
('view_resources', 'resources', 'view', '查看资源信息'),
('manage_resources', 'resources', 'manage', '管理资源信息'),
('view_devices', 'devices', 'view', '查看设备信息'),
('manage_devices', 'devices', 'manage', '管理设备信息'),
('view_reports', 'reports', 'view', '查看报表信息'),
('manage_reports', 'reports', 'manage', '管理报表信息'),
('view_statistics', 'statistics', 'view', '查看统计信息'),
('manage_statistics', 'statistics', 'manage', '管理统计信息'),
('view_users', 'users', 'view', '查看用户信息'),
('manage_users', 'users', 'manage', '管理用户信息'),
('view_roles', 'roles', 'view', '查看角色信息'),
('manage_roles', 'roles', 'manage', '管理角色信息'),
('view_permissions', 'permissions', 'view', '查看权限信息'),
('manage_permissions', 'permissions', 'manage', '管理权限信息');

-- 插入用户角色关联数据
INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1), -- admin 是系统管理员
(2, 2), -- data_admin 是数据管理员
(3, 3), -- forest_guard_1 是区域护林员
(4, 3), -- forest_guard_2 是区域护林员
(5, 4), -- public_user_1 是公众用户
(6, 4), -- public_user_2 是公众用户
(7, 5), -- inspector_1 是监管人员
(8, 5), -- inspector_2 是监管人员
(9, 3), -- maintainer_1 是区域护林员（原注释错误，已修正）
(10, 3), -- maintainer_2 是区域护林员（原注释错误，已修正）
(11, 3), -- forest_guard_3 是区域护林员
(12, 3), -- forest_guard_4 是区域护林员
(13, 3), -- forest_guard_5 是区域护林员
(14, 3), -- forest_guard_6 是区域护林员
(15, 3), -- forest_guard_7 是区域护林员
(16, 3), -- forest_guard_8 是区域护林员
(17, 3), -- forest_guard_9 是区域护林员
(18, 3), -- forest_guard_10 是区域护林员
(19, 2), -- data_admin_2 是数据管理员
(20, 1); -- system_admin_2 是系统管理员

-- 插入角色权限关联数据
INSERT INTO role_permissions (role_id, permission_id) VALUES
-- 系统管理员：拥有所有权限
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9), (1, 10), (1, 11), (1, 12), (1, 13), (1, 14), (1, 15), (1, 16), (1, 17), (1, 18), (1, 19), (1, 20),
-- 数据管理员：录入、校验林草资源基础信息，处理监测数据的异常情况，生成并发布统计报表
(2, 1), (2, 3), (2, 5), (2, 7), (2, 8), (2, 9), (2, 11), (2, 12), (2, 13), (2, 14), (2, 15), (2, 16), (2, 17), (2, 18), (2, 19),
-- 区域护林员：查看负责区域的实时监测数据与历史记录，接收并处理相关灾害预警，记录设备巡检与维护情况，更新区域内林草资源变动信息
(3, 1), (3, 3), (3, 5), (3, 7), (3, 8), (3, 9), (3, 11), (3, 13), (3, 15), (3, 17),
-- 公众用户：查看公开的林草资源统计数据，浏览非涉密的环境监测信息，提交林草资源异常情况反馈
(4, 1), (4, 3), (4, 7), (4, 11), (4, 13), (4, 15),
-- 监管人员：查看全系统业务数据与操作记录，监督预警处理流程的及时性与规范性，审核资源变动与设备维护记录的真实性
(5, 1), (5, 3), (5, 5), (5, 7), (5, 9), (5, 11), (5, 13), (5, 15), (5, 17), (5, 19);

-- 插入测试区域数据
INSERT INTO regions (region_name, region_type, latitude, longitude, manager_id) VALUES
('常绿森林A区', 'FOREST', 39.9042, 116.4074, 3),
('常绿森林B区', 'FOREST', 39.9142, 116.4174, 4),
('草地保护区1', 'GRASSLAND', 39.9242, 116.4274, 3),
('草地保护区2', 'GRASSLAND', 39.9342, 116.4374, 4),
('落叶森林C区', 'FOREST', 39.9442, 116.4474, 3),
('落叶森林D区', 'FOREST', 39.9542, 116.4574, 4),
('湿地保护区', 'GRASSLAND', 39.9642, 116.4674, 3),
('山地森林E区', 'FOREST', 39.9742, 116.4774, 4),
('山地森林F区', 'FOREST', 39.9842, 116.4874, 3),
('城市绿地1', 'GRASSLAND', 39.9942, 116.4974, 4),
('城市绿地2', 'GRASSLAND', 40.0042, 116.5074, 3),
('自然保护区1', 'FOREST', 40.0142, 116.5174, 4),
('自然保护区2', 'GRASSLAND', 40.0242, 116.5274, 3),
('生态公园1', 'FOREST', 40.0342, 116.5374, 4),
('生态公园2', 'GRASSLAND', 40.0442, 116.5474, 3),
('水源保护区', 'FOREST', 40.0542, 116.5574, 4),
('防风林带1', 'FOREST', 40.0642, 116.5674, 3),
('防风林带2', 'FOREST', 40.0742, 116.5774, 4),
('水土保持区1', 'GRASSLAND', 40.0842, 116.5874, 3),
('水土保持区2', 'GRASSLAND', 40.0942, 116.5974, 4);

-- 插入测试传感器数据
INSERT INTO sensors (sensor_id, region_id, model, monitor_type, install_time, protocol, status) VALUES
(1001, 1, '型号A', 'TEMPERATURE', '2023-01-01', 'MQTT', 'ACTIVE'),
(1002, 1, '型号B', 'HUMIDITY', '2023-01-02', 'MQTT', 'ACTIVE'),
(1003, 2, '型号A', 'TEMPERATURE', '2023-01-03', 'MQTT', 'ACTIVE'),
(1004, 2, '型号B', 'HUMIDITY', '2023-01-04', 'MQTT', 'ACTIVE'),
(1005, 3, '型号A', 'TEMPERATURE', '2023-01-05', 'MQTT', 'ACTIVE'),
(1006, 3, '型号B', 'HUMIDITY', '2023-01-06', 'MQTT', 'ACTIVE'),
(1007, 4, '型号A', 'TEMPERATURE', '2023-01-07', 'MQTT', 'ACTIVE'),
(1008, 4, '型号B', 'HUMIDITY', '2023-01-08', 'MQTT', 'ACTIVE'),
(1009, 5, '型号A', 'TEMPERATURE', '2023-01-09', 'MQTT', 'ACTIVE'),
(1010, 5, '型号B', 'HUMIDITY', '2023-01-10', 'MQTT', 'ACTIVE'),
(1011, 6, '型号A', 'TEMPERATURE', '2023-01-11', 'MQTT', 'ACTIVE'),
(1012, 6, '型号B', 'HUMIDITY', '2023-01-12', 'MQTT', 'ACTIVE'),
(1013, 7, '型号A', 'TEMPERATURE', '2023-01-13', 'MQTT', 'ACTIVE'),
(1014, 7, '型号B', 'HUMIDITY', '2023-01-14', 'MQTT', 'ACTIVE'),
(1015, 8, '型号A', 'TEMPERATURE', '2023-01-15', 'MQTT', 'ACTIVE'),
(1016, 8, '型号B', 'HUMIDITY', '2023-01-16', 'MQTT', 'ACTIVE'),
(1017, 9, '型号A', 'TEMPERATURE', '2023-01-17', 'MQTT', 'ACTIVE'),
(1018, 9, '型号B', 'HUMIDITY', '2023-01-18', 'MQTT', 'ACTIVE'),
(1019, 10, '型号A', 'TEMPERATURE', '2023-01-19', 'MQTT', 'ACTIVE'),
(1020, 10, '型号B', 'HUMIDITY', '2023-01-20', 'MQTT', 'ACTIVE'),
(1021, 11, '型号A', 'TEMPERATURE', '2023-01-21', 'MQTT', 'ACTIVE'),
(1022, 11, '型号B', 'HUMIDITY', '2023-01-22', 'MQTT', 'ACTIVE'),
(1023, 12, '型号A', 'TEMPERATURE', '2023-01-23', 'MQTT', 'ACTIVE'),
(1024, 12, '型号B', 'HUMIDITY', '2023-01-24', 'MQTT', 'ACTIVE'),
(1025, 13, '型号A', 'TEMPERATURE', '2023-01-25', 'MQTT', 'ACTIVE'),
(1026, 13, '型号B', 'HUMIDITY', '2023-01-26', 'MQTT', 'ACTIVE'),
(1027, 14, '型号A', 'TEMPERATURE', '2023-01-27', 'MQTT', 'ACTIVE'),
(1028, 14, '型号B', 'HUMIDITY', '2023-01-28', 'MQTT', 'ACTIVE'),
(1029, 15, '型号A', 'TEMPERATURE', '2023-01-29', 'MQTT', 'ACTIVE'),
(1030, 15, '型号B', 'HUMIDITY', '2023-01-30', 'MQTT', 'ACTIVE');

-- 插入测试传感器读数数据
-- 为每个传感器生成最近7天的读数
DELIMITER //
CREATE PROCEDURE generate_sensor_readings()
BEGIN
    DECLARE sensor_id INT DEFAULT 1001;
    DECLARE i INT DEFAULT 0;
    DECLARE reading_time DATETIME;
    DECLARE temperature DECIMAL(10,2);
    DECLARE humidity DECIMAL(10,2);
    
    WHILE sensor_id <= 1030 DO
        SET i = 0;
        WHILE i < 168 DO -- 7天 * 24小时
            SET reading_time = NOW() - INTERVAL i HOUR;
            
            -- 生成温度数据
            IF sensor_id % 2 = 1 THEN -- 奇数传感器是温度传感器
                SET temperature = 15 + RAND() * 20; -- 15-35度
                INSERT INTO sensor_readings (sensor_id, collected_at, reading_type, numeric_value, value_unit, data_status)
                VALUES (sensor_id, reading_time, 'TEMPERATURE', temperature, '℃', 'VALID');
            ELSE -- 偶数传感器是湿度传感器
                SET humidity = 30 + RAND() * 50; -- 30-80%
                INSERT INTO sensor_readings (sensor_id, collected_at, reading_type, numeric_value, value_unit, data_status)
                VALUES (sensor_id, reading_time, 'HUMIDITY', humidity, '%', 'VALID');
            END IF;
            
            SET i = i + 1;
        END WHILE;
        
        SET sensor_id = sensor_id + 1;
    END WHILE;
END //
DELIMITER ;

-- 调用存储过程生成传感器读数数据
CALL generate_sensor_readings();

-- 删除临时存储过程
DROP PROCEDURE generate_sensor_readings;

-- 插入测试预警规则数据
INSERT INTO alert_rules (alert_type, condition_expr, severity_level, is_active) VALUES
('FIRE', 'temperature > 35 AND humidity < 20', 'SEVERE', 1),
('FIRE', 'temperature > 30 AND humidity < 25', 'MODERATE', 1),
('DROUGHT', 'humidity < 15', 'SEVERE', 1),
('DROUGHT', 'humidity < 20', 'MODERATE', 1),
('PEST', 'other_sensor = "pest_detected"', 'GENERAL', 1),
('OTHER', 'system_error = true', 'CRITICAL', 1);

-- 插入测试预警数据
INSERT INTO alerts (rule_id, region_id, triggered_at, content, status, handler_id, handle_result) VALUES
(1, 1, NOW() - INTERVAL 1 HOUR, '常绿森林A区温度超过35℃，湿度低于20%，可能发生火灾', 'PROCESSING', 3, NULL),
(2, 2, NOW() - INTERVAL 2 HOUR, '常绿森林B区温度超过30℃，湿度低于25%，请注意防火', 'PENDING', NULL, NULL),
(3, 3, NOW() - INTERVAL 3 HOUR, '草地保护区1湿度低于15%，严重干旱', 'CLOSED', 4, '已安排浇水作业'),
(4, 4, NOW() - INTERVAL 4 HOUR, '草地保护区2湿度低于20%，轻度干旱', 'PENDING', NULL, NULL),
(5, 5, NOW() - INTERVAL 5 HOUR, '落叶森林C区检测到病虫害', 'PROCESSING', 3, NULL),
(6, 6, NOW() - INTERVAL 6 HOUR, '落叶森林D区系统错误', 'CLOSED', 1, '已修复系统故障');

-- 插入测试通知数据
INSERT INTO notifications (alert_id, receiver_id, channel, sent_at, receive_status) VALUES
(1, 3, 'SMS', NOW() - INTERVAL 1 HOUR, 'DELIVERED'),
(1, 4, 'SMS', NOW() - INTERVAL 1 HOUR, 'DELIVERED'),
(1, 7, 'SYSTEM', NOW() - INTERVAL 1 HOUR, 'DELIVERED'),
(2, 3, 'SMS', NOW() - INTERVAL 2 HOUR, 'DELIVERED'),
(2, 4, 'SYSTEM', NOW() - INTERVAL 2 HOUR, 'DELIVERED'),
(3, 3, 'SMS', NOW() - INTERVAL 3 HOUR, 'DELIVERED'),
(3, 4, 'SMS', NOW() - INTERVAL 3 HOUR, 'DELIVERED'),
(4, 3, 'SYSTEM', NOW() - INTERVAL 4 HOUR, 'DELIVERED'),
(5, 3, 'SMS', NOW() - INTERVAL 5 HOUR, 'DELIVERED'),
(5, 7, 'SYSTEM', NOW() - INTERVAL 5 HOUR, 'DELIVERED'),
(6, 1, 'SMS', NOW() - INTERVAL 6 HOUR, 'DELIVERED'),
(6, 7, 'SYSTEM', NOW() - INTERVAL 6 HOUR, 'DELIVERED');

-- 插入测试资源数据
INSERT INTO resources (resource_type, region_id, species_name, quantity, area, growth_stage, planted_at) VALUES
('TREE', 1, '油松', 1000, NULL, 'MATURE', '2010-03-01'),
('TREE', 1, '落叶松', 800, NULL, 'MATURE', '2010-03-01'),
('TREE', 2, '油松', 1200, NULL, 'MATURE', '2010-03-01'),
('TREE', 2, '落叶松', 900, NULL, 'MATURE', '2010-03-01'),
('TREE', 3, '杨树', 1500, NULL, 'GROWING', '2015-03-01'),
('TREE', 3, '柳树', 1300, NULL, 'GROWING', '2015-03-01'),
('TREE', 4, '杨树', 1600, NULL, 'GROWING', '2015-03-01'),
('TREE', 4, '柳树', 1400, NULL, 'GROWING', '2015-03-01'),
('GRASS', 5, '早熟禾', NULL, 5000.00, 'MATURE', '2018-04-01'),
('GRASS', 5, '黑麦草', NULL, 4000.00, 'MATURE', '2018-04-01'),
('GRASS', 6, '早熟禾', NULL, 5500.00, 'MATURE', '2018-04-01'),
('GRASS', 6, '黑麦草', NULL, 4500.00, 'MATURE', '2018-04-01'),
('TREE', 7, '柏树', 800, NULL, 'MATURE', '2008-03-01'),
('TREE', 7, '松树', 700, NULL, 'MATURE', '2008-03-01'),
('TREE', 8, '柏树', 900, NULL, 'MATURE', '2008-03-01'),
('TREE', 8, '松树', 800, NULL, 'MATURE', '2008-03-01'),
('GRASS', 9, '野牛草', NULL, 3000.00, 'GROWING', '2019-04-01'),
('GRASS', 9, '结缕草', NULL, 2500.00, 'GROWING', '2019-04-01'),
('GRASS', 10, '野牛草', NULL, 3500.00, 'GROWING', '2019-04-01'),
('GRASS', 10, '结缕草', NULL, 3000.00, 'GROWING', '2019-04-01');

-- 插入测试资源变动数据
INSERT INTO resource_changes (resource_id, change_type, change_reason, change_amount, change_area, operator_id) VALUES
(1, 'ADD', '新种植油松', 100, NULL, 2),
(2, 'ADD', '新种植落叶松', 80, NULL, 2),
(3, 'ADD', '新种植油松', 120, NULL, 2),
(4, 'ADD', '新种植落叶松', 90, NULL, 2),
(5, 'ADD', '新种植杨树', 150, NULL, 2),
(6, 'ADD', '新种植柳树', 130, NULL, 2),
(7, 'ADD', '新种植杨树', 160, NULL, 2),
(8, 'ADD', '新种植柳树', 140, NULL, 2),
(9, 'REDUCE', '干旱导致部分草地死亡', NULL, 500.00, 3),
(10, 'REDUCE', '干旱导致部分草地死亡', NULL, 400.00, 3),
(11, 'REDUCE', '干旱导致部分草地死亡', NULL, 550.00, 4),
(12, 'REDUCE', '干旱导致部分草地死亡', NULL, 450.00, 4),
(13, 'STATUS_UPDATE', '柏树进入成熟期', NULL, NULL, 2),
(14, 'STATUS_UPDATE', '松树进入成熟期', NULL, NULL, 2),
(15, 'STATUS_UPDATE', '柏树进入成熟期', NULL, NULL, 2),
(16, 'STATUS_UPDATE', '松树进入成熟期', NULL, NULL, 2),
(17, 'ADD', '新种植野牛草', NULL, 300.00, 2),
(18, 'ADD', '新种植结缕草', NULL, 250.00, 2),
(19, 'ADD', '新种植野牛草', NULL, 350.00, 2),
(20, 'ADD', '新种植结缕草', NULL, 300.00, 2);

-- 插入测试设备数据
INSERT INTO devices (device_id, device_name, device_type, model_spec, purchased_at, install_region_id, installer_id, warranty_until) VALUES
(2001, '温度传感器1', 'SENSOR', '型号A', '2023-01-01', 1, 9, '2026-01-01'),
(2002, '湿度传感器1', 'SENSOR', '型号B', '2023-01-02', 1, 9, '2026-01-02'),
(2003, '温度传感器2', 'SENSOR', '型号A', '2023-01-03', 2, 9, '2026-01-03'),
(2004, '湿度传感器2', 'SENSOR', '型号B', '2023-01-04', 2, 9, '2026-01-04'),
(2005, '摄像头1', 'CAMERA', '型号C', '2023-02-01', 1, 10, '2026-02-01'),
(2006, '摄像头2', 'CAMERA', '型号C', '2023-02-02', 2, 10, '2026-02-02'),
(2007, '预警器1', 'ALERTER', '型号D', '2023-03-01', 1, 9, '2026-03-01'),
(2008, '预警器2', 'ALERTER', '型号D', '2023-03-02', 2, 9, '2026-03-02'),
(2009, '温度传感器3', 'SENSOR', '型号A', '2023-01-05', 3, 9, '2026-01-05'),
(2010, '湿度传感器3', 'SENSOR', '型号B', '2023-01-06', 3, 9, '2026-01-06'),
(2011, '温度传感器4', 'SENSOR', '型号A', '2023-01-07', 4, 9, '2026-01-07'),
(2012, '湿度传感器4', 'SENSOR', '型号B', '2023-01-08', 4, 9, '2026-01-08'),
(2013, '摄像头3', 'CAMERA', '型号C', '2023-02-03', 3, 10, '2026-02-03'),
(2014, '摄像头4', 'CAMERA', '型号C', '2023-02-04', 4, 10, '2026-02-04'),
(2015, '预警器3', 'ALERTER', '型号D', '2023-03-03', 3, 9, '2026-03-03'),
(2016, '预警器4', 'ALERTER', '型号D', '2023-03-04', 4, 9, '2026-03-04'),
(2017, '温度传感器5', 'SENSOR', '型号A', '2023-01-09', 5, 9, '2026-01-09'),
(2018, '湿度传感器5', 'SENSOR', '型号B', '2023-01-10', 5, 9, '2026-01-10'),
(2019, '温度传感器6', 'SENSOR', '型号A', '2023-01-11', 6, 9, '2026-01-11'),
(2020, '湿度传感器6', 'SENSOR', '型号B', '2023-01-12', 6, 9, '2026-01-12');

-- 插入测试设备状态日志数据
-- 为每个设备生成最近24小时的状态日志
DELIMITER //
CREATE PROCEDURE generate_device_status_logs()
BEGIN
    DECLARE device_id INT DEFAULT 2001;
    DECLARE i INT DEFAULT 0;
    DECLARE status_time DATETIME;
    DECLARE run_status VARCHAR(20);
    DECLARE battery_percent INT;
    DECLARE signal_strength INT;
    
    WHILE device_id <= 2020 DO
        SET i = 0;
        WHILE i < 24 DO -- 最近24小时，每小时一次
            SET status_time = NOW() - INTERVAL i HOUR;
            
            -- 随机生成设备状态
            SET run_status = CASE 
                WHEN RAND() < 0.9 THEN 'NORMAL'
                WHEN RAND() < 0.95 THEN 'FAULT'
                ELSE 'OFFLINE'
            END;
            
            -- 生成电池电量（0-100%）
            SET battery_percent = FLOOR(RAND() * 100);
            
            -- 生成信号强度（0-100）
            SET signal_strength = FLOOR(RAND() * 100);
            
            INSERT INTO device_status_logs (device_id, collected_at, run_status, battery_percent, signal_strength)
            VALUES (device_id, status_time, run_status, battery_percent, signal_strength);
            
            SET i = i + 1;
        END WHILE;
        
        SET device_id = device_id + 1;
    END WHILE;
END //
DELIMITER ;

-- 调用存储过程生成设备状态日志数据
CALL generate_device_status_logs();

-- 删除临时存储过程
DROP PROCEDURE generate_device_status_logs;

-- 插入测试维护记录数据
INSERT INTO maintenance_records (device_id, maintenance_type, maintenance_time, maintainer_id, content, result) VALUES
(2001, 'INSPECTION', NOW() - INTERVAL 1 DAY, 9, '定期巡检温度传感器1', '正常'),
(2002, 'INSPECTION', NOW() - INTERVAL 1 DAY, 9, '定期巡检湿度传感器1', '正常'),
(2003, 'INSPECTION', NOW() - INTERVAL 1 DAY, 9, '定期巡检温度传感器2', '正常'),
(2004, 'INSPECTION', NOW() - INTERVAL 1 DAY, 9, '定期巡检湿度传感器2', '正常'),
(2005, 'INSPECTION', NOW() - INTERVAL 2 DAY, 10, '定期巡检摄像头1', '正常'),
(2006, 'INSPECTION', NOW() - INTERVAL 2 DAY, 10, '定期巡检摄像头2', '正常'),
(2007, 'INSPECTION', NOW() - INTERVAL 3 DAY, 9, '定期巡检预警器1', '正常'),
(2008, 'INSPECTION', NOW() - INTERVAL 3 DAY, 9, '定期巡检预警器2', '正常'),
(2009, 'REPAIR', NOW() - INTERVAL 4 DAY, 9, '修复温度传感器3故障', '已修复'),
(2010, 'REPAIR', NOW() - INTERVAL 5 DAY, 9, '修复湿度传感器3故障', '已修复'),
(2011, 'REPLACE', NOW() - INTERVAL 6 DAY, 9, '更换损坏的温度传感器4', '已更换'),
(2012, 'REPLACE', NOW() - INTERVAL 7 DAY, 9, '更换损坏的湿度传感器4', '已更换'),
(2013, 'REPAIR', NOW() - INTERVAL 8 DAY, 10, '修复摄像头3故障', '已修复'),
(2014, 'REPAIR', NOW() - INTERVAL 9 DAY, 10, '修复摄像头4故障', '已修复'),
(2015, 'INSPECTION', NOW() - INTERVAL 10 DAY, 9, '定期巡检预警器3', '正常'),
(2016, 'INSPECTION', NOW() - INTERVAL 11 DAY, 9, '定期巡检预警器4', '正常'),
(2017, 'INSPECTION', NOW() - INTERVAL 12 DAY, 9, '定期巡检温度传感器5', '正常'),
(2018, 'INSPECTION', NOW() - INTERVAL 13 DAY, 9, '定期巡检湿度传感器5', '正常'),
(2019, 'INSPECTION', NOW() - INTERVAL 14 DAY, 9, '定期巡检温度传感器6', '正常'),
(2020, 'INSPECTION', NOW() - INTERVAL 15 DAY, 9, '定期巡检湿度传感器6', '正常');

-- 插入测试报表模板数据
INSERT INTO report_templates (template_name, dimension, metrics, cycle, is_active) VALUES
('区域传感器统计', 'REGION', 'sensor_count, avg_temperature, avg_humidity', 'DAY', 1),
('设备状态统计', 'TYPE', 'device_count, normal_count, fault_count, offline_count', 'DAY', 1),
('预警统计', 'TIME', 'alert_count, fire_alerts, drought_alerts', 'WEEK', 1),
('资源统计', 'REGION', 'tree_count, grass_area, growth_stage_distribution', 'MONTH', 1),
('环境趋势分析', 'TIME', 'temperature_trend, humidity_trend', 'MONTH', 1);

-- 插入测试报表数据
INSERT INTO reports (template_id, period, generated_at, file_path, data_source) VALUES
(1, '2025-01-01', NOW() - INTERVAL 1 DAY, '/reports/daily_region_sensors_20250101.pdf', '自动生成'),
(1, '2025-01-02', NOW() - INTERVAL 2 DAY, '/reports/daily_region_sensors_20250102.pdf', '自动生成'),
(2, '2025-01-01', NOW() - INTERVAL 1 DAY, '/reports/daily_device_status_20250101.pdf', '自动生成'),
(2, '2025-01-02', NOW() - INTERVAL 2 DAY, '/reports/daily_device_status_20250102.pdf', '自动生成'),
(3, '2025-W01', NOW() - INTERVAL 7 DAY, '/reports/weekly_alerts_2025W01.pdf', '自动生成'),
(3, '2025-W02', NOW() - INTERVAL 14 DAY, '/reports/weekly_alerts_2025W02.pdf', '自动生成'),
(4, '2024-12', NOW() - INTERVAL 30 DAY, '/reports/monthly_resources_202412.pdf', '自动生成'),
(4, '2024-11', NOW() - INTERVAL 60 DAY, '/reports/monthly_resources_202411.pdf', '自动生成'),
(5, '2024-12', NOW() - INTERVAL 30 DAY, '/reports/monthly_environment_202412.pdf', '自动生成'),
(5, '2024-11', NOW() - INTERVAL 60 DAY, '/reports/monthly_environment_202411.pdf', '自动生成');

COMMIT;
SET FOREIGN_KEY_CHECKS = 1; -- 恢复外键校验

-- 显示插入结果
SELECT '测试数据生成完成' AS result;