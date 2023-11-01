package main

import (
	"context"
	// "fmt"
	"github.com/Bennu-Li/aws-ec2-manager/controllers"
	"github.com/Bennu-Li/aws-ec2-manager/docs"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"os"
)

// @title                      Manage EC2 API
// @version                    1.0
// @description                This API is used to manager aws ec2.
// @host                       localhost:8080
// @BasePath                   /v1/ec2
// @securityDefinitions.apikey Bearer
// @in                         header
// @name                       Authorization
func main() {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	docHost := os.Getenv("DOCHOST")
	if docHost != "" {
		docs.SwaggerInfo.Host = docHost
	}

	router := gin.Default()
	group := router.Group("/v1/ec2")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Welcome to use this api to manager ec2!",
		})
	})

	group.POST("/create", func(c *gin.Context) {
		controllers.CreateInstance(c, &cfg)
	})

	group.POST("/key", func(c *gin.Context) {
		controllers.CreateKeyPair(c, &cfg)
	})

	group.POST("/stop", func(c *gin.Context) {
		controllers.StopInstance(c, &cfg)
	})

	group.POST("/start", func(c *gin.Context) {
		controllers.StartInstance(c, &cfg)
	})

	group.POST("/reboot", func(c *gin.Context) {
		controllers.RebootInstance(c, &cfg)
	})

	group.POST("/terminate", func(c *gin.Context) {
		controllers.TerminateInstance(c, &cfg)
	})

	group.POST("/describe", func(c *gin.Context) {
		controllers.DescribeInstance(c, &cfg)
	})

	group.POST("/list", func(c *gin.Context) {
		controllers.ListInstance(c, &cfg)
	})

	// group.POST("/create", controllers.JWTAuthMiddleware(), func(c *gin.Context) {
	// 	controllers.CreateInstance(c, &cfg)
	// })

	router.Run() // 0.0.0.0:8080
}
