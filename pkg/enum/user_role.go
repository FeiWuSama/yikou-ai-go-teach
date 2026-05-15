package enum

// UserRoleEnum 用户角色枚举
type UserRoleEnum string

const (
	UserRole  UserRoleEnum = "user"
	AdminRole UserRoleEnum = "admin"
)

var roleTextMap = map[UserRoleEnum]string{
	UserRole:  "用户",
	AdminRole: "管理员",
}

// GetRoleText 获取角色文本
func (e UserRoleEnum) GetRoleText() string {
	if text, ok := roleTextMap[e]; ok {
		return text
	}
	return "未知角色"
}
