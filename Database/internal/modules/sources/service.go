package sources

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSource(ctx context.Context, source *Source) error {
	return s.repo.CreateSource(ctx, source)
}

func (s *Service) GetUsersIDList(ctx context.Context) ([]string, error) {
	return s.repo.GetUsersIDList(ctx)
}

func (s *Service) GetUserSourcesByID(ctx context.Context, userID int) ([]string, error) {
	return s.repo.GetUserSourcesByID(ctx, userID)
}

func (s *Service) DeleteSourceByLink(ctx context.Context, source *Source) error {
	return s.repo.DeleteSourceByLink(ctx, source)
}
