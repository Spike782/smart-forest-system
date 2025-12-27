DROP DATABASE IF EXISTS smart_forest;
CREATE DATABASE smart_forest DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE smart_forest;

-- -----------------------------
-- 1) 用户/角色（用于登录+权限控制）
-- -----------------------------
CREATE TABLE t_user (
  user_id      VARCHAR(20) PRIMARY KEY,
  username     VARCHAR(50) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role         VARCHAR(20) NOT NULL COMMENT 'admin / ranger / auditor / public',
  failed_count INT NOT NULL DEFAULT 0,
  locked_until DATETIME NULL,
  create_time  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表（简化RBAC）';

-- -----------------------------
-- 2) 区域（预警记录需要 area_id；区域也可绑定负责人）
-- -----------------------------
CREATE TABLE t_area (
  area_id     VARCHAR(20) PRIMARY KEY,
  area_name   VARCHAR(50) NOT NULL,
  area_type   VARCHAR(20) NOT NULL COMMENT 'forest/grassland',
  manager_id  VARCHAR(20) NOT NULL,
  FOREIGN KEY (manager_id) REFERENCES t_user(user_id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='区域信息';

-- -----------------------------
-- 3) 监测数据（插入后触发预警）
-- -----------------------------
CREATE TABLE t_monitoring_data (
  data_id     BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  area_id     VARCHAR(20) NOT NULL,
  collect_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  temperature DECIMAL(6,2) NULL,
  humidity    DECIMAL(6,2) NULL,
  data_status TINYINT NOT NULL DEFAULT 1 COMMENT '1-有效 0-无效',
  INDEX idx_area_time(area_id, collect_time),
  FOREIGN KEY (area_id) REFERENCES t_area(area_id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='环境监测数据（简化：温湿度）';

-- =========================================================
-- 4) 你文档里的三张核心表（规则/预警记录/通知）
--    (结构按你提供的DDL：外键/索引/字段含义一致)
-- =========================================================

-- 4.1 预警规则表 t_warning_rule
CREATE TABLE `t_warning_rule` (
  `rule_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '规则编号（主键）',
  `warning_type` VARCHAR(20) NOT NULL COMMENT '预警类型（火灾/旱情/病虫害等）',
  `trigger_condition` VARCHAR(100) NOT NULL COMMENT '触发条件（建议存JSON字符串，如：{"temp_gte":38,"hum_lte":20}）',
  `warning_level` TINYINT NOT NULL COMMENT '预警级别（1-一般/2-较重/3-严重/4-特别严重）',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '生效状态（0-禁用/1-启用）',
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`rule_id`),
  INDEX `idx_warning_type` (`warning_type`) COMMENT '按预警类型查询索引',
  INDEX `idx_status` (`status`) COMMENT '按生效状态查询索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='预警规则配置表';

-- 4.2 预警记录表 t_warning_record
CREATE TABLE `t_warning_record` (
  `warning_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '预警编号（主键）',
  `rule_id` INT UNSIGNED NOT NULL COMMENT '触发规则编号（外键关联t_warning_rule）',
  `area_id` VARCHAR(20) NOT NULL COMMENT '涉及区域编号（关联区域信息表）',
  `trigger_time` DATETIME NOT NULL COMMENT '触发时间',
  `warning_content` VARCHAR(200) NOT NULL COMMENT '预警内容',
  `handle_status` TINYINT NOT NULL DEFAULT 0 COMMENT '处理状态（0-未处理/1-处理中/2-已结案）',
  `handler_id` VARCHAR(20) DEFAULT NULL COMMENT '处理人ID',
  `handle_result` VARCHAR(500) DEFAULT NULL COMMENT '处理结果',
  `handle_time` DATETIME DEFAULT NULL COMMENT '处理完成时间',
  PRIMARY KEY (`warning_id`),
  FOREIGN KEY (`rule_id`) REFERENCES `t_warning_rule` (`rule_id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  INDEX `idx_area_time` (`area_id`, `trigger_time`) COMMENT '按区域+时间查询索引',
  INDEX `idx_rule_id` (`rule_id`) COMMENT '按规则编号查询索引',
  INDEX `idx_handle_status` (`handle_status`) COMMENT '按处理状态查询索引',
  FOREIGN KEY (`area_id`) REFERENCES `t_area` (`area_id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='预警触发记录表';

-- 4.3 通知记录表 t_notification_record
CREATE TABLE `t_notification_record` (
  `notification_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '通知编号（主键）',
  `warning_id` BIGINT UNSIGNED NOT NULL COMMENT '关联预警编号（外键关联t_warning_record）',
  `receiver_id` VARCHAR(20) NOT NULL COMMENT '接收人ID',
  `notify_method` TINYINT NOT NULL COMMENT '通知方式（1-短信/2-系统消息）',
  `send_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
  `receive_status` TINYINT NOT NULL DEFAULT 0 COMMENT '接收状态（0-未接收/1-已接收/2-已读）',
  `receive_time` DATETIME DEFAULT NULL COMMENT '接收时间',
  PRIMARY KEY (`notification_id`),
  FOREIGN KEY (`warning_id`) REFERENCES `t_warning_record` (`warning_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  INDEX `idx_warning_id` (`warning_id`) COMMENT '按预警编号查询通知',
  INDEX `idx_receiver_status` (`receiver_id`, `receive_status`) COMMENT '按接收人+状态查询索引',
  FOREIGN KEY (`receiver_id`) REFERENCES `t_user` (`user_id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='预警通知记录表';


-- =========================================================
-- 5) 视图（灾害预警业务线至少3个）
-- =========================================================

-- 5.1 启用规则视图
CREATE OR REPLACE VIEW v_active_rules AS
SELECT rule_id, warning_type, trigger_condition, warning_level, status, update_time
FROM t_warning_rule
WHERE status = 1;

-- 5.2 最近预警视图（含规则信息）
CREATE OR REPLACE VIEW v_recent_warnings AS
SELECT
  wr.warning_id, wr.area_id, a.area_name,
  wr.trigger_time, wr.warning_content,
  wr.handle_status, wr.handler_id, wr.handle_time,
  r.warning_type, r.warning_level
FROM t_warning_record wr
JOIN t_warning_rule r ON r.rule_id = wr.rule_id
JOIN t_area a ON a.area_id = wr.area_id;

-- 5.3 未读/未接收通知视图
CREATE OR REPLACE VIEW v_pending_notifications AS
SELECT
  n.notification_id, n.warning_id, n.receiver_id,
  n.notify_method, n.send_time, n.receive_status
FROM t_notification_record n
WHERE n.receive_status IN (0,1);

-- =========================================================
-- 6) 存储过程：插入监测数据后，自动比对规则并写预警+通知
--    说明：先实现“火灾”类型：温度>=temp_gte 且 湿度<=hum_lte
-- =========================================================
DELIMITER $$

DROP PROCEDURE IF EXISTS sp_check_and_create_warning $$
CREATE PROCEDURE sp_check_and_create_warning(
  IN p_area_id VARCHAR(20),
  IN p_temp DECIMAL(6,2),
  IN p_hum  DECIMAL(6,2),
  IN p_time DATETIME
)
BEGIN
  DECLARE done INT DEFAULT 0;
  DECLARE v_rule_id INT;
  DECLARE v_type VARCHAR(20);
  DECLARE v_cond VARCHAR(100);
  DECLARE v_level TINYINT;

  DECLARE v_json JSON;
  DECLARE v_temp_gte DECIMAL(6,2);
  DECLARE v_hum_lte  DECIMAL(6,2);

  DECLARE v_warning_id BIGINT;

  DECLARE cur CURSOR FOR
    SELECT rule_id, warning_type, trigger_condition, warning_level
    FROM t_warning_rule
    WHERE status = 1;

  DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;

  OPEN cur;

  read_loop: LOOP
    FETCH cur INTO v_rule_id, v_type, v_cond, v_level;
    IF done = 1 THEN
      LEAVE read_loop;
    END IF;

    -- 仅演示火灾：{"temp_gte":38,"hum_lte":20}
    IF v_type = '火灾' THEN
      SET v_json = CAST(v_cond AS JSON);
      SET v_temp_gte = CAST(JSON_UNQUOTE(JSON_EXTRACT(v_json,'$.temp_gte')) AS DECIMAL(6,2));
      SET v_hum_lte  = CAST(JSON_UNQUOTE(JSON_EXTRACT(v_json,'$.hum_lte'))  AS DECIMAL(6,2));

      IF p_temp IS NOT NULL AND p_hum IS NOT NULL AND p_temp >= v_temp_gte AND p_hum <= v_hum_lte THEN

        INSERT INTO t_warning_record(rule_id, area_id, trigger_time, warning_content, handle_status)
        VALUES(
          v_rule_id,
          p_area_id,
          p_time,
          CONCAT('区域', p_area_id, ' 温度=', p_temp, '℃ 湿度=', p_hum, '% 触发【火灾】预警，级别=', v_level),
          0
        );
        SET v_warning_id = LAST_INSERT_ID();

        -- 通知：发给区域负责人 + 所有admin
        INSERT INTO t_notification_record(warning_id, receiver_id, notify_method, receive_status)
        SELECT v_warning_id, a.manager_id, 2, 0
        FROM t_area a WHERE a.area_id = p_area_id;

        INSERT INTO t_notification_record(warning_id, receiver_id, notify_method, receive_status)
        SELECT v_warning_id, u.user_id, 2, 0
        FROM t_user u
        WHERE u.role = 'admin';

      END IF;
    END IF;

  END LOOP;

  CLOSE cur;
END $$

-- =========================================================
-- 7) 触发器：监测数据插入后自动调用存储过程
-- =========================================================
DROP TRIGGER IF EXISTS trg_monitor_insert $$
CREATE TRIGGER trg_monitor_insert
AFTER INSERT ON t_monitoring_data
FOR EACH ROW
BEGIN
  IF NEW.data_status = 1 THEN
    CALL sp_check_and_create_warning(NEW.area_id, NEW.temperature, NEW.humidity, NEW.collect_time);
  END IF;
END $$

DELIMITER ;

-- =========================================================
-- 8) 测试数据（账号/区域/规则）
--    账号密码见下：admin123 / ranger123
-- =========================================================


INSERT INTO t_user(user_id, username, password_hash, role) VALUES
('u_admin',  'admin',  '$2b$12$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin'),
('u_ranger', 'ranger', '$2b$12$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'ranger');

INSERT INTO t_area(area_id, area_name, area_type, manager_id) VALUES
('A001','北区森林','forest','u_ranger'),
('A002','南区草地','grassland','u_ranger');

-- 火灾规则：温度>=38 且 湿度<=20
INSERT INTO t_warning_rule(warning_type, trigger_condition, warning_level, status)
VALUES
('火灾', '{"temp_gte":38,"hum_lte":20}', 3, 1),
('火灾', '{"temp_gte":41,"hum_lte":18}', 4, 1);

-- =========================================================
-- 用来快速验证触发器是否工作：
-- INSERT INTO t_monitoring_data(area_id, temperature, humidity, data_status) VALUES ('A001', 39, 19, 1);
-- 然后查：
-- SELECT * FROM t_warning_record ORDER BY warning_id DESC;
-- SELECT * FROM t_notification_record ORDER BY notification_id DESC;
-- =========================================================
