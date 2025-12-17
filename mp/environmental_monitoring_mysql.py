#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
智慧林草系统 - 环境监测和统计分析模块MySQL实现

该脚本使用MySQL数据库实现环境监测和统计分析模块的主要功能，
包括数据的增删改查、统计分析、异常检测和报表生成等。
"""

import pymysql
import datetime
import random
from collections import defaultdict
import sys
import os

class EnvironmentalMonitoringSystem:
    """环境监测系统MySQL实现类"""
    
    def __init__(self, host='localhost', user='root', password='mp050410', database='smart_forest_grass', port=3306):
        """初始化系统"""
        print("正在连接MySQL数据库...")
        try:
            # 先连接到MySQL服务器（不指定数据库）
            self.conn = pymysql.connect(
                host=host,
                user=user,
                password=password,
                port=port,
                charset='utf8mb4',
                cursorclass=pymysql.cursors.DictCursor
            )
            print("MySQL服务器连接成功")
            self.cursor = self.conn.cursor()
            
            # 初始化数据库结构
            self._initialize_database(database)
            
            print("智慧林草系统 - 环境监测和统计分析模块启动")
            print("=" * 60)
        except pymysql.Error as e:
            print(f"MySQL连接错误: {e}")
            print("请确保MySQL服务已启动，并且连接参数正确")
            print(f"连接参数: host={host}, user={user}, password=******, database={database}, port={port}")
            exit(1)
    
    def _initialize_database(self, database):
        """初始化数据库结构"""
        print("正在初始化数据库结构...")
        
        # 创建数据库（如果不存在）
        self.cursor.execute(f"CREATE DATABASE IF NOT EXISTS {database} CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
        self.cursor.execute(f"USE {database}")
        # 手动更新数据库连接的db属性
        self.conn.db = database
        
        # 1. 环境监测站表
        self.cursor.execute('''
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
        '''),
        
        # 2. 环境数据采集表
        self.cursor.execute('''
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
        '''),
        
        # 3. 统计分析数据表
        self.cursor.execute('''
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
        ''')
            
        
    
        # 4. 报表模板表
        self.cursor.execute('''
        CREATE TABLE IF NOT EXISTS report_template (
            template_id INT PRIMARY KEY AUTO_INCREMENT COMMENT '模板ID',
            template_name VARCHAR(100) NOT NULL COMMENT '模板名称',
            report_type VARCHAR(50) NOT NULL COMMENT '报表类型',
            template_content TEXT NOT NULL COMMENT '模板内容',
            create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
            update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
            UNIQUE KEY uk_template_name (template_name)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='报表模板表'
        ''')
        
        # 5. 生成的报表表
        self.cursor.execute('''
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
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='生成的报表表'
        ''')
        
        # 6. 异常数据记录表
        self.cursor.execute('''
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
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='异常数据记录表'
        ''')
        
        # 7. 系统用户表（用于权限控制）
        self.cursor.execute('''
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
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户表'
        ''')
        
        # 8. 操作日志表
        self.cursor.execute('''
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
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表'
        ''')
        
        # 创建视图
        self._create_views()
        
        # 创建存储过程
        self._create_stored_procedures()
        
        # 创建触发器
        self._create_triggers()
        
        # 初始化数据
        self._initialize_data()
        
        self.conn.commit()
        print("数据库初始化完成\n")
    
    def _create_views(self):
        """创建系统视图"""
        print("正在创建视图...")
        
        # 环境数据概览视图
        self.cursor.execute('''
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
        JOIN monitoring_station ms ON ed.station_id = ms.station_id
        ''')
        
        # 站点统计数据视图
        self.cursor.execute('''
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
        JOIN monitoring_station ms ON sd.station_id = ms.station_id
        ''')
        
        # 异常数据视图
        self.cursor.execute('''
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
        JOIN monitoring_station ms ON ed.station_id = ms.station_id
        ''')
        
        # 管理员视图
        self.cursor.execute('''
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
        JOIN monitoring_station ms ON ed.station_id = ms.station_id
        ''')
        
        self.cursor.execute('''
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
        JOIN monitoring_station ms ON sd.station_id = ms.station_id
        ''')
        
        self.cursor.execute('''
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
        JOIN monitoring_station ms ON ed.station_id = ms.station_id
        ''')
        
        # 普通用户视图
        self.cursor.execute('''
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
        WHERE ed.collection_time >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
        ''')
        
        # 访客视图
        self.cursor.execute('''
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
        GROUP BY ms.station_name, DATE(ed.collection_time)
        ''')
        
        # 实时数据视图
        self.cursor.execute('''
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
        AND ms.status = '正常'
        ''')
        
        # 环境趋势视图
        self.cursor.execute('''
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
        ORDER BY ms.station_name, collection_date
        ''')
        
        # 空气质量视图
        self.cursor.execute('''
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
        ORDER BY ms.station_name, collection_date
        ''')
        
        print("视图创建完成\n")
    
    def _create_stored_procedures(self):
        """创建存储过程"""
        print("正在创建存储过程...")
        
        # 移除所有DELIMITER命令，确保每个存储过程可以单独执行
        
        # 按日统计环境数据存储过程
        self.cursor.execute('''
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
        END
        ''')
        
        # 检测异常数据存储过程
        self.cursor.execute('''
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
        END
        ''')
        
        # 处理异常数据存储过程
        self.cursor.execute('''
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
        END
        ''')
        
        # 生成日报存储过程
        self.cursor.execute('''
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
            SELECT LAST_INSERT_ID() AS report_id;
        END
        ''')
        
        # 获取系统状态存储过程
        self.cursor.execute('''
        CREATE PROCEDURE IF NOT EXISTS sp_get_system_status()
        BEGIN
            SELECT 
                (SELECT COUNT(*) FROM monitoring_station) AS station_count,
                (SELECT COUNT(*) FROM monitoring_station WHERE status = '正常') AS normal_station_count,
                (SELECT COUNT(*) FROM environmental_data WHERE DATE(collection_time) = CURDATE()) AS today_data_count,
                (SELECT COUNT(*) FROM abnormal_data WHERE processing_status = '未处理') AS unprocessed_abnormal_count;
        END
        ''')
        
        print("存储过程创建完成\n")
    
    def _create_triggers(self):
        """创建触发器"""
        print("正在创建触发器...")
        
        # 环境数据插入后自动检测异常触发器
        self.cursor.execute('''
        CREATE TRIGGER IF NOT EXISTS tr_after_insert_environmental_data
        AFTER INSERT ON environmental_data
        FOR EACH ROW
        CALL sp_detect_abnormal_data();
        ''')
        
        # 环境数据更新后自动检测异常触发器
        self.cursor.execute('''
        CREATE TRIGGER IF NOT EXISTS tr_after_update_environmental_data
        AFTER UPDATE ON environmental_data
        FOR EACH ROW
        BEGIN
            DELETE FROM abnormal_data WHERE data_id = NEW.data_id;
            CALL sp_detect_abnormal_data();
        END
        ''')
        
        # 环境数据删除后自动删除相关异常记录触发器
        self.cursor.execute('''
        CREATE TRIGGER IF NOT EXISTS tr_after_delete_environmental_data
        AFTER DELETE ON environmental_data
        FOR EACH ROW
        DELETE FROM abnormal_data WHERE data_id = OLD.data_id;
        ''')
        
        print("触发器创建完成\n")
    
    def _initialize_data(self):
        """初始化系统数据"""
        print("正在初始化系统数据...")
        
        # 初始化环境监测站
        stations = [
            ("监测站1", 39.9042, 116.4074, 50.0, "气象站", "2023-01-01", "正常"),
            ("监测站2", 39.9142, 116.4174, 60.0, "土壤监测站", "2023-02-01", "正常"),
            ("监测站3", 39.9242, 116.4274, 70.0, "空气质量监测站", "2023-03-01", "正常"),
            ("监测站4", 39.9342, 116.4374, 80.0, "综合监测站", "2023-04-01", "维护"),
            ("监测站5", 39.9442, 116.4474, 90.0, "气象站", "2023-05-01", "正常")
        ]
        self.cursor.executemany('''
        INSERT IGNORE INTO monitoring_station (
            station_name, latitude, longitude, altitude, station_type, installation_date, status
        ) VALUES (%s, %s, %s, %s, %s, %s, %s)
        ''', stations)
        
        # 初始化报表模板
        templates = [
            ("每日环境监测报表", "日报", "每日环境监测数据汇总报表模板"),
            ("每周环境统计报表", "周报", "每周环境统计数据汇总报表模板"),
            ("每月环境分析报表", "月报", "每月环境分析数据汇总报表模板"),
            ("季度环境趋势报表", "季报", "季度环境趋势分析报表模板"),
            ("年度环境总结报表", "年报", "年度环境总结分析报表模板")
        ]
        self.cursor.executemany('''
        INSERT IGNORE INTO report_template (template_name, report_type, template_content) VALUES (%s, %s, %s)
        ''', templates)
        
        # 初始化系统用户
        users = [
            ("admin", "123456", "管理员", "admin@example.com", "13800138000"),
            ("user1", "123456", "普通用户", "user1@example.com", "13800138001"),
            ("user2", "123456", "普通用户", "user2@example.com", "13800138002"),
            ("guest", "123456", "访客", "guest@example.com", "13800138003")
        ]
        self.cursor.executemany('''
        INSERT IGNORE INTO system_user (username, password, role, email, phone) VALUES (%s, %s, %s, %s, %s)
        ''', users)
        
        # 生成模拟环境数据
        self.generate_sample_data()
        
        self.conn.commit()
        print("系统数据初始化完成\n")
    
    def generate_sample_data(self):
        """生成模拟环境数据"""
        # 获取所有正常状态的监测站
        self.cursor.execute('SELECT station_id, station_type FROM monitoring_station WHERE status = "正常"')
        stations = self.cursor.fetchall()
        
        # 生成最近7天的模拟数据
        for i in range(7):
            # 每天生成24条数据（每小时一条）
            for hour in range(24):
                collection_time = datetime.datetime.now() - datetime.timedelta(days=i, hours=23-hour)
                collection_time_str = collection_time.strftime("%Y-%m-%d %H:%M:%S")
                
                for station in stations:
                    station_id = station['station_id']
                    station_type = station['station_type']
                    
                    # 生成随机环境数据
                    temperature = round(random.uniform(-10, 35), 2)
                    humidity = round(random.uniform(30, 90), 2)
                    wind_speed = round(random.uniform(0, 20), 2)
                    wind_direction = random.choice(["东", "南", "西", "北", "东南", "西南", "东北", "西北"])
                    rainfall = round(random.uniform(0, 50), 2)
                    sunshine_duration = round(random.uniform(0, 12), 2)
                    soil_temperature = round(random.uniform(0, 30), 2) if station_type in ["土壤监测站", "综合监测站"] else None
                    soil_humidity = round(random.uniform(20, 80), 2) if station_type in ["土壤监测站", "综合监测站"] else None
                    soil_ph = round(random.uniform(4.5, 9.5), 2) if station_type in ["土壤监测站", "综合监测站"] else None
                    pm25 = round(random.uniform(0, 500), 2) if station_type in ["空气质量监测站", "综合监测站"] else None
                    pm10 = round(random.uniform(0, 1000), 2) if station_type in ["空气质量监测站", "综合监测站"] else None
                    
                    # 插入环境数据
                    self.cursor.execute('''
                    INSERT INTO environmental_data (
                        station_id, collection_time, temperature, humidity, wind_speed, wind_direction,
                        rainfall, sunshine_duration, soil_temperature, soil_humidity, soil_ph, pm25, pm10
                    ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                    ''', (
                        station_id, collection_time_str, temperature, humidity, wind_speed, wind_direction,
                        rainfall, sunshine_duration, soil_temperature, soil_humidity, soil_ph, pm25, pm10
                    ))
        
        # 调用异常数据检测
        self.detect_abnormal_data()
    
    def detect_abnormal_data(self):
        """检测异常数据"""
        print("正在检测异常数据...")
        
        # 先清空现有异常数据
        self.cursor.execute('DELETE FROM abnormal_data')
        
        # 检测温度异常（>40℃或< -20℃）
        self.cursor.execute('''
        INSERT INTO abnormal_data (data_id, abnormal_type, abnormal_description)
        SELECT 
            data_id,
            '温度异常',
            CONCAT('温度超出正常范围：', temperature, '℃')
        FROM environmental_data
        WHERE (temperature > 40 OR temperature < -20)
        ''')
        
        # 检测湿度异常（>100%或< 0%）
        self.cursor.execute('''
        INSERT INTO abnormal_data (data_id, abnormal_type, abnormal_description)
        SELECT 
            data_id,
            '湿度异常',
            CONCAT('湿度超出正常范围：', humidity, '%')
        FROM environmental_data
        WHERE (humidity > 100 OR humidity < 0)
        ''')
        
        # 检测风速异常（>30 m/s）
        self.cursor.execute('''
        INSERT INTO abnormal_data (data_id, abnormal_type, abnormal_description)
        SELECT 
            data_id,
            '风速异常',
            CONCAT('风速超出正常范围：', wind_speed, 'm/s')
        FROM environmental_data
        WHERE wind_speed > 30
        ''')
        
        # 检测PM2.5异常（>500 μg/m³）
        self.cursor.execute('''
        INSERT INTO abnormal_data (data_id, abnormal_type, abnormal_description)
        SELECT 
            data_id,
            'PM2.5异常',
            CONCAT('PM2.5超出正常范围：', pm25, 'μg/m³')
        FROM environmental_data
        WHERE pm25 > 500
        ''')
        
        # 检测PM10异常（>1000 μg/m³）
        self.cursor.execute('''
        INSERT INTO abnormal_data (data_id, abnormal_type, abnormal_description)
        SELECT 
            data_id,
            'PM10异常',
            CONCAT('PM10超出正常范围：', pm10, 'μg/m³')
        FROM environmental_data
        WHERE pm10 > 1000
        ''')
        
        # 检测土壤pH值异常（<4.5或>9.5）
        self.cursor.execute('''
        INSERT INTO abnormal_data (data_id, abnormal_type, abnormal_description)
        SELECT 
            data_id,
            '土壤pH值异常',
            CONCAT('土壤pH值超出正常范围：', soil_ph)
        FROM environmental_data
        WHERE (soil_ph < 4.5 OR soil_ph > 9.5)
        ''')
        
        self.conn.commit()
        
        # 获取检测到的异常数据数量
        self.cursor.execute('SELECT COUNT(*) AS abnormal_count FROM abnormal_data')
        abnormal_count = self.cursor.fetchone()['abnormal_count']
        print(f"检测到 {abnormal_count} 条异常数据\n")
    
    def get_realtime_data(self):
        """获取实时数据"""
        query = '''
        SELECT 
            ms.station_name, ms.station_type, ed.collection_time, 
            ed.temperature, ed.humidity, ed.wind_speed, ed.wind_direction, ed.rainfall
        FROM environmental_data ed
        JOIN monitoring_station ms ON ed.station_id = ms.station_id
        JOIN (
            SELECT station_id, MAX(collection_time) AS max_time
            FROM environmental_data
            GROUP BY station_id
        ) latest ON ed.station_id = latest.station_id AND ed.collection_time = latest.max_time
        WHERE ms.status = "正常"
        ORDER BY ms.station_name
        '''
        self.cursor.execute(query)
        return self.cursor.fetchall()
    
    def get_daily_statistics(self, stat_date):
        """获取日统计数据"""
        query = '''
        SELECT 
            ms.station_name,
            COUNT(ed.data_id) AS data_count,
            AVG(ed.temperature) AS avg_temperature,
            MAX(ed.temperature) AS max_temperature,
            MIN(ed.temperature) AS min_temperature,
            AVG(ed.humidity) AS avg_humidity,
            SUM(ed.rainfall) AS total_rainfall,
            AVG(ed.wind_speed) AS avg_wind_speed
        FROM environmental_data ed
        JOIN monitoring_station ms ON ed.station_id = ms.station_id
        WHERE DATE(ed.collection_time) = %s
        GROUP BY ms.station_id
        ORDER BY ms.station_name
        '''
        self.cursor.execute(query, (stat_date,))
        return self.cursor.fetchall()
    
    def get_weekly_statistics(self, year, week):
        """获取周统计数据"""
        query = '''
        SELECT 
            ms.station_name,
            COUNT(ed.data_id) AS data_count,
            AVG(ed.temperature) AS avg_temperature,
            MAX(ed.temperature) AS max_temperature,
            MIN(ed.temperature) AS min_temperature,
            AVG(ed.humidity) AS avg_humidity,
            SUM(ed.rainfall) AS total_rainfall,
            AVG(ed.wind_speed) AS avg_wind_speed
        FROM environmental_data ed
        JOIN monitoring_station ms ON ed.station_id = ms.station_id
        WHERE YEAR(ed.collection_time) = %s AND WEEK(ed.collection_time) = %s
        GROUP BY ms.station_id
        ORDER BY ms.station_name
        '''
        self.cursor.execute(query, (year, week))
        return self.cursor.fetchall()
    
    def get_monthly_statistics(self, year, month):
        """获取月统计数据"""
        query = '''
        SELECT 
            ms.station_name,
            COUNT(ed.data_id) AS data_count,
            AVG(ed.temperature) AS avg_temperature,
            MAX(ed.temperature) AS max_temperature,
            MIN(ed.temperature) AS min_temperature,
            AVG(ed.humidity) AS avg_humidity,
            SUM(ed.rainfall) AS total_rainfall,
            AVG(ed.wind_speed) AS avg_wind_speed
        FROM environmental_data ed
        JOIN monitoring_station ms ON ed.station_id = ms.station_id
        WHERE YEAR(ed.collection_time) = %s AND MONTH(ed.collection_time) = %s
        GROUP BY ms.station_id
        ORDER BY ms.station_name
        '''
        self.cursor.execute(query, (year, month))
        return self.cursor.fetchall()
    
    def get_quarterly_statistics(self, year, quarter):
        """获取季度统计数据"""
        query = '''
        SELECT 
            ms.station_name,
            COUNT(ed.data_id) AS data_count,
            AVG(ed.temperature) AS avg_temperature,
            MAX(ed.temperature) AS max_temperature,
            MIN(ed.temperature) AS min_temperature,
            AVG(ed.humidity) AS avg_humidity,
            SUM(ed.rainfall) AS total_rainfall,
            AVG(ed.wind_speed) AS avg_wind_speed
        FROM environmental_data ed
        JOIN monitoring_station ms ON ed.station_id = ms.station_id
        WHERE YEAR(ed.collection_time) = %s AND QUARTER(ed.collection_time) = %s
        GROUP BY ms.station_id
        ORDER BY ms.station_name
        '''
        self.cursor.execute(query, (year, quarter))
        return self.cursor.fetchall()
    
    def get_annual_statistics(self, year):
        """获取年度统计数据"""
        query = '''
        SELECT 
            ms.station_name,
            COUNT(ed.data_id) AS data_count,
            AVG(ed.temperature) AS avg_temperature,
            MAX(ed.temperature) AS max_temperature,
            MIN(ed.temperature) AS min_temperature,
            AVG(ed.humidity) AS avg_humidity,
            SUM(ed.rainfall) AS total_rainfall,
            AVG(ed.wind_speed) AS avg_wind_speed
        FROM environmental_data ed
        JOIN monitoring_station ms ON ed.station_id = ms.station_id
        WHERE YEAR(ed.collection_time) = %s
        GROUP BY ms.station_id
        ORDER BY ms.station_name
        '''
        self.cursor.execute(query, (year,))
        return self.cursor.fetchall()
    
    def generate_daily_report(self, report_date):
        """生成日报"""
        # 确保报表模板存在
        self.cursor.execute('SELECT COUNT(*) FROM report_template WHERE template_name = "每日环境监测报表"')
        if self.cursor.fetchone()['COUNT(*)'] == 0:
            # 创建默认模板
            self.cursor.execute('''
            INSERT INTO report_template (template_name, report_type, template_content)
            VALUES ("每日环境监测报表", "日报", "默认日报模板")
            ''')
            self.conn.commit()
        
        # 获取日报模板
        self.cursor.execute('SELECT template_id FROM report_template WHERE template_name = "每日环境监测报表"')
        template_id = self.cursor.fetchone()['template_id']
        
        # 获取日统计数据
        statistics = self.get_daily_statistics(report_date)
        
        # 生成报表内容
        report_content = f"智慧林草系统 - 环境监测日报\n"
        report_content += f"报表日期：{report_date}\n"
        report_content += f"生成时间：{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n"
        report_content += "\n========================================\n"
        report_content += "\n站点环境数据汇总：\n"
        report_content += "\n{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}\n".format(
            "站点名称", "数据条数", "平均温度", "最高温度", "最低温度", "平均湿度", "总降雨量"
        )
        report_content += "-" * 85 + "\n"
        
        for stat in statistics:
            report_content += "{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}\n".format(
                stat['station_name'],
                stat['data_count'],
                f"{round(stat['avg_temperature'], 2) if stat['avg_temperature'] else 0}℃",
                f"{round(stat['max_temperature'], 2) if stat['max_temperature'] else 0}℃",
                f"{round(stat['min_temperature'], 2) if stat['min_temperature'] else 0}℃",
                f"{round(stat['avg_humidity'], 2) if stat['avg_humidity'] else 0}%",
                f"{round(stat['total_rainfall'], 2) if stat['total_rainfall'] else 0}mm"
            )
        
        # 获取当日异常数据数量
        self.cursor.execute('''
        SELECT COUNT(*) AS abnormal_count FROM abnormal_data
        WHERE DATE(discovery_time) = %s
        ''', (report_date,))
        abnormal_count = self.cursor.fetchone()['abnormal_count']
        report_content += f"\n当日异常数据条数：{abnormal_count}\n"
        
        # 保存报表
        report_name = f"环境监测日报_{report_date}"
        self.cursor.execute('''
        INSERT INTO generated_report (template_id, report_name, report_content, generated_by)
        VALUES (%s, %s, %s, "system")
        ''', (template_id, report_name, report_content))
        
        self.conn.commit()
        return report_content
    
    def generate_weekly_report(self, year, week):
        """生成周报"""
        # 确保报表模板存在
        self.cursor.execute('SELECT COUNT(*) FROM report_template WHERE template_name = "每周环境监测报表"')
        if self.cursor.fetchone()['COUNT(*)'] == 0:
            # 创建默认模板
            self.cursor.execute('''
            INSERT INTO report_template (template_name, report_type, template_content)
            VALUES ("每周环境监测报表", "周报", "默认周报模板")
            ''')
            self.conn.commit()
        
        # 获取周报模板
        self.cursor.execute('SELECT template_id FROM report_template WHERE template_name = "每周环境监测报表"')
        template_id = self.cursor.fetchone()['template_id']
        
        # 获取周统计数据
        statistics = self.get_weekly_statistics(year, week)
        
        # 生成报表内容
        report_content = f"智慧林草系统 - 环境监测周报\n"
        report_content += f"报表年份：{year}\n"
        report_content += f"报表周数：第{week}周\n"
        report_content += f"生成时间：{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n"
        report_content += "\n========================================\n"
        report_content += "\n站点环境数据汇总：\n"
        report_content += "\n{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}\n".format(
            "站点名称", "数据条数", "平均温度", "最高温度", "最低温度", "平均湿度", "总降雨量"
        )
        report_content += "-" * 85 + "\n"
        
        for stat in statistics:
            report_content += "{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}\n".format(
                stat['station_name'],
                stat['data_count'],
                f"{round(stat['avg_temperature'], 2) if stat['avg_temperature'] else 0}℃",
                f"{round(stat['max_temperature'], 2) if stat['max_temperature'] else 0}℃",
                f"{round(stat['min_temperature'], 2) if stat['min_temperature'] else 0}℃",
                f"{round(stat['avg_humidity'], 2) if stat['avg_humidity'] else 0}%",
                f"{round(stat['total_rainfall'], 2) if stat['total_rainfall'] else 0}mm"
            )
        
        # 保存报表
        report_name = f"环境监测周报_{year}_第{week}周"
        self.cursor.execute('''
        INSERT INTO generated_report (template_id, report_name, report_content, generated_by)
        VALUES (%s, %s, %s, "system")
        ''', (template_id, report_name, report_content))
        
        self.conn.commit()
        return report_content
    
    def generate_monthly_report(self, year, month):
        """生成月报"""
        # 确保报表模板存在
        self.cursor.execute('SELECT COUNT(*) FROM report_template WHERE template_name = "每月环境监测报表"')
        if self.cursor.fetchone()['COUNT(*)'] == 0:
            # 创建默认模板
            self.cursor.execute('''
            INSERT INTO report_template (template_name, report_type, template_content)
            VALUES ("每月环境监测报表", "月报", "默认月报模板")
            ''')
            self.conn.commit()
        
        # 获取月报模板
        self.cursor.execute('SELECT template_id FROM report_template WHERE template_name = "每月环境监测报表"')
        template_id = self.cursor.fetchone()['template_id']
        
        # 获取月统计数据
        statistics = self.get_monthly_statistics(year, month)
        
        # 生成报表内容
        report_content = f"智慧林草系统 - 环境监测月报\n"
        report_content += f"报表年份：{year}\n"
        report_content += f"报表月份：{month}\n"
        report_content += f"生成时间：{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n"
        report_content += "\n========================================\n"
        report_content += "\n站点环境数据汇总：\n"
        report_content += "\n{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}\n".format(
            "站点名称", "数据条数", "平均温度", "最高温度", "最低温度", "平均湿度", "总降雨量"
        )
        report_content += "-" * 85 + "\n"
        
        for stat in statistics:
            report_content += "{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}\n".format(
                stat['station_name'],
                stat['data_count'],
                f"{round(stat['avg_temperature'], 2) if stat['avg_temperature'] else 0}℃",
                f"{round(stat['max_temperature'], 2) if stat['max_temperature'] else 0}℃",
                f"{round(stat['min_temperature'], 2) if stat['min_temperature'] else 0}℃",
                f"{round(stat['avg_humidity'], 2) if stat['avg_humidity'] else 0}%",
                f"{round(stat['total_rainfall'], 2) if stat['total_rainfall'] else 0}mm"
            )
        
        # 保存报表
        report_name = f"环境监测月报_{year}_{month}"
        self.cursor.execute('''
        INSERT INTO generated_report (template_id, report_name, report_content, generated_by)
        VALUES (%s, %s, %s, "system")
        ''', (template_id, report_name, report_content))
        
        self.conn.commit()
        return report_content
    
    def generate_quarterly_report(self, year, quarter):
        """生成季报"""
        # 确保报表模板存在
        self.cursor.execute('SELECT COUNT(*) FROM report_template WHERE template_name = "季度环境监测报表"')
        if self.cursor.fetchone()['COUNT(*)'] == 0:
            # 创建默认模板
            self.cursor.execute('''
            INSERT INTO report_template (template_name, report_type, template_content)
            VALUES ("季度环境监测报表", "季报", "默认季报模板")
            ''')
            self.conn.commit()
        
        # 获取季报模板
        self.cursor.execute('SELECT template_id FROM report_template WHERE template_name = "季度环境监测报表"')
        template_id = self.cursor.fetchone()['template_id']
        
        # 获取季度统计数据
        statistics = self.get_quarterly_statistics(year, quarter)
        
        # 生成报表内容
        report_content = f"智慧林草系统 - 环境监测季报\n"
        report_content += f"报表年份：{year}\n"
        report_content += f"报表季度：第{quarter}季度\n"
        report_content += f"生成时间：{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n"
        report_content += "\n========================================\n"
        report_content += "\n站点环境数据汇总：\n"
        report_content += "\n{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}\n".format(
            "站点名称", "数据条数", "平均温度", "最高温度", "最低温度", "平均湿度", "总降雨量"
        )
        report_content += "-" * 85 + "\n"
        
        for stat in statistics:
            report_content += "{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}\n".format(
                stat['station_name'],
                stat['data_count'],
                f"{round(stat['avg_temperature'], 2) if stat['avg_temperature'] else 0}℃",
                f"{round(stat['max_temperature'], 2) if stat['max_temperature'] else 0}℃",
                f"{round(stat['min_temperature'], 2) if stat['min_temperature'] else 0}℃",
                f"{round(stat['avg_humidity'], 2) if stat['avg_humidity'] else 0}%",
                f"{round(stat['total_rainfall'], 2) if stat['total_rainfall'] else 0}mm"
            )
        
        # 保存报表
        report_name = f"环境监测季报_{year}_第{quarter}季度"
        self.cursor.execute('''
        INSERT INTO generated_report (template_id, report_name, report_content, generated_by)
        VALUES (%s, %s, %s, "system")
        ''', (template_id, report_name, report_content))
        
        self.conn.commit()
        return report_content
    
    def generate_annual_report(self, year):
        """生成年报"""
        # 确保报表模板存在
        self.cursor.execute('SELECT COUNT(*) FROM report_template WHERE template_name = "年度环境监测报表"')
        if self.cursor.fetchone()['COUNT(*)'] == 0:
            # 创建默认模板
            self.cursor.execute('''
            INSERT INTO report_template (template_name, report_type, template_content)
            VALUES ("年度环境监测报表", "年报", "默认年报模板")
            ''')
            self.conn.commit()
        
        # 获取年报模板
        self.cursor.execute('SELECT template_id FROM report_template WHERE template_name = "年度环境监测报表"')
        template_id = self.cursor.fetchone()['template_id']
        
        # 获取年度统计数据
        statistics = self.get_annual_statistics(year)
        
        # 生成报表内容
        report_content = f"智慧林草系统 - 环境监测年报\n"
        report_content += f"报表年份：{year}\n"
        report_content += f"生成时间：{datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n"
        report_content += "\n========================================\n"
        report_content += "\n站点环境数据汇总：\n"
        report_content += "\n{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}\n".format(
            "站点名称", "数据条数", "平均温度", "最高温度", "最低温度", "平均湿度", "总降雨量"
        )
        report_content += "-" * 85 + "\n"
        
        for stat in statistics:
            report_content += "{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}\n".format(
                stat['station_name'],
                stat['data_count'],
                f"{round(stat['avg_temperature'], 2) if stat['avg_temperature'] else 0}℃",
                f"{round(stat['max_temperature'], 2) if stat['max_temperature'] else 0}℃",
                f"{round(stat['min_temperature'], 2) if stat['min_temperature'] else 0}℃",
                f"{round(stat['avg_humidity'], 2) if stat['avg_humidity'] else 0}%",
                f"{round(stat['total_rainfall'], 2) if stat['total_rainfall'] else 0}mm"
            )
        
        # 保存报表
        report_name = f"环境监测年报_{year}"
        self.cursor.execute('''
        INSERT INTO generated_report (template_id, report_name, report_content, generated_by)
        VALUES (%s, %s, %s, "system")
        ''', (template_id, report_name, report_content))
        
        self.conn.commit()
        return report_content
    
    def get_environmental_trend(self, station_name, days=7):
        """获取环境趋势数据"""
        query = '''
        SELECT 
            DATE(ed.collection_time) AS trend_date,
            AVG(ed.temperature) AS avg_temperature,
            AVG(ed.humidity) AS avg_humidity,
            SUM(ed.rainfall) AS total_rainfall
        FROM environmental_data ed
        JOIN monitoring_station ms ON ed.station_id = ms.station_id
        WHERE ms.station_name = %s AND DATE(ed.collection_time) >= DATE_SUB(CURDATE(), INTERVAL %s DAY)
        GROUP BY DATE(ed.collection_time)
        ORDER BY trend_date
        '''
        self.cursor.execute(query, (station_name, days))
        return self.cursor.fetchall()
    
    def get_air_quality(self, station_name, days=7):
        """获取空气质量数据"""
        query = '''
        SELECT 
            DATE(ed.collection_time) AS air_date,
            AVG(ed.pm25) AS avg_pm25,
            AVG(ed.pm10) AS avg_pm10
        FROM environmental_data ed
        JOIN monitoring_station ms ON ed.station_id = ms.station_id
        WHERE ms.station_name = %s AND DATE(ed.collection_time) >= DATE_SUB(CURDATE(), INTERVAL %s DAY)
        GROUP BY DATE(ed.collection_time)
        ORDER BY air_date
        '''
        self.cursor.execute(query, (station_name, days))
        return self.cursor.fetchall()
    
    def get_abnormal_data_statistics(self, days=30):
        """获取异常数据统计"""
        query = '''
        SELECT 
            ms.station_name,
            ad.abnormal_type,
            COUNT(ad.abnormal_id) AS abnormal_count
        FROM abnormal_data ad
        JOIN environmental_data ed ON ad.data_id = ed.data_id
        JOIN monitoring_station ms ON ed.station_id = ms.station_id
        WHERE DATE(ad.discovery_time) >= DATE_SUB(CURDATE(), INTERVAL %s DAY)
        GROUP BY ms.station_name, ad.abnormal_type
        ORDER BY ms.station_name, abnormal_count DESC
        '''
        self.cursor.execute(query, (days,))
        return self.cursor.fetchall()
    
    def print_menu(self):
        """打印菜单"""
        print("=" * 60)
        print("智慧林草系统 - 环境监测和统计分析模块")
        print("=" * 60)
        print("1. 查看实时数据")
        print("2. 查看日统计数据")
        print("3. 生成日报")
        print("4. 查看周统计数据")
        print("5. 生成周报")
        print("6. 查看月统计数据")
        print("7. 生成月报")
        print("8. 查看季度统计数据")
        print("9. 生成季报")
        print("10. 查看年度统计数据")
        print("11. 生成年报")
        print("12. 查看环境趋势")
        print("13. 查看空气质量")
        print("14. 查看异常数据统计")
        print("15. 执行自定义SQL查询")
        print("16. 查看环境监测站列表")
        print("17. 查看系统用户列表")
        print("0. 退出系统")
        print("=" * 60)
    
    def run(self):
        """运行系统"""
        while True:
            self.print_menu()
            choice = input("请输入您的选择：")
            
            if choice == "0":
                print("\n感谢使用智慧林草系统，再见！")
                self.conn.close()
                break
            elif choice == "1":
                print("\n=== 实时数据 ===")
                realtime_data = self.get_realtime_data()
                if realtime_data:
                    print("\n{:<15} {:<15} {:<19} {:<10} {:<10} {:<10} {:<10} {:<15}".format(
                        "站点名称", "站点类型", "采集时间", "温度", "湿度", "风速", "风向", "降雨量"
                    ))
                    print("-" * 110)
                    for data in realtime_data:
                        # 格式化日期时间
                        collection_time = data['collection_time'].strftime('%Y-%m-%d %H:%M:%S') if data['collection_time'] else "-"
                        print("{:<15} {:<15} {:<19} {:<10} {:<10} {:<10} {:<10} {:<15}".format(
                            data['station_name'], data['station_type'], collection_time,
                            f"{data['temperature']:.2f}℃" if data['temperature'] else "-",
                            f"{data['humidity']:.2f}%" if data['humidity'] else "-",
                            f"{data['wind_speed']:.2f}m/s" if data['wind_speed'] else "-",
                            data['wind_direction'] if data['wind_direction'] else "-",
                            f"{data['rainfall']:.2f}mm" if data['rainfall'] else "-"
                        ))
                else:
                    print("\n暂无实时数据")
            elif choice == "2":
                date = input("请输入要查询的日期（格式：YYYY-MM-DD）：")
                print(f"\n=== {date} 日统计数据 ===")
                statistics = self.get_daily_statistics(date)
                if statistics:
                    print("\n{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}".format(
                        "站点名称", "数据条数", "平均温度", "最高温度", "最低温度", "平均湿度", "总降雨量"
                    ))
                    print("-" * 85)
                    for stat in statistics:
                        print("{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}".format(
                            stat['station_name'], stat['data_count'],
                            f"{round(stat['avg_temperature'], 2) if stat['avg_temperature'] else 0}℃",
                            f"{round(stat['max_temperature'], 2) if stat['max_temperature'] else 0}℃",
                            f"{round(stat['min_temperature'], 2) if stat['min_temperature'] else 0}℃",
                            f"{round(stat['avg_humidity'], 2) if stat['avg_humidity'] else 0}%",
                            f"{round(stat['total_rainfall'], 2) if stat['total_rainfall'] else 0}mm"
                        ))
                else:
                    print(f"\n暂无 {date} 的统计数据")
            elif choice == "3":
                date = input("请输入要生成报表的日期（格式：YYYY-MM-DD）：")
                print(f"\n正在生成 {date} 日报...")
                report = self.generate_daily_report(date)
                print(f"\n=== 环境监测日报_{date} ===")
                print(report)
            elif choice == "4":
                year = input("请输入要查询的年份：")
                week = input("请输入要查询的周数：")
                print(f"\n=== {year}年第{week}周统计数据 ===")
                statistics = self.get_weekly_statistics(year, week)
                if statistics:
                    print("\n{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}".format(
                        "站点名称", "数据条数", "平均温度", "最高温度", "最低温度", "平均湿度", "总降雨量"
                    ))
                    print("-" * 85)
                    for stat in statistics:
                        print("{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}".format(
                            stat['station_name'], stat['data_count'],
                            f"{round(stat['avg_temperature'], 2) if stat['avg_temperature'] else 0}℃",
                            f"{round(stat['max_temperature'], 2) if stat['max_temperature'] else 0}℃",
                            f"{round(stat['min_temperature'], 2) if stat['min_temperature'] else 0}℃",
                            f"{round(stat['avg_humidity'], 2) if stat['avg_humidity'] else 0}%",
                            f"{round(stat['total_rainfall'], 2) if stat['total_rainfall'] else 0}mm"
                        ))
                else:
                    print(f"\n暂无 {year}年第{week}周的统计数据")
            elif choice == "5":
                year = input("请输入要生成报表的年份：")
                week = input("请输入要生成报表的周数：")
                print(f"\n正在生成 {year}年第{week}周周报...")
                report = self.generate_weekly_report(year, week)
                print(f"\n=== 环境监测周报_{year}_第{week}周 ===")
                print(report)
            elif choice == "6":
                year = input("请输入要查询的年份：")
                month = input("请输入要查询的月份：")
                print(f"\n=== {year}年{month}月统计数据 ===")
                statistics = self.get_monthly_statistics(year, month)
                if statistics:
                    print("\n{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}".format(
                        "站点名称", "数据条数", "平均温度", "最高温度", "最低温度", "平均湿度", "总降雨量"
                    ))
                    print("-" * 85)
                    for stat in statistics:
                        print("{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}".format(
                            stat['station_name'], stat['data_count'],
                            f"{round(stat['avg_temperature'], 2) if stat['avg_temperature'] else 0}℃",
                            f"{round(stat['max_temperature'], 2) if stat['max_temperature'] else 0}℃",
                            f"{round(stat['min_temperature'], 2) if stat['min_temperature'] else 0}℃",
                            f"{round(stat['avg_humidity'], 2) if stat['avg_humidity'] else 0}%",
                            f"{round(stat['total_rainfall'], 2) if stat['total_rainfall'] else 0}mm"
                        ))
                else:
                    print(f"\n暂无 {year}年{month}月的统计数据")
            elif choice == "7":
                year = input("请输入要生成报表的年份：")
                month = input("请输入要生成报表的月份：")
                print(f"\n正在生成 {year}年{month}月月报...")
                report = self.generate_monthly_report(year, month)
                print(f"\n=== 环境监测月报_{year}_{month} ===")
                print(report)
            elif choice == "8":
                year = input("请输入要查询的年份：")
                quarter = input("请输入要查询的季度（1-4）：")
                print(f"\n=== {year}年第{quarter}季度统计数据 ===")
                statistics = self.get_quarterly_statistics(year, quarter)
                if statistics:
                    print("\n{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}".format(
                        "站点名称", "数据条数", "平均温度", "最高温度", "最低温度", "平均湿度", "总降雨量"
                    ))
                    print("-" * 85)
                    for stat in statistics:
                        print("{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}".format(
                            stat['station_name'], stat['data_count'],
                            f"{round(stat['avg_temperature'], 2) if stat['avg_temperature'] else 0}℃",
                            f"{round(stat['max_temperature'], 2) if stat['max_temperature'] else 0}℃",
                            f"{round(stat['min_temperature'], 2) if stat['min_temperature'] else 0}℃",
                            f"{round(stat['avg_humidity'], 2) if stat['avg_humidity'] else 0}%",
                            f"{round(stat['total_rainfall'], 2) if stat['total_rainfall'] else 0}mm"
                        ))
                else:
                    print(f"\n暂无 {year}年第{quarter}季度的统计数据")
            elif choice == "9":
                year = input("请输入要生成报表的年份：")
                quarter = input("请输入要生成报表的季度（1-4）：")
                print(f"\n正在生成 {year}年第{quarter}季度季报...")
                report = self.generate_quarterly_report(year, quarter)
                print(f"\n=== 环境监测季报_{year}_第{quarter}季度 ===")
                print(report)
            elif choice == "10":
                year = input("请输入要查询的年份：")
                print(f"\n=== {year}年统计数据 ===")
                statistics = self.get_annual_statistics(year)
                if statistics:
                    print("\n{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}".format(
                        "站点名称", "数据条数", "平均温度", "最高温度", "最低温度", "平均湿度", "总降雨量"
                    ))
                    print("-" * 85)
                    for stat in statistics:
                        print("{:<15} {:<10} {:<10} {:<10} {:<10} {:<15} {:<15}".format(
                            stat['station_name'], stat['data_count'],
                            f"{round(stat['avg_temperature'], 2) if stat['avg_temperature'] else 0}℃",
                            f"{round(stat['max_temperature'], 2) if stat['max_temperature'] else 0}℃",
                            f"{round(stat['min_temperature'], 2) if stat['min_temperature'] else 0}℃",
                            f"{round(stat['avg_humidity'], 2) if stat['avg_humidity'] else 0}%",
                            f"{round(stat['total_rainfall'], 2) if stat['total_rainfall'] else 0}mm"
                        ))
                else:
                    print(f"\n暂无 {year}年的统计数据")
            elif choice == "11":
                year = input("请输入要生成报表的年份：")
                print(f"\n正在生成 {year}年年报...")
                report = self.generate_annual_report(year)
                print(f"\n=== 环境监测年报_{year} ===")
                print(report)
            elif choice == "12":
                station_name = input("请输入要查询的站点名称：")
                days = input("请输入要查询的天数（默认7天）：") or "7"
                print(f"\n=== {station_name} 近{days}天环境趋势 ===")
                trend = self.get_environmental_trend(station_name, int(days))
                if trend:
                    print("\n{:<15} {:<10} {:<10} {:<15}".format(
                        "日期", "平均温度", "平均湿度", "总降雨量"
                    ))
                    print("-" * 50)
                    for data in trend:
                        print("{:<15} {:<10} {:<10} {:<15}".format(
                            data['trend_date'],
                            f"{round(data['avg_temperature'], 2) if data['avg_temperature'] else 0}℃",
                            f"{round(data['avg_humidity'], 2) if data['avg_humidity'] else 0}%",
                            f"{round(data['total_rainfall'], 2) if data['total_rainfall'] else 0}mm"
                        ))
                else:
                    print(f"\n暂无 {station_name} 的环境趋势数据")
            elif choice == "13":
                station_name = input("请输入要查询的站点名称：")
                days = input("请输入要查询的天数（默认7天）：") or "7"
                print(f"\n=== {station_name} 近{days}天空气质量 ===")
                air_quality = self.get_air_quality(station_name, int(days))
                if air_quality:
                    print("\n{:<15} {:<10} {:<10}".format(
                        "日期", "平均PM2.5", "平均PM10"
                    ))
                    print("-" * 35)
                    for data in air_quality:
                        print("{:<15} {:<10} {:<10}".format(
                            data['air_date'],
                            f"{round(data['avg_pm25'], 2) if data['avg_pm25'] else 0}μg/m³",
                            f"{round(data['avg_pm10'], 2) if data['avg_pm10'] else 0}μg/m³"
                        ))
                else:
                    print(f"\n暂无 {station_name} 的空气质量数据")
            elif choice == "14":
                days = input("请输入要查询的天数（默认30天）：") or "30"
                print(f"\n=== 近{days}天异常数据统计 ===")
                abnormal_stats = self.get_abnormal_data_statistics(int(days))
                if abnormal_stats:
                    print("\n{:<15} {:<20} {:<15}".format(
                        "站点名称", "异常类型", "异常数量"
                    ))
                    print("-" * 50)
                    for data in abnormal_stats:
                        print("{:<15} {:<20} {:<15}".format(
                            data['station_name'], data['abnormal_type'], data['abnormal_count']
                        ))
                else:
                    print(f"\n近{days}天暂无异常数据")
            elif choice == "15":
                print("\n=== 执行自定义SQL查询 ===")
                print("注意：请确保输入的SQL查询安全，避免删除或修改关键数据")
                print("示例查询：SELECT * FROM monitoring_station LIMIT 5;")
                sql = input("请输入SQL查询语句：")
                try:
                    self.cursor.execute(sql)
                    results = self.cursor.fetchall()
                    if results:
                        # 打印列名
                        columns = results[0].keys()
                        print("\n{:<20} " * len(columns).format(*columns))
                        print("-" * (20 * len(columns)))
                        # 打印数据
                        for row in results:
                            print("{:<20} " * len(row).format(*row.values()))
                    else:
                        print("\n查询结果为空")
                except pymysql.Error as e:
                    print(f"\nSQL执行错误：{e}")
            elif choice == "16":
                print("\n=== 环境监测站列表 ===")
                self.cursor.execute('SELECT * FROM monitoring_station ORDER BY station_id')
                stations = self.cursor.fetchall()
                print("\n{:<10} {:<15} {:<10} {:<15} {:<15} {:<10}".format(
                    "站点ID", "站点名称", "站点类型", "安装日期", "纬度", "经度"
                ))
                print("-" * 80)
                for station in stations:
                    print("{:<10} {:<15} {:<10} {:<15} {:<15} {:<10}".format(
                        station['station_id'], station['station_name'], station['station_type'], 
                        station['installation_date'], station['latitude'], station['longitude']
                    ))
            elif choice == "17":
                print("\n=== 系统用户列表 ===")
                self.cursor.execute('SELECT * FROM system_user ORDER BY user_id')
                users = self.cursor.fetchall()
                print("\n{:<10} {:<15} {:<10} {:<20} {:<15}".format(
                    "用户ID", "用户名", "角色", "邮箱", "电话"
                ))
                print("-" * 70)
                for user in users:
                    print("{:<10} {:<15} {:<10} {:<20} {:<15}".format(
                        user['user_id'], user['username'], user['role'], 
                        user['email'], user['phone']
                    ))
            else:
                print("\n无效的选择，请重新输入")
            
            input("\n按回车键继续...")

def main():
    """主函数"""
    print("=== 智慧林草系统 - 环境监测和统计分析模块 ===")
    print(f"Python版本: {sys.version}")
    print(f"当前目录: {os.getcwd()}")
    print(f"脚本路径: {os.path.abspath(__file__)}")
    print("========================================")
    
    try:
        # 创建系统实例，使用MySQL数据库
        # 注意：请根据实际情况修改数据库连接参数
        system = EnvironmentalMonitoringSystem(
            host='localhost',
            user='root',
            password='mp050410',
            database='smart_forest_grass',
            port=3306
        )
        
        # 运行系统
        system.run()
    except Exception as e:
        print(f"系统运行错误: {e}")
        print(f"错误类型: {type(e).__name__}")
        import traceback
        traceback.print_exc()
        exit(1)

if __name__ == "__main__":
    main()
