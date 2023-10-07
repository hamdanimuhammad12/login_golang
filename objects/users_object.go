package objects

type UserLogin struct {
	UserID         string `json:"user_id"`
	RoleID         string `json:"role_id"`
	Email          string `json:"email"`
	OrganizationID string `json:"organization_id"`
	Username       string `json:"username" `
}

type User struct {
	UserID string `json:"userID"`
	Email  string `json:"email,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Name   string `json:"name,omitempty"`
}

type Users struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Phone       string `json:"phone" form:"phone" binding:"required"`
	Email       string `json:"email" form:"email"`
	GenderID    int    `json:"gender_id" form:"gender_id" `
	DateOfBirth string `json:"date_of_birth" form:"date_of_birth"`
	Password    string `json:"password" form:"password"`
	Username    string `json:"username" `
	BaseUrl     string `json:"base_url" form:"base_url"`
}

type Login struct {
	Email    string `form:"email" json:"email" `
	Phone    string `form:"phone" json:"phone" `
	Username string `form:"username" json:"username" `
	Password string `form:"password" json:"password" `
	OtpCode  string `form:"otp_code" json:"otp_code" `
}
