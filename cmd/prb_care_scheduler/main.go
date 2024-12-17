package main

import (
	"context"
	"github.com/robfig/cron/v3"
	"golang.org/x/exp/slog"
	"log"
	"prb_care_scheduler/internal/business"
	"prb_care_scheduler/internal/config"
)

func main() {
	conf := config.NewViper()
	db := config.NewDatabase(conf)
	client := config.NewFirebase()
	ctx := context.Background()

	logger := cron.VerbosePrintfLogger(log.New(log.Writer(), "", log.LstdFlags))
	c := cron.New(cron.WithLogger(logger))

	_, err := c.AddFunc(conf.GetString("time.notify.kontrol"), func() {
		if err := business.NotifyStatusKontrolBalikMenunggu(ctx, db, client); err != nil {
			slog.Warn("failed to execute NotifyStatusKontrolBalikMenunggu(): " + err.Error())
		}
	})
	if err != nil {
		log.Fatalln("failed to add NotifyStatusKontrolBalikMenunggu(): " + err.Error())
	}

	_, err = c.AddFunc(conf.GetString("time.notify.obat"), func() {
		if err := business.NotifyStatusPengambilanObatMenunggu(ctx, db, client); err != nil {
			slog.Warn("failed to execute NotifyStatusPengambilanObatMenunggu(): " + err.Error())
		}
	})
	if err != nil {
		log.Fatalln("failed to add NotifyStatusPengambilanObatMenunggu(): " + err.Error())
	}

	_, err = c.AddFunc(conf.GetString("time.cancel"), func() {
		if err := business.BatalkanStatusKontrolBalikMenunggu(ctx, db); err != nil {
			slog.Warn("failed to execute BatalkanStatusKontrolBalikMenunggu(): " + err.Error())
		}
	})
	if err != nil {
		log.Fatalln("failed to add BatalkanStatusKontrolBalikMenunggu(): " + err.Error())
	}

	_, err = c.AddFunc(conf.GetString("time.cancel"), func() {
		if err := business.BatalkanStatusPengambilanObatMenunggu(ctx, db); err != nil {
			slog.Warn("failed to execute BatalkanStatusPengambilanObatMenunggu(): " + err.Error())
		}
	})
	if err != nil {
		log.Fatalln("failed to add BatalkanStatusPengambilanObatMenunggu(): " + err.Error())
	}

	_, err = c.AddFunc(conf.GetString("time.notify.prolanis"), func() {
		if err := business.NotifyProlanis(ctx, db, client); err != nil {
			slog.Warn("failed to execute NotifyProlanis(): " + err.Error())
		}
	})
	if err != nil {
		log.Fatalln("failed to add NotifyProlanis(): " + err.Error())
	}

	_, err = c.AddFunc(conf.GetString("time.cancel"), func() {
		if err := business.BatalkanStatusProlanisAktif(ctx, db); err != nil {
			slog.Warn("failed to execute BatalkanStatusProlanisAktif(): " + err.Error())
		}
	})
	if err != nil {
		log.Fatalln("failed to add BatalkanStatusKontrolBalikMenunggu(): " + err.Error())
	}

	c.Start()

	select {}
}
