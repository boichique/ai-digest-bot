package sources

import (
	"digest_bot_database/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	Client     *Client
	Handler    *Handler
	Service    *Service
	Repository *Repository
}

func NewModule(db *pgxpool.Pool, cfg *config.Config) *Module {
	client := NewClient(cfg)
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service, client)

	return &Module{
		Client:     client,
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}
