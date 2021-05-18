package main

import (
	// Log items to the terminal
	"log"
	// Import gin for route definition
	"github.com/gin-gonic/gin"
	// Import godotenv for .env variables
	"github.com/joho/godotenv"
	// Import our app controllers
	"sarvbooksapi/controllers"
)

// init gets called before the main function
func init() {
	// Log error if .env file does not exist
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}
}

func main() {
	router := gin.Default()
	v1 := router.Group("/")
	{
		user := new(controllers.UserController)
		// Create the login endpoint
		v1.POST("/login", user.Login)
		// Create the signup phone endpoint
		v1.POST("/signupphone", user.SignupRequest)
		// Create the login phone endpoint
		v1.POST("/loginphone", user.LoginRequest)
		// Create the otp verification endpoint
		v1.POST("/signupOTPCheck", user.SignupOTPCheck)
		// Create the otp verification endpoint
		v1.POST("/loginOTPCheck", user.LoginOTPCheck)
	}
	// Handle error response when a route is not defined
	router.NoRoute(func(c *gin.Context) {
		// In gin this is how you return a JSON response
		c.JSON(404, gin.H{"message": "Not found"})
	})
	// Init our server
	router.Run(":8080")
}

//{"name":"Manan", "email" : "abcx", "phone":"7638872795", "password":"aaa"} // SignUp Request
//{"phone":"7638872795","otp" : 1111} //SignupOTP Check
//{"email": "xyz", "password": 1234}
