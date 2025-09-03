package repository

import (
	"context"
	"tz/internal/models"

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
