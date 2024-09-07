package utils

import (
	crand "crypto/rand"
	"errors"
	"io"
	"math/big"
	"math/rand"
)

func RandFloat64(max float64) float64 {
	return rand.Float64() * max
}

func RangeFloat64(min, max float64) float64 {
	if min > max {
		return RangeFloat64(max, min)
	}

	return min + RandFloat64(max)
}

func RandFloat64Positive(max float64) float64 {
	lhx := RandFloat64(max)
	rhx := RandFloat64(max)
	return lhx - rhx
}

func RandPosition(maxX, maxY float64) (float64, float64) {
	x := rand.Float64() * maxX
	y := rand.Float64() * maxY
	return x, y
}

func RandInt(min, max int, seed int64) int {
	if min == max {
		return min
	}

	if min > max {
		max, min = min, max
	}

	return min + rand.Intn(max-min)
}

func GenerateToken(nBytes int, encodeFunc func([]byte) string) (string, error) {
	if nBytes <= 0 {
		return "", errors.New("ensure byte count > 0")
	}

	token := make([]byte, nBytes)
	_, err := io.ReadFull(crand.Reader, token)
	if err != nil {
		return "", err
	}

	var tokenStr string
	if encodeFunc != nil {
		tokenStr = encodeFunc(token)
	} else {
		tokenStr = string(token)
	}

	return tokenStr, nil
}

func RandIntByCrypto(min, max int) int {
	if min == max {
		return 0
	}

	if min > max {
		max, min = min, max
	}
	result, _ := crand.Int(crand.Reader, big.NewInt(int64(max-min)))
	return min + int(result.Int64())
}

// 返回 不重复的数组
func RandIntsNoRepeat(origin []int, count int) []int {
	tmpOrigin := make([]int, len(origin))
	copy(tmpOrigin, origin)

	rand.Shuffle(len(tmpOrigin), func(i int, j int) {
		tmpOrigin[i], tmpOrigin[j] = tmpOrigin[j], tmpOrigin[i]
	})

	result := make([]int, 0, count)
	for index, value := range tmpOrigin {
		if index == count {
			break
		}
		result = append(result, value)
	}

	return result
}
