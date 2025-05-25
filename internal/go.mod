module github.com/raphaeldiscky/go-food-micro/internal

go 1.21

require (
	github.com/raphaeldiscky/go-food-micro/internal/pkg v0.0.0-00010101000000-000000000000
	github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice v0.0.0-00010101000000-000000000000
)

replace (
	github.com/raphaeldiscky/go-food-micro/internal/pkg => ./pkg
	github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice => ./services/catalogreadservice
) 