package controllers

import (
	"sarvbooksapi/forms"
	"sarvbooksapi/helpers"
	"sarvbooksapi/models"
	"sarvbooksapi/services"

	"gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"
)

// Import the userModel from the models
var userModel = new(models.UserModel)

// UserController defines the user controller methods
type UserController struct{}

// Login allows a user to login a user and get
// access token
func (u *UserController) Login(c *gin.Context) {
	var data forms.LoginEmailCommand

	// Bind the request body data to var data and check if all details are provided
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide required details"})
		c.Abort()
		return
	}

	result, err := userModel.GetUserByEmail(data.Email)

	if result.Email == "" {
		c.JSON(404, gin.H{"message": "User account was not found"})
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem logging into your account"})
		c.Abort()
		return
	}

	// Get the hashed password from the saved document
	hashedPassword := []byte(result.Password)
	// Get the password provided in the request.body
	password := []byte(data.Password)

	err = helpers.PasswordCompare(password, hashedPassword)

	if err != nil {
		c.JSON(403, gin.H{"message": "Invalid user credentials"})
		c.Abort()
		return
	}

	jwtToken, err2 := services.GenerateToken(data.Email)

	// If we fail to generate token for access
	if err2 != nil {
		c.JSON(403, gin.H{"message": "There was a problem logging you in, try again later"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "Log in success", "token": jwtToken})
}

//SignupRequest for OTP
func (u *UserController) SignupRequest(c *gin.Context) {
	var data forms.SignupUserCommand

	// Bind the request body data to var data and check if all details are provided
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide required details"})
		c.Abort()
		return
	}

	_, err := userModel.GetUserByPhone(data.Phone)

	// If there happens to be a result respond with a
	// descriptive mesage
	if err != mgo.ErrNotFound {
		c.JSON(403, gin.H{"message": "Phone is already in use"})
		c.Abort()
		return
	}

	result, _ := userModel.GetUserByEmail(data.Email)

	// If there happens to be a result respond with a descriptive mesage
	if result.Email != "" {
		c.JSON(403, gin.H{"message": "Email is already in use"})
		c.Abort()
		return
	}

	err = userModel.GenerateOTP(data.Phone)
	// Check if there was an error when saving user
	if err != nil {
		c.JSON(400, gin.H{"message": "Problem creating an account"})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "OTP Sent to " + data.Phone})
}

// SignupOTPCheck OTP
func (u *UserController) SignupOTPCheck(c *gin.Context) {
	var data forms.SignupUserCommand

	// Bind the request body data to var data and check if all details are provided
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide required details"})
		c.Abort()
		return
	}

	issue, msg, err := userModel.VerifyOTP(data.Phone, data.OTP)

	if err == mgo.ErrNotFound {
		c.JSON(400, gin.H{"message": "OTP not Found"})
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem verifying OTP"})
		c.Abort()
		return
	}

	if issue == false && msg == "Expired OTP" {
		c.JSON(400, gin.H{"message": "Expired OTP"})
		c.Abort()
		return
	}

	if issue == true && msg == "Incorrect OTP" {
		c.JSON(400, gin.H{"message": "Incorrect OTP, Enter Once more"})
		c.Abort()
		return
	}
	err = userModel.AddUser(data)
	c.JSON(201, gin.H{"message": "OTP Verified"})
	c.Abort()
	return
}

// LoginRequest OTP
func (u *UserController) LoginRequest(c *gin.Context) {
	var data forms.PhoneCommand

	// Bind the request body data to var data and check if all details are provided
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide required details"})
		c.Abort()
		return
	}

	err := userModel.GenerateOTP(data.Phone)
	// Check if there was an error when saving user
	if err != nil {
		c.JSON(400, gin.H{"message": "Problem during login!"})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "OTP Sent to " + data.Phone})
}

// LoginOTPCheck OTP
func (u *UserController) LoginOTPCheck(c *gin.Context) {
	var data forms.PhoneCommand

	// Bind the request body data to var data and check if all details are provided
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide required details"})
		c.Abort()
		return
	}

	issue, msg, err := userModel.VerifyOTP(data.Phone, data.OTP)

	if err == mgo.ErrNotFound {
		c.JSON(400, gin.H{"message": "OTP not Found"})
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem verifying OTP"})
		c.Abort()
		return
	}

	if issue == false && msg == "Expired OTP" {
		c.JSON(400, gin.H{"message": "Expired OTP"})
		c.Abort()
		return
	}

	if issue == true && msg == "Incorrect OTP" {
		c.JSON(400, gin.H{"message": "Incorrect OTP, Enter Once more"})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "OTP Verified"})
	c.Abort()
	return
}
