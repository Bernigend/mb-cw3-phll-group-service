package endpoint

import (
	"context"

	uuid "github.com/satori/go.uuid"

	customErrors "github.com/Bernigend/mb-cw3-phll-group-service/internal/app/custom-errors"

	api "github.com/Bernigend/mb-cw3-phll-group-service/pkg/group-service-api"
)

// Необходимые для работы API методы сервиса
type IService interface {
	GetGroupByName(ctx context.Context, groupName string) (*api.GetGroup_Response, error)
	GetGroupByUUID(ctx context.Context, groupUUID uuid.UUID) (*api.GetGroup_Response, error)

	GetGroupList(ctx context.Context, filter *api.GetGroupList_Request) (*api.GetGroupList_Response, error)

	AddGroups(ctx context.Context, groupsList []*api.AddGroups_GroupItem) (*api.AddGroups_Response, error)
}

type Endpoint struct {
	service IService
}

func NewEndpoint(service IService) *Endpoint {
	return &Endpoint{service: service}
}

// [Метод API] Возвращает группу
func (e Endpoint) GetGroup(ctx context.Context, request *api.GetGroup_Request) (*api.GetGroup_Response, error) {
	if len(request.GetGroupUuid()) > 0 {
		parsedGroupUuid, err := uuid.FromString(request.GetGroupUuid())
		if err != nil {
			return nil, customErrors.InvalidArgument.New(ctx, "invalid group UUID")
		}

		result, err := e.service.GetGroupByUUID(ctx, parsedGroupUuid)
		return result, customErrors.ToGRPC(err)
	}

	if len(request.GetGroupName()) > 0 {
		result, err := e.service.GetGroupByName(ctx, request.GetGroupName())
		return result, customErrors.ToGRPC(err)
	}

	return nil, customErrors.InvalidArgument.New(ctx, "expected group name or group uuid")
}

// [Метод API] Добавляет группу
func (e Endpoint) AddGroups(ctx context.Context, request *api.AddGroups_Request) (*api.AddGroups_Response, error) {
	if groupsList := request.GetGroupsList(); groupsList != nil {
		result, err := e.service.AddGroups(ctx, groupsList)
		return result, customErrors.ToGRPC(err)
	}

	return nil, customErrors.InvalidArgument.New(ctx, "expected groups list")
}

// [Метод API] Возвращает список групп
func (e Endpoint) GetGroupList(ctx context.Context, request *api.GetGroupList_Request) (*api.GetGroupList_Response, error) {
	result, err := e.service.GetGroupList(ctx, request)
	return result, customErrors.ToGRPC(err)
}
