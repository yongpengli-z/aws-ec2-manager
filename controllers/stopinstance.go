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

type EC2StopInstancesAPI interface {
	StopInstances(ctx context.Context,
		params *ec2.StopInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.StopInstancesOutput, error)
}

type StopParams struct {
	DryRun     bool   `json:"dryRun" form:"dryRun"`
	InstanceId string `json:"instanceId" form:"instanceId"`
}

// Ec2Manager godoc
// @Summary     Stop a EC2 Instance
// @Description Stop a EC2 Instance.
// @Tags        Instance
// @Accept      json
// @Produce     json
// @Param       instanceId          query    string   true    "instance Id"
// @Success     200               {object} map[string]any
// @Router      /stop              [post]
// @Security    Bearer
func StopInstance(c *gin.Context, cfg *aws.Config) {
	client := ec2.NewFromConfig(*cfg)
	s := &StopParams{
		DryRun: false,
	}

	err := c.ShouldBind(s)
	if err != nil {
		log.Println(err)
		ReturnErrorBody(c, 1, "Request parameter invalid.", err)
		return
	}

	log.Println("params: ", s)

	input := &ec2.StopInstancesInput{
		InstanceIds: []string{
			s.InstanceId,
		},
		DryRun: aws.Bool(s.DryRun),
	}

	result, err := MakeStop(context.TODO(), client, input)
	if err != nil {
		log.Println("Got an error stopping instance: ", err)
		ReturnErrorBody(c, 1, "Got an error stopping instance.", err)
		return
	}

	data := make(map[string]string)
	data["instanceId"] = *result.StoppingInstances[0].InstanceId

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func MakeStop(c context.Context, api EC2StopInstancesAPI, input *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
	resp, err := api.StopInstances(c, input)

	// var apiErr smithy.APIError
	// if errors.As(err, &apiErr) && apiErr.ErrorCode() == "DryRunOperation" {
	// 	fmt.Println("User has permission to stop instances.")
	// 	input.DryRun = aws.Bool(false)
	// 	return api.StopInstances(c, input)
	// }

	return resp, err
}
