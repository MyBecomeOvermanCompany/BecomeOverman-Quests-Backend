// Package grpcclient Example Usage
//
// This file demonstrates how to use the RecommendationClient in your Go services.

package grpcclient

import (
	"context"
	"log"
	"time"

	"BecomeOverMan/internal/models"
)

// Example 1: Searching for quests
func ExampleSearchQuests(client *RecommendationClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &models.RecommendationService_SearchQuest_Request{
		Query:    "morning routine",
		TopK:     10,
		Category: "health",
		Status:   "all",
	}

	resp, err := client.SearchQuests(ctx, req)
	if err != nil {
		log.Printf("SearchQuests failed: %v", err)
		return
	}

	log.Printf("Found %d quests", len(resp.Results))
	for _, result := range resp.Results {
		log.Printf("Quest %d: %s (score: %.2f)", result.ID, result.Title, result.SimilarityScore)
	}
}

// Example 2: Adding quests to recommendation system
func ExampleAddQuests(client *RecommendationClient) {
	ctx := context.Background()

	req := &models.RecommendationService_AddQuests_Request{
		Quests: []models.RecommendationService_questToAdd{
			{
				ID:          1,
				Title:       "Champion Wake Up",
				Description: "Wake up early every day",
				Category:    "health",
			},
			{
				ID:          2,
				Title:       "Meditation Master",
				Description: "Daily meditation practice",
				Category:    "mindfulness",
			},
		},
	}

	err := client.AddQuests(ctx, req)
	if err != nil {
		log.Printf("AddQuests failed: %v", err)
		return
	}

	log.Println("Quests added successfully")
}

// Example 3: Adding users to recommendation system
func ExampleAddUsers(client *RecommendationClient) {
	ctx := context.Background()

	req := &models.RecommendationService_AddUsers_Request{
		Users: []models.UserWithQuestIDS{
			{
				UserID:   123,
				QuestIDs: []int{1, 2, 5, 10},
			},
			{
				UserID:   456,
				QuestIDs: []int{2, 3, 7},
			},
		},
	}

	err := client.AddUsers(ctx, req)
	if err != nil {
		log.Printf("AddUsers failed: %v", err)
		return
	}

	log.Println("Users added successfully")
}

// Example 4: Getting user recommendations
func ExampleRecommendUsers(client *RecommendationClient) {
	ctx := context.Background()

	req := &models.RecommendationService_RecommendUsers_Request{
		UserID: 123,
		TopK:   5,
	}

	resp, err := client.RecommendUsers(ctx, req)
	if err != nil {
		log.Printf("RecommendUsers failed: %v", err)
		return
	}

	log.Printf("Recommended users for user %d:", resp.UserID)
	for _, rec := range resp.Results {
		log.Printf("User %d (score: %.2f) - %v", rec.UserID, rec.SimilarityScore, rec.Explanation)
	}
}

// Example 5: Getting quest recommendations
func ExampleRecommendQuests(client *RecommendationClient) {
	ctx := context.Background()

	req := &models.RecommendationService_RecommendQuests_Req{
		UserQuestIDs: []int{1, 2, 5, 10},
		TopK:         5,
		Category:     "health",
	}

	resp, err := client.RecommendQuests(ctx, req)
	if err != nil {
		log.Printf("RecommendQuests failed: %v", err)
		return
	}

	log.Printf("Quest recommendations (profile: %v):", resp.UserProfileInfo)
	for _, rec := range resp.Recommendations {
		log.Printf("Quest %d: %s (score: %.2f) - %s",
			rec.ID, rec.Title, rec.SimilarityScore, rec.Explanation)
	}
}

// Example 6: Using in a service
type QuestService struct {
	grpcClient *RecommendationClient
	// ... other dependencies
}

func (s *QuestService) GetSimilarQuests(questID int) ([]models.RecommendQuests_Result, error) {
	ctx := context.Background()

	req := &models.RecommendationService_RecommendQuests_Req{
		UserQuestIDs: []int{questID},
		TopK:         10,
	}

	resp, err := s.grpcClient.RecommendQuests(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Recommendations, nil
}

// Example 7: Error handling with fallback
func ExampleWithFallback(client *RecommendationClient) []models.RecommendationService_SearchQuest_Result {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req := &models.RecommendationService_SearchQuest_Request{
		Query: "workout",
		TopK:  5,
	}

	resp, err := client.SearchQuests(ctx, req)
	if err != nil {
		log.Printf("gRPC call failed, using fallback: %v", err)
		// Return default/cached results
		return []models.RecommendationService_SearchQuest_Result{}
	}

	return resp.Results
}
