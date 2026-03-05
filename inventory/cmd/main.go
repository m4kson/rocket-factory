package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	inventoryAPI "github.com/m4kson/rocket-factory/inventory/internal/api/inventory/v1"
	inventoryRepository "github.com/m4kson/rocket-factory/inventory/internal/repository/part"
	inventoryService "github.com/m4kson/rocket-factory/inventory/internal/service/part"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = "50052"

//type InventoryService struct {
//	inventoryV1.UnimplementedInventoryServiceServer
//
//	mu    sync.RWMutex
//	parts map[string]*inventoryV1.Part
//}
//
//func (s *InventoryService) GetPart(ctx context.Context, request *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	part, ok := s.parts[request.GetPartUuid()]
//	if !ok {
//		return nil, status.Errorf(codes.NotFound, "part with uuid %s not found", request.PartUuid)
//	}
//
//	return &inventoryV1.GetPartResponse{
//		Part: part,
//	}, nil
//}
//
//func (s *InventoryService) ListParts(ctx context.Context, request *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	filter := request.GetFilter()
//
//	results := make([]*inventoryV1.Part, 0, len(s.parts))
//	for _, p := range s.parts {
//		if p == nil {
//			continue
//		}
//
//		// Пропускаем деталь, если она не соответствует какому‑либо активному условию фильтра.
//		if filter != nil {
//			// uuids
//			if len(filter.Uuids) > 0 {
//				matched := false
//				for _, u := range filter.Uuids {
//					if p.Uuid == u {
//						matched = true
//						break
//					}
//				}
//				if !matched {
//					continue
//				}
//			}
//
//			// names
//			if len(filter.Names) > 0 {
//				matched := false
//				for _, n := range filter.Names {
//					if p.Name == n {
//						matched = true
//						break
//					}
//				}
//				if !matched {
//					continue
//				}
//			}
//
//			// categories
//			if len(filter.Categories) > 0 {
//				matched := false
//				for _, c := range filter.Categories {
//					if p.Category == c {
//						matched = true
//						break
//					}
//				}
//				if !matched {
//					continue
//				}
//			}
//
//			// manufacturer_countries
//			if len(filter.ManufacturerCountries) > 0 {
//				matched := false
//				for _, mc := range filter.ManufacturerCountries {
//					if p.Manufacturer.Country == mc {
//						matched = true
//						break
//					}
//				}
//				if !matched {
//					continue
//				}
//			}
//
//			// tags (ищем пересечение тегов)
//			if len(filter.Tags) > 0 {
//				matched := false
//				for _, ft := range filter.Tags {
//					for _, pt := range p.Tags {
//						if pt == ft {
//							matched = true
//							break
//						}
//					}
//					if matched {
//						break
//					}
//				}
//				if !matched {
//					continue
//				}
//			}
//		}
//
//		results = append(results, p)
//	}
//
//	return &inventoryV1.ListPartsResponse{
//		Parts: results,
//	}, nil
//}

//// initializeTestData инициализирует сервис тестовыми данными о деталях
//func initializeTestData(service *InventoryService) {
//	now := timestamppb.Now()
//
//	testParts := []*inventoryV1.Part{
//		{
//			Uuid:          uuid.New().String(),
//			Name:          "Турбокомпрессор",
//			Description:   "Высокопроизводительный турбокомпрессор для двигателя",
//			Price:         15999.99,
//			StockQuantity: 45,
//			Category: &inventoryV1.Category{
//				Category: &inventoryV1.Category_Engine{Engine: "engine"},
//			},
//			Dimensions: &inventoryV1.Dimensions{
//				Length: 250.0,
//				Width:  180.0,
//				Height: 200.0,
//				Weight: 12.5,
//			},
//			Manufacturer: &inventoryV1.Manufacturer{
//				Name:    "TurboTech Industries",
//				Country: "Germany",
//				Website: "https://turbotech.de",
//			},
//			Tags:      []string{"performance", "engine", "boost"},
//			CreatedAt: now,
//			UpdatedAt: now,
//		},
//		{
//			Uuid:          uuid.New().String(),
//			Name:          "Топливный насос",
//			Description:   "Электрический топливный насос высокого давления",
//			Price:         8499.50,
//			StockQuantity: 120,
//			Category: &inventoryV1.Category{
//				Category: &inventoryV1.Category_Fuel{Fuel: "fuel"},
//			},
//			Dimensions: &inventoryV1.Dimensions{
//				Length: 150.0,
//				Width:  100.0,
//				Height: 120.0,
//				Weight: 2.3,
//			},
//			Manufacturer: &inventoryV1.Manufacturer{
//				Name:    "FuelFlow Corp",
//				Country: "Japan",
//				Website: "https://fuelflow.jp",
//			},
//			Tags:      []string{"fuel", "pump", "electrical"},
//			CreatedAt: now,
//			UpdatedAt: now,
//		},
//		{
//			Uuid:          uuid.New().String(),
//			Name:          "Радиатор охлаждения",
//			Description:   "Алюминиевый радиатор охлаждения двигателя",
//			Price:         5299.00,
//			StockQuantity: 60,
//			Category: &inventoryV1.Category{
//				Category: &inventoryV1.Category_Engine{Engine: "engine"},
//			},
//			Dimensions: &inventoryV1.Dimensions{
//				Length: 600.0,
//				Width:  400.0,
//				Height: 50.0,
//				Weight: 8.2,
//			},
//			Manufacturer: &inventoryV1.Manufacturer{
//				Name:    "CoolTech Solutions",
//				Country: "USA",
//				Website: "https://cooltech.com",
//			},
//			Tags:      []string{"cooling", "radiator", "engine"},
//			CreatedAt: now,
//			UpdatedAt: now,
//		},
//		{
//			Uuid:          uuid.New().String(),
//			Name:          "Крыло левое",
//			Description:   "Алюминиевое крыло с аэродинамическим профилем",
//			Price:         12500.00,
//			StockQuantity: 25,
//			Category: &inventoryV1.Category{
//				Category: &inventoryV1.Category_Wing{Wing: "wing"},
//			},
//			Dimensions: &inventoryV1.Dimensions{
//				Length: 8500.0,
//				Width:  2500.0,
//				Height: 150.0,
//				Weight: 180.0,
//			},
//			Manufacturer: &inventoryV1.Manufacturer{
//				Name:    "AeroWings Ltd",
//				Country: "UK",
//				Website: "https://aerowings.co.uk",
//			},
//			Tags:      []string{"wing", "aerodynamic", "aluminum"},
//			CreatedAt: now,
//			UpdatedAt: now,
//		},
//		{
//			Uuid:          uuid.New().String(),
//			Name:          "Иллюминатор",
//			Description:   "Прочный многослойный иллюминатор с защитой от удара",
//			Price:         3750.25,
//			StockQuantity: 80,
//			Category: &inventoryV1.Category{
//				Category: &inventoryV1.Category_Porthole{Porthole: "porthole"},
//			},
//			Dimensions: &inventoryV1.Dimensions{
//				Length: 800.0,
//				Width:  600.0,
//				Height: 80.0,
//				Weight: 15.5,
//			},
//			Manufacturer: &inventoryV1.Manufacturer{
//				Name:    "ClearView Optics",
//				Country: "Switzerland",
//				Website: "https://clearview.ch",
//			},
//			Tags:      []string{"window", "porthole", "safety"},
//			CreatedAt: now,
//			UpdatedAt: now,
//		},
//	}
//
//	service.mu.Lock()
//	defer service.mu.Unlock()
//
//	for _, part := range testParts {
//		service.parts[part.Uuid] = part
//		log.Printf("✅ Инициализирована деталь: %s (UUID: %s)", part.Name, part.Uuid)
//	}
//
//	log.Printf("📦 Всего деталей инициализировано: %d", len(service.parts))
//}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()

	repo := inventoryRepository.NewPartRepository()
	service := inventoryService.NewPartService(repo)
	api := inventoryAPI.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("grpc server listening on %s\n", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped")
}
