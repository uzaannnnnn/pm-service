package main

var AdultProfile = map[string]int{
	"L1": 5, "L2": 5, "L3": 5, "L4": 4, "L5": 4,
	"C1": 5, "C2": 5, "C3": 5, "C4": 5, "C5": 4,
	"P1": 5, "P2": 5, "P3": 4, "P4": 5, "P5": 5,
	"N1": 5, "N2": 5, "N3": 4, "N4": 4, "N5": 4,
}

var CF = []string{
	"L1", "L2", "L3",
	"C1", "C2", "C3",
	"P1", "P2", "P4",
	"N1", "N3",
}

var SF = []string{
	"L4", "L5",
	"C4", "C5",
	"P3", "P5",
	"N2", "N4", "N5",
}

var ParameterWeight = map[string]float64{
	"L": 0.40,
	"C": 0.30,
	"P": 0.20,
	"N": 0.10,
}

func gapWeight(gap int) float64 {
	switch gap {
	case 0:
		return 5
	case 1:
		return 4.5
	case -1:
		return 4
	case 2:
		return 3.5
	case -2:
		return 3
	case 3:
		return 2.5
	case -3:
		return 2
	case 4:
		return 1.5
	case -4:
		return 1
	default:
		return 0
	}
}

func contains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func getParameter(k string) string {
	return string(k[0])
}

func calculatePM(scores map[string]int) float64 {
	type agg struct {
		cfSum   float64
		cfCount float64
		sfSum   float64
		sfCount float64
	}

	paramAgg := map[string]*agg{
		"L": {}, "C": {}, "P": {}, "N": {},
	}

	for k, v := range scores {
		param := getParameter(k)
		gap := v - AdultProfile[k]
		weight := gapWeight(gap)

		if contains(CF, k) {
			paramAgg[param].cfSum += weight
			paramAgg[param].cfCount++
		} else {
			paramAgg[param].sfSum += weight
			paramAgg[param].sfCount++
		}
	}

	finalScore := 0.0

	for p, a := range paramAgg {
		var ncf, nsf float64

		if a.cfCount > 0 {
			ncf = a.cfSum / a.cfCount
		}
		if a.sfCount > 0 {
			nsf = a.sfSum / a.sfCount
		}

		totalParam := (0.6 * ncf) + (0.4 * nsf)
		finalScore += ParameterWeight[p] * totalParam
	}

	return finalScore
}
