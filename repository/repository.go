package repository

import (
	"errors"
	"time"

	"github.com/mohan7-code/url-shortener/models"
	context "github.com/mohan7-code/url-shortener/utils/context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IURLRepository interface {
	Create(ctx *context.Context, url *models.URL) error
	GetUrlByShortCode(ctx *context.Context, shortCode string) (*models.URL, error)
	GetByOriginalURL(ctx *context.Context, originalURL string) (*models.URL, error)
	IncrementClickCount(ctx *context.Context, id string) error
	IncrementClickCountByShortCode(ctx *context.Context, code string) error
	ListURLs(ctx *context.Context, limit, offset int) ([]*models.URL, int64, error)
}

type urlRepository struct {
}

func NewURLRepository() IURLRepository {
	return &urlRepository{}
}

func (r *urlRepository) getTable() string {
	return "url_shortner"
}

func (r *urlRepository) Create(ctx *context.Context, url *models.URL) error {
	err := ctx.DB.WithContext(ctx).Table(r.getTable()).Save(url).Error
	if err != nil {
		ctx.Log.Error("failed to create short url", zap.Error(err))
		return err
	}
	return nil
}

func (r *urlRepository) GetUrlByShortCode(ctx *context.Context, shortCode string) (*models.URL, error) {
	var url models.URL
	err := ctx.DB.Debug().WithContext(ctx).Table(r.getTable()).Where("short_code = ?", shortCode).First(&url).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.Log.Warn("short code not found", zap.Any("short_code", shortCode))
			return nil, nil
		}
		ctx.Log.Error("failed to get short url by short code", zap.Any("short_code", shortCode), zap.Error(err))
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) GetByOriginalURL(ctx *context.Context, originalURL string) (*models.URL, error) {
	var url models.URL
	err := ctx.DB.Table(r.getTable()).Where("original_url = ?", originalURL).First(&url).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		ctx.Log.Error("failed to find by original url", zap.String("original_url", originalURL), zap.Error(err))
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) IncrementClickCount(ctx *context.Context, id string) error {
	err := ctx.DB.WithContext(ctx).Table(r.getTable()).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"click_count":      gorm.Expr("click_count + ?", 1),
			"last_accessed_at": time.Now(),
		}).Error

	if err != nil {
		ctx.Log.Error("failed to increment click count", zap.String("id", id), zap.Error(err))
		return err
	}
	return nil
}

func (r *urlRepository) IncrementClickCountByShortCode(ctx *context.Context, code string) error {

	err := ctx.DB.WithContext(ctx).Table(r.getTable()).
		Where("short_code = ?", code).
		Updates(map[string]interface{}{
			"click_count":      gorm.Expr("click_count + ?", 1),
			"last_accessed_at": time.Now(),
		}).Error

	if err != nil {
		ctx.Log.Error("failed to increment click count", zap.String("code", code), zap.Error(err))
		return err
	}
	return nil
}

func (r *urlRepository) ListURLs(ctx *context.Context, limit, offset int) ([]*models.URL, int64, error) {
	var urls []*models.URL
	var total int64

	query := ctx.DB.WithContext(ctx).Table(r.getTable())

	if err := query.Count(&total).Error; err != nil {
		ctx.Log.Error("failed to count urls", zap.Error(err))
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&urls).Error
	if err != nil {
		ctx.Log.Error("failed to list urls", zap.Error(err))
		return nil, 0, err
	}

	return urls, total, nil
}
