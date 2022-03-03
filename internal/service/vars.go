package service

import "os"

var (
	apiUrl  = os.Getenv("API_URL")
	apiPort = os.Getenv("API_PORT")
)
