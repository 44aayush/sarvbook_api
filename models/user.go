package models

import (
	"math/rand"
	"sarvbooksapi/forms"
	"sarvbooksapi/helpers"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//User struct
type User struct {
	ID       bson.ObjectId `json:"_id" bson:"_id"`
	Name     string        `json:"name" bson:"name"`
	Email    string        `json:"email" bson:"email"`
	Phone    string        `json:"phone" bson:"phone"`
	Password string        `json:"password" bson:"password"`
}

//OTP Struct
type OTP struct {
	Phone   string    `json:"phone" bson:"phone"`
	Otp     int       `json:"otp" bson:"otp"`
	Expires time.Time `json:"expires" bson:"expires"`
}

//UserModel struct
type UserModel struct{}

// GetUserByEmail handles fetching user by email
func (u *UserModel) GetUserByEmail(email string) (user User, err error) {
	// Connect to the user collection
	collection := dbConnect.Use(databaseName, "user")
	// Assign result to error object while saving user
	err = collection.Find(bson.M{"email": email}).One(&user)
	return user, err
}

// AddUser handles registering a user
func (u *UserModel) AddUser(data forms.SignupUserCommand) error {

	// Connect to the user collection
	collection := dbConnect.Use(databaseName, "user")

	// Assign result to error object while saving user
	err := collection.Insert(bson.M{
		"name":     data.Name,
		"email":    data.Email,
		"phone":    data.Phone,
		"password": helpers.GeneratePasswordHash([]byte(data.Password)),
	})
	return err
}

// GetUserByPhone handles fetching user by OTP
func (u *UserModel) GetUserByPhone(phone string) (user User, err error) {
	// Connect to the user collection
	collection := dbConnect.Use(databaseName, "user")
	// Assign result to error object while saving user
	err = collection.Find(bson.M{"phone": phone}).One(&user)
	return user, err
}

// GenerateOTP handles registering a user
func (u *UserModel) GenerateOTP(phone string) error {

	// Connect to the user collection
	collection := dbConnect.Use(databaseName, "otp")

	// Assign result to error object while saving user
	err := collection.Insert(bson.M{
		"otp":     rand.Intn(9999-1000) + 1000,
		"phone":   phone,
		"expires": time.Now().Add(15 * time.Minute),
	})
	return err
}

// VerifyOTP , function receives phone number and otp, reads the database and checks if the otp matches the phone number
func (u *UserModel) VerifyOTP(phone string, otp int) (bool, string, error) {

	var otpdata OTP
	collection := dbConnect.Use(databaseName, "otp")
	err := collection.Find(bson.M{"phone": phone}).One(&otpdata)

	if err != nil {
		return false, "OTP not sent", err
	}
	err = u.RemoveOldOTP(phone)

	if otpdata.Otp == otp {
		if time.Now().Before(otpdata.Expires) {
			return true, "Verified OTP", err
		} else {
			return false, "Expired OTP", err
		}
	} else {
		return false, "Incorrect OTP", err
	}

}

//// OTP expires after 15 mins, so check if .the otp has been created within 15 mins or not.
//// If created within 15 mins, all good. else reply with expired message.
//// Remove otp once verified

//RemoveOldOTP fucntion removes all otp related to a phone as soon as the above 2 functions are called
func (u *UserModel) RemoveOldOTP(phone string) error {
	collection := dbConnect.Use(databaseName, "otp")
	err := collection.Remove(bson.M{"phone": phone})
	return err
}
