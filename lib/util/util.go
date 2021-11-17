package util

import (
	"errors"
	"math/rand"
)

func GetRandomInt(startInclusive, endExclusive int) (int, error) {
	if startInclusive >= endExclusive {
		return -1, errors.New("start should be < than end")
	}

	return rand.Intn(endExclusive-startInclusive) + startInclusive, nil
}
