package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpcPort = "50052"

type InventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (s *InventoryService) GetPart(ctx context.Context, request *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[request.GetPartUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with uuid %s not found", request.PartUuid)
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *InventoryService) ListParts(ctx context.Context, request *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filter := request.GetFilter()

	results := make([]*inventoryV1.Part, 0, len(s.parts))
	for _, p := range s.parts {
		if p == nil {
			continue
		}

		// Пропускаем деталь, если она не соответствует какому‑либо активному условию фильтра.
		if filter != nil {
			// uuids
			if len(filter.Uuids) > 0 {
				matched := false
				for _, u := range filter.Uuids {
					if p.Uuid == u {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}

			// names
			if len(filter.Names) > 0 {
				matched := false
				for _, n := range filter.Names {
					if p.Name == n {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}

			// categories
			if len(filter.Categories) > 0 {
				matched := false
				for _, c := range filter.Categories {
					if p.Category == c {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}

			// manufacturer_countries
			if len(filter.ManufacturerCountries) > 0 {
				matched := false
				for _, mc := range filter.ManufacturerCountries {
					if p.Manufacturer.Country == mc {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}

			// tags (ищем пересечение тегов)
			if len(filter.Tags) > 0 {
				matched := false
				for _, ft := range filter.Tags {
					for _, pt := range p.Tags {
						if pt == ft {
							matched = true
							break
						}
					}
					if matched {
						break
					}
				}
				if !matched {
					continue
				}
			}
		}

		results = append(results, p)
	}

	return &inventoryV1.ListPartsResponse{
		Parts: results,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()

	service := &InventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

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
