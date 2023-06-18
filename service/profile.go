package service

import (
	"context"
	"youapp/entity"
	"youapp/repository"
)

type ProfileService struct {
	bRepo *repository.ProfileRepository
}

func NewProfileService(bRepo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{
		bRepo: bRepo,
	}
}

func (b *ProfileService) GetAllProfile(ctx context.Context, userid uint) ([]entity.Profile, error) {
	return b.bRepo.GetAllProfile(ctx, userid)
}

func (b *ProfileService) AddProfile(ctx context.Context, profile entity.Profile) (entity.Profile, error) {

	res, err := b.bRepo.AddProfile(ctx, profile)
	if err != nil {
		return entity.Profile{}, err
	}
	return res, nil
}

func (b *ProfileService) GetProfileByID(ctx context.Context, id int) (entity.Profile, error) {
	return b.bRepo.GetProfileByID(ctx, id)
}

func (b *ProfileService) UpdateProfile(ctx context.Context, profile entity.Profile) (entity.Profile, error) {
	err := b.bRepo.UpdateProfile(ctx, profile)
	if err != nil {
		return entity.Profile{}, err
	}

	return profile, nil
}

func (b *ProfileService) DeleteProfile(ctx context.Context, id int) error {
	return b.bRepo.DeleteProfile(ctx, id)
}
