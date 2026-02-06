package services

import (
	"BecomeOverMan/internal/models"
	"BecomeOverMan/internal/repositories"
	"context"
	"fmt"
	"log/slog"
	"slices"

	grpcclient "BecomeOverMan/internal/grpc"
)

type QuestService struct {
	questRepo  *repositories.QuestRepository
	userRepo   *repositories.UserRepository
	grpcClient *grpcclient.RecommendationClient
}

func NewQuestService(
	questRepo *repositories.QuestRepository,
	userRepo *repositories.UserRepository,
	grpcClient *grpcclient.RecommendationClient,
) *QuestService {
	return &QuestService{
		questRepo:  questRepo,
		userRepo:   userRepo,
		grpcClient: grpcClient,
	}
}

// GetAvailableQuests returns quests available for the user
func (s *QuestService) GetAvailableQuests(ctx context.Context, userID int) ([]models.Quest, error) {
	return s.questRepo.GetAvailableQuests(ctx, userID)
}

func (s *QuestService) GetQuestShop(ctx context.Context, userID int) ([]models.Quest, error) {
	return s.questRepo.GetQuestShop(ctx, userID)
}

func (s *QuestService) GetMyActiveQuests(ctx context.Context, userID int) ([]models.Quest, error) {
	return s.questRepo.GetMyActiveQuests(ctx, userID)
}

func (s *QuestService) GetMyCompletedQuests(ctx context.Context, userID int) ([]models.Quest, error) {
	return s.questRepo.GetMyCompletedQuests(ctx, userID)
}

func (s *QuestService) GetMyAllQuestsWithDetails(ctx context.Context, userID int) ([]models.Quest, error) {
	return s.questRepo.GetMyAllQuestsWithDetails(ctx, userID)
}

// PurchaseQuest handles the purchase of a quest by a user
func (s *QuestService) PurchaseQuest(ctx context.Context, userID, questID int) error {
	err := s.questRepo.PurchaseQuest(ctx, userID, questID)
	if err != nil {
		slog.Error("Failed to purchase quest", "error", err)
		return err
	}

	go func() {
		questIDS, err := s.getUserQuestIDs(userID)
		if err != nil {
			slog.Error("Failed to get user quest IDs", "error", err, "user_id", userID)
			return
		}

		if len(questIDS) == 0 {
			slog.Info("User has no quests", "user_id", userID)
		}

		req := models.RecommendationService_AddUsers_Request{
			Users: []models.UserWithQuestIDS{
				{
					UserID:   userID,
					QuestIDs: questIDS,
				},
			},
		}

		response, err := s.sendUserQuestToRecommendationService(req)
		if err != nil {
			slog.Error("Failed to send user quest to recommendation service", "error", err, "user_id", userID)
		}

		slog.Info("User quest sent to recommendation service", "user_id", userID, "response", response)
	}()

	return nil
}

func (s *QuestService) getUserQuestIDs(userID int) ([]int, error) {
	return s.questRepo.GetUserQuestIDs(userID)
}

func (s *QuestService) sendUserQuestToRecommendationService(req models.RecommendationService_AddUsers_Request) (map[string]any, error) {
	if s.grpcClient == nil {
		return nil, fmt.Errorf("gRPC client is not initialized")
	}

	ctx := context.Background()
	err := s.grpcClient.AddUsers(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error calling gRPC AddUsers: %v", err)
	}

	response := map[string]any{
		"status": "success",
	}
	return response, nil
}

// StartQuest begins the execution of a purchased quest
func (s *QuestService) StartQuest(ctx context.Context, userID, questID int) error {
	return s.questRepo.StartQuest(ctx, userID, questID)
}

// CompleteTask marks a task as completed by the user
func (s *QuestService) CompleteTask(ctx context.Context, userID, questID, taskID int) error {
	return s.questRepo.CompleteTask(ctx, userID, questID, taskID)
}

// CompleteQuest finalizes the quest completion
func (s *QuestService) CompleteQuest(ctx context.Context, userID, questID int) error {
	return s.questRepo.CompleteQuest(ctx, userID, questID)
}

func (s *QuestService) GetQuestDetails(ctx context.Context, questID int, userID int) (*models.Quest, error) {
	return s.questRepo.GetQuestDetails(ctx, questID, userID)
}

func (s *QuestService) CreateSharedQuest(user1ID, user2ID, questID int) error {
	return s.questRepo.CreateSharedQuest(user1ID, user2ID, questID)
}

func (s *QuestService) SaveQuestToDB(quest *models.Quest, tasks []models.Task) (int, error) {
	return s.questRepo.SaveQuestToDB(quest, tasks)
}

func (s *QuestService) SearchQuests(
	ctx context.Context,
	req models.RecommendationService_SearchQuest_Request,
	userID int,
) (models.SearchQuestsResponse, error) {
	if s.grpcClient == nil {
		return nil, fmt.Errorf("gRPC client is not initialized")
	}

	// Call gRPC service
	response, err := s.grpcClient.SearchQuests(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error calling gRPC SearchQuests: %v", err)
	}

	// Extract quest IDs
	questsIDS := make([]int, len(response.Results))
	for i, result := range response.Results {
		questsIDS[i] = result.ID
	}

	var questsWithDetails []models.Quest

	// Get quests with details from DB
	questsWithDetails, err = s.questRepo.SearchQuestsWithDetailsByIDs(ctx, questsIDS)
	if err != nil {
		slog.ErrorContext(ctx, "ошибка получения квестов из БД с указанными ids во время поиска",
			"error", err,
			"ids", questsIDS,
		)
		return nil, fmt.Errorf("В поиске квестов по запросу произошла внутренняя ошибка: %w", err)
	}

	// Return result
	questsWithDetailsAndSimilarityResponse := models.NewSearchQuestsResponse(questsWithDetails, *response)

	return questsWithDetailsAndSimilarityResponse, nil
}

func (s *QuestService) RecommendFriends(
	ctx context.Context,
	req models.RecommendationService_RecommendUsers_Request,
) ([]models.UserProfileWithSimilarityScore, error) {
	if s.grpcClient == nil {
		return nil, fmt.Errorf("gRPC client is not initialized")
	}

	// Call gRPC service
	response, err := s.grpcClient.RecommendUsers(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error calling gRPC RecommendUsers: %v", err)
	}

	// Extract user IDs
	userIDs := make([]int, len(response.Results))
	userIDsWithSimilarityScore := make(map[int]float64, len(response.Results))
	explanations := make(map[int]map[string]any, len(response.Results))
	for i, result := range response.Results {
		userIDs[i] = result.UserID
		userIDsWithSimilarityScore[result.UserID] = result.SimilarityScore
		explanations[result.UserID] = result.Explanation
	}

	// Get recommended profiles (excluding existing friends)
	recommendedProfiles, err := s.userRepo.GetProfiles(userIDs)
	if err != nil {
		slog.ErrorContext(ctx, "ошибка получения профилей из БД с указанными ids во время рекомендации друзей",
			"error", err,
			"ids", userIDs,
		)
		return nil, fmt.Errorf("В рекомендации друзей по запросу произошла внутренняя ошибка: %w", err)
	}

	// Exclude users who are already friends
	friendsIDS, err := s.userRepo.GetAllAcceptedFriends(req.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "ошибка получения друзей из БД с указанными ids во время рекомендации друзей",
			"error", err,
			"user_id", req.UserID,
		)
		return nil, fmt.Errorf("В рекомендации друзей по запросу произошла внутренняя ошибка: %w", err)
	}

	recommendedNotFriendsProfiles := make([]models.UserProfile, 0, len(recommendedProfiles))
	for _, profile := range recommendedProfiles {
		if !slices.Contains(friendsIDS, profile.ID) {
			recommendedNotFriendsProfiles = append(recommendedNotFriendsProfiles, profile)
		}
	}

	// Return result
	recommendedProfilesAndSimilarityResponse := make([]models.UserProfileWithSimilarityScore, 0, len(recommendedNotFriendsProfiles))
	for _, profile := range recommendedNotFriendsProfiles {
		recommendedProfilesAndSimilarityResponse = append(recommendedProfilesAndSimilarityResponse, models.UserProfileWithSimilarityScore{
			UserProfile:     profile,
			SimilarityScore: userIDsWithSimilarityScore[profile.ID],
			Explanation:     explanations[profile.ID],
		})
	}

	return recommendedProfilesAndSimilarityResponse, nil
}

// quests/recommend
func (s *QuestService) RecommendQuests(ctx context.Context, userID int) (*models.RecommendationService_RecommendQuests_Resp, error) {
	questIDS, err := s.questRepo.GetUserQuestIDs(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user quest ids (while recommend quests): %w", err)
	}

	if len(questIDS) == 0 {
		return &models.RecommendationService_RecommendQuests_Resp{}, nil
	}

	req := models.RecommendationService_RecommendQuests_Req{
		UserQuestIDs: questIDS,
	}

	return s.recommendQuests(ctx, req)
}

func (s *QuestService) recommendQuests(ctx context.Context, req models.RecommendationService_RecommendQuests_Req) (*models.RecommendationService_RecommendQuests_Resp, error) {
	if s.grpcClient == nil {
		return nil, fmt.Errorf("gRPC client is not initialized")
	}

	// Call gRPC service
	response, err := s.grpcClient.RecommendQuests(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error calling gRPC RecommendQuests: %v", err)
	}

	return response, nil
}
