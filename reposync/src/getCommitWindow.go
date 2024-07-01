package reposync

import "math"

func GetCommitWindow(lenCommitList int) int {
	return int(math.Max(1, math.Floor(float64(lenCommitList)*0.05)))
}
