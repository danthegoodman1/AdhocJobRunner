package template_funcs

import (
	"os"
)

func Env(envVar string) string {
	// fmt.Println("getting", envVar)
	return os.Getenv(envVar)
}
