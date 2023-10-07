package controllers

import (
	"fmt"
	"net/http"

	"login_golang/models"
	"login_golang/objects"
	"login_golang/tables"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

// interface
type UserController interface {
	RegisterHandler(c *gin.Context)
}

type userController struct {
	userMod models.UserModels
}

func NewUserController(userModel models.UserModels) UserController {
	return &userController{
		userMod: userModel,
	}
}

func (ctr *userController) RegisterHandler(c *gin.Context) {

	var reqData objects.Users
	err := c.ShouldBindJSON(&reqData)
	if err != nil {
		fmt.Println(err)
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error validate %s, condition: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorMessages,
		})
		return

	}

	var whreAdmPhone tables.Users
	whreAdmPhone.Phone = reqData.Phone

	var whreAdmUserName tables.Users
	whreAdmUserName.Username = reqData.Username
	whreStrAdmUserName := "users.role_id in (1,2)"
	checkAdmEmailResp, _ := ctr.userMod.GetUserWhereRow(whreAdmUserName, whreStrAdmUserName)
	if checkAdmEmailResp.ID == 0 {
		hash, err := Hash(reqData.Password)
		if err != nil {
			fmt.Println(err)
			return
		}

		// lowName := strings.Replace(strings.ToLower(reqData.Name), " ", "", -1)
		// encryptCode := lowName[0:4] + helpers.EncodeToString(4)

		var userData tables.Users
		userData.Name = reqData.Name
		userData.Phone = reqData.Phone
		userData.Email = reqData.Email
		userData.Username = reqData.Username
		userData.Password = hash
		userData.IsActive = true
		userData.RoleID = 1 // owner

		resUser, err := ctr.userMod.InsertUser(userData)
		if err != nil {
			fmt.Println(err)
			return
		}

		if resUser.ID > 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Register Success",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Register not Success",
			})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "This username or phone has registered before. Please use another account.",
			"data":    nil,
		})
		return
	}

}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
