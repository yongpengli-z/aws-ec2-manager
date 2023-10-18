package controllers

import (
	"context"
	// "bytes"
	// "encoding/json"
	// "fmt"
	"github.com/gin-gonic/gin"
	// "io"
	// "io/ioutil"
	"net/http"
	// "os"

	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type CreateParams struct {
	DryRun          bool   `json:"dryRun" form:"dryRun"`
	InstanceName    string `json:"instanceName" form:"instanceName" binding:"required"`
	UserName        string `json:"userName" form:"userName" binding:"required"`
	Source          string `json:"source" form:"source"`
	InstanceImageId string `json:"instanceImageId" form:"instanceImageId"`
	InstanceType    string `json:"instanceType" form:"instanceType"`
	CuSize          string `json:"cuSize" form:"cuSize"`
	DiskSize        string `json:"diskSize" form:"diskSize"`
}

type EC2CreateInstanceAPI interface {
	RunInstances(ctx context.Context,
		params *ec2.RunInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)

	CreateTags(ctx context.Context,
		params *ec2.CreateTagsInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)
}

// SendNotification godoc
// @Summary     Create a ec2 instance
// @Description Cretae a rc2 instance
// @Tags        Instance
// @Accept      json
// @Produce     json
// @Param       instanceName      query    string true    "instance name"
// @Param       userName          query    string true    "user name"
// @Param       source            query    string false   "source"
// @Param       instanceImageId   query    string false   "instance image id"
// @Param       instanceType      query    string false   "instance type"
// @Param       cuSize            query    string false   "instance"
// @Param       diskSize          query    string false   "instance"
// @Param       dryRun            query    bool   fale    "if dry run"
// @Success     200       {object} map[string]any
// @Router      /create         [post]
// @Security    Bearer
func CreateInstance(c *gin.Context, cfg *aws.Config) {
	// dryRun := true
	minMaxCount := int32(1)

	client := ec2.NewFromConfig(*cfg)

	createParam := &CreateParams{
		Source:          "QTP",
		InstanceImageId: "ami-0556fb70e2e8f34b7",
		InstanceType:    "t2.micro",
		DryRun:          true,
	}

	var typeInstance types.InstanceType
	typeInstance = "t2.micro"

	err := c.ShouldBind(createParam)
	if err != nil {
		log.Println(err)
		ReturnErrorBody(c, 1, "Request parameter invalid.", err)
		return
	}

	log.Println(createParam)

	input := &ec2.RunInstancesInput{
		ImageId:      aws.String(createParam.InstanceImageId),
		InstanceType: typeInstance,
		MinCount:     &minMaxCount,
		MaxCount:     &minMaxCount,
		KeyName:      aws.String(createParam.InstanceName),
		DryRun:       &createParam.DryRun,
	}

	log.Println(input.ImageId)

	result, err := MakeInstance(context.TODO(), client, input)
	if err != nil {
		log.Println("Got an error creating an instance:")
		log.Println(err)
		ReturnErrorBody(c, 1, "Got an error creating an instance.", err)
		return
	}

	log.Println(result)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})

	return
}

func MakeInstance(ctx context.Context, api EC2CreateInstanceAPI, input *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {
	return api.RunInstances(ctx, input)
}

// func Email(c *gin.Context, db *sql.DB) {
// 	e := &EmailParams{}
// 	err := c.ShouldBind(e)
// 	if err != nil {
// 		fmt.Println(err)
// 		ReturnErrorBody(c, 1, "Your request parameter invalid.", err)
// 		return
// 	}

// 	reader, err := e.generateRequestBody()
// 	if err != nil {
// 		fmt.Println(err)
// 		ReturnErrorBody(c, 1, "faild to generate request body.", err)
// 		return
// 	}

// 	responce, err := Post(os.Getenv("NOTIFICATIONSERVER"), "application/json", reader)

// 	// Record send message
// 	status := fmt.Sprintf("%v", responce["Status"])
// 	errRecord := RecordBehavior(c, db, "email", e.Message, e.Receiver, status)
// 	if errRecord != nil {
// 		fmt.Println("record error: ", errRecord)
// 	}

// 	if err != nil {
// 		fmt.Println(err)
// 		ReturnErrorBody(c, 1, "faild to send message.", err)
// 		return
// 	}
// 	if status != "200" {
// 		ReturnErrorBody(c, 1, "faild to send message.", fmt.Errorf("%v", responce["Message"]))
// 		return
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{
// 			"code": 0,
// 			"msg":  "success",
// 		})
// 	}
// 	return
// }

// func (e *EmailParams) generateRequestBody() (io.Reader, error) {
// 	var requestBody map[string]interface{}

// 	requestBody, err := ReadJson("./alert/to_email.json")
// 	if err != nil {
// 		return nil, err
// 	}
// 	email := requestBody["receiver"].(map[string]interface{})["spec"].(map[string]interface{})["email"].(map[string]interface{})
// 	email["to"] = []string{e.Receiver}
// 	if e.Format != "html" {
// 		email["tmplType"] = "text"
// 	} else {
// 		email["tmplType"] = "html"
// 	}

// 	alerts := requestBody["alert"].(map[string]interface{})["alerts"].([]interface{})[0].(map[string]interface{})
// 	alerts["annotations"].(map[string]interface{})["message"] = e.Message
// 	alerts["annotations"].(map[string]interface{})["subject"] = e.Subject

// 	bytesData, _ := json.Marshal(requestBody)
// 	reader := bytes.NewReader(bytesData)

// 	return reader, nil
// }

// func ReadJson(filename string) (map[string]interface{}, error) {
// 	var requestBody map[string]interface{}

// 	// box := packr.NewBox("alert")
// 	// byteValue := box.String("./alert/alert.json")

// 	jsonFile, err := os.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer jsonFile.Close()

// 	byteValue, _ := ioutil.ReadAll(jsonFile)
// 	json.Unmarshal([]byte(byteValue), &requestBody)
// 	// fmt.Println(requestBody)

// 	return requestBody, nil

// }

// func Post(url string, contentType string, jsonFile io.Reader) (map[string]interface{}, error) {
// 	client := http.Client{}
// 	rsp, err := client.Post(url, contentType, jsonFile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rsp.Body.Close()

// 	body, err := ioutil.ReadAll(rsp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var responce map[string]interface{}
// 	json.Unmarshal([]byte(body), &responce)
// 	// fmt.Println(responce["Status"], responce["Message"])
// 	fmt.Println("RSP:", string(body))
// 	return responce, nil
// }

// func RecordBehavior(c *gin.Context, db *sql.DB, mess_type, message, receiver, status string) error {
// 	sqlStr := "INSERT INTO userBehavior(user, application, mess_type, message, receiver, status) values (?, ?, ?, ?, ?, ?);"
// 	userName, ok := c.Get("username")
// 	user := fmt.Sprintf("%v", userName)
// 	if !ok {
// 		return fmt.Errorf("The requested user name is not recognized")
// 	}
// 	appName, ok := c.Get("appname")
// 	app := fmt.Sprintf("%v", appName)
// 	if !ok {
// 		return fmt.Errorf("The requested app name is not recognized")
// 	}
// 	err := models.UserBehavior(db, sqlStr, user, app, mess_type, message, receiver, status)
// 	return err
// }
