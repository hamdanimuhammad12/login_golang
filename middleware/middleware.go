package middleware

import (
	"fmt"
	"login_golang/models"
	"login_golang/objects"
	"login_golang/tables"
	"net/mail"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

var identityKey = "id"
var roleIDKey = "role_id"
var organizationIDKey = "organization_id"

func SchemaMstr(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("mstr" + "." + tableName)
	}
}

func SchemaUsr(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("public" + "." + tableName)
	}
}

func SetupMiddlewareUserName(db *gorm.DB) *jwt.GinJWTMiddleware {

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "jwt",
		Key:         []byte("#freelance#"),
		Timeout:     time.Duration(24*365) * time.Hour,
		MaxRefresh:  time.Duration(24*365) * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			// simpan data login (save token)
			fmt.Println("PayloadFunc ")

			if v, ok := data.(*objects.UserLogin); ok {

				tokenResult := jwt.MapClaims{
					identityKey: v.UserID,
					roleIDKey:   v.RoleID,
				}

				return tokenResult
			}

			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			fmt.Println("IdentityHandler ----- ")
			claims := jwt.ExtractClaims(c)

			fmt.Println("extraxt claims---", claims, len(claims))

			return &objects.UserLogin{
				UserID: claims[identityKey].(string),
				RoleID: claims[roleIDKey].(string),
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			//pengecekan token yg sudah disimpan di DB
			fmt.Println("Authorizator ----- ")
			fmt.Println("data tables user------->>", data.(*objects.UserLogin).UserID, data.(*objects.UserLogin).OrganizationID)

			// if data.(*objects.UserLogin).OrganizationID == "" {
			// 	return false
			// }

			if v, ok := data.(*objects.UserLogin); ok {

				fmt.Println("v.UserID------>>>>>>", v.UserID)
				var userData tables.Users

				errc := db.Debug().Scopes(models.SchemaUsr("users")).Where("status is true").First(&userData, "id = ? ", v.UserID).Error
				if errc != nil {
					fmt.Println(errc)
					return false
				}

				fmt.Println("return userData.ID------>>>>>>", userData.ID)
				if userData.ID > 0 {
					return true
				}
			}

			fmt.Println("---false---->>", data)

			return false
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			// pengecekan akun login
			var loginVals objects.Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			fmt.Println("Authenticator ----- ", loginVals)

			var userData tables.Users
			errc := db.Debug().Scopes(models.SchemaUsr("users")).Where("status is true AND role_id in (1, 2)").First(&userData, "lower(email) = lower(?) and is_email_verified=true", loginVals.Email).Error
			if errc != nil {
				fmt.Println(errc)
			}
			fmt.Println("Authenticator userData.RoleID ----------->", userData.RoleID)
			// jika user admin tidak di dalam organization manapunn then is not allowed

			checkPassword := VerifyPassword(loginVals.Password, userData.Password)
			fmt.Println("checkPassword guest ::::", loginVals.Password, userData.Password, checkPassword)
			if checkPassword {
				fmt.Println("getUserData guest---", userData)

				// save tokeN here
				return &objects.UserLogin{
					UserID: strconv.Itoa(userData.ID),
					Email:  userData.Email,
					RoleID: strconv.Itoa(userData.RoleID),
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			fmt.Println("Unauthorized ----- ", code)

			c.JSON(code, gin.H{
				"code":    code,
				"status":  false,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		fmt.Println("Err: ", err)
		return nil
	}

	return authMiddleware
}

func SetupMiddleware(db *gorm.DB) *jwt.GinJWTMiddleware {

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "jwt",
		Key:         []byte("#snapin-new#"),
		Timeout:     time.Duration(24*365) * time.Hour,
		MaxRefresh:  time.Duration(24*365) * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			// simpan data login (save token)
			fmt.Println("PayloadFunc ")

			if v, ok := data.(*objects.UserLogin); ok {

				tokenResult := jwt.MapClaims{
					identityKey:       v.UserID,
					roleIDKey:         v.RoleID,
					organizationIDKey: v.OrganizationID,
				}

				// save token
				// saveToken := db.Debug().Scopes(models.SchemaPublic("users"))
				// var userData models.Users
				// db.Model(&userData).Update("remember_token", tokenResult)

				fmt.Println("dataaaa payload----- ", v.UserID, v.Email, v.RoleID, tokenResult)

				return tokenResult
			}

			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			fmt.Println("IdentityHandler ----- ")
			claims := jwt.ExtractClaims(c)

			fmt.Println("extraxt claims---", claims, len(claims))

			return &objects.UserLogin{
				UserID: claims[identityKey].(string),
				RoleID: claims[roleIDKey].(string),
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			//pengecekan token yg sudah disimpan di DB
			fmt.Println("Authorizator ----- ")
			fmt.Println("data tables user------->>", data.(*objects.UserLogin).UserID, data.(*objects.UserLogin).OrganizationID)

			// if data.(*objects.UserLogin).OrganizationID == "" {
			// 	return false
			// }

			if v, ok := data.(*objects.UserLogin); ok {

				fmt.Println("v.UserID------>>>>>>", v.UserID)
				var userData tables.Users

				errc := db.Debug().Scopes(models.SchemaUsr("users")).Where("is_active is true").First(&userData, "id = ? ", v.UserID).Error
				if errc != nil {
					fmt.Println(errc)
					return false
				}

				fmt.Println("return userData.ID------>>>>>>", userData.ID)
				if userData.ID > 0 {
					return true
				}
			}

			fmt.Println("---false---->>", data)

			return false
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			// pengecekan akun login
			var loginVals objects.Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			fmt.Println("Authenticator ----- ", loginVals)

			var userData tables.Users
			errc := db.Debug().Scopes(models.SchemaUsr("users")).Where("is_active is true AND role_id in (1)").First(&userData, "username = ?", loginVals.Username).Error
			if errc != nil {
				fmt.Println(errc)
			}
			fmt.Println("Authenticator userData.RoleID ----------->", userData.RoleID)
			// jika user admin tidak di dalam organization manapunn then is not allowed

			checkPassword := VerifyPassword(loginVals.Password, userData.Password)
			fmt.Println("checkPassword ::::", loginVals.Password, userData.Password, checkPassword)
			if checkPassword {
				fmt.Println("getUserData---", userData)

				// save tokeN here
				return &objects.UserLogin{
					UserID: strconv.Itoa(userData.ID),
					Email:  userData.Username,
					RoleID: strconv.Itoa(userData.RoleID),
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			fmt.Println("Unauthorized ----- ", code)

			c.JSON(code, gin.H{
				"code":    code,
				"status":  false,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		fmt.Println("Err: ", err)
		return nil
	}

	return authMiddleware
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
