package business

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"prb_care_scheduler/internal/constant"
	"prb_care_scheduler/internal/entity"
	"prb_care_scheduler/internal/helper"
	"strconv"
	"time"
)

func NotifyProlanis(ctx context.Context, db *gorm.DB, client *messaging.Client) error {
	tx := db.WithContext(ctx).Begin()
	defer tx.Rollback()

	var prolanis []entity.Prolanis
	t := time.Now()
	now := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	tomorrow := now + 86400
	overmorrow := now + (2 * 86400)
	err := tx.
		Where("status = ?", constant.StatusProlanisAktif).
		Where("waktu_mulai  BETWEEN ? AND ?", tomorrow, overmorrow).
		Preload("AdminPuskesmas").
		Find(&prolanis).Error

	if err != nil {
		return err
	}

	if len(prolanis) == 0 {
		slog.Info("no prolanis status meets the condition for notification")
		return nil
	}

	var tokenPerangkat []string
	err = tx.Model(&entity.Pengguna{}).
		Where("token_perangkat != ''").
		Pluck("token_perangkat", &tokenPerangkat).Error
	if err != nil {
		return err
	}

	batchSize := 500
	totalTokens := len(tokenPerangkat)

	for _, p := range prolanis {
		data := map[string]string{
			"namaPuskesmas":  p.AdminPuskesmas.NamaPuskesmas,
			"tanggalMulai":   strconv.FormatInt(p.WaktuMulai, 10),
			"tanggalSelesai": strconv.FormatInt(p.WaktuSelesai, 10),
		}

		for i := 0; i < totalTokens; i += batchSize {
			end := i + batchSize
			if end > totalTokens {
				end = totalTokens
			}
			tokenBatch := tokenPerangkat[i:end]

			if err := helper.SendNotificationBroadcastData(ctx, client, data, tokenBatch); err != nil {
				slog.Info(fmt.Sprintf("failed to send notification data for batch: %v, error: %s", tokenBatch, err.Error()))
			} else {
				slog.Info(fmt.Sprintf("successfully sent notification for batch: %v", tokenBatch))
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func BatalkanStatusProlanisAktif(ctx context.Context, db *gorm.DB) error {
	tx := db.WithContext(ctx).Begin()
	defer tx.Rollback()

	t := time.Now()
	now := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()

	var prolanis []entity.Prolanis
	if err := tx.Where("status = ? AND waktu_selesai < ?", constant.StatusProlanisAktif, now).Find(&prolanis).Error; err != nil {
		return err
	}

	if len(prolanis) == 0 {
		slog.Info("no prolanis status meets the condition for cancellation")
		return nil
	}

	for _, p := range prolanis {
		if err := tx.Model(&p).Update("status", constant.StatusProlanisSelesai).Error; err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
