package usecase

import (
	"context"
	"time"

	"github.com/fajarachmadyusup13/url-generator/model"
	"github.com/fajarachmadyusup13/url-generator/utils"
	"github.com/sirupsen/logrus"
)

type urlUsecase struct {
	urlRepo model.UrlRepository
}

func NewUrlUsecase(urlRepo model.UrlRepository) model.UrlUsecase {
	return &urlUsecase{
		urlRepo: urlRepo,
	}
}

func (uu *urlUsecase) GenerateShortURL(ctx context.Context, urlStr string) (*model.Url, error) {
	url := &model.Url{
		ID:        utils.GenerateID(),
		Url:       urlStr,
		ExpiredAt: time.Now().Add(1 * time.Hour),
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"url": url,
	})

	res, err := uu.urlRepo.Create(ctx, url)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return res, nil
}

func (uu *urlUsecase) UpdateURL(ctx context.Context, url *model.Url) (*model.Url, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"url": url,
	})

	oldUrl, err := uu.urlRepo.FindByID(ctx, url.ID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if oldUrl == nil {
		return nil, ErrRecordNotFound
	}

	urlBySlug, err := uu.urlRepo.FindBySlug(ctx, url.Slug)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if urlBySlug != nil && oldUrl.ID != urlBySlug.ID {
		return nil, ErrURLAlreadyExists
	}

	res, err := uu.urlRepo.UpdateByID(ctx, url)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

func (uu *urlUsecase) FindByID(ctx context.Context, id int64) (*model.Url, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"id":  id,
	})

	res, err := uu.urlRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if res == nil {
		return nil, ErrRecordNotFound
	}

	return res, nil
}

func (uu *urlUsecase) FindBySlug(ctx context.Context, slug string) (*model.Url, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":  ctx,
		"slug": slug,
	})

	res, err := uu.urlRepo.FindBySlug(ctx, slug)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if res == nil {
		return nil, ErrRecordNotFound
	}

	return res, nil
}
