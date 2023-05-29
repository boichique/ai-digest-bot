package sources

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSource(ctx context.Context, user *User) error {
	return s.repo.CreateSource(ctx, user)
}

func (s *Service) GetUserSourcesByID(ctx context.Context, userID int) (*User, error) {
	return s.repo.GetUserSourcesByID(ctx, userID)
}
