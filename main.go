package main

import (
	"tg-mc/services"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	services.Run()
}
