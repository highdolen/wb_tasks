package repository

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"shortener/internal/config"

	"github.com/jackc/pgconn"
	wbfdb "github.com/wb-go/wbf/dbpg"
	wbfredis "github.com/wb-go/wbf/redis"
)

var (
	ErrCodeAlreadyExists = errors.New("short link code already exists")
	ErrNotFound          = errors.New("not found")
)

type Repository struct {
	db    *wbfdb.DB
	redis *wbfredis.Client
	ttl   time.Duration
}

type VisitStats struct {
	Date      time.Time `json:"date"`
	UserAgent string    `json:"user_agent"`
	Count     int       `json:"count"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func New(
	cfg *config.AppConfig,
	dbConn *wbfdb.DB,
	redisClient *wbfredis.Client,
) *Repository {
	return &Repository{
		db:    dbConn,
		redis: redisClient,
		ttl:   cfg.Redis.TTL,
	}
}

// CREATE SHORT LINK
func (r *Repository) CreateShortLink(
	ctx context.Context,
	originalURL string,
	customCode *string,
) (string, error) {

	code := ""
	if customCode != nil && *customCode != "" {
		code = *customCode
	} else {
		code = generateRandomCode(6)
	}

	query := `
		INSERT INTO short_links (short_code, original_url)
		VALUES ($1, $2)
	`
	_, err := r.db.ExecContext(ctx, query, code, originalURL)
	if err != nil {
		if isUniqueViolation(err) {
			return "", ErrCodeAlreadyExists
		}
		return "", err
	}

	if r.redis != nil {
		_ = r.redis.SetWithExpiration(ctx, code, originalURL, r.ttl)
	}

	return code, nil
}

// GET ORIGINAL URL
func (r *Repository) GetOriginalURL(
	ctx context.Context,
	code string,
) (string, error) {

	if r.redis != nil {
		if val, err := r.redis.Get(ctx, code); err == nil {
			return val, nil
		}
	}

	query := `
		SELECT original_url
		FROM short_links
		WHERE short_code = $1 AND is_active = true
	`

	var original string
	err := r.db.QueryRowContext(ctx, query, code).Scan(&original)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}

	if r.redis != nil {
		_ = r.redis.SetWithExpiration(ctx, code, original, r.ttl)
	}

	return original, nil
}

// SAVE VISIT (БЕЗ IP)
func (r *Repository) SaveVisit(
	ctx context.Context,
	code string,
	userAgent string,
) error {

	query := `
		INSERT INTO visits (short_url_id, user_agent, created_at)
		VALUES (
			(SELECT id FROM short_links WHERE short_code = $1),
			$2,
			NOW()
		)
	`
	_, err := r.db.ExecContext(ctx, query, code, userAgent)
	return err
}

// ANALYTICS
func (r *Repository) GetAnalytics(
	ctx context.Context,
	code string,
) ([]VisitStats, error) {

	query := `
		SELECT
			DATE(v.created_at) AS date,
			v.user_agent,
			COUNT(*)
		FROM visits v
		JOIN short_links s ON s.id = v.short_url_id
		WHERE s.short_code = $1
		GROUP BY DATE(v.created_at), v.user_agent
		ORDER BY DATE(v.created_at) DESC
	`

	rows, err := r.db.QueryContext(ctx, query, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []VisitStats
	for rows.Next() {
		var s VisitStats
		if err := rows.Scan(&s.Date, &s.UserAgent, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}

func generateRandomCode(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
