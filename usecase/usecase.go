package usecase

import (
	"NORSI-TRANS/appErrors"
	"NORSI-TRANS/models"
	"NORSI-TRANS/repository"
	"context"
)

type NotionUseCase interface {
	InsertNotion(ctx context.Context, notion models.Notion, userId int64) (int64, error)
	GetNotion(ctx context.Context, id, userId int64) (models.Notion, error)
	DeleteNotion(ctx context.Context, id, userId int64) error
	GetUserNotions(ctx context.Context, userId int64) ([]models.Notion, error)
}

type notionUC struct {
	notionRepo repository.NotionRepo
}

func NewNotionUseCase(notionRepo repository.NotionRepo) NotionUseCase {
	return notionUC{notionRepo: notionRepo}
}

func (n notionUC) InsertNotion(ctx context.Context, notion models.Notion, userId int64) (int64, error) {
	notion.UserId = userId

	return n.notionRepo.InsertNotion(ctx, notion)
}

func (n notionUC) GetNotion(ctx context.Context, id, userId int64) (models.Notion, error) {
	notion, err := n.notionRepo.GetNotion(ctx, id)
	if err != nil {
		return notion, err
	}

	if notion.UserId != userId {
		return notion, appErrors.ErrForbidden
	}

	return notion, nil
}

func (n notionUC) DeleteNotion(ctx context.Context, id, userId int64) error {
	if _, err := n.GetNotion(ctx, id, userId); err != nil {
		return err
	}

	return n.notionRepo.DeleteNotion(ctx, id)
}

func (n notionUC) GetUserNotions(ctx context.Context, userId int64) ([]models.Notion, error) {
	return n.notionRepo.GetUserNotions(ctx, userId)
}
