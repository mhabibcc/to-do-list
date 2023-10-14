package env

import "os"

const (
	DevelopmentEnv = "development"
	StagingEnv     = "staging"
	ProductionEnv  = "production"
)

func ServiceEnv() string {
	e := os.Getenv("ENV")
	if e == "" {
		e = DevelopmentEnv
	}
	return e
}
