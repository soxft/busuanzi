package dbutil

import (
	"context"
)

// DBInterface defines the database operations needed for busuanzi
type DBInterface interface {
	// Increment site page views and return new count
	IncrSitePv(ctx context.Context, siteKey string) (int64, error)
	
	// Increment page views for a specific path and return new count
	IncrPagePv(ctx context.Context, siteKey, pathKey string) (float64, error)
	
	// Add user to site unique visitors
	AddSiteUv(ctx context.Context, siteKey, userIdentity string) error
	
	// Add user to page unique visitors
	AddPageUv(ctx context.Context, siteKey, pathKey, userIdentity string) error
	
	// Get site page views
	GetSitePv(ctx context.Context, siteKey string) (int64, error)
	
	// Get page views for a specific path
	GetPagePv(ctx context.Context, siteKey, pathKey string) (float64, error)
	
	// Count site unique visitors
	CountSiteUv(ctx context.Context, siteKey string) (int64, error)
	
	// Count page unique visitors
	CountPageUv(ctx context.Context, siteKey, pathKey string) (int64, error)
	
	// Set expire for keys (only applicable for Redis)
	SetExpire(keys ...string)
}

var DB DBInterface