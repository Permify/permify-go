package permifygrpc

import (
	"context"
	"io"
	"testing"

	"google.golang.org/protobuf/types/known/anypb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	pservice "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
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
			wr, err := client.Schema.Write(context.Background(), &pservice.SchemaWriteRequest{
				TenantId: "t1",
				Schema: `
            entity user {}
            
            entity document {
               relation viewer @user
               
               action view = viewer
            }`,
			})

			Expect(err).ShouldNot(HaveOccurred())

			cr, err := client.Permission.Check(context.Background(), &pservice.PermissionCheckRequest{
				TenantId: "t1",
				Metadata: &pservice.PermissionCheckRequestMetadata{
					SnapToken:     "",
					SchemaVersion: wr.SchemaVersion,
					Depth:         50,
				},
				Entity: &pservice.Entity{
					Type: "document",
					Id:   "1",
				},
				Permission: "view",
				Subject: &pservice.Subject{
					Type: "user",
					Id:   "3",
				},
			})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(cr.Can).Should(Equal(pservice.CheckResult_CHECK_RESULT_DENIED))
		})
	})

	Context("Write Data", func() {
		It("abac schema", func() {
			wr, err := client.Schema.Write(context.Background(), &pservice.SchemaWriteRequest{
				TenantId: "instagram",
				Schema: `
					entity user {}
					
					entity account {
						relation owner @user
						relation following @user
						relation follower @user
					
						attribute public boolean
					
						action view = (owner or follower) or public
						
					}`,
			})

			Expect(err).ShouldNot(HaveOccurred())

			// Convert the wrapped attribute value into Any proto message
			value, err := anypb.New(&pservice.BooleanValue{
				Data: true,
			})
			if err != nil {
				Expect(err).ShouldNot(HaveOccurred())
			}

			cr, err := client.Data.Write(context.Background(), &pservice.DataWriteRequest{
				TenantId: "instagram",
				Metadata: &pservice.DataWriteRequestMetadata{
					SchemaVersion: wr.SchemaVersion,
				},
				Attributes: []*pservice.Attribute{
					{
						Entity: &pservice.Entity{
							Type: "account",
							Id:   "1",
						},
						Attribute: "public",
						Value:     value,
					},
				},
			})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(cr.SnapToken).ShouldNot(BeNil())
		})
	})

	Context("Lookup", func() {
		It("Lookup entity request", func() {
			wr, err := client.Schema.Write(context.Background(), &pservice.SchemaWriteRequest{
				TenantId: "t1",
				Schema: `
				entity user {}
            
				entity document {
				   relation viewer @user
				   
				   action view = viewer
				}`,
			})

			rr, err := client.Data.WriteRelationships(context.Background(), &pservice.RelationshipWriteRequest{
				TenantId: "t1",
				Metadata: &pservice.RelationshipWriteRequestMetadata{
					SchemaVersion: wr.SchemaVersion,
				},
				Tuples: []*pservice.Tuple{
					{
						Entity: &pservice.Entity{
							Type: "document",
							Id:   "1",
						},
						Relation: "viewer",
						Subject: &pservice.Subject{
							Type: "user",
							Id:   "1",
						},
					},
					{
						Entity: &pservice.Entity{
							Type: "document",
							Id:   "3",
						},
						Relation: "viewer",
						Subject: &pservice.Subject{
							Type: "user",
							Id:   "1",
						},
					},
					{
						Entity: &pservice.Entity{
							Type: "document",
							Id:   "4",
						},
						Relation: "viewer",
						Subject: &pservice.Subject{
							Type: "user",
							Id:   "1",
						},
					},
				},
			})

			Expect(err).ShouldNot(HaveOccurred())

			cr, err := client.Permission.LookupEntityStream(context.Background(), &pservice.PermissionLookupEntityRequest{
				TenantId: "t1",
				Metadata: &pservice.PermissionLookupEntityRequestMetadata{
					SnapToken:     rr.SnapToken,
					SchemaVersion: wr.SchemaVersion,
					Depth:         50,
				},
				EntityType: "document",
				Permission: "view",
				Subject: &pservice.Subject{
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
