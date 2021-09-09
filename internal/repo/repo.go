package repo

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"ova-conversation-api/internal/domain"
)

type Repo interface {
	AddEntity(entity domain.Conversation) (uint64, error)
	AddEntities(entities []domain.Conversation) error
	ListEntities(limit, offset uint64) ([]domain.Conversation, error)
	DescribeEntity(entityId uint64) (*domain.Conversation, error)
	RemoveEntity(entityId uint64) (uint64, error)
	UpdateEntity(entity domain.Conversation) (uint64, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{db: db}
}

func (r *repo) AddEntity(entity domain.Conversation) (uint64, error) {
	query := squirrel.
		Insert("conversations").
		Columns("user_id", "text").
		Values(entity.UserID, entity.Text).
		Suffix("RETURNING \"id\"").
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	var id int64
	err := query.QueryRow().Scan(&id)
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (r *repo) AddEntities(entities []domain.Conversation) error {
	sq := squirrel.
		Insert("conversations").
		Columns("user_id", "text")

	for i := range entities {
		sq = sq.Values(entities[i].UserID, entities[i].Text)
	}

	query, args, err := sq.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, args...)

	return err
}

func (r *repo) ListEntities(limit uint64, offset uint64) ([]domain.Conversation, error) {
	query, args, err := squirrel.
		Select("id", "user_id", "text", "date").
		From("conversations").
		OrderBy("id").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	res := make([]domain.Conversation, 0, limit)
	for rows.Next() {
		var id, userId uint64
		var text string
		var date time.Time
		if err := rows.Scan(&id, &userId, &text, &date); err != nil {
			return nil, err
		}

		res = append(res, domain.Conversation{ID: id, UserID: userId, Text: text, Date: date})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *repo) DescribeEntity(entityId uint64) (*domain.Conversation, error) {
	query, args, err := squirrel.
		Select("id", "user_id", "text", "date").
		From("conversations").
		Where("id = ?", entityId).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRowx(query, args...)
	var id, userId uint64
	var text string
	var date time.Time
	err = row.Scan(&id, &userId, &text, &date)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	result := domain.Conversation{ID: id, UserID: userId, Text: text, Date: date}

	return &result, nil
}

func (r *repo) RemoveEntity(entityId uint64) (uint64, error) {
	query, args, err := squirrel.
		Delete("conversations").
		Where(squirrel.Eq{"id": entityId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return entityId, nil
}

func (r *repo) UpdateEntity(entity domain.Conversation) (uint64, error) {
	query, args, err := squirrel.
		Update("conversations").
		Set("text", entity.Text).
		Where(squirrel.Eq{"id": entity.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, err
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if ra == 0 {
		return 0, nil
	}

	return entity.ID, nil
}
