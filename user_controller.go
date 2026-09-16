package controllers

import (
	"net/http"
	"strings"

	"12middleware/models"
	"12middleware/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
	Age      int    `json:"age" validate:"gte=18,lte=100"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
	Age      int    `json:"age" validate:"gte=18,lte=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func validationErrorResponse(c *gin.Context, err error) {

	validationErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	errors := make(map[string]string)

	for _, validationError := range validationErrors {

		field := strings.ToLower(validationError.Field())
		tag := validationError.Tag()

		if tag == "required" {
			errors[field] = field + " is required"
			continue
		}

		if tag == "email" {
			errors[field] = field + " must be a valid email address"
			continue
		}

		if tag == "min" {
			errors[field] = field + " must be at least " + validationError.Param() + " characters"
			continue
		}

		if tag == "max" {
			errors[field] = field + " must be at most " + validationError.Param() + " characters"
			continue
		}

		if tag == "gte" {
			errors[field] = field + " must be at least " + validationError.Param()
			continue
		}

		if tag == "lte" {
			errors[field] = field + " must be at most " + validationError.Param()
			continue
		}

		errors[field] = "Invalid value for " + field
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"errors": errors,
	})
}

func CreateUser(c *gin.Context) {

	var request CreateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)

	if err := validate.Struct(request); err != nil {
		validationErrorResponse(c, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
		Role:     request.Role,
		Age:      request.Age,
	}

	models.DB.Create(&user)

	c.JSON(http.StatusCreated, user)
}

func GetUsers(c *gin.Context) {

	var users []models.User

	if result := models.DB.Find(&users); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {

	id := c.Param("id")

	var user models.User

	if result := models.DB.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {

	id := c.Param("id")

	var user models.User

	if result := models.DB.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	var request UpdateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)

	if err := validate.Struct(request); err != nil {
		validationErrorResponse(c, err)
		return
	}

	user.Name = request.Name
	user.Email = request.Email
	user.Role = request.Role
	user.Age = request.Age

	if request.Password != "" {

		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(request.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash password",
			})
			return
		}

		user.Password = string(hashedPassword)
	}

	models.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}

func LoginUser(c *gin.Context) {

	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	request.Email = strings.TrimSpace(request.Email)

	if err := validate.Struct(request); err != nil {
		validationErrorResponse(c, err)
		return
	}

	var user models.User

	if result := models.DB.Where("email = ?", request.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func DeleteUser(c *gin.Context) {

	id := c.Param("id")

	var user models.User

	if result := models.DB.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	models.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
