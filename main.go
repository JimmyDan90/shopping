package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"shopping/api"
	_ "shopping/docs"
	"shopping/utils/graceful"
	"time"
)

// @title 电商练手demo
// @description 电商练手demo
// @version 1.0
// @contact.name golang gin mysql
func main() {
	r := gin.Default()
	registerMiddleware(r)
	api.RegisterHandlers(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()
	graceful.ShutdownGin(srv, time.Second*5)
}

func registerMiddleware(r *gin.Engine) {
	r.Use(
		gin.LoggerWithFormatter(
			func(params gin.LogFormatterParams) string {
				return fmt.Sprintf(
					"%s - [%s] \"%s %s %s %d %s %s\" \n",
					params.ClientIP,
					params.TimeStamp.Format(time.RFC3339),
					params.Method,
					params.Request.Proto,
					params.StatusCode,
					params.Latency,
					params.ErrorMessage,
				)
			}))
	r.Use(gin.Recovery())
}
