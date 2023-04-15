package handler

import (
	"NORSI-TRANS/appErrors"
	"NORSI-TRANS/models"
	"NORSI-TRANS/usecase"
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type NotionHandler struct {
	notionUC usecase.NotionUseCase
}

func NewHandler(userUseCase usecase.NotionUseCase) NotionHandler {
	return NotionHandler{notionUC: userUseCase}
}

// CreateNotion godoc
// @Summary создание заметки
// @Tags notion
// @Accept  json
// @Produce  json
// @Param userId query string true "id пользователя"
// @Param data body models.Notion true "заметка"
// @Success 200 {object} string
// @Failure 400 {object} appErrors.AppError
// @Failure 500 {object} appErrors.AppError
// @Router /notion [post]
func (h NotionHandler) CreateNotion(ctx *fiber.Ctx) error {
	notion := models.Notion{}
	if err := json.Unmarshal(ctx.Body(), &notion); err != nil {
		return err
	}

	userId, err := parseId(ctx.Query("userId", "0"))
	if err != nil {
		return err
	}

	id, err := h.notionUC.InsertNotion(context.Background(), notion, userId)
	if err != nil {
		return err
	}

	response := make(map[string]int64)
	response["notionId"] = id
	responseByte, err := json.Marshal(response)
	if err != nil {

	}

	ctx.Context().Success(fiber.MIMEApplicationJSON, responseByte)
	return nil
}

// GetNotion godoc
// @Summary получение заметки
// @Tags notion
// @Accept  json
// @Produce  json
// @Param userId query string true "id пользователя"
// @Param id path string true "id заметки"
// @Success 200 {object} models.Notion
// @Failure 400 {object} appErrors.AppError
// @Failure 403 {object} appErrors.AppError
// @Failure 404 {object} appErrors.AppError
// @Failure 500 {object} appErrors.AppError
// @Router /notion/{id} [get]
func (h NotionHandler) GetNotion(ctx *fiber.Ctx) error {
	notionId, err := parseId(ctx.Params("id", "0"))
	if err != nil {
		return err
	}

	userId, err := parseId(ctx.Query("userId", "0"))
	if err != nil {
		return err
	}

	if notionId == 0 {
		return appErrors.ErrBadRequest
	}

	model, err := h.notionUC.GetNotion(context.Background(), notionId, userId)
	if err != nil {
		return err
	}

	responseByte, err := json.Marshal(model)
	if err != nil {
		return appErrors.ErrBaseApp
	}

	ctx.Context().Success(fiber.MIMEApplicationJSON, responseByte)

	return nil
}

// DeleteNotion godoc
// @Summary удаление заметки
// @Tags notion
// @Accept  json
// @Produce  json
// @Param userId query string true "id пользователя"
// @Param id path string true "id заметки"
// @Success 200 {object} string
// @Failure 400 {object} appErrors.AppError
// @Failure 404 {object} appErrors.AppError
// @Failure 403 {object} appErrors.AppError
// @Failure 500 {object} appErrors.AppError
// @Router /notion/{id} [delete]
func (h NotionHandler) DeleteNotion(ctx *fiber.Ctx) error {
	notionId, err := parseId(ctx.Params("id", "0"))
	if err != nil {
		return err
	}

	userId, err := parseId(ctx.Query("userId", "0"))
	if err != nil {
		return err
	}

	if err = h.notionUC.DeleteNotion(context.Background(), notionId, userId); err != nil {
		return err
	}

	return nil
}

// GetUserNotions godoc
// @Summary получение всех заметок пользователя
// @Tags notion
// @Accept  json
// @Produce  json
// @Param userId query string true "id пользователя"
// @Success 200 {object} []models.Notion
// @Success 400 {object} appErrors.AppError
// @Failure 500 {object} appErrors.AppError
// @Router /notions [get]
func (h NotionHandler) GetUserNotions(ctx *fiber.Ctx) error {
	userId, err := parseId(ctx.Query("userId", "0"))
	if err != nil {
		return err
	}

	notions, err := h.notionUC.GetUserNotions(context.Background(), userId)
	if err != nil {
		return err
	}

	responseByte, err := json.Marshal(notions)
	if err != nil {
		return appErrors.ErrBaseApp
	}

	ctx.Context().Success(fiber.MIMEApplicationJSON, responseByte)

	return nil
}

func parseId(id string) (int64, error) {
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if errors.Is(err, strconv.ErrSyntax) {
		return parsedId, appErrors.ErrBadRequest
	} else if err != nil {
		return parsedId, appErrors.ErrBaseApp
	}

	if parsedId == 0 {
		return parsedId, appErrors.ErrBadRequest
	}

	return parsedId, nil
}
