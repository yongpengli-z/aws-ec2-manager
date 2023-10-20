package controllers

import (
	"context"
	// "fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type CreateKeyParams struct {
	DryRun   bool   `json:"dryRun" form:"dryRun"`
	UserName string `json:"username" form:"username" binding:"required"`
	Source   string
}

// Ec2Manager godoc
// @Summary     Create a keyPair
// @Description Cretae a rsa keyPair for the user.
// @Tags        Instance
// @Accept      json
// @Produce     json
// @Param       username          query    string   true    "user name"
// @Success     200               {object} map[string]any
// @Router      /key              [post]
// @Security    Bearer
func CreateKeyPair(c *gin.Context, cfg *aws.Config) {

	client := ec2.NewFromConfig(*cfg)

	k := &CreateKeyParams{
		Source: "qtp",
		DryRun: false,
	}

	err := c.ShouldBind(k)
	if err != nil {
		log.Println(err)
		ReturnErrorBody(c, 1, "Request parameter invalid.", err)
		return
	}

	log.Println("params: ", k)

	keyName := k.UserName + "-" + k.Source

	input := &ec2.CreateKeyPairInput{
		KeyName:   aws.String(keyName),
		DryRun:    aws.Bool(k.DryRun),
		KeyFormat: "pem",
		KeyType:   "rsa",
	}

	// result, err := MakeKeyPairs(context.TODO(), client, input)
	result, err := client.CreateKeyPair(context.TODO(), input)
	if err != nil {
		log.Println("Got an error creating a key pair: ", err)
		ReturnErrorBody(c, 1, "Got an error creating a key pair.", err)
		return
	}

	data := make(map[string]string)
	data["keyName"] = *result.KeyName
	data["keyMaterial"] = *result.KeyMaterial

	log.Println("Key name: ", *result.KeyName)
	log.Println("KeyMaterial: ", *result.KeyMaterial)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}
