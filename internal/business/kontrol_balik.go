package business

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"prb_care_scheduler/internal/constant"
	"prb_care_scheduler/internal/entity"
	"prb_care_scheduler/internal/helper"
	"strconv"
	"time"
)

func NotifyStatusKontrolBalikMenunggu(ctx context.Context, db *gorm.DB, client *messaging.Client) error {
	tx := db.WithContext(ctx).Begin()
	defer tx.Rollback()

	var kontrolBaliks []entity.KontrolBalik
	t := time.Now()
	now := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	lastMonth := now - (30 * 86400)
	nextWeek := now + (7 * 86400)

	err := tx.
		Where("status = ?", constant.StatusKontrolBalikMenunggu).
		Where("tanggal_kontrol BETWEEN ? AND ?", lastMonth, nextWeek).
		Preload("Pasien.Pengguna").
		Preload("Pasien.AdminPuskesmas").
		Find(&kontrolBaliks).Error

	if err != nil {
		return err
	}

	if len(kontrolBaliks) == 0 {
		slog.Info("no kontrol balik status meets the condition for notification")
		return nil
	}

	for _, k := range kontrolBaliks {
		data := map[string]string{
			"namaPuskesmas":  k.Pasien.AdminPuskesmas.NamaPuskesmas,
			"namaLengkap":    k.Pasien.Pengguna.NamaLengkap,
			"tanggalKontrol": strconv.FormatInt(k.TanggalKontrol, 10),
		}
		if err := helper.SendNotificationData(ctx, client, data, k.Pasien.Pengguna.TokenPerangkat); err != nil {
			slog.Info(fmt.Sprintf("failed to send notification data for %s: %s", k.Pasien.Pengguna.TokenPerangkat, err.Error()))
		} else {
			slog.Info("success send notification data for " + k.Pasien.Pengguna.TokenPerangkat)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func BatalkanStatusKontrolBalikMenunggu(ctx context.Context, db *gorm.DB) error {
	tx := db.WithContext(ctx).Begin()
	defer tx.Rollback()

	t := time.Now()
	now := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	lastMonth := now - (30 * 86400)
	var kontrolBaliks []entity.KontrolBalik
	if err := tx.Where("status = ? AND tanggal_kontrol < ?", constant.StatusKontrolBalikMenunggu, lastMonth).Find(&kontrolBaliks).Error; err != nil {
		return err
	}

	if len(kontrolBaliks) == 0 {
		slog.Info("no kontrol balik status meets the condition for cancellation")
		return nil
	}

	for _, k := range kontrolBaliks {
		if err := tx.Model(&k).Update("status", constant.StatusKontrolBalikBatal).Error; err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
