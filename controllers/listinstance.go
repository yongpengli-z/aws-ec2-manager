package controllers

import (
	"context"
	// "fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// type EC2DescribeInstancesAPI interface {
// 	DescribeInstances(ctx context.Context,
// 		params *ec2.DescribeInstancesInput,
// 		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
// }

type ListParams struct {
	DryRun   bool   `json:"dryRun" form:"dryRun"`
	UserName string `json:"userName" form:"userName"`
}

// Ec2Manager godoc
// @Summary     List EC2 Instance
// @Description DList EC2 Instance.
// @Tags        Instance
// @Accept      json
// @Produce     json
// @Param       userName          query    string   false    "UserName"
// @Success     200               {object} map[string]any
// @Router      /list              [post]
// @Security    Bearer
func ListInstance(c *gin.Context, cfg *aws.Config) {
	client := ec2.NewFromConfig(*cfg)
	l := &ListParams{
		DryRun:   false,
		UserName: "",
	}

	err := c.ShouldBind(l)
	if err != nil {
		log.Println(err)
		ReturnErrorBody(c, 1, "Request parameter invalid.", err)
		return
	}

	// log.Println("params: ", d)

	filter1 := &types.Filter{
		Name:   aws.String("tag:cost/business/qtp"),
		Values: []string{"true"},
	}

	filter := []types.Filter{*filter1}

	if l.UserName != "" {
		filter2 := &types.Filter{
			Name:   aws.String("tag:cost/owner"),
			Values: []string{l.UserName},
		}
		filter = append(filter, *filter2)
	}

	input := &ec2.DescribeInstancesInput{
		Filters: filter,
		DryRun:  aws.Bool(l.DryRun),
	}

	result, err := MakeDescribe(context.TODO(), client, input)
	if err != nil {
		log.Println("Got an error listing instance: ", err)
		ReturnErrorBody(c, 1, "Got an error listing instance.", err)
		return
	}

	var datas []map[string]interface{}

	log.Println("length of Reservations: ", len(result.Reservations))

	if len(result.Reservations) > 0 {
		for _, res := range result.Reservations {
			log.Println("length of instance: ", len(res.Instances))
			for _, instance := range res.Instances {
				data := make(map[string]interface{})
				data["instanceId"] = instance.InstanceId
				data["state"] = instance.State.Code
				datas = append(datas, data)
			}
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": "None instance found",
		})
		return
	}

	// data := make(map[string]interface{})

	// if len(result.Reservations) > 0 {
	// 	data["instanceId"] = *result.Reservations[0].Instances[0].InstanceId
	// 	data["state"] = *result.Reservations[0].Instances[0].State.Code
	// 	if result.Reservations[0].Instances[0].PrivateIpAddress != nil {
	// 		data["instanceIP"] = *result.Reservations[0].Instances[0].PrivateIpAddress
	// 	} else {
	// 		data["instanceIP"] = "none"
	// 	}
	// } else {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code": 1,
	// 		"msg":  "None instance found",
	// 	})
	// 	return
	// }

	// log.Println("mess: ", data)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": datas,
	})
}

// func MakeDescribe(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
// 	resp, err := api.DescribeInstances(c, input)
// 	return resp, err
// }
