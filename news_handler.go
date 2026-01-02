package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type NewsItem struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Link       string    `json:"link"`
	ImageURL   string    `json:"image_url"`
	Content    string    `json:"contentSnippet"`
	FinalScore float64   `json:"final_score"`
	AgeClass   string    `json:"age_class"`
	CreatedAt  time.Time `json:"created_at"`
}

func parsePagination(c *gin.Context) (int, int) {
	page := 1
	limit := 20

	if value := strings.TrimSpace(c.Query("page")); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if value := strings.TrimSpace(c.Query("limit")); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			if parsed < 1 {
				parsed = 1
			}
			if parsed > 50 {
				parsed = 50
			}
			limit = parsed
		}
	}

	offset := (page - 1) * limit
	return limit, offset
}

func handleGetNews(c *gin.Context) {
	age := strings.ToUpper(strings.TrimSpace(c.Param("age")))
	switch age {
	case "ADULT", "TEENAGER", "CHILDREN":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid age"})
		return
	}

	dbURL := strings.TrimSpace(os.Getenv("SUPABASE_DB_URL"))
	if dbURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "missing database url"})
		return
	}

	cfg, err := pgx.ParseConfig(dbURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid database config"})
		return
	}

	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	ctx := context.Background()
	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection failed"})
		return
	}
	defer conn.Close(ctx)

	limit, offset := parsePagination(c)
	search := strings.TrimSpace(c.Query("q"))
	query := `
		SELECT id, title, link, image_url, final_score, age_class, created_at
		FROM news_result
		WHERE age_class = $1
	`
	args := []any{age}
	if search != "" {
		query += " AND (title ILIKE $2 OR content_snippet ILIKE $2)"
		args = append(args, "%"+search+"%")
	}
	query += " ORDER BY created_at DESC"
	query += " LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	defer rows.Close()

	items := make([]NewsItem, 0, 20)
	for rows.Next() {
		var item NewsItem
		err = rows.Scan(
			&item.ID,
			&item.Title,
			&item.Link,
			&item.ImageURL,
			&item.FinalScore,
			&item.AgeClass,
			&item.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan failed"})
			return
		}
		items = append(items, item)
	}
	if rows.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rows error"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func handleGetNewsDetail(c *gin.Context) {
	age := strings.ToUpper(strings.TrimSpace(c.Param("age")))
	switch age {
	case "ADULT", "TEENAGER", "CHILDREN":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid age"})
		return
	}

	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	dbURL := strings.TrimSpace(os.Getenv("SUPABASE_DB_URL"))
	if dbURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "missing database url"})
		return
	}

	cfg, err := pgx.ParseConfig(dbURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid database config"})
		return
	}

	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	ctx := context.Background()
	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection failed"})
		return
	}
	defer conn.Close(ctx)

	var item NewsItem
	err = conn.QueryRow(ctx, `
		SELECT id, title, link, image_url, content_snippet, final_score, age_class, created_at
		FROM news_result
		WHERE age_class = $1 AND id = $2
		LIMIT 1
	`, age, id).Scan(
		&item.ID,
		&item.Title,
		&item.Link,
		&item.ImageURL,
		&item.Content,
		&item.FinalScore,
		&item.AgeClass,
		&item.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}

	c.JSON(http.StatusOK, item)
}
