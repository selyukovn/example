package kafapi

import (
	"context"
	"fmt"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------

// handle
//
// Ошибки:
//   - ErrorHandling
func handle(service ServiceInterface, ctx context.Context, meta *Meta, data any) error {
	assert.NotNilDeepMust(service)
	assert.NotNilDeepMust(ctx)
	assert.NotNilDeepMust(meta)
	assert.NotNilDeepMust(data)

	_om_ := "kafapi" + "." + "service"

	var err error
	switch vData := data.(type) {
	case *DataAccountCreatedV1:
		err = service.AccountCreatedV1(ctx, meta, vData)
	case *DataAccountDeactivatedV1:
		err = service.AccountDeactivatedV1(ctx, meta, vData)
	case *DataIpWhitelistChangedV1:
		err = service.IpWhitelistChangedV1(ctx, meta, vData)
	case *DataSessionCreatedV1:
		err = service.SessionCreatedV1(ctx, meta, vData)
	case *DataSessionClosedV1:
		err = service.SessionClosedV1(ctx, meta, vData)
	default:
		// добавлен тип сообщения, но не вызван метод обработки -- ошибка разработчика.
		panic(fmt.Errorf("%s не знает, как обработать сообщение %T", _om_, data))
	}

	if err != nil {
		return newErrorHandling(err)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

type ServiceInterface interface {
	AccountCreatedV1(context.Context, *Meta, *DataAccountCreatedV1) error
	AccountDeactivatedV1(context.Context, *Meta, *DataAccountDeactivatedV1) error
	IpWhitelistChangedV1(context.Context, *Meta, *DataIpWhitelistChangedV1) error
	SessionCreatedV1(context.Context, *Meta, *DataSessionCreatedV1) error
	SessionClosedV1(context.Context, *Meta, *DataSessionClosedV1) error
}

// ---------------------------------------------------------------------------------------------------------------------

var _ ServiceInterface = UnimplementedService{}

type UnimplementedService struct{}

func (s UnimplementedService) AccountCreatedV1(context.Context, *Meta, *DataAccountCreatedV1) error {
	return std.NewErrorRuntimeFf("%s не реализован!", "DataAccountCreatedV1")
}
func (s UnimplementedService) AccountDeactivatedV1(context.Context, *Meta, *DataAccountDeactivatedV1) error {
	return std.NewErrorRuntimeFf("%s не реализован!", "DataAccountDeactivatedV1")
}
func (s UnimplementedService) IpWhitelistChangedV1(context.Context, *Meta, *DataIpWhitelistChangedV1) error {
	return std.NewErrorRuntimeFf("%s не реализован!", "DataIpWhitelistChangedV1")
}
func (s UnimplementedService) SessionCreatedV1(context.Context, *Meta, *DataSessionCreatedV1) error {
	return std.NewErrorRuntimeFf("%s не реализован!", "DataSessionCreatedV1")
}
func (s UnimplementedService) SessionClosedV1(context.Context, *Meta, *DataSessionClosedV1) error {
	return std.NewErrorRuntimeFf("%s не реализован!", "DataSessionClosedV1")
}

// ---------------------------------------------------------------------------------------------------------------------

var _ ServiceInterface = EmptyService{}

type EmptyService struct{}

func (s EmptyService) AccountCreatedV1(context.Context, *Meta, *DataAccountCreatedV1) error {
	return nil
}
func (s EmptyService) AccountDeactivatedV1(context.Context, *Meta, *DataAccountDeactivatedV1) error {
	return nil
}
func (s EmptyService) IpWhitelistChangedV1(context.Context, *Meta, *DataIpWhitelistChangedV1) error {
	return nil
}
func (s EmptyService) SessionCreatedV1(context.Context, *Meta, *DataSessionCreatedV1) error {
	return nil
}
func (s EmptyService) SessionClosedV1(context.Context, *Meta, *DataSessionClosedV1) error {
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

var _ ServiceInterface = ServiceDecoratorDefault{}

type ServiceDecoratorDefault struct {
	origin ServiceInterface
}

// NewServiceDecoratorDefault
//
// Паникует при нулевых аргументах.
func NewServiceDecoratorDefault(origin ServiceInterface) ServiceDecoratorDefault {
	assert.NotNilDeepMust(origin)

	return ServiceDecoratorDefault{
		origin: origin,
	}
}

func (s ServiceDecoratorDefault) AccountCreatedV1(ctx context.Context, meta *Meta, data *DataAccountCreatedV1) error {
	return s.origin.AccountCreatedV1(ctx, meta, data)
}
func (s ServiceDecoratorDefault) AccountDeactivatedV1(ctx context.Context, meta *Meta, data *DataAccountDeactivatedV1) error {
	return s.origin.AccountDeactivatedV1(ctx, meta, data)
}
func (s ServiceDecoratorDefault) IpWhitelistChangedV1(ctx context.Context, meta *Meta, data *DataIpWhitelistChangedV1) error {
	return s.origin.IpWhitelistChangedV1(ctx, meta, data)
}
func (s ServiceDecoratorDefault) SessionCreatedV1(ctx context.Context, meta *Meta, data *DataSessionCreatedV1) error {
	return s.origin.SessionCreatedV1(ctx, meta, data)
}
func (s ServiceDecoratorDefault) SessionClosedV1(ctx context.Context, meta *Meta, data *DataSessionClosedV1) error {
	return s.origin.SessionClosedV1(ctx, meta, data)
}

// ---------------------------------------------------------------------------------------------------------------------
