package api

import (
	"github.com/gin-gonic/gin"
	"log"
	userApi "shopping/api/user"
	"shopping/config"
	"shopping/domain/user"
	"shopping/utils/database_handler"
)

// Database 结构体
type Database struct {
	userRepository *user.Repository
}

var AppConfig = &config.Configuration{}

// CreateDBs 根据配置文件创建数据库
func CreateDBs() *Database {
	cfgFile := "./config/config.yaml"
	conf, err := config.GetAllConfigValues(cfgFile)
	AppConfig = conf
	if err != nil {
		return nil
	}
	if err != nil {
		log.Fatalf("读取配置文件失败，%v", err.Error())
	}
	db := database_handler.NewMySqlDb(AppConfig.DatabaseSettings.DatabaseURI)
	return &Database{
		userRepository: user.NewUserRepository(db),
	}
}

func RegisterHandlers(r *gin.Engine) {
	dbs := *CreateDBs()
	RegisterUserHandlers(r, dbs)
}

func RegisterUserHandlers(r *gin.Engine, dbs Database) {
	userService := user.NewUserService(*dbs.userRepository)
	userController := userApi.NewUserController(userService, AppConfig)
	userGroup := r.Group("/user")
	userGroup.POST("", userController.CreateUser)
	userGroup.POST("/login", userController.Login)
}
