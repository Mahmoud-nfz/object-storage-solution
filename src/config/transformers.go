package config

import (
	"log"
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
