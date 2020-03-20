package applyLogin

import (
	"RAS/myUtils"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
)

func SendSMS(msg *MessageContent) (resp string, err error) {
	//Load IdKey from file
	filename := "/IdKey/SendSMS/config.json"
	idKey, err := myUtils.LoadIdKey(filename)
	if err != nil {
		return "", fmt.Errorf("queryIP Load IdKey from file error: %v", err)
	}

	credential := common.NewCredential(
		idKey.SecretId,
		idKey.SecretKey,
	)
	/////////////////////////////////////////////
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := sms.NewClient(credential, "", cpf)

	request := sms.NewSendSmsRequest()

	request.PhoneNumberSet = msg.PhoneNumberSet
	request.TemplateID = msg.TemplateID
	request.Sign = msg.Sign
	request.TemplateParamSet = msg.TemplateParamSet
	request.SmsSdkAppid = msg.SmsSdkAppid

	response, err := client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return "", fmt.Errorf("client.SendSms An API error has returned: %v", err)
	}
	if err != nil {
		return "", fmt.Errorf("client.SendSms unknown type error has returned: %v", err)
	}
	return response.ToJsonString(), nil
}

type MessageContent struct {
	PhoneNumberSet   []*string
	TemplateID       *string
	Sign             *string
	TemplateParamSet []*string
	SmsSdkAppid      *string
}

func DefaultMessageContent() *MessageContent {
	templateID := "549292"
	sign := "海天新亚"
	smsSdkAppid := "1400327700"

	msg := MessageContent{
		PhoneNumberSet:   nil,
		TemplateID:       &templateID,
		Sign:             &sign,
		TemplateParamSet: nil,
		SmsSdkAppid:      &smsSdkAppid,
	}
	return &msg
}
