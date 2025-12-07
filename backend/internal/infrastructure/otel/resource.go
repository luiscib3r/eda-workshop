package otel

import (
	"go.opentelemetry.io/otel/sdk/resource"
)

func Resource() (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.Environment(),
	)
}
