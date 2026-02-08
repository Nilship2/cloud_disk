// internal/constant/error_code.go
package constant

// 错误码定义
const (
	SuccessCode = 0

	// 通用错误 1000-1999
	ErrParamInvalid      = 1001
	ErrDatabase          = 1002
	ErrInternal          = 1003
	ErrNotFound          = 1004
	ErrUnauthorized      = 1005
	ErrPermissionDenied  = 1006
	ErrResourceExhausted = 1007

	// 用户相关错误 2000-2999
	ErrUserExists     = 2001
	ErrUserNotExists  = 2002
	ErrPasswordWrong  = 2003
	ErrEmailExists    = 2004
	ErrUsernameExists = 2005
	ErrUserInactive   = 2006

	// 文件相关错误 3000-3999
	ErrFileTooLarge     = 3001
	ErrFileTypeInvalid  = 3002
	ErrFileUploadFailed = 3003
	ErrFileNotFound     = 3004
	ErrFileExists       = 3005
	ErrStorageFull      = 3006
	ErrFileCorrupted    = 3007

	// 分享相关错误 4000-4999
	ErrShareNotFound      = 4001
	ErrShareExpired       = 4002
	ErrSharePasswordWrong = 4003
	ErrShareLimitExceeded = 4004
)

// 错误消息映射
var ErrorMessages = map[int]string{
	SuccessCode: "success",

	ErrParamInvalid:      "参数无效",
	ErrDatabase:          "数据库错误",
	ErrInternal:          "服务器内部错误",
	ErrNotFound:          "资源不存在",
	ErrUnauthorized:      "未授权访问",
	ErrPermissionDenied:  "权限不足",
	ErrResourceExhausted: "资源不足",

	ErrUserExists:     "用户已存在",
	ErrUserNotExists:  "用户不存在",
	ErrPasswordWrong:  "密码错误",
	ErrEmailExists:    "邮箱已存在",
	ErrUsernameExists: "用户名已存在",
	ErrUserInactive:   "用户未激活",

	ErrFileTooLarge:     "文件太大",
	ErrFileTypeInvalid:  "文件类型不支持",
	ErrFileUploadFailed: "文件上传失败",
	ErrFileNotFound:     "文件不存在",
	ErrFileExists:       "文件已存在",
	ErrStorageFull:      "存储空间不足",
	ErrFileCorrupted:    "文件已损坏",

	ErrShareNotFound:      "分享不存在",
	ErrShareExpired:       "分享已过期",
	ErrSharePasswordWrong: "分享密码错误",
	ErrShareLimitExceeded: "分享次数超限",
}

func GetErrorMessage(code int) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
