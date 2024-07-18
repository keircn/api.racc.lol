package utils

import (
	"math/rand"
)

func GetRandomIndex() int {
	randomIndex := rand.Intn(NUMBER_OF_IMAGES)
	if randomIndex == 0 {
		randomIndex = 1
	}

	return randomIndex
}

func GetRandomIndexVideo() int {
	randomIndex := rand.Intn(NUMBER_OF_VIDEOS)
	if randomIndex == 0 {
		randomIndex = 1
	}

	return randomIndex
}

func GetRandomIndexTrans() int {
	randomIndex := rand.Intn(NUMBER_OF_TRANS)
	if randomIndex == 0 {
		randomIndex = 1
	}

	return randomIndex
}