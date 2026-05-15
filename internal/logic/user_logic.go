package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"gorm.io/gorm"
	"net/url"
	"yikou-ai-go-teach/internal/api"
	"yikou-ai-go-teach/internal/dal/model"
	"yikou-ai-go-teach/internal/dal/query"
	"yikou-ai-go-teach/internal/dal/vo"
	"yikou-ai-go-teach/pkg/constants"
	"yikou-ai-go-teach/pkg/enum"
	"yikou-ai-go-teach/pkg/errorutil"
	"yikou-ai-go-teach/pkg/response"
	"yikou-ai-go-teach/pkg/snowflake"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) Logout(ctx context.Context, c *app.RequestContext) error {
	// 清除Cookie
	c.SetCookie(constants.UserLoginState, "", 0, "/", "", protocol.CookieSameSiteLaxMode, false, true)
	return nil
}

func (s *UserService) GetEncryptPassword(ctx context.Context, password string) string {
	h := md5.New()
	h.Write([]byte("feiwu" + password)) // 加盐
	return hex.EncodeToString(h.Sum(nil))
}

func (s *UserService) UserRegister(ctx context.Context, req *api.YiKouUserRegisterRequest) (int64, error) {
	// 1. 校验参数
	if req.UserAccount == "" || req.UserPassword == "" || req.CheckPassword == "" {
		return 0, errorutil.ParamsError
	}
	if len(req.UserAccount) < 4 || len(req.UserAccount) > 12 {
		return 0, errorutil.ParamsError.WithMessage("用户账号长度必须在4到12之间")
	}
	if len(req.UserPassword) < 8 || len(req.UserPassword) > 12 {
		return 0, errorutil.ParamsError.WithMessage("用户密码长度必须在8到12之间")
	}
	if req.UserPassword != req.CheckPassword {
		return 0, errorutil.ParamsError.WithMessage("两次输入密码不一致")
	}
	// 2. 校验用户名是否已被注册
	count, _ := query.Use(s.db).User.Where(query.User.UserAccount.Eq(req.UserAccount)).Count()
	if count > 0 {
		return 0, errorutil.ParamsError.WithMessage("用户名已被注册")
	}
	// 3. 密码加密
	encryptPassword := s.GetEncryptPassword(ctx, req.UserPassword)
	userId, err := snowflake.GenerateSnowFlakeId()
	if err != nil {
		return 0, err
	}
	// 4. 创建用户
	newUser := &model.User{
		ID:           userId,
		UserAccount:  req.UserAccount,
		UserPassword: encryptPassword,
		UserName:     "无名",
		UserRole:     string(enum.UserRole),
	}
	err = query.Use(s.db).User.Create(newUser)
	if err != nil {
		return 0, err
	}
	return newUser.ID, nil
}

func (s *UserService) GetLoginUserVo(ctx context.Context, c *app.RequestContext) (vo.UserVo, error) {
	// 1. 获取user
	userJson := c.Request.Header.Cookie(constants.UserLoginState)
	if userJson == nil || string(userJson) == "" {
		return vo.UserVo{}, errorutil.ParamsError.WithMessage("登录已过期，请重新登录")
	}
	decodeUser, err := url.QueryUnescape(string(userJson))
	if err != nil {
		return vo.UserVo{}, err
	}
	var user model.User
	err = json.Unmarshal([]byte(decodeUser), &user)
	if err != nil {
		return vo.UserVo{}, err
	}
	// 2. 校验用户是否存在
	_, err = query.Use(s.db).User.Where(query.User.ID.Eq(user.ID), query.User.IsDelete.Eq(0)).First()
	if err != nil {
		return vo.UserVo{}, err
	}
	// 3. 构建 User
	loginUserVo := vo.UserVo{
		ID:          user.ID,
		UserAccount: user.UserAccount,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
		UserRole:    user.UserRole,
		CreateTime:  user.CreateTime,
		UpdateTime:  user.UpdateTime,
	}
	return loginUserVo, nil
}

// AddUser 新增用户
func (s *UserService) AddUser(ctx context.Context, req *api.YiKouUserAddRequest) (int64, error) {
	// 1. 校验参数
	if req.UserAccount == "" || req.UserPassword == "" {
		return 0, errorutil.ParamsError
	}
	if len(req.UserAccount) < 4 || len(req.UserAccount) > 12 {
		return 0, errorutil.ParamsError.WithMessage("用户账号长度必须在4到12之间")
	}
	if len(req.UserPassword) < 8 || len(req.UserPassword) > 12 {
		return 0, errorutil.ParamsError.WithMessage("用户密码长度必须在8到12之间")
	}

	// 2. 校验用户名是否已被注册
	count, _ := query.Use(s.db).User.Where(query.User.UserAccount.Eq(req.UserAccount)).Count()
	if count > 0 {
		return 0, errorutil.ParamsError.WithMessage("用户名已被注册")
	}

	// 3. 密码加密
	encryptPassword := s.GetEncryptPassword(ctx, req.UserPassword)

	// 4. 创建用户
	newUser := &model.User{
		UserAccount:  req.UserAccount,
		UserPassword: encryptPassword,
		UserName:     req.UserAccount, // 默认使用账号作为用户名
		UserAvatar:   req.UserAvatar,
		UserProfile:  req.UserProfile,
		UserRole:     req.UserRole,
	}
	if req.UserRole == "" {
		newUser.UserRole = string(enum.UserRole) // 默认角色
	}

	err := query.Use(s.db).User.Create(newUser)
	if err != nil {
		return 0, err
	}
	return newUser.ID, nil
}

// GetUser 根据ID获取用户
func (s *UserService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	user, err := query.Use(s.db).User.Where(query.User.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserVo 根据ID获取用户VO
func (s *UserService) GetUserVo(ctx context.Context, id int64) (vo.UserVo, error) {
	user, err := s.GetUser(ctx, id)
	if err != nil {
		return vo.UserVo{}, err
	}

	userVo := vo.UserVo{
		ID:          user.ID,
		UserAccount: user.UserAccount,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
		UserRole:    user.UserRole,
		CreateTime:  user.CreateTime,
		UpdateTime:  user.UpdateTime,
	}
	return userVo, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, id int64) (bool, error) {
	// 软删除
	_, err := query.Use(s.db).User.Where(query.User.ID.Eq(id)).Update(query.User.IsDelete, 1)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, req *api.YiKouUserUpdateRequest) (bool, error) {
	// 1. 检查用户是否存在
	_, err := query.Use(s.db).User.Where(query.User.ID.Eq(int64(req.Id))).First()
	if err != nil {
		return false, err
	}

	// 2. 更新用户信息
	updateMap := make(map[string]interface{})
	updateMap["user_name"] = req.UserName
	updateMap["user_avatar"] = req.UserAvatar
	updateMap["user_profile"] = req.UserProfile
	updateMap["user_role"] = req.UserRole
	_, err = query.Use(s.db).User.Where(query.User.ID.Eq(int64(req.Id))).Updates(updateMap)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ListUserVoByPage 分页获取用户VO列表
func (s *UserService) ListUserVoByPage(ctx context.Context, req *api.YiKouUserQueryRequest) (*response.PageResponse[vo.UserVo], error) {
	// 1. 设置默认分页参数
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 2. 构建查询条件
	queryBuilder := query.Use(s.db).User.Where(query.User.IsDelete.Eq(0))

	if req.UserAccount != "" {
		queryBuilder = queryBuilder.Where(query.User.UserAccount.Like("%" + req.UserAccount + "%"))
	}
	if req.UserName != "" {
		queryBuilder = queryBuilder.Where(query.User.UserName.Like("%" + req.UserName + "%"))
	}
	if req.UserProfile != "" {
		queryBuilder = queryBuilder.Where(query.User.UserProfile.Like("%" + req.UserProfile + "%"))
	}
	if req.UserRole != "" {
		queryBuilder = queryBuilder.Where(query.User.UserRole.Eq(req.UserRole))
	}

	// 3. 查询总数
	totalCount, err := queryBuilder.Count()
	if err != nil {
		return nil, err
	}

	// 4. 计算分页
	totalPage := int((totalCount + int64(req.PageSize) - 1) / int64(req.PageSize))
	offset := (req.PageNum - 1) * req.PageSize

	// 5. 排序
	if req.SortField != "" {
		if orderExpr, ok := query.User.GetFieldByName(req.SortField); ok {
			if req.SortOrder == "desc" {
				queryBuilder = queryBuilder.Order(orderExpr.Desc())
			} else {
				queryBuilder = queryBuilder.Order(orderExpr)
			}
		} else {
			// 如果字段不存在，使用默认排序
			queryBuilder = queryBuilder.Order(query.User.CreateTime.Desc())
		}
	} else {
		queryBuilder = queryBuilder.Order(query.User.CreateTime.Desc())
	}

	// 6. 分页查询
	users, err := queryBuilder.Offset(offset).Limit(req.PageSize).Find()
	if err != nil {
		return nil, err
	}

	// 7. 转换为UserVo列表
	var userVoList []vo.UserVo
	for _, user := range users {
		userVo := vo.UserVo{
			ID:          user.ID,
			UserAccount: user.UserAccount,
			UserName:    user.UserName,
			UserAvatar:  user.UserAvatar,
			UserProfile: user.UserProfile,
			UserRole:    user.UserRole,
			CreateTime:  user.CreateTime,
			UpdateTime:  user.UpdateTime,
		}
		userVoList = append(userVoList, userVo)
	}

	// 8. 构建分页响应
	pageResponse := &response.PageResponse[vo.UserVo]{
		Records:            userVoList,
		PageNum:            req.PageNum,
		PageSize:           req.PageSize,
		TotalPage:          totalPage,
		TotalRow:           int(totalCount),
		OptimizeCountQuery: false,
	}

	return pageResponse, nil
}

func (s *UserService) UserLogin(ctx context.Context, req *api.YiKouUserLoginRequest, c *app.RequestContext) (vo.UserVo, error) {
	// 1. 校验参数
	if req.UserAccount == "" || req.UserPassword == "" {
		return vo.UserVo{}, errorutil.ParamsError
	}
	// 2. 校验用户是否存在
	user, err := query.Use(s.db).User.Where(query.User.UserAccount.Eq(req.UserAccount)).First()
	if err != nil {
		return vo.UserVo{}, err
	}
	// 3. 校验密码是否正确
	encryptPassword := s.GetEncryptPassword(ctx, req.UserPassword)
	if user.UserPassword != encryptPassword {
		return vo.UserVo{}, errorutil.ParamsError.WithMessage("密码错误")
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		return vo.UserVo{}, err
	}
	// 4. 保存sessionId到cookie
	c.SetCookie(constants.UserLoginState, string(userJson),
		86400, "/", "", protocol.CookieSameSiteLaxMode, false, true)
	// 5. 构建userVo对象
	loginUserVo := vo.UserVo{
		ID:          user.ID,
		UserAccount: user.UserAccount,
		UserName:    user.UserName,
		UserAvatar:  user.UserAvatar,
		UserProfile: user.UserProfile,
		UserRole:    user.UserRole,
		CreateTime:  user.CreateTime,
		UpdateTime:  user.UpdateTime,
	}
	return loginUserVo, nil
}
