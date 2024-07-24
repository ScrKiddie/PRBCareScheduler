package business

import (
	"context"
	"firebase.google.com/go/messaging"
	"gorm.io/gorm"
	"log"
	"prb_care_scheduler/internal/entity"
	"prb_care_scheduler/internal/helper"
	"strconv"
	"time"
)

func NotifyStatusKontrolBalikMenunggu(ctx context.Context, db *gorm.DB, client *messaging.Client) error {
	tx := db.WithContext(ctx).Begin()
	defer tx.Rollback()

	var kontrolBaliks []entity.KontrolBalik
	now := time.Now().Unix()
	tomorrow := now + 86400
	threeDaysAgo := now - 259200

	err := tx.Preload("Pasien.Pengguna").Preload("Pasien.AdminPuskesmas").
		Where("status = ?", "menunggu").
		Where("tanggal_kontrol BETWEEN ? AND ?", threeDaysAgo, tomorrow).
		Find(&kontrolBaliks).Error

	if err != nil {
		return err
	}
	if len(kontrolBaliks) == 0 {
		log.Println("Tidak ada kontrol balik dengan status 'menunggu' yang memenuhi kondisi untuk notifikasi.")
		return nil
	}
	for _, k := range kontrolBaliks {
		data := map[string]string{
			"title":          "PRB Care - Kontrol Balik",
			"namaPuskesmas":  k.Pasien.AdminPuskesmas.NamaPuskesmas,
			"namaLengkap":    k.Pasien.Pengguna.NamaLengkap,
			"tanggalKontrol": strconv.FormatInt(k.TanggalKontrol, 10),
			"tanggalBatal":   strconv.FormatInt(k.TanggalKontrol+259200, 10),
		}
		if err := helper.SendNotificationData(ctx, client, data, k.Pasien.Pengguna.TokenPerangkat); err != nil {
			log.Printf("Failed to send notification data for %s : %s", k.Pasien.Pengguna.TokenPerangkat, err.Error())
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

	now := time.Now().Unix()
	threeDaysAgo := now - 259200

	var kontrolBaliks []entity.KontrolBalik
	if err := tx.Where("status = ? AND tanggal_kontrol < ?", "menunggu", threeDaysAgo).Find(&kontrolBaliks).Error; err != nil {
		return err
	}
	if len(kontrolBaliks) == 0 {
		log.Println("Tidak ada kontrol balik dengan status 'menunggu' yang memenuhi kondisi untuk dibatalkan.")
		return nil
	}
	for _, k := range kontrolBaliks {
		if err := tx.Model(&k).Update("status", "batal").Error; err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
