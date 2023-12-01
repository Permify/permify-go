package permifygotest

import (
	"google.golang.org/grpc"

	v1 "github.com/Permify/permify-go/generated/base/v1"
)

// Client - Permify client
type Client struct {
	Permission v1.PermissionClient
	Schema     v1.SchemaClient
	Data       v1.DataClient
	Bundle     v1.BundleClient
	Tenancy    v1.TenancyClient
	Watch      v1.WatchClient
}

// Config - Permify client configuration
type Config struct {
	Endpoint string
	Cert     byte
}

// NewClient - Creates new Permify client
func NewClient(c Config, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.Dial(
		c.Endpoint,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		Permission: v1.NewPermissionClient(conn),
		Schema:     v1.NewSchemaClient(conn),
		Data:       v1.NewDataClient(conn),
		Bundle:     v1.NewBundleClient(conn),
		Tenancy:    v1.NewTenancyClient(conn),
		Watch:      v1.NewWatchClient(conn),
	}, nil
}
