// internal/model/dto/response/user.go
package response

import "time"

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}

type UserInfo struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Capacity  int64     `json:"capacity"`
	Used      int64     `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

type StorageInfo struct {
	Capacity    int64   `json:"capacity"`
	Used        int64   `json:"used"`
	Available   int64   `json:"available"`
	UsageRate   float64 `json:"usage_rate"`
	FileCount   int64   `json:"file_count"`
	FolderCount int64   `json:"folder_count"`
}
