package e2e

import (
	"context"

	inventory_v1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ = Describe("Inventory", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventory_v1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)
		Expect(env).ToNot(BeNil(), "ожидали поднятое тестовое окружение")
		Expect(env.App).ToNot(BeNil(), "ожидали поднятый контейнер приложения")

		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventory_v1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции")

		cancel()
	})

	Describe("Get", func() {
		var partId string

		BeforeEach(func() {
			var err error
			partId, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали")
		})

		It("должно успешно вернуть деталь по UUID", func() {
			resp, err := inventoryClient.GetPart(ctx,
				&inventory_v1.GetPartRequest{
					PartUuid: partId,
				})

			Expect(err).ToNot(HaveOccurred())
			part := resp.GetPart()
			Expect(part).ToNot(BeNil())
			Expect(part.GetUuid()).To(Equal(partId))
			Expect(part.GetName()).ToNot(BeEmpty())
			Expect(part.GetDescription()).ToNot(BeEmpty())
			Expect(part.GetCategory()).ToNot(BeNil())
			Expect(part.GetCategory().GetEngine()).ToNot(BeEmpty())
			Expect(part.GetPrice()).To(BeNumerically(">", 0))
			Expect(part.GetManufacturer()).ToNot(BeNil())
			Expect(part.GetManufacturer().GetName()).ToNot(BeEmpty())
			Expect(part.GetUpdatedAt()).ToNot(BeNil())
		})
	})

	Describe("Get List", func() {
		BeforeEach(func() {
			var err error
			_, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали")

			_, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали")
		})

		It("должно успешно вернуть деталь по UUID", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventory_v1.ListPartsRequest{})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(len(resp.Parts)).To(BeNumerically(">", 0))
		})

	})

})
