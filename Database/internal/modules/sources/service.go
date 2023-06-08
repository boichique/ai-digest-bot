package sources

import (
	"context"

	"digest_bot_database/internal/log"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSource(ctx context.Context, source *Source) error {
	if err := s.repo.CreateSource(ctx, source); err != nil {
		return err
	}

	log.FromContext(ctx).Info(
		"source created",
		"source", source.Source,
		"forUser", source.UserID,
	)

	return nil
}

func (s *Service) GetUsersIDList(ctx context.Context) ([]string, error) {
	return s.repo.GetUsersIDList(ctx)
}

func (s *Service) GetUserSourcesByUserID(ctx context.Context, userID int) ([]string, error) {
	return s.repo.GetUserSourcesByID(ctx, userID)
}

func (s *Service) DeleteSourceByLink(ctx context.Context, source *Source) error {
	if err := s.repo.DeleteSourceByLink(ctx, source); err != nil {
		log.FromContext(ctx).Error(err.Error())
		return err
	}

	log.FromContext(ctx).Info(
		"source deleted",
		"source", source.Source,
		"forUser", source.UserID,
	)

	return nil
}
