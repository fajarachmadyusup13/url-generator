package repository

import (
	"context"
	"errors"

	"github.com/fajarachmadyusup13/url-generator/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type urlRepository struct {
	db *gorm.DB
}

func NewUrlRepository(db *gorm.DB) model.UrlRepository {
	return &urlRepository{
		db: db,
	}
}

func (ur *urlRepository) Create(ctx context.Context, url *model.Url) (*model.Url, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"url": url,
	})

	tx := ur.db.WithContext(ctx).Begin()
	url.GenerateSlug()
	err := tx.Create(url).Error
	if err != nil {
		logger.Error(err)
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return url, err
}

func (ur *urlRepository) UpdateByID(ctx context.Context, url *model.Url) (*model.Url, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"url": url,
	})

	oldUrl, err := ur.FindByID(ctx, url.ID)
	if err != nil {
		return nil, err
	}

	if oldUrl == nil {
		return nil, nil
	}

	tx := ur.db.WithContext(ctx).Begin()
	err = tx.Model(url).Omit(url.ImmutableColumns()...).Save(url).Error
	if err != nil {
		tx.Rollback()
		logger.Error(err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return ur.FindByID(ctx, url.ID)
}

func (ur *urlRepository) FindByID(ctx context.Context, id int64) (*model.Url, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"id":  id,
	})

	var url model.Url
	err := ur.db.WithContext(ctx).Take(&url, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Error(err)
		return nil, err
	}
	return &url, err
}

func (ur *urlRepository) FindBySlug(ctx context.Context, slug string) (*model.Url, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":  ctx,
		"slug": slug,
	})

	var url model.Url

	err := ur.db.Where(&model.Url{Slug: slug}).Take(&url).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Error(err)
		return nil, err
	}

	return &url, nil
}

func (ur *urlRepository) DeleteByID(ctx context.Context, id int64) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"id":  id,
	})

	var url model.Url
	tx := ur.db.WithContext(ctx).Begin()
	err := tx.Find(&url, id).Delete(&url).Error
	if err != nil {
		logger.Error(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Error(err)
		return err
	}

	err = ur.db.WithContext(ctx).Unscoped().Find(&url, id).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (ur *urlRepository) DeleteBySlug(ctx context.Context, slug string) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":  ctx,
		"slug": slug,
	})

	var url model.Url
	tx := ur.db.WithContext(ctx).Begin()
	err := tx.Where(&model.Url{Slug: slug}).Delete(&url).Error
	if err != nil {
		logger.Error(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Error(err)
		return err
	}

	err = ur.db.WithContext(ctx).Unscoped().Where(&model.Url{Slug: slug}).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
