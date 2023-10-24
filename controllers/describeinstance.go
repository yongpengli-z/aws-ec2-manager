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

type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

type DescribeParams struct {
	DryRun     bool   `json:"dryRun" form:"dryRun"`
	InstanceId string `json:"instanceId" form:"instanceId" binding:"required"`
}

// Ec2Manager godoc
// @Summary     Describe a EC2 Instance
// @Description Describe a EC2 Instance.
// @Tags        Instance
// @Accept      json
// @Produce     json
// @Param       instanceId          query    string   true    "instance Id"
// @Success     200               {object} map[string]any
// @Router      /describe              [post]
// @Security    Bearer
func DescribeInstance(c *gin.Context, cfg *aws.Config) {
	client := ec2.NewFromConfig(*cfg)
	d := &DescribeParams{
		DryRun: false,
	}

	err := c.ShouldBind(d)
	if err != nil {
		log.Println(err)
		ReturnErrorBody(c, 1, "Request parameter invalid.", err)
		return
	}

	log.Println("params: ", d)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{
			d.InstanceId,
		},
		DryRun: aws.Bool(d.DryRun),
	}

	result, err := MakeDescribe(context.TODO(), client, input)
	if err != nil {
		log.Println("Got an error describeing instance: ", err)
		ReturnErrorBody(c, 1, "Got an error describeing instance.", err)
		return
	}

	data := make(map[string]interface{})

	if len(result.Reservations) > 0 {
		data["instanceId"] = *result.Reservations[0].Instances[0].InstanceId
		data["state"] = *result.Reservations[0].Instances[0].State.Code
		if result.Reservations[0].Instances[0].PrivateIpAddress != nil {
			data["instanceIP"] = *result.Reservations[0].Instances[0].PrivateIpAddress
		} else {
			data["instanceIP"] = "none"
		}

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "None instance found",
		})
		return
	}

	log.Println("mess: ", data)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func MakeDescribe(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	resp, err := api.DescribeInstances(c, input)
	return resp, err
}
