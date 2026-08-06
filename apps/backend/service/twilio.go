package service

import (
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"resuming/systemconfig"
)

func SendSMS(receipient_phone_number, subject, body string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: systemconfig.TWILIOAccountSID,
		Password: systemconfig.TWILIOAuthToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(receipient_phone_number)
	params.SetFrom(systemconfig.TWILIOPhoneNumber)
	params.SetBody(body)

	_, err := client.Api.CreateMessage(params)
	return err
}
