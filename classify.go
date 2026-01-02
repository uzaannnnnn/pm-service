package main

func classify(scores map[string]int) (string, float64) {
	score := calculatePM(scores)

	switch {
	case score > 3.15:
		return "ADULT", score
	case score > 1.49:
		return "TEENAGER", score
	default:
		return "CHILDREN", score
	}
}
