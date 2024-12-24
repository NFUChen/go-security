package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type INotificationApproachRepository interface {
	GetNotificationApproachesByUserID(ctx context.Context, userID uint) ([]NotificationApproach, error)
	UpdateNotificationApproaches(ctx context.Context, userID uint, approaches []NotificationApproach) error
	GetNumberOfApproachesByUserID(ctx context.Context, userID uint) int
	SaveNotificationApproaches(ctx context.Context, userID uint, approaches []NotificationApproach) error
}

type NotificationApproachRepository struct {
	Engine *redis.Client
	Key    string
}

func NewNotificationApproachRepository(engine *redis.Client) *NotificationApproachRepository {
	return &NotificationApproachRepository{
		Engine: engine,
		Key:    "notification_approaches",
	}
}

func (repo NotificationApproachRepository) createKey(userID uint) string {
	return fmt.Sprintf("%s:%d", repo.Key, userID)
}

func (repo NotificationApproachRepository) ApproachesAsJson(approaches []NotificationApproach) ([]byte, error) {
	value, err := json.Marshal(approaches)
	if err != nil {
		return []byte{}, err
	}
	return value, nil
}

func (repo NotificationApproachRepository) GetNotificationApproachesByUserID(ctx context.Context, userID uint) ([]NotificationApproach, error) {
	value, err := repo.Engine.Get(ctx, repo.createKey(userID)).Result()
	if err != nil {
		return nil, err
	}

	var approaches []NotificationApproach
	err = json.Unmarshal([]byte(value), &approaches)
	if err != nil {
		return nil, err
	}

	return approaches, nil
}

func (repo NotificationApproachRepository) GetNumberOfApproachesByUserID(ctx context.Context, userID uint) int {
	approaches, err := repo.GetNotificationApproachesByUserID(ctx, userID)
	if err != nil {
		return 0
	}
	return len(approaches)
}

func (repo NotificationApproachRepository) UpdateNotificationApproaches(ctx context.Context, userID uint, approaches []NotificationApproach) error {
	value, err := repo.ApproachesAsJson(approaches)
	if err != nil {
		return err
	}

	return repo.Engine.Set(ctx, repo.createKey(userID), value, 0).Err()
}

func (repo NotificationApproachRepository) SaveNotificationApproaches(ctx context.Context, userID uint, approaches []NotificationApproach) error {
	value, err := repo.ApproachesAsJson(approaches)
	if err != nil {
		return err
	}

	return repo.Engine.Set(ctx, repo.createKey(userID), value, 0).Err()
}
