package models

type UserCreateDTO struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

type UserLoginResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserUpdatePasswordDTO struct {
	Password string `json:"password"`
}

type UserUpdateDTO struct {
	Email       string  `json:"email" db:"email"`
	DisplayName string  `json:"display_name" db:"display_name"`
	AvatarURL   *string `json:"avatar_url" db:"avatar_url"`
	ProfileText *string `json:"profile_text" db:"profile_text"`
	UserTypeKey string  `json:"user_type_key" db:"user_type_key"`
}

type UserBanDTO struct {
	Reason string `json:"reason"`
}
