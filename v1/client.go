package permifygotest

import (
	"google.golang.org/grpc"

	v1 `github.com/Permify/permify-go-test/generated/base/v1`
)

// Client - Permify client
type Client struct {
	Permission   v1.PermissionClient
	Schema       v1.SchemaClient
	Relationship v1.RelationshipClient
}

// Config - Permify client configuration
type Config struct {
	endpoint string
	cert     byte
}

// NewClient - Creates new Permify client
func NewClient(c Config, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.Dial(
		c.endpoint,
		opts...,
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		Permission:   v1.NewPermissionClient(conn),
		Schema:       v1.NewSchemaClient(conn),
		Relationship: v1.NewRelationshipClient(conn),
	}, nil
}
