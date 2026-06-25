package feedback

import (
	"context"
	"fmt"

	"github.com/danyele/podp/internal/shared/redis"
	"github.com/google/uuid"
)

type SaveFeedbackUsecase struct {
	Redis redis.Cache
}

func (u *SaveFeedbackUsecase) Execute(ctx context.Context, name, email, message string) error {
	id := uuid.New().String()
	key := fmt.Sprintf("feedback:%s", id)
	data := fmt.Sprintf("%s|%s|%s", name, email, message)
	return u.Redis.Set(ctx, key, data)

}
