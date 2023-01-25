<h1 align="center">
    <img src="https://raw.githubusercontent.com/Permify/permify/master/assets/permify-logo.svg" alt="Permify logo" width="336px" /><br />
    Permify Golang Client
</h1>

<p align="center">
    <a href="https://github.com/Permify/permify" target="_blank"><img src="https://img.shields.io/github/package-json/v/permify/permify-node?style=for-the-badge" alt="GitHub package.json version" /></a>&nbsp;
    <a href="https://github.com/Permify/permify" target="_blank"><img src="https://img.shields.io/github/license/Permify/permify?style=for-the-badge" alt="Permify Licence" /></a>&nbsp;
    <a href="https://discord.gg/MJbUjwskdH" target="_blank"><img src="https://img.shields.io/discord/950799928047833088?style=for-the-badge&logo=discord&label=DISCORD" alt="Permify Discord Channel" /></a>&nbsp;
</p>

# Installation

```shell
go get github.com/Permify/permify-go
```

# How to use

### Import Permify.

```go
import permify `github.com/Permify/permify-go`
```

### Initialize the new Permify client.

```go
import permify `github.com/Permify/permify-go`

// generate new client
client, err = permify.NewClient(
    Config{
	    endpoint: `localhost:3478`,
    },
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
```

### Create a new tenant

```go
ct, err := client.Tenancy.Create(context.Background(), &v1.TenantCreateRequest{
    Id:   "t1",
    Name: "tenant 1",
})
```

### Write Schema

```go
sr, err: = client.Schema.Write(context.Background(), &v1.SchemaWriteRequest {
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

rr, err: = client.Relationship.Write(context.Background(), & v1.RelationshipWriteRequest {
    TenantId: "t1",
    Metadata: & v1.RelationshipWriteRequestMetadata {
        SchemaVersion: sr.SchemaVersion, // sr --> schema write response
    },
    Tuples: [] * v1.Tuple {
        {
            Entity: & v1.Entity {
                Type: "document",
                Id: "1",
            },
            Relation: "viewer",
            Subject: & v1.Subject {
                Type: "user",
                Id: "1",
            },
        }, {
            Entity: & v1.Entity {
                Type: "document",
                Id: "3",
            },
            Relation: "viewer",
            Subject: & v1.Subject {
                Type: "user",
                Id: "1",
            },
        }
    },
})
```

### Check

```go
cr, err: = client.Permission.Check(context.Background(), & v1.PermissionCheckRequest {
    TenantId: "t1",
	Metadata: & v1.PermissionCheckRequestMetadata {
        SnapToken: rr.SnapToken, // rr --> relationship write response
        SchemaVersion: sr.SchemaVersion, // sr --> schema write response
        Depth: 50,
    },
    Entity: & v1.Entity {
        Type: "document",
        Id: "1",
    },
    Permission: "view",
    Subject: & v1.Subject {
        Type: "user",
        Id: "3",
    },

    if (cr.can === PermissionCheckResponse_Result.RESULT_ALLOWED) {
        // RESULT_ALLOWED
    } else {
        // RESULT_DENIED
    }
})
```

### Streaming Calls

```go
str, err: = client.Permission.LookupEntityStream(context.Background(), & v1.PermissionLookupEntityRequest {
    TenantId: "t1",
	Metadata: & v1.PermissionLookupEntityRequestMetadata {
        SnapToken: rr.SnapToken, // rr --> relationship write response
        SchemaVersion: sr.SchemaVersion, // sr --> schema write response
        Depth: 50,
    },
    EntityType: "document",
    Permission: "view",
    Subject: & v1.Subject {
        Type: "user",
        Id: "1",
    },
})

// handle stream response
for {
    res, err: = str.Recv()

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

## Community & Support

Join our [Discord channel](https://discord.gg/MJbUjwskdH) for issues, feature requests, feedbacks or anything else. We love to talk about authorization and access control :heart:

<p align="left">
<a href="https://discord.gg/MJbUjwskdH">
 <img height="70px" width="70px" alt="permify | Discord" src="https://user-images.githubusercontent.com/39353278/187209316-3d01a799-c51b-4eaa-8f52-168047078a14.png" />
</a>
<a href="https://twitter.com/GetPermify">
  <img height="70px" width="70px" alt="permify | Twitter" src="https://user-images.githubusercontent.com/39353278/187209323-23f14261-d406-420d-80eb-1aa707a71043.png"/>
</a>
<a href="https://www.linkedin.com/company/permifyco">
  <img height="70px" width="70px" alt="permify | Linkedin" src="https://user-images.githubusercontent.com/39353278/187209321-03293a24-6f63-4321-b362-b0fc89fdd879.png" />
</a>
</p>
