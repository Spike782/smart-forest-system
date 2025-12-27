智慧林草系统 - 灾害预警业务线

项目介绍

项目背景

森林和草地作为核心生态资源，传统管理模式依赖人工巡查与线下记录，存在数据滞后、预警不及时、管理效率低下等问题。智慧林草系统通过整合物联网、数据库与 API 技术，实现环境监测、灾害预警、资源管理的数字化闭环，为林草管理提供精准决策支持。

核心目标

实时比对监测数据与预警规则，自动触发灾害预警

动态管理预警规则（启用 / 禁用、多梯度配置）

自动向相关人员发送预警通知，跟踪接收状态

支持预警记录的全生命周期管理（未处理 / 处理中 / 已结案）

核心功能

基于角色的权限控制（管理员、巡护员、审计员等）

预警规则 CRUD 与状态管理

环境监测数据采集与存储

规则触发自动生成预警记录

多方式预警通知与状态跟踪

预警记录查询与处理结果反馈

环境准备

依赖软件

Python 3.8+

MySQL 8.0+

pip（Python 包管理工具）

Python 依赖包

bash

运行

pip install fastapi uvicorn mysql-connector-python python-jose passlib\[bcrypt]

快速部署

1\. 数据库初始化

登录 MySQL，执行schema.sql脚本创建数据库与测试数据：

bash

运行

mysql -u root -p < schema.sql

脚本执行完成后将自动创建：

数据库smart\_forest及所有业务表

3 个核心视图（启用规则、最近预警、未读通知）

1 个存储过程（规则校验与预警生成）

1 个触发器（监测数据插入后自动触发预警）

测试账号与测试数据

2\. 系统配置修改

编辑main.py中的核心配置（根据实际环境调整）：

python

运行

\# 数据库配置

DB\_CONFIG = {

&nbsp;   "host": "127.0.0.1",  # 数据库地址

&nbsp;   "port": 3306,         # 数据库端口

&nbsp;   "user": "root",       # 数据库用户名

&nbsp;   "password": "123456", # 数据库密码

&nbsp;   "database": "smart\_forest",

}



\# JWT密钥（务必修改为随机长字符串）

JWT\_SECRET = "CHANGE\_ME\_TO\_A\_RANDOM\_LONG\_SECRET"

3\. 启动服务

bash

运行

uvicorn main:app --reload

服务默认运行在 http://127.0.0.1:8000

--reload 开启热重载（开发环境使用）

API 文档

FastAPI 自带交互式 API 文档，启动服务后可访问：

Swagger UI：http://127.0.0.1:8000/docs（推荐，支持在线调试）

ReDoc：http://127.0.0.1:8000/redoc

核心 API 列表

接口路径	方法	功能描述	权限要求

/api/auth/login	POST	用户登录（返回 JWT 令牌）	无

/api/me	GET	获取当前登录用户信息	已认证

/api/admin/reset\_password	POST	重置用户密码	admin 角色

/api/areas	GET	查询所有区域信息	已认证

/api/rules	GET	查询所有预警规则	已认证

/api/rules	POST	创建预警规则	admin 角色

/api/monitoring/reading	POST	写入监测数据（自动触发预警）	已认证

数据库设计

核心业务表

1\. 预警规则表（t\_warning\_rule）

存储预警规则配置，支持动态启用 / 禁用与多梯度配置：

核心字段：规则编号、预警类型（火灾 / 旱情 / 病虫害）、触发条件（JSON 格式）、预警级别（1-4 级）、生效状态

索引：按预警类型、生效状态优化查询

2\. 预警记录表（t\_warning\_record）

记录每次预警触发详情：

核心字段：预警编号、触发规则 ID、区域 ID、触发时间、预警内容、处理状态、处理人、处理结果

外键关联：规则表（t\_warning\_rule）、区域表（t\_area）

索引：按区域 + 时间、处理状态优化查询

3\. 通知记录表（t\_notification\_record）

跟踪预警通知发送与接收状态：

核心字段：通知编号、预警 ID、接收人 ID、通知方式（短信 / 系统消息）、发送时间、接收状态

外键关联：预警表（t\_warning\_record）、用户表（t\_user）

索引：按预警 ID、接收人 + 状态优化查询

关键数据库对象

视图：

v\_active\_rules：仅查询启用状态的预警规则

v\_recent\_warnings：关联预警记录、规则、区域信息的联合查询

v\_pending\_notifications：查询未接收 / 未读的通知

存储过程：

sp\_check\_and\_create\_warning：校验监测数据是否满足规则，生成预警记录与通知

触发器：

trg\_monitor\_insert：监测数据插入后自动调用存储过程，实现预警自动触发

安全策略

1\. 身份认证与授权

JWT 认证：登录成功后生成有效期 8 小时的令牌，所有需授权接口均需携带令牌

基于角色的权限控制（RBAC）：

admin：全权限（创建规则、重置密码等）

ranger：巡护员（查看预警、处理预警等）

auditor：审计员（只读权限）

public：公共角色（有限访问权限）

2\. 登录安全

密码加密：使用pbkdf2\_sha256算法加密存储，不存储明文

登录保护：支持 3 次登录失败锁定 10 分钟（代码已预留逻辑，需启用）

会话超时：令牌有效期 8 小时，超时需重新登录

3\. 数据完整性

外键约束：核心表之间通过外键强关联，避免数据孤岛

事务支持：数据库操作默认开启事务，确保操作原子性

数据备份：建议定期备份smart\_forest数据库（推荐每日全量备份 + 增量备份）

测试指南

测试账号

用户名	密码	角色	权限说明

admin	admin123	admin	全权限操作

ranger	ranger123	ranger	查看预警、处理预警

功能测试步骤

登录获取令牌：

访问/api/auth/login，传入username=admin、password=admin123

记录返回的token（后续接口需在请求头携带Authorization: Bearer {token}）

创建预警规则：

访问/api/rules（POST），传入参数：

json

{

&nbsp; "warning\_type": "火灾",

&nbsp; "trigger\_condition": "{\\"temp\_gte\\":38,\\"hum\_lte\\":20}",

&nbsp; "warning\_level": 3,

&nbsp; "status": 1

}

写入监测数据触发预警：

访问/api/monitoring/reading（POST），传入参数：

json

{

&nbsp; "area\_id": "A001",

&nbsp; "temperature": 39.5,

&nbsp; "humidity": 18.2,

&nbsp; "data\_status": 1

}

验证预警结果：

执行 SQL 查询预警记录：

sql

SELECT \* FROM t\_warning\_record ORDER BY warning\_id DESC LIMIT 1;

执行 SQL 查询通知记录：

sql

SELECT \* FROM t\_notification\_record ORDER BY notification\_id DESC LIMIT 10;

注意事项

安全配置：

生产环境务必修改JWT\_SECRET为随机长字符串（建议至少 32 位）

生产环境禁用--reload参数，避免代码泄露

数据库密码建议使用环境变量注入，避免硬编码

未完成功能：

登录失败锁定逻辑已预留，需在/api/auth/login中启用注释代码

预警记录处理、通知状态更新接口需根据业务扩展

目前仅支持 "火灾" 类型预警，需扩展其他类型（旱情、病虫害）需修改存储过程

性能优化：

生产环境建议增大数据库连接池（pool\_size）

监测数据量较大时，建议对t\_monitoring\_data表按时间分区

扩展方向

支持更多预警类型（旱情、病虫害），扩展存储过程的规则校验逻辑

增加预警处理接口（更新处理状态、提交处理结果）

扩展通知方式（短信、邮件），集成第三方消息服务

增加数据可视化功能（预警趋势图、区域预警统计）

实现数据库定时备份脚本，支持一键恢复

增加日志记录功能，跟踪系统操作与异常



