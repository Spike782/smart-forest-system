package middleware

import (
	"database/sql"
	"net/http"
	"strings"

	"smart-forest-system/backend/util"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization头中获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证令牌"})
			c.Abort()
			return
		}

		// 检查令牌格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌格式"})
			c.Abort()
			return
		}

		// 解析令牌
		claims, err := util.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// CheckPermissionMiddleware 权限检查中间件
func CheckPermissionMiddleware(db *sql.DB, resource string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证的用户"})
			c.Abort()
			return
		}

		// 检查用户是否有指定资源的指定操作权限
		var hasPermission bool
		query := `
			SELECT EXISTS(
				SELECT 1
				FROM user_roles ur
				JOIN role_permissions rp ON ur.role_id = rp.role_id
				JOIN permissions p ON rp.permission_id = p.permission_id
				WHERE ur.user_id = ? AND p.resource = ? AND p.action = ?
			) AS has_permission`
		err := db.QueryRow(query, userID, resource, action).Scan(&hasPermission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有操作权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}
