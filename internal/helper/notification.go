package helper

import (
	"context"
	"firebase.google.com/go/v4/messaging"
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

func SendNotificationBroadcastData(ctx context.Context, client *messaging.Client, data map[string]string, token []string) error {
	message := &messaging.MulticastMessage{
		Data:   data,
		Tokens: token,
	}

	_, err := client.SendEachForMulticast(ctx, message)
	if err != nil {
		return err
	}
	return nil
}
