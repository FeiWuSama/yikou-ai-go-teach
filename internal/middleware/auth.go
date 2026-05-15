package middleware

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
	"net/url"
	"yikou-ai-go-teach/internal/dal/model"
	"yikou-ai-go-teach/internal/dal/query"
	"yikou-ai-go-teach/pkg/constants"
	"yikou-ai-go-teach/pkg/enum"
	"yikou-ai-go-teach/pkg/errorutil"
)

// AuthMiddleware 鉴权中间件
func AuthMiddleware(roleEnum enum.UserRoleEnum, db *gorm.DB) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 校验权限
		var decodeUser []byte
		if roleEnum != "" {
			// 2. 获取存储在cookie的用户信息
			userJson := c.Request.Header.Cookie(constants.UserLoginState)
			if userJson == nil {
				c.JSON(200, errorutil.NotLoginError)
				c.Abort()
				return
			}
			// 3. URL解码
			decodeUserJson, err := url.QueryUnescape(string(userJson))
			if err != nil {
				c.JSON(200, errorutil.NotLoginError)
				c.Abort()
				return
			}
			decodeUser = []byte(decodeUserJson)
		}

		// 4. 解析用户信息
		var user model.User
		err := json.Unmarshal(decodeUser, &user)
		if err != nil {
			c.JSON(200, errorutil.SystemError.WithMessage(err.Error()))
			c.Abort()
			return
		}

		// 5. 校验用户权限等级是否符合要求
		dbUser, err := query.Use(db).User.Where(query.User.ID.Eq(user.ID), query.User.IsDelete.Eq(0)).First()
		if err != nil {
			c.JSON(200, errorutil.NotAuthError)
			c.Abort()
			return
		}
		if roleEnum == enum.AdminRole && enum.UserRoleEnum(dbUser.UserRole) != roleEnum {
			c.JSON(200, errorutil.NotAuthError)
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}
