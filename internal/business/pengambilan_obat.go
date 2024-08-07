package business

import (
	"context"
	"firebase.google.com/go/messaging"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"prb_care_scheduler/internal/entity"
	"prb_care_scheduler/internal/helper"
	"strconv"
	"time"
)

func NotifyStatusPengambilanObatMenunggu(ctx context.Context, db *gorm.DB, client *messaging.Client) error {
	tx := db.WithContext(ctx).Begin()
	defer tx.Rollback()

	var pengambilanObats []entity.PengambilanObat
	now := time.Now().Unix()
	tomorrow := now + 86400
	threeDaysAgo := now - 259200

	err := tx.Preload("Obat").Preload("Pasien.Pengguna").Preload("Obat.AdminApotek").
		Where("status = ?", "menunggu").
		Where("tanggal_pengambilan BETWEEN ? AND ?", threeDaysAgo, tomorrow).
		Find(&pengambilanObats).Error
	if err != nil {
		return err
	}

	if len(pengambilanObats) == 0 {
		log.Println("Tidak ada pengambilan obat dengan status 'menunggu' yang memenuhi kondisi untuk notifikasi.")
		return nil
	}

	for _, p := range pengambilanObats {
		data := map[string]string{
			"title":              "PRB Care - Pengambilan Obat",
			"namaApotek":         p.Obat.AdminApotek.NamaApotek,
			"namaLengkap":        p.Pasien.Pengguna.NamaLengkap,
			"tanggalPengambilan": strconv.FormatInt(p.TanggalPengambilan, 10),
			"tanggalBatal":       strconv.FormatInt(p.TanggalPengambilan+259200, 10),
		}
		if err := helper.SendNotificationData(ctx, client, data, p.Pasien.Pengguna.TokenPerangkat); err != nil {
			log.Printf("Failed to send notification data for %s : %s", p.Pasien.Pengguna.TokenPerangkat, err.Error())
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

	now := time.Now().Unix()
	threeDaysAgo := now - 259200

	var pengambilanObats []entity.PengambilanObat
	if err := tx.Where("status = ? AND tanggal_pengambilan < ?", "menunggu", threeDaysAgo).Find(&pengambilanObats).Error; err != nil {
		return err
	}

	if len(pengambilanObats) == 0 {
		log.Println("Tidak ada pengambilan obat dengan status 'menunggu' yang memenuhi kondisi untuk dibatalkan.")
		return nil
	}

	for _, p := range pengambilanObats {
		if err := tx.Model(&p).Update("status", "batal").Error; err != nil {
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
