<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/Permify/permify/raw/master/assets/logo-permify-dark.png">
    <img alt="Permify logo" src="https://github.com/Permify/permify/raw/master/assets/logo-permify-light.png" width="40%">
  </picture>
<h1 align="center">
   Permify Golang Client
</h1>
</div>

<p align="center">
    <a href="https://github.com/Permify/permify" target="_blank"><img src="https://img.shields.io/github/license/Permify/permify-go?style=for-the-badge" alt="Permify Go Licence" /></a>&nbsp;
    <a href="https://discord.gg/MJbUjwskdH" target="_blank"><img src="https://img.shields.io/discord/950799928047833088?style=for-the-badge&logo=discord&label=DISCORD" alt="Permify Discord Channel" /></a>&nbsp;
</p>

# Installation

```shell
go get buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1
go get github.com/Permify/permify-go
```

# How to use

### Import Permify.

```go
import (
	permify_payload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	permify_grpc "github.com/Permify/permify-go/grpc"
)
```

### Initialize the new Permify client.

```go
// generate new client
client, err := permify_grpc.NewClient(
	permify_grpc.Config{
		Endpoint: `localhost:3478`,
	},
	grpc.WithTransportCredentials(insecure.NewCredentials()),
)
```

### Create a new tenant

```go
ct, err := client.Tenancy.Create(context.Background(), &permify_payload.TenantCreateRequest{
	Id:   "t1",
	Name: "tenant 1",
})
```

### Write Schema

```go
sr, err := client.Schema.Write(context.Background(), &permify_payload.SchemaWriteRequest {
    TenantId: "t1",
    Schema: `
        entity user {}
            
        entity document {
    
        relation viewer @user
        action view = viewer
    }`,
})
```

### Write Relationships

```go

rr, err := client.Data.WriteRelationships(context.Background(), & permify_payload.RelationshipWriteRequest {
    TenantId: "t1",
    Metadata: & permify_payload.RelationshipWriteRequestMetadata {
        SchemaVersion: sr.SchemaVersion, // sr --> schema write response
    },
    Tuples: [] * permify_payload.Tuple {
        {
            Entity: & permify_payload.Entity {
                Type: "document",
                Id: "1",
            },
            Relation: "viewer",
            Subject: & permify_payload.Subject {
                Type: "user",
                Id: "1",
            },
        }, {
            Entity: & permify_payload.Entity {
                Type: "document",
                Id: "3",
            },
            Relation: "viewer",
            Subject: & permify_payload.Subject {
                Type: "user",
                Id: "1",
            },
        },
    },
})
```

### Check

```go
cr, err := client.Permission.Check(context.Background(), & permify_payload.PermissionCheckRequest {
    TenantId: "t1",
	Metadata: & permify_payload.PermissionCheckRequestMetadata {
        SnapToken: rr.SnapToken, // rr --> relationship write response
        SchemaVersion: sr.SchemaVersion, // sr --> schema write response
        Depth: 50,
    },
    Entity: & permify_payload.Entity {
        Type: "document",
        Id: "1",
    },
    Permission: "view",
    Subject: & permify_payload.Subject {
        Type: "user",
        Id: "3",
    },
})

if (cr.Can == permify_payload.CheckResult_CHECK_RESULT_ALLOWED) {
    // RESULT_ALLOWED
} else {
    // RESULT_DENIED
}
```

### Streaming Calls

```go
str, err := client.Permission.LookupEntityStream(context.Background(), & permify_payload.PermissionLookupEntityRequest {
    TenantId: "t1",
	Metadata: & permify_payload.PermissionLookupEntityRequestMetadata {
        SnapToken: rr.SnapToken, // rr --> relationship write response
        SchemaVersion: sr.SchemaVersion, // sr --> schema write response
        Depth: 50,
    },
    EntityType: "document",
    Permission: "view",
    Subject: & permify_payload.Subject {
        Type: "user",
        Id: "1",
    },
})

// handle stream response
for {
    res, err := str.Recv()

    if err == io.EOF {
        break
    }

    // res.EntityId
}
```

Permify is an **open-source authorization service** for creating and maintaining fine-grained authorizations across your individual applications and services.

* [Permify website](https://permify.co)
* [Permify documentation](https://docs.permify.co/docs/)
* [Permify playground](https://play.permify.co)
* [Permify GitHub Repository](https://github.com/Permify/permify)

## Communication Channels

If you like Permify, please consider giving us a :star:

[![Discord](https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white)](https://discord.gg/n6KfzYxhPp)&nbsp;&nbsp;
[![X (Twitter)](https://img.shields.io/badge/X-000000?style=for-the-badge&logo=x&logoColor=white)](https://twitter.com/GetPermify)&nbsp;&nbsp;
[![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/company/permifyco)
