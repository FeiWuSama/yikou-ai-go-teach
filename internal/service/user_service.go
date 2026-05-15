package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"yikou-ai-go-teach/internal/api"
	"yikou-ai-go-teach/internal/dal/model"
	"yikou-ai-go-teach/internal/dal/vo"
	"yikou-ai-go-teach/pkg/response"
)

type IUserService interface {
	UserRegister(ctx context.Context, req *api.YiKouUserRegisterRequest) (int64, error)
	GetEncryptPassword(ctx context.Context, password string) string
	GetLoginUserVo(ctx context.Context, c *app.RequestContext) (vo.UserVo, error)
	UserLogin(ctx context.Context, req *api.YiKouUserLoginRequest, c *app.RequestContext) (vo.UserVo, error)
	Logout(ctx context.Context, c *app.RequestContext) error
	AddUser(ctx context.Context, req *api.YiKouUserAddRequest) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	GetUserVo(ctx context.Context, id int64) (vo.UserVo, error)
	DeleteUser(ctx context.Context, id int64) (bool, error)
	UpdateUser(ctx context.Context, req *api.YiKouUserUpdateRequest) (bool, error)
	ListUserVoByPage(ctx context.Context, req *api.YiKouUserQueryRequest) (*response.PageResponse[vo.UserVo], error)
}
