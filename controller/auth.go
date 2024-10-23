package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sahilq312/workly/initializer"
	"github.com/sahilq312/workly/model"
	"github.com/sahilq312/workly/utils"
)

// Login function to authenticate a user
func Login(c *gin.Context) {
	// Get the request body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind the JSON body to the struct
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Check if user exists
	var user model.User
	result := initializer.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "User does not exist"})
		return
	}

	// Compare hashed password
	match, err := utils.CompareHashedPassword(body.Password, user.Password)
	if err != nil || !match {
		c.JSON(400, gin.H{"error": "Invalid password"})
		return
	}
	// Return user details
	c.JSON(200, gin.H{
		"user": user,
	})
}

// Register function to register a new user
func Register(c *gin.Context) {
	//struct for request
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	//Get the request body
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid Request"})
		return
	}

	//Check if user already exist
	userExist := model.User{Email: body.Email}
	result := initializer.DB.Where(&userExist).First(&userExist)
	if result.Error == nil {
		c.JSON(400, gin.H{
			"error": "User already exist",
		})
		return
	}

	//hash the password
	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Error in hashing the password",
		})
		return
	}

	//Create a new user
	user := model.User{Name: body.Name, Email: body.Email, Password: hashedPassword}
	result = initializer.DB.Create(&user)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": "Error in creating the user",
		})
		return
	}

	//Return User and session
	c.JSON(200, gin.H{
		"user": user,
	})
}

func GetUserById(c *gin.Context) {
	// Get id from params
	id := c.Params.ByName("id")

	//find user with the given id
	var user model.User
	if err := initializer.DB.Where("id = ?", id).First(&user).Error; err != nil {
		// If no user is found, return a 404 error
		c.JSON(404, gin.H{
			"error": "No user found",
		})
		return
	}
	// return user
	c.JSON(200, gin.H{
		"user": user,
	})
}

func GetUser(c *gin.Context) {
	//Get User from Session

	c.JSON(200, gin.H{
		"user": "user",
	})
	//Return User
}

func Logout(c *gin.Context) {
	//Delete Session

	//Return Success
}
