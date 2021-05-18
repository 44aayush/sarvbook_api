package forms

// SignupUserCommand defines user form struct
type SignupUserCommand struct {
	// binding:"required" ensures that the field is provided
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	OTP      int    `json:"otp"`
}

// LoginEmailCommand defines user login form struct
type LoginEmailCommand struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//OTPCommand defines user form struct
/*type OTPCommand struct {
	OTP   int    `json:"otp" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}*/

// PhoneCommand defines user login form struct
type PhoneCommand struct {
	Phone string `json:"phone" binding:"required"`
	OTP   int    `json:"otp"`
}
