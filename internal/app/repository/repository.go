package repository

import (
	"context"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	customError "github.com/Bernigend/mb-cw3-phll-group-service/internal/app/custom-errors"

	"github.com/Bernigend/mb-cw3-phll-group-service/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

// Создаёт новое подключение к базе данных
func NewRepository(dsn string) (*Repository, error) {
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	return &Repository{db: conn}, nil
}

// Закрывает соединение с базой данных, если оно было установлено
func (r Repository) Close() error {
	if r.db == nil {
		return nil
	}

	db, err := r.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

// Выполняет автоматические миграции в базе данных
func (r Repository) AutoMigrate() error {
	if r.db == nil {
		return nil
	}

	return r.db.AutoMigrate(&ds.Group{})
}

func (r Repository) GetGroup(ctx context.Context, filter *ds.Group) (*ds.Group, error) {
	var group ds.Group

	err := r.db.Where(filter).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customError.NotFound.NewWrap(ctx, "группа не найдена", err)
		} else {
			return nil, customError.Internal.NewWrap(ctx, "произошла непредвиденная ошибка", err)
		}
	}

	return &group, nil
}

func (r Repository) GetGroupList(ctx context.Context, filter *ds.Group) (ds.GroupList, error) {
	var groupsList []*ds.Group

	err := r.db.Where(filter).Find(&groupsList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customError.NotFound.NewWrap(ctx, "по указанным фильтрам групп не найдено", err)
		} else {
			return nil, customError.Internal.NewWrap(ctx, "произошла непредвиденная ошибка", err)
		}
	}

	return groupsList, nil
}

func (r Repository) AddGroup(ctx context.Context, group *ds.Group) (uuid.UUID, error) {
	err := r.db.Create(&group).Error
	if err != nil {
		return uuid.Nil, customError.Internal.NewWrap(ctx, "произошла непредвиденная ошибка", err)
	}

	return group.UUID, nil
}
