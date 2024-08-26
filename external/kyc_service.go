package external

import (
	"log"
	"math/rand"
	"time"
)

type ExternalKYCService struct{}

type ExternalKYCRequest struct {
	FullName string
	Email    string
	Phone    string
	Address  string
}

type ExternalKYCResponse struct {
	Status string
}

func (s *ExternalKYCService) Validate(request *ExternalKYCRequest) (*ExternalKYCResponse, error) {
	return &ExternalKYCResponse{Status: "pending"}, nil
}

func (s *ExternalKYCService) Verify(request *ExternalKYCRequest) (*ExternalKYCResponse, error) {
	// Simulate network delay
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

	// Simulate random approval or rejection
	status := "approved"
	if rand.Float32() < 0.1 { // 10% chance of rejection
		status = "rejected"
	}

	log.Println("KYC validation for:", request.FullName, "Status:", status)

	return &ExternalKYCResponse{Status: "approved"}, nil
}
