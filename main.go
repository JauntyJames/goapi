package main

import (
	"os"

	"github.com/joho/godotenv"
)

func main() {
	a := App{}
	godotenv.Load()
	a.Initialize(
		os.Getenv("APT_TEST_USERNAME"),
		os.Getenv("ATP_TEST_PASSWORD"),
		os.Getenv("ATP_TEST_NAME"))
	a.Run(":8000")
}
