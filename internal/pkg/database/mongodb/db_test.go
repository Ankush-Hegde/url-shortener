package mongodb

import (
	"context"
	"testing"
)

func TestConnectRequiresURI(t *testing.T) {
	_, err := Connect(context.Background(), " ")
	if err == nil {
		t.Fatal("Connect() error = nil, want an error for an empty URI")
	}
}
