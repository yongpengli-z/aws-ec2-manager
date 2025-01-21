package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type CreateParams struct {
	DryRun bool `json:"dryRun" form:"dryRun"`
	// InstanceName string `json:"instanceName" form:"instanceName"`
	UserName       string `json:"username" form:"username" binding:"required"`
	Source         string
	ImageId        string `json:"imageId" form:"imageId"`
	InstanceType   string `json:"instanceType" form:"instanceType"`
	RootDiskSize   int32  `json:"rootDiskSize" form:"rootDiskSize"`
	DeviceDiskSize int32  `json:"deviceDiskSize" form:"deviceDiskSize"`
	DeviceName     string
	Department     string `json:"department" form:"department" binding:"required"`
}

type EC2CreateInstanceAPI interface {
	RunInstances(ctx context.Context,
		params *ec2.RunInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)

	CreateTags(ctx context.Context,
		params *ec2.CreateTagsInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)
}

// Ec2Manager godoc
// @Summary     Create a ec2 instance
// @Description Cretae a rc2 instance
// @Tags        Instance
// @Accept      json
// @Produce     json
// @Param       username          query    string   true    "user name"
// @Param       imageId           query    string   false   "instance image id"
// @Param       instanceType      query    string   false   "instance type"
// @Param       rootDiskSize      query    int32    false   "instance rootDiskSize"
// @Param       deviceDiskSize    query    int32    false   "instance deviceDiskSize"
// @Param       department        query    string   true    "user department"
// @Success     200               {object} map[string]any
// @Router      /create           [post]
// @Security    Bearer
func CreateInstance(c *gin.Context, cfg *aws.Config) {

	client := ec2.NewFromConfig(*cfg)

	createParam := &CreateParams{
		Source:         "qtp",
		ImageId:        "ami-09ac7e749b0a8d2a1",
		InstanceType:   "t2.micro",
		DryRun:         false,
		RootDiskSize:   20,
		DeviceDiskSize: 0,
		DeviceName:     "/dev/sdb",
	}

	err := c.ShouldBind(createParam)
	if err != nil {
		log.Println(err)
		ReturnErrorBody(c, 1, "Request parameter invalid.", err)
		return
	}

	log.Println("params: ", createParam)

	keyName := createParam.UserName + "-" + createParam.Source

	rootDeviceName, err := getAmiImageDevice(client, createParam.ImageId)
	if err != nil {
		log.Println(err)
		ReturnErrorBody(c, 1, "Got an error getting image root device name.", err)
		return
	}

	rootDevice := &types.BlockDeviceMapping{
		DeviceName: aws.String(rootDeviceName),
		Ebs: &types.EbsBlockDevice{
			DeleteOnTermination: aws.Bool(true),
			VolumeSize:          aws.Int32(createParam.RootDiskSize),
			VolumeType:          "gp3",
		},
	}

	// blockDeviceMapping[0] = *rootDevice
	blockDeviceMapping := []types.BlockDeviceMapping{*rootDevice}

	if createParam.DeviceDiskSize != 0 {
		blockDevice := &types.BlockDeviceMapping{
			DeviceName: aws.String(createParam.DeviceName),
			Ebs: &types.EbsBlockDevice{
				DeleteOnTermination: aws.Bool(true),
				VolumeSize:          aws.Int32(createParam.DeviceDiskSize),
				VolumeType:          "gp3",
			},
		}
		blockDeviceMapping = append(blockDeviceMapping, *blockDevice)
	}

	input := &ec2.RunInstancesInput{
		ImageId:             aws.String(createParam.ImageId),
		InstanceType:        types.InstanceType(createParam.InstanceType),
		MinCount:            aws.Int32(1),
		MaxCount:            aws.Int32(1),
		BlockDeviceMappings: blockDeviceMapping,
		SecurityGroupIds:    []string{"sg-0f044b5adec791eb3"},
		// SecurityGroups:   []string{"ec2-default-sg"},
		SubnetId: aws.String("subnet-020e66b15e0965ace"),
		KeyName:  aws.String(keyName),
		DryRun:   &createParam.DryRun,
	}

	data := make(map[string]string)

	result, err := MakeInstance(context.TODO(), client, input)
	if err != nil {
		log.Println("Got an error creating an instance:", err)
		ReturnErrorBody(c, 1, "Got an error creating an instance.", err)
		return
	} else {
		data["instanceId"] = fmt.Sprintf("%v", *result.Instances[0].InstanceId)
		data["instanceIP"] = fmt.Sprintf("%v", *result.Instances[0].PrivateIpAddress)
		log.Println("Instance created, instanceId: ", data["instanceId"])
	}

	tagInput := &ec2.CreateTagsInput{
		Resources: []string{data["instanceId"]},
		Tags: []types.Tag{
			{
				Key:   aws.String("cost/owner"),
				Value: aws.String(createParam.UserName),
			},
			{
				Key:   aws.String("cost/business-name"),
				Value: aws.String(createParam.Source),
			},
			{
				Key:   aws.String("cost/service-name"),
				Value: aws.String("ec2"),
			},
			{
				Key:   aws.String("cost/env"),
				Value: aws.String("uat"),
			},
			{
				Key:   aws.String("cost/org/" + createParam.Department),
				Value: aws.String("true"),
			},
		},
	}

	_, err = MakeTags(context.TODO(), client, tagInput)

	if err != nil {
		log.Println("Got an error tagging the instance: ", err)
		message := "Instance created, but tagging the instance failed. instanceId: " + data["instanceId"] + "IP: " + data["instanceIP"]
		ReturnErrorBody(c, 2, message, err)
		return
	} else {
		log.Println("tagging the instance success")
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})

	return
}

func getAmiImageDevice(client *ec2.Client, imageId string) (string, error) {
	input := &ec2.DescribeImagesInput{
		ImageIds: []string{imageId},
	}
	res, err := client.DescribeImages(context.TODO(), input)
	if err != nil {
		return "", err
	}
	if len(res.Images) > 0 {
		imageImformation := res.Images[0]
		if len(imageImformation.BlockDeviceMappings) > 0 {
			fmt.Println("deviceName: ", *imageImformation.BlockDeviceMappings[0].DeviceName)
			return *imageImformation.BlockDeviceMappings[0].DeviceName, nil
		}
	}
	// fmt.Println(res)
	return "", fmt.Errorf("Get no imformation for the image")
}

func MakeInstance(ctx context.Context, api EC2CreateInstanceAPI, input *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {
	return api.RunInstances(ctx, input)
}

func MakeTags(ctx context.Context, api EC2CreateInstanceAPI, input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return api.CreateTags(ctx, input)
}
