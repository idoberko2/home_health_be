package main

import (
	"context"

	"github.com/idoberko2/home_health_be/app"
)

func main() {
	app.New().Run(context.Background())
}
