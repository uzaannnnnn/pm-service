package main

import "github.com/gin-gonic/gin"

type Request struct {
	Title          string         `json:"title"`
	Link           string         `json:"link"`
	IsoDate        string         `json:"iso_date"`
	Image          string         `json:"image"`
	ContentSnippet string         `json:"contentSnippet"`
	Scores         map[string]int `json:"scores"`
}

func handlePM(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	age, finalScore := classify(req.Scores)

	saveToDB(req, age, finalScore)

	c.JSON(200, gin.H{
		"age_class":   age,
		"final_score": finalScore,
	})
}
