package business

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"prb_care_scheduler/internal/constant"
	"prb_care_scheduler/internal/entity"
	"prb_care_scheduler/internal/helper"
	"strconv"
	"time"
)

func NotifyStatusPengambilanObatMenunggu(ctx context.Context, db *gorm.DB, client *messaging.Client) error {
	tx := db.WithContext(ctx).Begin()
	defer tx.Rollback()

	var pengambilanObats []entity.PengambilanObat
	var pengambilanObat entity.PengambilanObat

	t := time.Now()
	now := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	tomorrow := now + 86400
	fourDaysAgo := now - (4 * 86400)

	var uniqueResis []string
	err := tx.
		Model(&pengambilanObat).
		Select("DISTINCT resi").
		Where("status = ?", constant.StatusPengambilanObatMenunggu).
		Where("tanggal_pengambilan BETWEEN ? AND ?", fourDaysAgo, tomorrow).
		Find(&uniqueResis).Error
	if err != nil {
		return err
	}

	for _, resi := range uniqueResis {
		pengambilanObat.ID = 0
		err = tx.
			Where("resi = ?", resi).
			Preload("Obat").
			Preload("Pasien.Pengguna").
			Preload("Obat.AdminApotek").
			First(&pengambilanObat).Error
		if err != nil {
			return err
		}
		pengambilanObats = append(pengambilanObats, pengambilanObat)
	}

	if len(pengambilanObats) == 0 {
		slog.Info("no waiting pengambilan obat meets the condition for notification")
		return nil
	}

	for _, p := range pengambilanObats {
		data := map[string]string{
			"namaApotek":         p.Obat.AdminApotek.NamaApotek,
			"namaLengkap":        p.Pasien.Pengguna.NamaLengkap,
			"tanggalPengambilan": strconv.FormatInt(p.TanggalPengambilan, 10),
			"tanggalBatal":       strconv.FormatInt(p.TanggalPengambilan+(4*86400), 10),
		}
		if err := helper.SendNotificationData(ctx, client, data, p.Pasien.Pengguna.TokenPerangkat); err != nil {
			slog.Info(fmt.Sprintf("failed to send notification data for %s: %s", p.Pasien.Pengguna.TokenPerangkat, err.Error()))
		} else {
			slog.Info("success send notification data for " + p.Pasien.Pengguna.TokenPerangkat)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func BatalkanStatusPengambilanObatMenunggu(ctx context.Context, db *gorm.DB) error {
	tx := db.WithContext(ctx).Begin()
	defer tx.Rollback()

	t := time.Now()
	now := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	fourDaysAgo := now - (4 * 86400)

	var pengambilanObats []entity.PengambilanObat
	if err := tx.Where("status = ? AND tanggal_pengambilan < ?", constant.StatusPengambilanObatMenunggu, fourDaysAgo).Find(&pengambilanObats).Error; err != nil {
		return err
	}

	if len(pengambilanObats) == 0 {
		slog.Info("no waiting pengambilan obat meets the condition for cancellation")
		return nil
	}

	for _, p := range pengambilanObats {
		if err := tx.Model(&p).Update("status", constant.StatusPengambilanObatBatal).Error; err != nil {
			return err
		}

		var obat entity.Obat
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&obat, p.IdObat).Error; err != nil {
			tx.Rollback()
			return err
		}

		newJumlah := obat.Jumlah + p.Jumlah
		if err := tx.Model(&obat).Update("jumlah", newJumlah).Error; err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
