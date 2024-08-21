package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DatabaseURLConnection string
	ServerPort            int
	JWTSecret             string
)

func Load() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ServerPort, err = strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		ServerPort = 9000
	}

	DatabaseURLConnection = fmt.Sprintf(
		"%s:%s@/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SCHEMA"),
	)

	JWTSecret = os.Getenv("JWT_SECRET")

	log.SetFlags(log.LstdFlags | log.LUTC | log.Llongfile)
}
