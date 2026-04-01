-- 视图1：v_DeviceStatusSummary（单独批处理）
CREATE VIEW v_DeviceStatusSummary
AS
SELECT
    da.device_type,
    COUNT(CASE WHEN ds.run_status = '正常' THEN 1 END) AS 正常设备数,
    COUNT(CASE WHEN ds.run_status = '故障' THEN 1 END) AS 故障设备数,
    COUNT(CASE WHEN ds.run_status = '离线' THEN 1 END) AS 离线设备数,
    COUNT(ds.device_id) AS 总设备数,
    -- 特殊字符列名用[]包裹
    ROUND(COUNT(CASE WHEN ds.run_status = '故障' THEN 1 END) * 100.0 / COUNT(ds.device_id), 2) AS [故障占比(%)]
FROM
    DeviceArchive da
LEFT JOIN DeviceStatus ds ON da.device_id = ds.device_id
    AND ds.collect_time = (SELECT MAX(collect_time) FROM DeviceStatus WHERE device_id = ds.device_id) -- 取最新状态
GROUP BY
    da.device_type;
GO -- 分隔批处理


-- 视图2：v_AreaDeviceDistribution（单独批处理）
CREATE VIEW v_AreaDeviceDistribution
AS
SELECT
    a.area_id,
    a.area_name,
    da.device_id,
    da.device_name,
    da.device_type,
    da.model,
    u.user_name AS 安装人员姓名,
    ds.run_status AS 最新运行状态,
    ds.signal_strength AS 最新信号强度,
    da.warranty_period AS 质保期截止日期
FROM
    Area a
JOIN DeviceArchive da ON a.area_id = da.area_id
JOIN [User] u ON da.installer_id = u.user_id
LEFT JOIN DeviceStatus ds ON da.device_id = ds.device_id
    AND ds.collect_time = (SELECT MAX(collect_time) FROM DeviceStatus WHERE device_id = ds.device_id);
GO -- 分隔批处理


-- 视图3：v_DeviceMaintenanceHistory（修正ORDER BY报错）
CREATE VIEW v_DeviceMaintenanceHistory
AS
SELECT TOP 100 PERCENT -- 新增这行，让ORDER BY在视图中生效
    da.device_id,
    da.device_name,
    dm.maintenance_id,
    dm.maintenance_type,
    dm.maintenance_time,
    u.user_name AS 维护人员姓名,
    dm.maintenance_content,
    dm.maintenance_result,
    ds.run_status AS 维护时设备状态
FROM
    DeviceArchive da
JOIN DeviceMaintenance dm ON da.device_id = dm.device_id
JOIN [User] u ON dm.maintainer_id = u.user_id
LEFT JOIN DeviceStatus ds ON da.device_id = ds.device_id
    AND ds.collect_time <= dm.maintenance_time
    AND ds.collect_time = (SELECT MAX(collect_time) FROM DeviceStatus WHERE device_id = da.device_id AND collect_time <= dm.maintenance_time)
ORDER BY
    dm.maintenance_time DESC;
GO