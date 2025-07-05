package sm2

import (
	"fmt"
	"math"
	"time"
)


type ReviewResult struct {
	NextReviewTime time.Time
	Interval       int // in minutes
	Easiness       float64
	Repetitions    int
}

func SM2(
	now time.Time,
	previousInterval int, // in minutes
	previousEasiness float64,
	repetitions int,
	grade int, // 0–5 scale
) ReviewResult {
	var interval int
	easiness := previousEasiness
	if grade < 3 {
		repetitions = 0
		interval = 1 // reset to 5 minutes
	} else {
		switch repetitions {
		case 0:
			interval = 5
		case 1:
			interval = 30
		default:
			minutes := float64(previousInterval) * easiness
			interval = int(math.Round(minutes))
		}
		repetitions++
	}

	easiness += 0.1 - float64(5-grade)*(0.08+float64(5-grade)*0.02)
	if easiness < 1.3 {
		easiness = 1.3
	}

	fmt.Println(time.Duration(interval))

	nextReviewTime := now.Add(time.Duration(interval) * time.Minute)

	return ReviewResult{
		NextReviewTime: nextReviewTime,
		Interval:       interval,
		Easiness:       easiness,
		Repetitions:    repetitions,
	}
}