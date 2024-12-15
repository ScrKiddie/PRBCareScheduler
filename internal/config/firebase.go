package config

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"

	"google.golang.org/api/option"
	"log"
)

func NewFirebase() *messaging.Client {
	opt := option.WithCredentialsFile("fcm.json")
	ctx := context.Context(context.Background())
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}
