package repository

import (
	"context"
	"youapp/entity"

	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db}
}

func (b *ProfileRepository) GetAllProfile(ctx context.Context, userId uint) ([]entity.Profile, error) {
	var profileResult []entity.Profile

	profile, err := b.db.
		WithContext(ctx).
		Table("profiles").
		Select("*").
		Where("user_id = ? AND deleted_at IS NULL", userId).
		Rows()
	if err != nil {
		return []entity.Profile{}, err
	}
	defer profile.Close()

	for profile.Next() {
		b.db.ScanRows(profile, &profileResult)
	}

	return profileResult, nil
}

func (b *ProfileRepository) AddProfile(ctx context.Context, profile entity.Profile) (entity.Profile, error) {
	err := b.db.
		WithContext(ctx).
		Create(&profile).Error
	return profile, err
}

func (b *ProfileRepository) GetProfileByID(ctx context.Context, id int) (entity.Profile, error) {
	var profileResult entity.Profile

	err := b.db.
		WithContext(ctx).
		Table("profiles").
		Where("id = ? AND deleted_at IS NULL", id).
		Find(&profileResult).Error
	if err != nil {
		return entity.Profile{}, err
	}

	return profileResult, nil
}

func (b *ProfileRepository) DeleteProfile(ctx context.Context, id int) error {
	err := b.db.
		WithContext(ctx).
		Delete(&entity.Profile{}, id).Error
	return err
}

func (b *ProfileRepository) UpdateProfile(ctx context.Context, profile entity.Profile) error {
	err := b.db.
		WithContext(ctx).
		Table("profiles").
		Where("id = ?", profile.ID).
		Updates(&profile).Error
	return err
}
