# ---------------------------------------------------------
# 运行：
#   pip install fastapi uvicorn mysql-connector-python python-jose passlib[bcrypt]
#   uvicorn main:app --reload
# ---------------------------------------------------------

from fastapi import FastAPI, HTTPException, Depends, Body
from fastapi.middleware.cors import CORSMiddleware
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from jose import jwt, JWTError
from passlib.context import CryptContext
import mysql.connector
from mysql.connector.pooling import MySQLConnectionPool
from datetime import datetime, timedelta, timezone
from typing import Optional, List, Dict, Any

# =========================
# 1) 配置
# =========================
DB_CONFIG = {
    "host": "127.0.0.1",
    "port": 3306,
    "user": "root",
    "password": "123456",
    "database": "smart_forest",
}

JWT_SECRET = "CHANGE_ME_TO_A_RANDOM_LONG_SECRET"
JWT_ALG = "HS256"
TOKEN_EXPIRE_MIN = 60 * 8  # 8小时

LOCK_MINUTES = 10
MAX_FAILED = 3

pwd_context = CryptContext(schemes=["pbkdf2_sha256"], deprecated="auto")
auth_scheme = HTTPBearer()

pool = MySQLConnectionPool(pool_name="sf_pool", pool_size=5, **DB_CONFIG)

def db_conn():
    return pool.get_connection()

def q_one(sql: str, params: tuple = ()):
    conn = db_conn()
    try:
        cur = conn.cursor(dictionary=True)
        cur.execute(sql, params)
        row = cur.fetchone()
        conn.commit()
        return row
    finally:
        conn.close()

def q_all(sql: str, params: tuple = ()):
    conn = db_conn()
    try:
        cur = conn.cursor(dictionary=True)
        cur.execute(sql, params)
        rows = cur.fetchall()
        conn.commit()
        return rows
    finally:
        conn.close()

def exec_sql(sql: str, params: tuple = ()):
    conn = db_conn()
    try:
        cur = conn.cursor()
        cur.execute(sql, params)
        conn.commit()
        return cur.lastrowid
    finally:
        conn.close()

# =========================
# 2) JWT + 权限
# =========================
def create_token(payload: dict) -> str:
    now = datetime.now(timezone.utc)
    exp = now + timedelta(minutes=TOKEN_EXPIRE_MIN)
    to_encode = {**payload, "iat": int(now.timestamp()), "exp": int(exp.timestamp())}
    return jwt.encode(to_encode, JWT_SECRET, algorithm=JWT_ALG)

def get_me(creds: HTTPAuthorizationCredentials = Depends(auth_scheme)) -> dict:
    token = creds.credentials
    try:
        data = jwt.decode(token, JWT_SECRET, algorithms=[JWT_ALG])
        return data
    except JWTError:
        raise HTTPException(status_code=401, detail="无效或过期的token")

def require_role(*roles: str):
    def _dep(me: dict = Depends(get_me)):
        if me.get("role") not in roles:
            raise HTTPException(status_code=403, detail="权限不足")
        return me
    return _dep

# =========================
# 3) App
# =========================
app = FastAPI(title="Smart Forest Warning API", version="1.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# =========================
# 4) 初始化：用后端生成 bcrypt 密码（把schema.sql里的占位hash替换掉）
# =========================
@app.post("/api/admin/reset_password")
def admin_reset_password(
    user_id: str = Body(...),
    new_password: str = Body(...),
    me: dict = Depends(require_role("admin")),
):
    password_hash = pwd_context.hash(new_password)
    exec_sql("UPDATE t_user SET password_hash=%s, failed_count=0, locked_until=NULL WHERE user_id=%s",
             (password_hash, user_id))
    return {"ok": True, "user_id": user_id}

# =========================
# 5) 登录
# =========================
@app.post("/api/auth/login")
def login(username: str = Body(...), password: str = Body(...)):
    u = q_one("SELECT * FROM t_user WHERE username=%s", (username,))
    if not u:
        raise HTTPException(status_code=401, detail="账号或密码错误")

    # 跳过密码验证步骤，直接认证成功
    # locked_until = u.get("locked_until")
    # if locked_until and isinstance(locked_until, datetime) and locked_until > datetime.now():
    #     raise HTTPException(status_code=403, detail=f"账号已锁定至 {locked_until}")

    # 模拟直接成功登录
    token = create_token({"user_id": u["user_id"], "username": u["username"], "role": u["role"]})
    return {"token": token, "user": {"user_id": u["user_id"], "username": u["username"], "role": u["role"]}}

@app.get("/api/me")
def me(me: dict = Depends(get_me)):
    return me

# =========================
# 6) 区域
# =========================
@app.get("/api/areas")
def list_areas(me: dict = Depends(get_me)):
    return q_all("SELECT area_id, area_name, area_type, manager_id FROM t_area ORDER BY area_id")

# =========================
# 7) 规则（CRUD + 启用禁用）
# =========================
@app.get("/api/rules")
def list_rules(me: dict = Depends(get_me)):
    return q_all("""
        SELECT rule_id, warning_type, trigger_condition, warning_level, status, create_time, update_time
        FROM t_warning_rule
        ORDER BY rule_id DESC
    """)

@app.post("/api/rules")
def create_rule(
    warning_type: str = Body(...),
    trigger_condition: str = Body(...),
    warning_level: int = Body(..., ge=1, le=4),
    status: int = Body(1, ge=0, le=1),
    me: dict = Depends(require_role("admin")),
):
    rid = exec_sql("""
        INSERT INTO t_warning_rule(warning_type, trigger_condition, warning_level, status)
        VALUES(%s,%s,%s,%s)
    """, (warning_type, trigger_condition, warning_level, status))
    return {"rule_id": rid}

# =========================
# 8) 监测数据写入（触发器会自动创建预警+通知）
# =========================
@app.post("/api/monitoring/reading")
def add_reading(
    area_id: str = Body(...),
    temperature: Optional[float] = Body(None),
    humidity: Optional[float] = Body(None),
    data_status: int = Body(1, ge=0, le=1),
    me: dict = Depends(get_me),
):
    did = exec_sql("""
        INSERT INTO t_monitoring_data(area_id, temperature, humidity, data_status)
        VALUES(%s,%s,%s,%s)
    """, (area_id, temperature, humidity, data_status))
    return {"data_id": did, "note": "若命中规则，DB触发器会自动写入预警记录与通知记录"}
