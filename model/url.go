package model

import (
	"context"
	"time"

	"github.com/fajarachmadyusup13/url-generator/utils"
	"github.com/gosimple/slug"
)

type Url struct {
	ID        int64      `json:"id"`
	Url       string     `json:"url"`
	Slug      string     `json:"slug"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ExpiredAt time.Time  `json:"expired_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UrlRepository interface {
	Create(ctx context.Context, url *Url) (*Url, error)
	UpdateByID(ctx context.Context, url *Url) (*Url, error)
	FindByID(ctx context.Context, id int64) (*Url, error)
	FindBySlug(ctx context.Context, slug string) (*Url, error)
	DeleteByID(ctx context.Context, id int64) error
	DeleteBySlug(ctx context.Context, slug string) error
}

type UrlUsecase interface {
	GenerateShortURL(ctx context.Context, url string) (*Url, error)
	UpdateURL(ctx context.Context, url *Url) (*Url, error)
	FindByID(ctx context.Context, id int64) (*Url, error)
	FindBySlug(ctx context.Context, slug string) (*Url, error)
}

func (ur *Url) ImmutableColumns() []string {
	return []string{"created_at", "expired_at"}
}

func (ur *Url) GenerateSlug() {
	ur.Slug = slug.Make(utils.GenerateString())
}
