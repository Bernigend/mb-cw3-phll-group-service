package ds

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Базовая модель включает в себя общие столбцы всех таблиц
type BaseModel struct {
	UUID      uuid.UUID `gorm:"type:uuid;primaryKey;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *BaseModel) BeforeCreate(_ *gorm.DB) error {
	b.UUID = uuid.NewV4()
	return nil
}

type GroupList []*Group

const (
	GroupNameMaxLength       = 128
	GroupDepartmentMaxLength = 5
	GroupFacultyMaxLength    = GroupDepartmentMaxLength

	GroupSemesterStartFormat = time.RFC3339
	GroupSemesterEndFormat   = time.RFC3339
)

type Group struct {
	BaseModel
	Name                 string    `gorm:"type:string;size:128;not null"`
	SemesterStart        time.Time `gorm:"type:time;not null"`
	SemesterEnd          time.Time `gorm:"type:time;not null"`
	IsFirstWeekNumerator bool      `gorm:"type:bool;not null"`
	Department           string    `gorm:"type:string;size:5;not null"`
	Faculty              string    `gorm:"type:string;size:5;not null"`
}
