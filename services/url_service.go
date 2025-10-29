package service

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mohan7-code/url-shortener/dtos"
	"github.com/mohan7-code/url-shortener/models"
	"github.com/mohan7-code/url-shortener/repository"
	"github.com/mohan7-code/url-shortener/utils/cache"
	context "github.com/mohan7-code/url-shortener/utils/context"
	helper "github.com/mohan7-code/url-shortener/utils/helpers"
	"go.uber.org/zap"
)

type IURLService interface {
	ShortenURL(ctx *context.Context, originalURL string) (*models.URL, error)
	GetOriginalURL(ctx *context.Context, shortCode string) (*models.URL, error)
	ListURLs(ctx *context.Context, page, limit int) (*dtos.ListResponse, error)
}

type urlServiceImpl struct {
	repo repository.IURLRepository
}

func NewURLService() IURLService {
	return &urlServiceImpl{
		repo: repository.NewURLRepository(),
	}
}

func (s *urlServiceImpl) ShortenURL(ctx *context.Context, originalURL string) (*models.URL, error) {
	if originalURL == "" {
		return nil, errors.New("original URL cannot be empty")
	}

	if !helper.IsValidURL(originalURL) {
		ctx.Log.Warn("invalid URL format", zap.String("url", originalURL))
		return nil, errors.New("invalid URL format — must be a valid  URL")
	}

	rdb := cache.New().Client
	// Check Redis cache first
	if cachedShortCode, err := rdb.Get(ctx, originalURL).Result(); err == nil && cachedShortCode != "" {
		ctx.Log.Info("cache hit for original URL", zap.String("short_code", cachedShortCode))
		return &models.URL{
			OriginalURL: originalURL,
			ShortCode:   cachedShortCode,
		}, nil
	}

	existing, err := s.repo.GetByOriginalURL(ctx, originalURL)
	if err != nil {
		ctx.Log.Error("error checking existing URL", zap.Error(err))
		return nil, err
	}

	if existing != nil && existing.ID != uuid.Nil {
		ctx.Log.Info("url already exists", zap.String("short_code", existing.ShortCode))
		return existing, nil
	}

	var shortCode string
	for {
		shortCode = generateShortCode(originalURL)

		existingCode, err := s.repo.GetUrlByShortCode(ctx, shortCode)
		if err != nil {
			return nil, err
		}

		if existingCode == nil {
			break
		}
	}
	url := &models.URL{
		ID:             uuid.New(),
		ShortCode:      shortCode,
		OriginalURL:    originalURL,
		ClickCount:     0,
		LastAccessedAt: time.Now(),
	}

	err = s.repo.Create(ctx, url)
	if err != nil {
		ctx.Log.Error("failed to create shortened URL", zap.Error(err))
		return nil, err
	}

	rdb.Set(ctx, shortCode, originalURL, 24*time.Hour).Err()
	rdb.Set(ctx, originalURL, shortCode, 24*time.Hour).Err()

	ctx.Log.Info("shortened URL created", zap.String("short_code", shortCode))
	return url, nil
}

func (s *urlServiceImpl) GetOriginalURL(ctx *context.Context, shortCode string) (*models.URL, error) {
	if strings.TrimSpace(shortCode) == "" {
		return nil, errors.New("short code cannot be empty")
	}

	rdb := cache.New().Client
	if cachedURL, err := rdb.Get(ctx, shortCode).Result(); err == nil && cachedURL != "" {
		ctx.Log.Info("cache hit for short code", zap.String("short_code", shortCode))
		return &models.URL{OriginalURL: cachedURL}, nil
	}

	url, err := s.repo.GetUrlByShortCode(ctx, shortCode)
	if err != nil {
		ctx.Log.Error("failed to fetch original URL", zap.Error(err))
		return nil, err
	}
	if url == nil {
		return nil, errors.New("short code not found")
	}

	rdb.Set(ctx, shortCode, url.OriginalURL, 24*time.Hour)

	if err := s.repo.IncrementClickCount(ctx, url.ID.String()); err != nil {
		ctx.Log.Warn("failed to increment click count", zap.String("short_code", shortCode))
	}

	return url, nil
}

func (s *urlServiceImpl) ListURLs(ctx *context.Context, page, limit int) (*dtos.ListResponse, error) {
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	if limit == 0 {
		offset = 0
	}

	urls, total, err := s.repo.ListURLs(ctx, limit, offset)
	if err != nil {
		ctx.Log.Error("failed to list URLs", zap.Error(err))
		return nil, err
	}

	totalPages := 1
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &dtos.ListResponse{
		Data:       urls,
		TotalCount: total,
		Pages:      totalPages,
	}, nil
}

func generateShortCode(originalURL string) string {
	hash := sha1.New()
	hash.Write([]byte(fmt.Sprintf("%s-%d", originalURL, time.Now().UnixNano())))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))[:8]
}
