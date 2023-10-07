package routes

import (
	"login_golang/config"
	"login_golang/controllers"
	m "login_golang/middleware"
	"login_golang/models"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// var (
// 	varconf        config.Configurations =
// 	db             *gorm.DB                   = config.ConnectDB(varconf)
// 	userModels     models.UserModels          = models.NewUserModels(db)
// 	userController controllers.UserController = controllers.NewUserController(userModels)
// )

func SetupRoutes(conf config.Configurations) *gin.Engine {

	r := gin.Default()

	var varconf config.Configurations = conf
	var db *gorm.DB = config.ConnectDB(varconf)

	// models
	var userModels models.UserModels = models.NewUserModels(db)

	//controller
	var UserController controllers.UserController = controllers.NewUserController(userModels)

	// CORS Handle
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,content-type,authorization, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Set("db", db)
		c.Next()
	})

	authMiddleware := m.SetupMiddleware(db) // dashboard access

	// test
	// Refresh time can be longer than token timeout
	// user.GET("/refresh_token", authMiddleware.RefreshHandler)

	// public endpoint
	r.GET("/", func(c *gin.Context) { c.JSON(404, gin.H{"code": 200, "message": "Wellcome to Aplikasi Freelace"}) })
	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/register", UserController.RegisterHandler)

	return r
}
