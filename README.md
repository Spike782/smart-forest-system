# 智慧林草系统 - 自动部署指南

## 项目简介

智慧林草系统是一个基于Go语言开发的现代化林草管理平台，集成了环境监测、灾害预警、资源管理、设备管理、统计分析等功能。

## 环境要求

- **操作系统**: Windows 10/11
- **Go版本**: 1.19+
- **MySQL**: 8.0+
- **浏览器**: Chrome 90+ / Firefox 88+ / Edge 90+

## 快速部署

### 步骤1: 修改配置文件

在运行部署脚本前，请先修改数据库配置文件：

**Windows:**
```bash
# 编辑数据库配置文件
notepad backend/db/config.yaml
```

**Linux/Mac:**
```bash
# 编辑数据库配置文件
nano backend/db/config.yaml
# 或
vim backend/db/config.yaml
```

配置文件内容示例：
```yaml
database:
  user: root                    # 数据库用户名
  password: your_password       # 数据库密码
  host: 127.0.0.1              # 数据库主机地址
  port: 3306                   # 数据库端口
  name: smart_forest           # 数据库名称
  extra_params: "charset=utf8mb4&parseTime=true&loc=Local"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600s
  conn_max_idle_time: 600s
```

### 步骤2: 部署数据库

登录`mysql`数据库，运行项目目录下`/database`文件夹中的`init.sql`文件，会自动创建数据库`smart_forest`

可选：运行项目目录下`/database`文件夹中的`test_data.sql`文件，会自动向数据库中插入测试数据

### 步骤4：启动后端服务

1.进入项目的`/backend`文件夹，执行`build.bat`脚本自动编译

2.执行编译成功后的server.exe启动后端服务器

### 步骤3: 访问系统

部署完成后，在浏览器中访问：
```
http://localhost:8080
```

## 功能模块

- **环境监测**: 传感器数据采集与展示
- **灾害预警**: 实时预警与通知管理
- **资源管理**: 林草资源统计与维护
- **设备管理**: 设备状态监控与维护记录
- **统计分析**: 数据报表生成与导出
- **用户管理**: 权限控制与用户认证

## 数据库结构

系统使用MySQL数据库，主要包含以下表：
- users - 用户表
- roles - 角色表
- permissions - 权限表
- regions - 区域表
- sensors - 传感器表
- sensor_readings - 传感器读数表
- alert_rules - 预警规则表
- alerts - 预警记录表
- resources - 资源表
- devices - 设备表

## 开发与调试

### 单独运行数据库初始化
```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE smart_forest CHARACTER SET utf8mb4;"

# 初始化表结构
Get-Content database/init.sql -Raw | mysql -u root -p --default-character-set=utf8mb4 smart_forest

# 导入测试数据
Get-Content database/test_data.sql -Raw | mysql -u root -p --default-character-set=utf8mb4 smart_forest
```

### 单独编译和运行
```bash
cd backend
go build -o server.exe api/server.go
server.exe
```

### 运行测试

持久层测试：

```bash
cd backend
go test ./test/persistence.go
```

## 注意事项

1. **数据库配置**: 确保MySQL服务正在运行，且连接信息正确
2. **端口占用**: 默认使用8080端口，确保该端口未被其他程序占用
3. **权限问题**: 确保当前用户有足够权限访问数据库和文件系统
4. **字符编码**: 系统使用UTF-8编码，确保终端支持中文显示

## 技术栈

- **后端**: Go + Gin框架
- **数据库**: MySQL 8.0
- **前端**: HTML5 + JavaScript + Chart.js
- **认证**: JWT
- **部署**: Batch脚本自动化

