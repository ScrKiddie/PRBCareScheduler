package business

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"prb_care_scheduler/internal/constant"
	"prb_care_scheduler/internal/entity"
	"time"
)

func BatalkanStatusProlanisAktif(ctx context.Context, db *gorm.DB) error {
	tx := db.WithContext(ctx).Begin()
	defer tx.Rollback()

	t := time.Now()
	now := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()

	var prolanis []entity.Prolanis
	fmt.Println(now)
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
