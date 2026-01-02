package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

func normalizeAgeClass(age string) string {
	switch strings.ToUpper(strings.TrimSpace(age)) {
	case "CHILD", "CHILDREN", "KIDS":
		return "CHILDREN"
	case "TEEN", "TEENAGER", "REMAJA":
		return "TEENAGER"
	case "ADULT", "DEWASA":
		return "ADULT"
	default:
		return "TEENAGER"
	}
}

func saveToDB(req Request, age string, score float64) {
	dbURL := strings.TrimSpace(os.Getenv("SUPABASE_DB_URL"))

	cfg, err := pgx.ParseConfig(dbURL)
	if err != nil {
		log.Println("DB parse config error:", err)
		return
	}

	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	conn, err := pgx.ConnectConfig(context.Background(), cfg)
	if err != nil {
		log.Println("DB connection error:", err)
		return
	}
	defer conn.Close(context.Background())

	rawJSON, err := json.Marshal(req.Scores)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}

	normalizedAge := normalizeAgeClass(age)

	_, err = conn.Exec(context.Background(), `
		INSERT INTO news_result
		(title, link, image_url, content_snippet, age_class, final_score, raw_scores)
		VALUES ($1,$2,$3,$4,$5,$6,$7::jsonb)
	`,
		req.Title,
		req.Link,
		req.Image,
		req.ContentSnippet,
		normalizedAge,
		score,
		string(rawJSON),
	)

	if err != nil {
		log.Println("DB insert error:", err)
	}
}
