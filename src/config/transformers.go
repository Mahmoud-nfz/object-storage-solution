package config

import (
	"log"
	"strconv"
	"time"

	"data-storage/src/utils"
)

type transformer[T any] func(key string, value string) T

func maximumAgeTransformer(key string, value string) time.Duration {
	maximumAge, err := utils.ParseDuration(value)
	if err != nil {
		log.Fatalln("Error parsing ", key, ": ", err)
	}

	return maximumAge
}

func jwtSecretTransformer(key string, value string) []byte {
	return []byte(value)
}

func chunkSizeTransformer(key string, value string) int64 {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatalln("Error parsing ", key, ": ", err)
	}
	return i
}
