package helper

import (
	"context"
	"firebase.google.com/go/messaging"
)

func SendNotificationData(ctx context.Context, client *messaging.Client, data map[string]string, token string) error {
	message := &messaging.Message{
		Data:  data,
		Token: token,
	}

	_, err := client.Send(ctx, message)
	if err != nil {
		return err
	}
	return nil
}
