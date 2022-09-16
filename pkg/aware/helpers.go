package aware

import (
	"math"
	"math/rand"
	"time"
)

func generateRandomFloat(min float64, max float64, decimals int) float64 {
    rand.Seed(time.Now().UnixNano())
    val := min + rand.Float64() * (max - min)
    return math.Round(val * math.Pow(10, float64(decimals)))/math.Pow(10, float64(decimals))
}

func generateRandomInt(min, max int) int {
    rand.Seed(time.Now().UnixNano())
    return min + rand.Int() * (max - min)
}

func generateRandomBool() bool {
    return generateRandomInt(0, 1) == 1
}

// TODO: generatedRandomBoolWithWeighting()

func formatTime(t time.Time) string {
    return t.Format(time.RFC3339)
}
