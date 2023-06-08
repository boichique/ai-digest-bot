package sources

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

const baseURL = "http://server:10000"

type Module struct {
	Handler    *Handler
	Service    *Service
	Repository *Repository
	Client     *Client
}

func NewModule(db *pgxpool.Pool) *Module {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)
	client := NewClient(baseURL)

	return &Module{
		Handler:    handler,
		Service:    service,
		Repository: repository,
		Client:     client,
	}
}
