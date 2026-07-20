package admin_auth_events

import (
	"example/admin/gateway/internal/adapt/infra/clients/auth/cachable"
	"example/admin/gateway/internal/api/kafcon/handlers/admin_auth_events/actions"
	"example/admin/gateway/internal/api/kafcon/handlers/admin_auth_events/kafapi"
	"example/admin/gateway/internal/api/kafcon/kernel"
)

func Register(
	router kernel.Router,
	sAuthCacher cachable.Cacher,
	sMws []func(kafapi.ServiceInterface) kafapi.ServiceInterface,
	hMws []func(kernel.HandlerInterface) kernel.HandlerInterface,
) {
	var service kafapi.ServiceInterface = newServiceDefault(
		actions.NewSessionClosedV1(sAuthCacher),
	)

	for i := len(sMws) - 1; i >= 0; i-- {
		service = sMws[i](service)
	}

	handler := newHandler(service)

	for i := len(hMws) - 1; i >= 0; i-- {
		handler = hMws[i](handler)
	}

	router.Register(TopicName, handler)
}
