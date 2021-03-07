package domain

// UserInfo Represents profile information on a user
type UserInfo struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func NewUserInfo() UserInfo {
	return UserInfo{}
}

func (u UserInfo) WithEmail(email string) UserInfo {
	return UserInfo{
		Email:    email,
		Username: u.Username,
	}
}

func (u UserInfo) WithUsername(username string) UserInfo {
	return UserInfo{
		Email:    u.Email,
		Username: username,
	}
}
