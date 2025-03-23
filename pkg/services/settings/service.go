package settings

import (
	"context"

	"github.com/opoccomaxao/mylittlefleet/pkg/models"
	"github.com/opoccomaxao/mylittlefleet/pkg/services/settings/repo"
)

type Service struct {
	repo *repo.Repo
}

func NewService(
	repo *repo.Repo,
) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) SetAll(
	ctx context.Context,
	settings map[models.SettingsKey]string,
) error {
	values := make([]*models.Settings, 0, len(settings))
	for key, value := range settings {
		values = append(values, &models.Settings{
			Key:   key,
			Value: value,
		})
	}

	return s.repo.UpdateSettings(ctx, values...)
}

func (s *Service) Set(
	ctx context.Context,
	key models.SettingsKey,
	value string,
) error {
	return s.repo.UpdateSettings(ctx, &models.Settings{
		Key:   key,
		Value: value,
	})
}

func (s *Service) SetBool(
	ctx context.Context,
	key models.SettingsKey,
	value bool,
) error {
	return s.Set(ctx, key, map[bool]string{true: "1", false: "0"}[value])
}

func (s *Service) GetAll(
	ctx context.Context,
	keys ...models.SettingsKey,
) (map[models.SettingsKey]string, error) {
	settings, err := s.repo.GetSettingsByKeys(ctx, keys...)
	if err != nil {
		return nil, err
	}

	res := make(map[models.SettingsKey]string)
	for _, setting := range settings {
		res[setting.Key] = setting.Value
	}

	return res, nil
}

func (s *Service) Get(
	ctx context.Context,
	key models.SettingsKey,
) (string, error) {
	settings, err := s.GetAll(ctx, key)
	if err != nil {
		return "", err
	}

	return settings[key], nil
}

func (s *Service) GetBool(
	ctx context.Context,
	key models.SettingsKey,
) (bool, error) {
	res, err := s.Get(ctx, key)
	if err != nil {
		return false, err
	}

	return res == "1", nil
}

func (s *Service) Delete(
	ctx context.Context,
	key ...models.SettingsKey,
) error {
	return s.repo.DeleteSettingsByKeys(ctx, key...)
}
