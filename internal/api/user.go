package api

import (
	"yikou-ai-go-teach/internal/dal/model"
	"yikou-ai-go-teach/internal/dal/vo"
	"yikou-ai-go-teach/pkg/request"
	"yikou-ai-go-teach/pkg/response"
)

type YiKouUserRegisterRequest struct {
	UserAccount   string `json:"userAccount"`
	UserPassword  string `json:"userPassword"`
	CheckPassword string `json:"checkPassword"`
}

type YiKouUserRegisterResponse response.BaseResponse[int64]

type YiKouUserLoginRequest struct {
	UserAccount  string `json:"userAccount"`
	UserPassword string `json:"userPassword"`
}

type YiKouUserLoginResponse response.BaseResponse[vo.UserVo]

type YiKouUserAddRequest struct {
	UserAccount  string `json:"userAccount"`
	UserPassword string `json:"userPassword"`
	UserAvatar   string `json:"userAvatar"`
	UserProfile  string `json:"userProfile"`
	UserRole     string `json:"userRole"`
}

type YiKouUserAddResponse response.BaseResponse[int64]

type YiKouUserGetResponse response.BaseResponse[model.User]

type YiKouUserGetVoResponse response.BaseResponse[vo.UserVo]

type YiKouUserDeleteResponse response.BaseResponse[bool]

type YiKouUserUpdateRequest struct {
	request.DeleteRequest
	UserName    string `json:"userName"`
	UserAvatar  string `json:"userAvatar"`
	UserProfile string `json:"userProfile"`
	UserRole    string `json:"userRole"`
}

type YiKouUserUpdateResponse response.BaseResponse[bool]

type YiKouUserQueryRequest struct {
	request.PageRequest
	UserAccount string `json:"userAccount"`
	UserProfile string `json:"userProfile"`
	UserName    string `json:"userName"`
	UserRole    string `json:"userRole"`
}

type YiKouUserPageVoResponse response.BaseResponse[response.PageResponse[vo.UserVo]]
