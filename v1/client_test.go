package permifygotest

import (
	"context"
	"io"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "github.com/Permify/permify-go/generated/base/v1"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "client-suite")
}

var _ = Describe("Client Test", func() {
	var client *Client

	BeforeEach(func() {
		var err error
		client, err = NewClient(
			Config{
				Endpoint: `localhost:3478`,
			},
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("Check request", func() {
		It("document schema", func() {
			wr, err := client.Schema.Write(context.Background(), &v1.SchemaWriteRequest{
				TenantId: "t1",
				Schema: `
            entity user {}
            
            entity document {
               relation viewer @user
               
               action view = viewer
            }`,
			})

			Expect(err).ShouldNot(HaveOccurred())

			cr, err := client.Permission.Check(context.Background(), &v1.PermissionCheckRequest{
				TenantId: "t1",
				Metadata: &v1.PermissionCheckRequestMetadata{
					SnapToken:     "",
					SchemaVersion: wr.SchemaVersion,
					Depth:         50,
				},
				Entity: &v1.Entity{
					Type: "document",
					Id:   "1",
				},
				Permission: "view",
				Subject: &v1.Subject{
					Type: "user",
					Id:   "3",
				},
			})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(cr.Can).Should(Equal(v1.CheckResult_RESULT_DENIED))
		})

		It("Lookup entity request", func() {
			wr, err := client.Schema.Write(context.Background(), &v1.SchemaWriteRequest{
				TenantId: "t1",
				Schema: `
				entity user {}
            
				entity document {
				   relation viewer @user
				   
				   action view = viewer
				}`,
			})

			rr, err := client.Relationship.Write(context.Background(), &v1.RelationshipWriteRequest{
				TenantId: "t1",
				Metadata: &v1.RelationshipWriteRequestMetadata{
					SchemaVersion: wr.SchemaVersion,
				},
				Tuples: []*v1.Tuple{
					{
						Entity: &v1.Entity{
							Type: "document",
							Id:   "1",
						},
						Relation: "viewer",
						Subject: &v1.Subject{
							Type: "user",
							Id:   "1",
						},
					},
					{
						Entity: &v1.Entity{
							Type: "document",
							Id:   "3",
						},
						Relation: "viewer",
						Subject: &v1.Subject{
							Type: "user",
							Id:   "1",
						},
					},
					{
						Entity: &v1.Entity{
							Type: "document",
							Id:   "4",
						},
						Relation: "viewer",
						Subject: &v1.Subject{
							Type: "user",
							Id:   "1",
						},
					},
				},
			})

			Expect(err).ShouldNot(HaveOccurred())

			cr, err := client.Permission.LookupEntityStream(context.Background(), &v1.PermissionLookupEntityRequest{
				TenantId: "t1",
				Metadata: &v1.PermissionLookupEntityRequestMetadata{
					SnapToken:     rr.SnapToken,
					SchemaVersion: wr.SchemaVersion,
					Depth:         50,
				},
				EntityType: "document",
				Permission: "view",
				Subject: &v1.Subject{
					Type: "user",
					Id:   "1",
				},
			})

			// handle(cr, ["1", "3", "4"])

			expected := map[string]struct{}{"1": {}, "3": {}, "4": {}}

			Expect(err).ShouldNot(HaveOccurred())

			for {
				res, err := cr.Recv()

				if err == io.EOF {
					break
				}

				Expect(err).ShouldNot(HaveOccurred())

				_, ok := expected[res.EntityId]
				Expect(ok).To(BeTrue())
			}
		})
	})
})
