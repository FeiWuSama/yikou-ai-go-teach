package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
	"yikou-ai-go-teach/internal/api"
	"yikou-ai-go-teach/internal/dal/model"
	"yikou-ai-go-teach/internal/dal/vo"
	"yikou-ai-go-teach/internal/service"
	"yikou-ai-go-teach/pkg/errorutil"
	"yikou-ai-go-teach/pkg/request"
	"yikou-ai-go-teach/pkg/response"
)

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// UserRegister 用户注册
// @Summary 用户注册
// @Description 用户注册
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param req body api.YiKouUserRegisterRequest true "用户注册请求"
// @Success 200 {object} api.YiKouUserRegisterResponse "用户ID"
// @Router /user/register [post]
func (u *UserHandler) UserRegister(ctx context.Context, c *app.RequestContext) {
	req := &api.YiKouUserRegisterRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	userId, err := u.userService.UserRegister(ctx, req)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	c.JSON(consts.StatusOK, response.NewSuccessResponse[int64](userId))
}

// UserLogin 用户登录
// @Summary 用户登录
// @Description 用户登录
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param req body api.YiKouUserLoginRequest true "用户登录请求"
// @Success 200 {object} api.YiKouUserLoginResponse "登录用户信息"
// @Router /user/login [post]
func (u *UserHandler) UserLogin(ctx context.Context, c *app.RequestContext) {
	req := &api.YiKouUserLoginRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	userVo, err := u.userService.UserLogin(ctx, req, c)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	c.JSON(consts.StatusOK, response.NewSuccessResponse[vo.UserVo](userVo))
}

// GetLoginUser 获取登录用户信息
// @Summary 获取登录用户信息
// @Description 获取登录用户信息
// @Tags 用户模块
// @Accept json
// @Produce json
// @Success 200 {object} api.YiKouUserLoginResponse "登录用户信息"
// @Router /user/get/login [get]
func (u *UserHandler) GetLoginUser(ctx context.Context, c *app.RequestContext) {
	userVo, err := u.userService.GetLoginUserVo(ctx, c)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	c.JSON(consts.StatusOK, response.NewSuccessResponse[vo.UserVo](userVo))
}

// Logout 退出登录
// @Summary 退出登录
// @Description 退出登录
// @Tags 用户模块
// @Accept json
// @Produce json
// @Success 200 {object} bool "退出登录成功"
// @Router /user/logout [get]
func (u *UserHandler) Logout(ctx context.Context, c *app.RequestContext) {
	err := u.userService.Logout(ctx, c)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	c.JSON(consts.StatusOK, response.NewSuccessResponse[any](true))
}

// AddUser 新增用户
// @Summary 新增用户
// @Description 新增用户
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param req body api.YiKouUserAddRequest true "新增用户请求"
// @Success 200 {object} api.YiKouUserAddResponse "用户ID"
// @Router /user/add [post]
func (u *UserHandler) AddUser(ctx context.Context, c *app.RequestContext) {
	req := &api.YiKouUserAddRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	userId, err := u.userService.AddUser(ctx, req)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	c.JSON(consts.StatusOK, response.NewSuccessResponse[int64](userId))
}

// GetUser 根据ID获取用户
// @Summary 根据ID获取用户
// @Description 根据ID获取用户
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param id query int true "用户ID"
// @Success 200 {object} api.YiKouUserGetResponse "用户信息"
// @Router /user/get [get]
func (u *UserHandler) GetUser(ctx context.Context, c *app.RequestContext) {
	id := c.Query("id")
	if id == "" {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](errorutil.ParamsError))
		return
	}
	idInt64, _ := strconv.ParseInt(id, 10, 64)
	user, err := u.userService.GetUser(ctx, idInt64)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	c.JSON(consts.StatusOK, response.NewSuccessResponse[model.User](*user))
}

// GetUserVo 根据ID获取用户VO
// @Summary 根据ID获取用户VO
// @Description 根据ID获取用户VO
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param id query int true "用户ID"
// @Success 200 {object} api.YiKouUserGetVoResponse "用户VO信息"
// @Router /user/get/vo [get]
func (u *UserHandler) GetUserVo(ctx context.Context, c *app.RequestContext) {
	id := c.Query("id")
	if id == "" {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](errorutil.ParamsError))
		return
	}
	idInt64, _ := strconv.ParseInt(id, 10, 64)
	userVo, err := u.userService.GetUserVo(ctx, idInt64)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	c.JSON(consts.StatusOK, response.NewSuccessResponse[vo.UserVo](userVo))
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除用户
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param req body request.DeleteRequest true "删除用户请求"
// @Success 200 {object} api.YiKouUserDeleteResponse "删除结果"
// @Router /user/delete [post]
func (u *UserHandler) DeleteUser(ctx context.Context, c *app.RequestContext) {
	req := &request.DeleteRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	success, err := u.userService.DeleteUser(ctx, int64(req.Id))
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	c.JSON(consts.StatusOK, response.NewSuccessResponse[bool](success))
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新用户
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param req body api.YiKouUserUpdateRequest true "更新用户请求"
// @Success 200 {object} api.YiKouUserUpdateResponse "更新结果"
// @Router /user/update [post]
func (u *UserHandler) UpdateUser(ctx context.Context, c *app.RequestContext) {
	req := &api.YiKouUserUpdateRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}

	if req.Id == 0 {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](errorutil.ParamsError.WithMessage("用户ID不能为空")))
		return
	}

	success, err := u.userService.UpdateUser(ctx, req)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	c.JSON(consts.StatusOK, response.NewSuccessResponse[bool](success))
}

// ListUserVoByPage 分页获取用户VO列表
// @Summary 分页获取用户VO列表
// @Description 分页获取用户VO列表
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param req body api.YiKouUserQueryRequest true "分页查询请求"
// @Success 200 {object} api.YiKouUserPageVoResponse "分页用户VO列表"
// @Router /user/list/page/vo [post]
func (u *UserHandler) ListUserVoByPage(ctx context.Context, c *app.RequestContext) {
	req := &api.YiKouUserQueryRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	pageResponse, err := u.userService.ListUserVoByPage(ctx, req)
	if err != nil {
		c.JSON(consts.StatusOK, response.NewErrorResponse[any](err))
		return
	}
	c.JSON(consts.StatusOK, response.NewSuccessResponse[*response.PageResponse[vo.UserVo]](pageResponse))
}
