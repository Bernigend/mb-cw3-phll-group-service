package service

import (
	"context"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"

	customErrors "github.com/Bernigend/mb-cw3-phll-group-service/internal/app/custom-errors"
	"github.com/Bernigend/mb-cw3-phll-group-service/internal/app/ds"
	api "github.com/Bernigend/mb-cw3-phll-group-service/pkg/group-service-api"
)

type IRepository interface {
	GetGroup(ctx context.Context, filter *ds.Group) (*ds.Group, error)
	GetGroupList(ctx context.Context, filter *ds.Group) (ds.GroupList, error)

	AddGroup(ctx context.Context, group *ds.Group) error
}

type Service struct {
	repository IRepository
}

func NewService(repository IRepository) *Service {
	return &Service{repository: repository}
}

func (s Service) GetGroupByName(ctx context.Context, groupName string) (*api.GetGroup_Response, error) {
	if len(groupName) == 0 {
		return nil, customErrors.InvalidArgument.New(ctx, "название группы не может быть пустым")
	}

	group, err := s.repository.GetGroup(ctx, &ds.Group{Name: groupName})
	if err != nil {
		return nil, err
	}

	return &api.GetGroup_Response{
		GroupUuid:            group.UUID.String(),
		GroupName:            group.Name,
		SemesterStartAt:      group.SemesterStart.Format(ds.GroupSemesterStartFormat),
		SemesterEndAt:        group.SemesterEnd.Format(ds.GroupSemesterEndFormat),
		IsFirstWeekNumerator: group.IsFirstWeekNumerator,
		Department:           group.Department,
		Faculty:              group.Faculty,
	}, nil
}

func (s Service) GetGroupByUUID(ctx context.Context, groupUUID uuid.UUID) (*api.GetGroup_Response, error) {
	if groupUUID == uuid.Nil {
		return nil, customErrors.InvalidArgument.New(ctx, "неверный UUID группы")
	}

	group, err := s.repository.GetGroup(ctx, &ds.Group{
		BaseModel: ds.BaseModel{UUID: groupUUID},
	})
	if err != nil {
		return nil, err
	}

	return &api.GetGroup_Response{
		GroupUuid:            group.UUID.String(),
		GroupName:            group.Name,
		SemesterStartAt:      group.SemesterStart.Format(ds.GroupSemesterStartFormat),
		SemesterEndAt:        group.SemesterEnd.Format(ds.GroupSemesterEndFormat),
		IsFirstWeekNumerator: group.IsFirstWeekNumerator,
		Department:           group.Department,
		Faculty:              group.Faculty,
	}, nil
}

func (s Service) GetGroupList(ctx context.Context, filter *api.GetGroupList_Request) (*api.GetGroupList_Response, error) {
	department := filter.GetDepartment()
	faculty := filter.GetFaculty()

	groupList, err := s.repository.GetGroupList(ctx, &ds.Group{
		Department: department,
		Faculty:    faculty,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*api.GetGroupList_GroupItem, 0, len(groupList))

	for _, group := range groupList {
		result = append(result, &api.GetGroupList_GroupItem{
			GroupUuid:            group.UUID.String(),
			GroupName:            group.Name,
			SemesterStartAt:      group.SemesterStart.Format(ds.GroupSemesterStartFormat),
			SemesterEndAt:        group.SemesterEnd.Format(ds.GroupSemesterEndFormat),
			IsFirstWeekNumerator: group.IsFirstWeekNumerator,
			Department:           group.Department,
			Faculty:              group.Faculty,
		})
	}

	return &api.GetGroupList_Response{GroupsList: result}, nil
}

func (s Service) AddGroups(ctx context.Context, groupsList []*api.AddGroups_GroupItem) (*api.AddGroups_Response, error) {
	var errorString string

	if groupsList == nil || len(groupsList) == 0 {
		return nil, customErrors.InvalidArgument.New(ctx, "список групп не может быть пустым")
	}

	result := make([]*api.AddGroups_ResultItem, 0, len(groupsList))

	for _, group := range groupsList {
		if len(group.GetGroupName()) > ds.GroupNameMaxLength {
			result = append(result, &api.AddGroups_ResultItem{
				Result: false,
				Error:  fmt.Sprintf("максимальная длина названия группы = %d", ds.GroupNameMaxLength),
			})
			continue
		}

		if _, err := s.GetGroupByName(ctx, group.GetGroupName()); err == nil {
			result = append(result, &api.AddGroups_ResultItem{
				Result: false,
				Error:  "группа с таким названием уже существует",
			})
			continue
		}

		semesterStart, err := time.Parse(ds.GroupSemesterStartFormat, group.GetSemesterStartAt())
		if err != nil {
			result = append(result, &api.AddGroups_ResultItem{
				Result: false,
				Error:  fmt.Sprintf("неверный формат времени начала семестра, требуется: %s и т.п.", ds.GroupSemesterStartFormat),
			})
			continue
		}

		semesterEnd, err := time.Parse(ds.GroupSemesterEndFormat, group.GetSemesterEndAt())
		if err != nil {
			result = append(result, &api.AddGroups_ResultItem{
				Result: false,
				Error:  fmt.Sprintf("неверный формат времени окончания семестра, требуется: %s и т.п.", ds.GroupSemesterEndFormat),
			})
			continue
		}

		if semesterStart.After(semesterEnd) {
			result = append(result, &api.AddGroups_ResultItem{
				Result: false,
				Error:  "время начала семестра не может быть позже времени окончания",
			})
			continue
		}

		if len(group.GetDepartment()) > ds.GroupDepartmentMaxLength {
			result = append(result, &api.AddGroups_ResultItem{
				Result: false,
				Error:  fmt.Sprintf("максимальная длина названия кафедры = %d", ds.GroupDepartmentMaxLength),
			})
			continue
		}

		if len(group.GetFaculty()) > ds.GroupFacultyMaxLength {
			result = append(result, &api.AddGroups_ResultItem{
				Result: false,
				Error:  fmt.Sprintf("максимальная длина названия факультета = %d", ds.GroupFacultyMaxLength),
			})
			continue
		}

		err = s.repository.AddGroup(ctx, &ds.Group{
			Name:                 group.GetGroupName(),
			SemesterStart:        semesterStart,
			SemesterEnd:          semesterEnd,
			IsFirstWeekNumerator: group.GetIsFirstWeekNumerator(),
			Department:           group.GetDepartment(),
			Faculty:              group.GetFaculty(),
		})

		if err != nil {
			errorString = customErrors.ToGRPC(err).Error()
		} else {
			errorString = ""
		}

		result = append(result, &api.AddGroups_ResultItem{
			Result: len(errorString) == 0,
			Error:  errorString,
		})
	}

	return &api.AddGroups_Response{ResultsList: result}, nil
}
