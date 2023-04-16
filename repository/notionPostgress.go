package repository

import (
	"NORSI-TRANS/appErrors"
	"NORSI-TRANS/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type NotionRepo interface {
	InsertNotion(ctx context.Context, notion models.Notion) (int64, error)
	GetNotion(ctx context.Context, id int64) (models.Notion, error)
	DeleteNotion(ctx context.Context, id int64) error
	GetUserNotions(ctx context.Context, userId int64) ([]models.Notion, error)
}

type notionPostgres struct {
	connection      *pgx.Conn
	notionTableName string
}

func NewNotionPostgresRepo(connection *pgx.Conn, notionTableName string) (NotionRepo, error) {
	_, err := connection.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS "+notionTableName+"("+
		"NotionId bigserial unique not null primary key,\n"+
		"UserId     bigserial,\n"+
		"Notion     text default ''\n"+
		");")
	if err != nil {
		return nil, appErrors.ErrBaseApp.Wrap(err, "create table failed")
	}

	return &notionPostgres{
		connection:      connection,
		notionTableName: notionTableName,
	}, nil
}

func (n *notionPostgres) InsertNotion(ctx context.Context, notion models.Notion) (int64, error) {
	sqlStatement := fmt.Sprintf("INSERT INTO %s (UserID, Notion) Values($1, $2) RETURNING notionId", n.notionTableName)

	var id int64
	if err := n.connection.QueryRow(ctx, sqlStatement, notion.UserId, notion.Information).Scan(&id); err != nil {
		return id, appErrors.ErrBaseApp.Wrap(err, "insert notion failed")
	}

	return id, nil
}

func (n *notionPostgres) GetNotion(ctx context.Context, id int64) (models.Notion, error) {
	notion := models.Notion{}

	sqlStatement := fmt.Sprintf("SELECT * FROM %s WHERE notionId = $1", n.notionTableName)
	if err := n.connection.QueryRow(ctx, sqlStatement, id).Scan(&notion.Id, &notion.UserId, &notion.Information); errors.Is(err, pgx.ErrNoRows) {
		return notion, appErrors.ErrNotFound
	} else if err != nil {
		return notion, appErrors.ErrBaseApp.Wrap(err, "get notion failed")
	}

	return notion, nil
}

func (n *notionPostgres) DeleteNotion(ctx context.Context, id int64) error {
	sqlStatement := fmt.Sprintf("DELETE FROM %s WHERE notionId = $1", n.notionTableName)
	if _, err := n.connection.Exec(ctx, sqlStatement, id); err != nil {
		return appErrors.ErrBaseApp.Wrap(err, "delete notion failed")
	}

	return nil
}

func (n *notionPostgres) GetUserNotions(ctx context.Context, userId int64) ([]models.Notion, error) {
	results := make([]models.Notion, 0)

	sqlStatement := fmt.Sprintf("SELECT * FROM %s WHERE userId = $1", n.notionTableName)
	rows, err := n.connection.Query(ctx, sqlStatement, userId)
	if err != nil {
		return results, appErrors.ErrBaseApp.Wrap(err, "get notions failed")
	}

	for rows.Next() {
		notion := models.Notion{}

		err = rows.Scan(&notion.Id, &notion.UserId, &notion.Information)
		if err != nil {
			return nil, appErrors.ErrBaseApp.Wrap(err, "")
		}

		results = append(results, notion)
	}

	return results, nil
}
