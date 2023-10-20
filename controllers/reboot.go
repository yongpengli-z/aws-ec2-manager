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

type EC2RebootInstancesAPI interface {
	RebootInstances(ctx context.Context,
		params *ec2.RebootInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.RebootInstancesOutput, error)
}

type RebootParams struct {
	DryRun     bool   `json:"dryRun" form:"dryRun"`
	InstanceId string `json:"instanceId" form:"instanceId" binding:"required"`
}

// Ec2Manager godoc
// @Summary     Reboot a EC2 Instance
// @Description Reboot a EC2 Instance.
// @Tags        Instance
// @Accept      json
// @Produce     json
// @Param       instanceId          query    string   true    "instance Id"
// @Success     200               {object} map[string]any
// @Router      /reboot              [post]
// @Security    Bearer
func RebootInstance(c *gin.Context, cfg *aws.Config) {
	client := ec2.NewFromConfig(*cfg)
	r := &RebootParams{
		DryRun: false,
	}

	err := c.ShouldBind(r)
	if err != nil {
		log.Println(err)
		ReturnErrorBody(c, 1, "Request parameter invalid.", err)
		return
	}

	log.Println("Reboot params: ", r)

	input := &ec2.RebootInstancesInput{
		InstanceIds: []string{
			r.InstanceId,
		},
		DryRun: aws.Bool(r.DryRun),
	}

	_, err = MakeReboot(context.TODO(), client, input)
	if err != nil {
		log.Println("Got an error rebooting instance: ", err)
		ReturnErrorBody(c, 1, "Got an error rebooting instance.", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": "",
	})
}

func MakeReboot(c context.Context, api EC2RebootInstancesAPI, input *ec2.RebootInstancesInput) (*ec2.RebootInstancesOutput, error) {
	resp, err := api.RebootInstances(c, input)
	return resp, err
}
