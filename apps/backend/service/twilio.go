package service

import (
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	systemconfig "resuming/system-config"
)

func SendSMS(receipient_phone_number, subject, body string) error {
	client := twilio.NewRestClient()

	params := &openapi.CreateMessageParams{}
	params.SetTo(receipient_phone_number)
	params.SetFrom(systemconfig.TWILIOPhoneNumber)
	params.SetBody(body)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return err
	}

	return nil
}
