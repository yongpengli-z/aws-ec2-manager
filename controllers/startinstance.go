package controllers

import (
	"context"
	// "fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	// "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2StartInstancesAPI interface {
	StartInstances(ctx context.Context,
		params *ec2.StartInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error)
}

type StartParams struct {
	DryRun     bool   `json:"dryRun" form:"dryRun"`
	InstanceId string `json:"instanceId" form:"instanceId" binding:"required"`
}

// Ec2Manager godoc
// @Summary     Start a EC2 Instance
// @Description Start a EC2 Instance.
// @Tags        Instance
// @Accept      json
// @Produce     json
// @Param       instanceId          query    string   true    "instance Id"
// @Success     200               {object} map[string]any
// @Router      /start              [post]
// @Security    Bearer
func StartInstance(c *gin.Context, cfg *aws.Config) {
	client := ec2.NewFromConfig(*cfg)
	s := &StartParams{
		DryRun: false,
	}

	err := c.ShouldBind(s)
	if err != nil {
		log.Println(err)
		ReturnErrorBody(c, 1, "Request parameter invalid.", err)
		return
	}

	log.Println("params: ", s)

	input := &ec2.StartInstancesInput{
		InstanceIds: []string{
			s.InstanceId,
		},
		DryRun: aws.Bool(s.DryRun),
	}

	result, err := MakeStart(context.TODO(), client, input)
	if err != nil {
		log.Println("Got an error startping instance: ", err)
		ReturnErrorBody(c, 1, "Got an error startping instance.", err)
		return
	}

	data := make(map[string]string)
	data["instanceId"] = *result.StartingInstances[0].InstanceId

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func MakeStart(c context.Context, api EC2StartInstancesAPI, input *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
	resp, err := api.StartInstances(c, input)
	return resp, err
}
