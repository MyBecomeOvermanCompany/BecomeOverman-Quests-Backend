package services

import (
	"BecomeOverMan/internal/integrations"
	"BecomeOverMan/internal/models"
	"BecomeOverMan/internal/repositories"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

type QuestService struct {
	questRepo         *repositories.QuestRepository
	questTreeRepo     *repositories.QuestTreeRepository
	habitTrackingRepo *repositories.HabitTrackingRepository
	userRepo          *repositories.UserRepository
}

func NewQuestService(
	questRepo *repositories.QuestRepository,
	userRepo *repositories.UserRepository,
) *QuestService {
	return &QuestService{
		questRepo:         questRepo,
		questTreeRepo:     repositories.NewQuestTreeRepository(questRepo.GetDB()),
		habitTrackingRepo: repositories.NewHabitTrackingRepository(questRepo.GetDB()),
		userRepo:          userRepo,
	}
}

// GetAvailableQuests returns quests available for the user (with tree check)
func (s *QuestService) GetAvailableQuests(ctx context.Context, userID int) ([]models.Quest, error) {
	// Получаем квесты с информацией о дереве
	questsWithTree, err := s.questTreeRepo.GetAvailableQuestsWithTree(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Фильтруем только разблокированные квесты
	availableQuests := make([]models.Quest, 0)
	for _, qwt := range questsWithTree {
		if qwt.IsUnlocked {
			availableQuests = append(availableQuests, qwt.Quest)
		}
	}

	return availableQuests, nil
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

// GetDevelopmentBranches - получить все ветки развития
func (s *QuestService) GetDevelopmentBranches(ctx context.Context) ([]models.DevelopmentBranch, error) {
	return s.questTreeRepo.GetDevelopmentBranches(ctx)
}

// GetQuestTree - получить дерево квестов для пользователя
func (s *QuestService) GetQuestTree(ctx context.Context, userID int) ([]models.QuestWithPrerequisites, error) {
	return s.questTreeRepo.GetAvailableQuestsWithTree(ctx, userID)
}

// GetUserPassiveBuffs - получить пассивные баффы пользователя
func (s *QuestService) GetUserPassiveBuffs(ctx context.Context, userID int) ([]models.PassiveBuff, error) {
	return s.questTreeRepo.GetUserPassiveBuffs(ctx, userID)
}

// PurchaseQuest handles the purchase of a quest by a user
func (s *QuestService) PurchaseQuest(ctx context.Context, userID, questID int) error {
	// Проверяем, разблокирован ли квест
	isUnlocked, _, _, err := s.questTreeRepo.CheckQuestUnlocked(ctx, userID, questID)
	if err != nil {
		return fmt.Errorf("error checking quest unlock: %w", err)
	}
	if !isUnlocked {
		return fmt.Errorf("quest is locked - complete prerequisites first")
	}

	err = s.questRepo.PurchaseQuest(ctx, userID, questID)
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
	// 1. Создаем URL
	url := integrations.Recommendation_Service_BASE_URL + "/users/add"

	// 2. Кодируем в JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// 3. Создаем io.Reader из JSON
	body := bytes.NewBuffer(jsonData)

	// 4. Делаем POST запрос с таймаутом
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Post(url, "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("error making POST request to recommendation service: %v", err)
	}

	defer resp.Body.Close()

	// 5. Проверяем статус
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recommendation service returned status %d", resp.StatusCode)
	}

	// 6. Читаем и парсим ответ
	var response map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
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

// GetQuestDetails - получает детали квеста
func (s *QuestService) GetQuestDetails(ctx context.Context, questID, userID int) (*models.Quest, error) {
	return s.questRepo.GetQuestDetails(ctx, questID, userID)
}

// CompleteQuest finalizes the quest completion
func (s *QuestService) CompleteQuest(ctx context.Context, userID, questID int) error {
	return s.questRepo.CompleteQuest(ctx, userID, questID)
}

// GetQuestDetailsForAI - получает детали квеста для AI (без userID)
func (s *QuestService) GetQuestDetailsForAI(ctx context.Context, questID int) (*models.Quest, error) {
	return s.questRepo.GetQuestDetails(ctx, questID, 0)
}

func (s *QuestService) CreateSharedQuest(user1ID, user2ID, questID int) error {
	return s.questRepo.CreateSharedQuest(user1ID, user2ID, questID)
}

func (s *QuestService) GetSharedQuests(ctx context.Context, userID int) ([]models.SharedQuest, error) {
	return s.questRepo.GetSharedQuests(ctx, userID)
}

func (s *QuestService) SaveQuestToDB(quest *models.Quest, tasks []models.Task) (int, error) {
	questID, err := s.questRepo.SaveQuestToDB(quest, tasks)
	if err != nil {
		return 0, err
	}

	// Связываем квест с ветками развития на основе category
	if quest.Category != "" {
		ctx := context.Background()
		// Получаем основную ветку по категории
		branches, _ := s.questTreeRepo.GetDevelopmentBranches(ctx)
		for _, branch := range branches {
			if branch.Level == 1 && branch.Name == quest.Category {
				// Связываем с основной веткой (вес 1.0)
				s.questTreeRepo.LinkQuestToBranch(ctx, questID, branch.ID, 1.0)
				break
			}
		}
	}

	return questID, nil
}

func (s *QuestService) SearchQuests(
	ctx context.Context,
	req models.RecommendationService_SearchQuest_Request,
	userID int,
) (models.SearchQuestsResponse, error) {
	// 1. Создаем URL
	url := integrations.Recommendation_Service_BASE_URL + "/search"

	// 2. Кодируем в JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// 3. Создаем io.Reader из JSON
	body := bytes.NewBuffer(jsonData)

	// 4. Делаем POST запрос с таймаутом
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Post(url, "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("error making POST request to recommendation service: %v", err)
	}

	defer resp.Body.Close()

	// 5. Проверяем статус
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recommendation service returned status %d", resp.StatusCode)
	}

	// 6. Читаем и парсим ответ
	var response models.RecommendationService_SearchQuests_Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 7. достаем ids
	questsIDS := make([]int, len(response.Results))
	for i, result := range response.Results {
		questsIDS[i] = result.ID
	}

	var questsWithDetails []models.Quest

	// 8. Достаем квесты с деталями из БД (тут сразу и те что есть у юзера и те что еще не куплены)
	questsWithDetails, err = s.questRepo.SearchQuestsWithDetailsByIDs(ctx, questsIDS)
	if err != nil {
		slog.ErrorContext(ctx, "ошибка получения квестов из БД с указанными ids во время поиска",
			"error", err,
			"ids", questsIDS,
		)
		return nil, fmt.Errorf("В поиске квестов по запросу произошла внутренняя ошибка: %w", err)
	}

	// 9. Возвращаем результат = []struct{questWithDetails, SimilaryScore}
	questsWithDetailsAndSimilarityResponse := models.NewSearchQuestsResponse(questsWithDetails, response)

	return questsWithDetailsAndSimilarityResponse, nil
}

func (s *QuestService) RecommendFriends(
	ctx context.Context,
	req models.RecommendationService_RecommendUsers_Request,
) ([]models.UserProfileWithSimilarityScore, error) {
	// 1. Создаем URL
	url := integrations.Recommendation_Service_BASE_URL + "/users/recommend"

	// 2. Кодируем в JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// 3. Создаем io.Reader из JSON
	body := bytes.NewBuffer(jsonData)

	// 4. Делаем POST запрос с таймаутом
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Post(url, "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("error making POST request to recommendation service: %v", err)
	}

	defer resp.Body.Close()

	// 5. Проверяем статус
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recommendation service returned status %d", resp.StatusCode)
	}

	// 6. Читаем и парсим ответ
	var response models.RecommendationService_RecommendUsers_Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 7. достаем ids
	userIDs := make([]int, len(response.Results))
	userIDsWithSimilarityScore := make(map[int]float64, len(response.Results))
	explanations := make(map[int]map[string]any, len(response.Results))
	for i, result := range response.Results {
		userIDs[i] = result.UserID
		userIDsWithSimilarityScore[result.UserID] = result.SimilarityScore
		explanations[result.UserID] = result.Explanation
	}

	// 8. Достаем профили потенциальных друзей (кроме пользователей с которыми уже дружба)
	recommendedProfiles, err := s.userRepo.GetProfiles(userIDs)
	if err != nil {
		slog.ErrorContext(ctx, "ошибка получения профилей из БД с указанными ids во время рекомендации друзей",
			"error", err,
			"ids", userIDs,
		)
		return nil, fmt.Errorf("В рекомендации друзей по запросу произошла внутренняя ошибка: %w", err)
	}

	// 9. Убираем оттуда пользователей с которыми уже дружба
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

	// 10. Возвращаем результат = []models.UserProfileWithSimilarityScore
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
	// 1. Создаем URL
	url := integrations.Recommendation_Service_BASE_URL + "/quests/recommend"

	// 2. Кодируем в JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// 3. Создаем io.Reader из JSON
	body := bytes.NewBuffer(jsonData)

	// 4. Делаем POST запрос с таймаутом
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Post(url, "application/json", body)
	if err != nil {
		return nil, fmt.Errorf("error making POST request to recommendation service: %v", err)
	}

	defer resp.Body.Close()

	// 5. Проверяем статус
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recommendation service returned status %d", resp.StatusCode)
	}

	// 6. Читаем и парсим ответ
	var response models.RecommendationService_RecommendQuests_Resp
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}
