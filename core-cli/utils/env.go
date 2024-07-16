package utils

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Env struct {
	WIKI_ROOT string
	TASK_ROOT string
}

func GetEnv() Env {
	env := Env{}
	if val, ok := os.LookupEnv("WIKI_ROOT"); ok {
		env.WIKI_ROOT = val
	}
	if val, ok := os.LookupEnv("TASK_ROOT"); ok {
		env.TASK_ROOT = val
	}
	return env
}
