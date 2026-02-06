package grpcclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "BecomeOverMan/proto"
	"BecomeOverMan/internal/models"
)

// RecommendationClient wraps the gRPC client for the recommendation service
type RecommendationClient struct {
	conn   *grpc.ClientConn
	client pb.RecommendationServiceClient
}

// NewRecommendationClient creates a new gRPC client connection to the recommendation service
func NewRecommendationClient(address string) (*RecommendationClient, error) {
	// Set up connection options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	// Create context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Establish connection
	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to recommendation service at %s: %w", address, err)
	}

	client := pb.NewRecommendationServiceClient(conn)
	log.Printf("✓ Connected to recommendation gRPC service at %s", address)

	return &RecommendationClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection
func (c *RecommendationClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// SearchQuests searches for quests based on query
func (c *RecommendationClient) SearchQuests(ctx context.Context, req *models.RecommendationService_SearchQuest_Request) (*models.RecommendationService_SearchQuests_Response, error) {
	// Convert Go model to proto message
	protoReq := &pb.SearchQuestsRequest{
		Query:    req.Query,
		TopK:     int32(req.TopK),
		Category: req.Category,
		Status:   req.GetStatus(),
	}

	// Call gRPC service
	resp, err := c.client.SearchQuests(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("SearchQuests gRPC call failed: %w", err)
	}

	// Convert proto response to Go model
	results := make([]models.RecommendationService_SearchQuest_Result, len(resp.Results))
	for i, r := range resp.Results {
		results[i] = models.RecommendationService_SearchQuest_Result{
			ID:              int(r.Id),
			Title:           r.Title,
			Description:     r.Description,
			Category:        r.Category,
			SimilarityScore: r.SimilarityScore,
		}
	}

	return &models.RecommendationService_SearchQuests_Response{
		Results: results,
	}, nil
}

// AddQuests adds quests to the recommendation system
func (c *RecommendationClient) AddQuests(ctx context.Context, req *models.RecommendationService_AddQuests_Request) error {
	// Convert Go model to proto message
	quests := make([]*pb.Quest, len(req.Quests))
	for i, q := range req.Quests {
		quests[i] = &pb.Quest{
			Id:          int32(q.ID),
			Title:       q.Title,
			Description: q.Description,
			Category:    q.Category,
		}
	}

	protoReq := &pb.AddQuestsRequest{
		Quests: quests,
	}

	// Call gRPC service
	_, err := c.client.AddQuests(ctx, protoReq)
	if err != nil {
		return fmt.Errorf("AddQuests gRPC call failed: %w", err)
	}

	return nil
}

// AddUsers adds users with their quest IDs to the recommendation system
func (c *RecommendationClient) AddUsers(ctx context.Context, req *models.RecommendationService_AddUsers_Request) error {
	// Convert Go model to proto message
	users := make([]*pb.UserWithQuestIDs, len(req.Users))
	for i, u := range req.Users {
		questIDs := make([]int32, len(u.QuestIDs))
		for j, id := range u.QuestIDs {
			questIDs[j] = int32(id)
		}

		users[i] = &pb.UserWithQuestIDs{
			UserId:   int32(u.UserID),
			QuestIds: questIDs,
		}
	}

	protoReq := &pb.AddUsersRequest{
		Users: users,
	}

	// Call gRPC service
	_, err := c.client.AddUsers(ctx, protoReq)
	if err != nil {
		return fmt.Errorf("AddUsers gRPC call failed: %w", err)
	}

	return nil
}

// RecommendUsers gets user recommendations
func (c *RecommendationClient) RecommendUsers(ctx context.Context, req *models.RecommendationService_RecommendUsers_Request) (*models.RecommendationService_RecommendUsers_Response, error) {
	// Convert Go model to proto message
	protoReq := &pb.RecommendUsersRequest{
		UserId: int32(req.UserID),
		TopK:   int32(req.TopK),
	}

	// Call gRPC service
	resp, err := c.client.RecommendUsers(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("RecommendUsers gRPC call failed: %w", err)
	}

	// Convert proto response to Go model
	results := make([]models.UserIDWithSimilarityScore, len(resp.Results))
	for i, r := range resp.Results {
		explanation := make(map[string]any)
		for k, v := range r.Explanation {
			explanation[k] = v
		}

		results[i] = models.UserIDWithSimilarityScore{
			UserID:          int(r.UserId),
			SimilarityScore: r.SimilarityScore,
			Explanation:     explanation,
		}
	}

	return &models.RecommendationService_RecommendUsers_Response{
		Status:  resp.Status,
		UserID:  int(resp.UserId),
		Results: results,
	}, nil
}

// RecommendQuests gets quest recommendations based on user profile
func (c *RecommendationClient) RecommendQuests(ctx context.Context, req *models.RecommendationService_RecommendQuests_Req) (*models.RecommendationService_RecommendQuests_Resp, error) {
	// Convert Go model to proto message
	questIDs := make([]int32, len(req.UserQuestIDs))
	for i, id := range req.UserQuestIDs {
		questIDs[i] = int32(id)
	}

	protoReq := &pb.RecommendQuestsRequest{
		UserQuestIds: questIDs,
		TopK:         int32(req.TopK),
		Category:     req.Category,
	}

	// Call gRPC service
	resp, err := c.client.RecommendQuests(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("RecommendQuests gRPC call failed: %w", err)
	}

	// Convert proto response to Go model
	recommendations := make([]models.RecommendQuests_Result, len(resp.Recommendations))
	for i, r := range resp.Recommendations {
		recommendations[i] = models.RecommendQuests_Result{
			RecommendationService_questToAdd: models.RecommendationService_questToAdd{
				ID:          int(r.Id),
				Title:       r.Title,
				Description: r.Description,
				Category:    r.Category,
			},
			SimilarityScore: r.SimilarityScore,
			Explanation:     r.Explanation,
		}
	}

	userProfileInfo := make(map[string]any)
	for k, v := range resp.UserProfileInfo {
		userProfileInfo[k] = v
	}

	return &models.RecommendationService_RecommendQuests_Resp{
		Recommendations: recommendations,
		UserProfileInfo: userProfileInfo,
	}, nil
}
