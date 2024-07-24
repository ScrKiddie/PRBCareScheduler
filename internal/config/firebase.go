package config

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
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
