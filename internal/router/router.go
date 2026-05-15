package router

import (
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
	"gorm.io/gorm"
	"time"
	"yikou-ai-go-teach/internal/handler"
	"yikou-ai-go-teach/internal/middleware"
	"yikou-ai-go-teach/pkg/enum"
	"yikou-ai-go-teach/pkg/errorutil"
	"yikou-ai-go-teach/pkg/response"
)

// RegisterRoutes 注册路由
func RegisterRoutes(h *server.Hertz, url func(config *swagger.Config), db *gorm.DB, userHandler *handler.UserHandler) {
	// 注册全局中间件
	// 处理跨域问题
	h.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	// 全局异常处理
	h.Use(recovery.Recovery(recovery.WithRecoveryHandler(CustomRecoveryHandler)))

	// 测试接口
	h.GET("/ping", handler.Ping)
	// swaggo文档
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))

	userRoute := h.Group("/user")
	{
		userRoute.POST("/register", userHandler.UserRegister)
		userRoute.POST("/login", userHandler.UserLogin)
		userRoute.GET("/get/vo", userHandler.GetUserVo)

		// 需要登录的接口
		userRoute.GET("/get/login", middleware.AuthMiddleware(enum.UserRole, db), userHandler.GetLoginUser)
		userRoute.POST("/logout", middleware.AuthMiddleware(enum.UserRole, db), userHandler.Logout)

		// 需要管理员权限的接口
		userRoute.POST("/add", middleware.AuthMiddleware(enum.AdminRole, db), userHandler.AddUser)
		userRoute.GET("/get", middleware.AuthMiddleware(enum.AdminRole, db), userHandler.GetUser)
		userRoute.POST("/delete", middleware.AuthMiddleware(enum.AdminRole, db), userHandler.DeleteUser)
		userRoute.POST("/update", middleware.AuthMiddleware(enum.AdminRole, db), userHandler.UpdateUser)
		userRoute.POST("/list/page/vo", middleware.AuthMiddleware(enum.AdminRole, db), userHandler.ListUserVoByPage)
	}
}

// CustomRecoveryHandler 全局异常处理器
func CustomRecoveryHandler(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
	logger.Errorf("panic recovered: %v\n%s", err, stack)
	c.JSON(consts.StatusOK, response.NewErrorResponse[any](errorutil.SystemError.WithMessage(fmt.Sprintf("%v", err))))
	c.Abort()
}
