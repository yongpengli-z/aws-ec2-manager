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

type EC2TerminateInstancesAPI interface {
	TerminateInstances(ctx context.Context,
		params *ec2.TerminateInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}

type TerminateParams struct {
	DryRun     bool   `json:"dryRun" form:"dryRun"`
	InstanceId string `json:"instanceId" form:"instanceId" binding:"required"`
}

// Ec2Manager godoc
// @Summary     Terminate a EC2 Instance
// @Description Terminate a EC2 Instance.
// @Tags        Instance
// @Accept      json
// @Produce     json
// @Param       instanceId          query    string   true    "instance Id"
// @Success     200                 {object} map[string]any
// @Router      /terminate          [post]
// @Security    Bearer
func TerminateInstance(c *gin.Context, cfg *aws.Config) {
	client := ec2.NewFromConfig(*cfg)
	t := &TerminateParams{
		DryRun: false,
	}

	err := c.ShouldBind(t)
	if err != nil {
		log.Println(err)
		ReturnErrorBody(c, 1, "Request parameter invalid.", err)
		return
	}

	log.Println("Terminate params: ", t)

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{
			t.InstanceId,
		},
		DryRun: aws.Bool(t.DryRun),
	}

	result, err := MakeTerminate(context.TODO(), client, input)
	// result, err := client.TerminateInstances(context.TODO(), input)
	if err != nil {
		log.Println("Got an error terminateing instance: ", err)
		ReturnErrorBody(c, 1, "Got an error terminateing instance.", err)
		return
	}

	data := make(map[string]string)
	data["instanceId"] = *result.TerminatingInstances[0].InstanceId
	// data["instanceId"] = *result.TerminatingInstances[0].CurrentState

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func MakeTerminate(c context.Context, api EC2TerminateInstancesAPI, input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	resp, err := api.TerminateInstances(c, input)
	return resp, err
}
