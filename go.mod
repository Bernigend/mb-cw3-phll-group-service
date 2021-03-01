module github.com/Bernigend/mb-cw3-phll-group-service

go 1.15

require (
	github.com/Bernigend/mb-cw3-phll-group-service/pkg/group-service-api v1.0.0
	github.com/satori/go.uuid v1.2.0
	google.golang.org/grpc v1.35.0
	gorm.io/driver/postgres v1.0.6
	gorm.io/gorm v1.20.11
)

replace (
	github.com/Bernigend/mb-cw3-phll-group-service/pkg/group-service-api => ./pkg/group-service-api
)
