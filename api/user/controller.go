package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"shopping/config"
	"shopping/domain/user"
	"shopping/utils/api_helper"
	jwtHelper "shopping/utils/jwt"
	"strconv"
	"time"
)

type Controller struct {
	userService *user.Service
	appConfig   *config.Configuration
}

// NewUserController 实例化
func NewUserController(service *user.Service, appConfig *config.Configuration) *Controller {
	return &Controller{
		userService: service,
		appConfig:   appConfig,
	}
}

// CreateUser godoc
// @Summary 根据给定的用户名和密码创建用户
// @Tags Auth
// @Accept json
// @Produce json
// @Param CreateUserRequest body CreateUserRequest true "user information"
// @Success 201 {object} CreateUserResponse
// @Failure 400 {object} api_helper.ErrResponse
// @Router /user [post]
func (c *Controller) CreateUser(g *gin.Context) {
	var req CreateUserRequest
	if err := g.ShouldBind(&req); err != nil {
		api_helper.HandlerError(g, api_helper.ErrInvalidBody)
		return
	}
	newUser := user.NewUser(req.Username, req.Password, req.Password2)
	err := c.userService.Create(newUser)
	if err != nil {
		api_helper.HandlerError(g, err)
		return
	}
	g.JSON(
		http.StatusCreated, CreateUserResponse{
			Username: req.Username,
		})
}

// Login godoc
// @Summary 根据用户名和密码登录
// @Tags Auth
// @Accept json
// @Product json
// @Param LoginRequest body LoginRequest true "user information"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} api_helper.ErrResponse
// Router /user/login [post]
func (c *Controller) Login(g *gin.Context) {
	var req LoginRequest
	if err := g.ShouldBind(&req); err != nil {
		api_helper.HandlerError(g, api_helper.ErrInvalidBody)
	}
	currentUser, err := c.userService.GetUser(req.Username, req.Password)
	if err != nil {
		api_helper.HandlerError(g, err)
		return
	}
	decodedClaims := jwtHelper.VerifyToken(currentUser.Token, c.appConfig.SecretKey)
	if decodedClaims == nil {
		jwtClaims := jwt.NewWithClaims(
			jwt.SigningMethodHS256, jwt.MapClaims{
				"userId":   strconv.FormatInt(int64(currentUser.ID), 10),
				"username": currentUser.Username,
				"iat":      time.Now().Unix(),
				"iss":      os.Getenv("ENV"),
				"exp": time.Now().Add(
					24 * time.Hour).Unix(),
				"isAdmin": currentUser.IsAdmin,
			})
		token := jwtHelper.GenerateToken(jwtClaims, c.appConfig.SecretKey)
		currentUser.Token = token
		err = c.userService.UpdateUser(&currentUser)
		if err != nil {
			api_helper.HandlerError(g, err)
			return
		}
	}
	g.JSON(http.StatusOK, LoginResponse{
		Username: currentUser.Username,
		UserId:   currentUser.ID,
		Token:    currentUser.Token,
	})
}

// VerifyToken 验证Token
func (c *Controller) VerifyToken(g *gin.Context) {
	token := g.GetHeader("Authorization")
	decodedClaims := jwtHelper.VerifyToken(token, c.appConfig.SecretKey)
	g.JSON(http.StatusOK, decodedClaims)
}
