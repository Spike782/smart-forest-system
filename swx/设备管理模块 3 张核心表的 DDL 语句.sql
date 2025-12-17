CREATE TABLE [User] (
    user_id VARCHAR(20) NOT NULL, 
    user_name VARCHAR(30) NOT NULL, 
    role VARCHAR(20) NOT NULL, 
    password VARCHAR(50) NOT NULL, 
    phone VARCHAR(20), 
    create_time DATETIME DEFAULT GETDATE(), 
    PRIMARY KEY (user_id),
    CHECK (role IN ('管理员', '护林员', '维护人员'))
);
GO

CREATE TABLE Area (
    area_id VARCHAR(20) NOT NULL, 
    area_name VARCHAR(50) NOT NULL, 
    area_type VARCHAR(20) NOT NULL, 
    manager_id VARCHAR(20) NOT NULL, 
    area_address VARCHAR(100), 
    PRIMARY KEY (area_id),
    -- 取消级联更新，消除循环警告
    FOREIGN KEY (manager_id) REFERENCES [User](user_id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CHECK (area_type IN ('森林', '草地'))
);
GO

-- 1.1 设备档案表（DeviceArchive）
CREATE TABLE DeviceArchive (
    device_id VARCHAR(20) NOT NULL, 
    device_name VARCHAR(50) NOT NULL, 
    device_type VARCHAR(20) NOT NULL, 
    model VARCHAR(30) NOT NULL, 
    purchase_time DATE NOT NULL, 
    area_id VARCHAR(20) NOT NULL, 
    installer_id VARCHAR(20) NOT NULL, 
    warranty_period DATE NOT NULL, 
    create_time DATETIME NOT NULL, 
    PRIMARY KEY (device_id),
    -- 关键：取消级联更新，用NO ACTION替代CASCADE，消除循环警告
    FOREIGN KEY (area_id) REFERENCES Area(area_id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    FOREIGN KEY (installer_id) REFERENCES [User](user_id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CHECK (device_type IN ('传感器', '摄像头', '预警器'))
);
GO

-- 1.2 设备状态数据表（DeviceStatus）
CREATE TABLE DeviceStatus (
    status_id INT IDENTITY(1,1) NOT NULL, 
    device_id VARCHAR(20) NOT NULL, 
    collect_time DATETIME NOT NULL, 
    run_status VARCHAR(20) NOT NULL, 
    battery_level INT, 
    signal_strength INT NOT NULL, 
    PRIMARY KEY (status_id),
    FOREIGN KEY (device_id) REFERENCES DeviceArchive(device_id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CHECK (run_status IN ('正常', '故障', '离线')),
    CHECK (battery_level IS NULL OR battery_level BETWEEN 0 AND 100),
    CHECK (signal_strength BETWEEN 0 AND 100)
);
GO

-- 1.3 设备维护记录表（DeviceMaintenance）
CREATE TABLE DeviceMaintenance (
    maintenance_id VARCHAR(20) NOT NULL, 
    device_id VARCHAR(20) NOT NULL, 
    maintenance_type VARCHAR(20) NOT NULL, 
    maintenance_time DATETIME NOT NULL, 
    maintainer_id VARCHAR(20) NOT NULL, 
    maintenance_content VARCHAR(500) NOT NULL, 
    maintenance_result VARCHAR(20) NOT NULL, 
    PRIMARY KEY (maintenance_id),
    -- 关键：取消级联更新
    FOREIGN KEY (device_id) REFERENCES DeviceArchive(device_id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    FOREIGN KEY (maintainer_id) REFERENCES [User](user_id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CHECK (maintenance_type IN ('巡检', '维修', '更换')),
    CHECK (maintenance_result IN ('成功', '待跟进', '失败'))
);
GO

-- 1.4 索引创建（最后执行）
CREATE NONCLUSTERED INDEX IX_DeviceArchive_AreaId ON DeviceArchive(area_id);
CREATE NONCLUSTERED INDEX IX_DeviceArchive_DeviceType ON DeviceArchive(device_type);
CREATE NONCLUSTERED INDEX IX_DeviceStatus_DeviceId ON DeviceStatus(device_id);
CREATE NONCLUSTERED INDEX IX_DeviceStatus_RunStatus ON DeviceStatus(run_status);
CREATE NONCLUSTERED INDEX IX_DeviceStatus_CollectTime ON DeviceStatus(collect_time);
CREATE NONCLUSTERED INDEX IX_DeviceMaintenance_DeviceId ON DeviceMaintenance(device_id);
CREATE NONCLUSTERED INDEX IX_DeviceMaintenance_MaintenanceTime ON DeviceMaintenance(maintenance_time);
GO
