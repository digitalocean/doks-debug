package connections

import (
	"context"
)

// Connector is a type that implements Connect.
type Connector interface {
	Connect(ctx context.Context) error
	Type() string
}
