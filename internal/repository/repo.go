package repository

import (
	"context"
	"database/sql"
	"time"
	"tz/internal/models"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	Db *gorm.DB
}

func NewSubscriptionRepository() *SubscriptionRepository {
	db, err := gorm.Open(postgres.Open("host=db user=postgres password=postgres dbname=subscriptions port=5432 sslmode=disable"))
	if err != nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil
	}
	if err := sqlDB.Ping(); err != nil {
		return nil
	}
	return &SubscriptionRepository{
		Db: db,
	}
}

func (sr *SubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	err := sr.Db.WithContext(ctx).Create(sub).Error
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (sr *SubscriptionRepository) Retrieve(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	err := sr.Db.WithContext(ctx).Where("id = ?", sub.ID).First(sub).Error
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (sr *SubscriptionRepository) Update(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	err := sr.Db.WithContext(ctx).Where("id = ?", sub.ID).Save(sub).Error
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (sr *SubscriptionRepository) Delete(ctx context.Context, sub *models.Subscription) error {
	err := sr.Db.WithContext(ctx).Where("id = ?", sub.ID).Delete(sub).Error
	if err != nil {
		return err
	}
	return nil
}

func (sr *SubscriptionRepository) List(ctx context.Context, subs *[]models.Subscription) (*[]models.Subscription, error) {
	err := sr.Db.WithContext(ctx).Find(subs).Error
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *SubscriptionRepository) CalculateTotal(
	ctx context.Context,
	startPeriodstr, endPeriodstr, userID, serviceName string,
) (int, error) {
	var total sql.NullInt64
	q := r.Db.WithContext(ctx).Model(&models.Subscription{})
	if startPeriodstr != "" {
		startPeriod, err := time.Parse("01-2006", startPeriodstr)
		if err != nil {
			return 0, err
		}
		q = q.Where("start_date > ?", startPeriod)
	}

	if endPeriodstr != "" {
		endPeriod, err := time.Parse("01-2006", endPeriodstr)
		if err != nil {
			return 0, err
		}
		q = q.Where("start_date < ?", endPeriod)
	}

	if userID != "" {
		q = q.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		q = q.Where("service_name = ?", serviceName)
	}
	err := q.Select("COALESCE(SUM(price), 0)").Scan(&total).Error

	if err != nil {
		log.Error().Err(err)
		return 0, err
	}

	if total.Valid {
		return int(total.Int64), nil
	}
	return 0, nil
}
