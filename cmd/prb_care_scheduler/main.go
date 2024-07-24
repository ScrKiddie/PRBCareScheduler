package main

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
	"prb_care_scheduler/internal/business"
	"prb_care_scheduler/internal/config"
)

func main() {
	conf := config.NewViper()
	db := config.NewDatabase(conf)
	client := config.NewFirebase()
	ctx := context.Background()

	logger := cron.VerbosePrintfLogger(log.New(log.Writer(), "cron: ", log.LstdFlags))
	c := cron.New(cron.WithLogger(logger))

	_, err := c.AddFunc("0 9,12,15 * * *", func() {
		if err := business.NotifyStatusKontrolBalikMenunggu(ctx, db, client); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Fatalln(err)
	}

	_, err = c.AddFunc("0 9,12,15 * * *", func() {
		if err := business.NotifyStatusPengambilanObatMenunggu(ctx, db, client); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Fatalln(err)
	}

	_, err = c.AddFunc("0 0 * * *", func() {
		if err := business.BatalkanStatusKontrolBalikMenunggu(ctx, db); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Fatalln(err)
	}

	_, err = c.AddFunc("0 0 * * *", func() {
		if err := business.BatalkanStatusPengambilanObatMenunggu(ctx, db); err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Fatalln(err)
	}

	c.Start()

	select {}
}
