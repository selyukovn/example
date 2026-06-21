package security

import (
	"context"
	"example/admin/gateway/internal/api/http/kernel"
	"example/admin/gateway/internal/infra/clients/auth"
	"fmt"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"net/http"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

const ctxKeyUser = "security.user"

type Security struct {
	authClient auth.ClientInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func New(authClient auth.ClientInterface) Security {
	assert.NotNilDeepMust(authClient)

	return Security{
		authClient: authClient,
	}
}

func (s Security) getSessId(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s Security) middlewareAsGuest(
	w http.ResponseWriter,
	r *http.Request,
	next http.Handler,
) {
	fromUag := kernel.UserAgent(r)
	if fromUag == "" {
		http.Error(w, "User-Agent обязателен", http.StatusBadRequest)
		return
	}

	user := newUserGuest(
		kernel.UserIp(r),
		fromUag,
	)

	ctx := r.Context()
	ctx = context.WithValue(ctx, ctxKeyUser, user)
	r = r.WithContext(ctx)

	logger.InfoFf(ctx, "Позвольтевотель не идетинфицирован")

	next.ServeHTTP(w, r)
}

func (s Security) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_m_ := "Middleware"

		// todo : некий guest-id ???

		sessId := s.getSessId(r)

		// guest
		if sessId == "" {
			s.middlewareAsGuest(w, r, next)
			return // !!!
		}

		// authenticated

		ctx := r.Context()

		fromIp := kernel.UserIp(r)
		fromUag := kernel.UserAgent(r)
		if fromUag == "" {
			http.Error(w, "User-Agent обязателен", http.StatusBadRequest)
			return
		}

		// Может показаться, что логично было бы удалить куку в исключительных ситуациях
		// (например, при отсутствии связанного аккаунта или ошибки валидации).
		// Однако, не стоит зависеть от предположений о работе клиента!
		// Пусть клиент сам разбирается с ситуацией, исходя из кода ответа обработчика.
		res, err := s.authClient.CheckSession(ctx, fromIp, fromUag, sessId)
		switch err.(type) {
		case nil:
		case std.ErrorNotFound:
			// Ошибка может быть как на стороне клиента (передан посторонний идентификатор или его подобие),
			// так и на стороне сервера (битые данные, ошибки, ...).
			// В любом случае дальнейшая обработка возможна только в неавторизованном виде.
			//
			// Ситуацию можно было бы приравнять к использованию протухшей сессии,
			// однако, 401-й код может быть неуместным, например, в обработчике,
			// который делает перенаправление в зависимости от состояния аутентификации пользователя.
			// Поэтому выбор кода ответа необходимо оставить на усмотрение обработчика.
			logger.WarnFf(ctx, "неопознанная сессия %q: %s", std.MaskStrNotFirstLast(sessId), err)
			s.middlewareAsGuest(w, r, next)
			return // !!!
		case auth.ErrorValidation:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return // !!!
		case auth.ErrorAccountAccessDenied:
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return // !!!
		case auth.ErrorSessionClosed:
			// В общем случае для обработчика запроса нет разницы, по какой причине не была пройдена аутентификация.
			// Если необходимость в различии все же появится, то можно использовать дополнительные флаги в User.
			//
			// Выбор кода ответа остается на усмотрение обработчика -- по аналогии с `case std.ErrorNotFound`.
			logger.WarnFf(ctx, "обращение к закрытой сессии %q: %s", std.MaskStrNotFirstLast(sessId), err)
			s.middlewareAsGuest(w, r, next)
			return // !!!
		case std.ErrorRuntime:
			logger.ErrorFf(ctx, std.WrapErrorToRuntime(err, s, _m_).Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return // !!!
		default:
			panic(err)
		}

		user := newUserAuthorized(
			fromIp,
			fromUag,
			sessId,
			res.SessionExpireAt,
			res.AccountId,
		)
		ctx = context.WithValue(ctx, ctxKeyUser, user)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// ---------------------------------------------------------------------------------------------------------------------

// AssociatedUser
//
// Паникует при нулевых аргументах.
// Паникует, если обработчик запроса не был обернут в Middleware.
func (s Security) AssociatedUser(ctx context.Context) *User {
	assert.NotNilDeepMust(ctx)

	_m_ := "AssociatedUser"

	vUser := ctx.Value(ctxKeyUser)
	if vUser == nil {
		panic(fmt.Errorf("%T.%s не нашел %s в контексте запроса", s, _m_, ctxKeyUser))
	}

	user, ok := vUser.(*User)
	if !ok {
		panic(fmt.Errorf("%T.%s нашел %s, но это не %T: %#v", s, _m_, ctxKeyUser, &User{}, vUser))
	}

	return user
}

// ---------------------------------------------------------------------------------------------------------------------

// Authenticate
//
// Паникует при нулевых аргументах.
// Паникует, если User.IsGuest() == false.
//
// Ошибки:
//   - std.ErrorRuntime
func (s Security) Authenticate(ctx context.Context, user *User, w http.ResponseWriter, sessId string) error {
	/* Может показаться, что один из аргументов `ctx` и `user` лишний, поскольку есть */ _ = s.AssociatedUser /**/
	// Это как минимум ухудшило бы интерфейс метода --> Authenticate(ctx, sessId, w) -- что тут аутентифицируется?

	assert.NotNilDeepMust(ctx)
	assert.NotNilDeepMust(user)
	assert.NotNilDeepMust(w)
	assert.Str().NotEmpty().Must(sessId)

	assert.TrueMust(user.IsGuest())

	_m_ := "Authenticate"

	res, err := s.authClient.CheckSession(ctx, user.Ip(), user.UserAgent(), sessId)
	switch err.(type) {
	case nil:
	case std.ErrorNotFound,
		auth.ErrorValidation,
		auth.ErrorAccountAccessDenied,
		auth.ErrorSessionClosed,
		std.ErrorRuntime:
		return std.WrapErrorToRuntime(err, s, _m_)
	default:
		panic(err)
	}

	user.authenticate(sessId, res.SessionExpireAt, res.AccountId)

	return nil
}

// UnAuthenticate
//
// Паникует при нулевых аргументах.
// Паникует, если User.IsAuthenticated() == false.
//
// Ошибки:
//   - std.ErrorRuntime
func (s Security) UnAuthenticate(ctx context.Context, user *User) error {
	/* Может показаться, что один из аргументов `ctx` и `user` лишний, поскольку есть */ _ = s.AssociatedUser /**/
	// Это как минимум ухудшило бы интерфейс метода --> Authenticate(ctx, sessId, w) -- что тут аутентифицируется?

	assert.NotNilDeepMust(ctx)
	assert.NotNilDeepMust(user)

	assert.TrueMust(user.IsAuthenticated())

	_m_ := "UnAuthenticate"

	// Поскольку пользователь аутентифицирован, то он прошел аутентификацию с теми же данными,
	// а значит все возможные ошибки клиента исключены -- любая ошибка, кроме `std.ErrorAlreadyDone`, есть бага.
	err := s.authClient.SignOut(ctx, user.Ip(), user.UserAgent(), user.sessionId())
	switch err.(type) {
	case nil:
	case std.ErrorAlreadyDone:
		// Наиболее вероятная ситуация -- нажатие кнопки выхода после истечения срока жизни сессии.
		// Но теоретически может быть и скрытая ошибка, поэтому нужно обязательно записать предупреждение.
		logger.WarnFf(ctx, "закрытие закрытой сессии %q: %s", std.MaskStrNotFirstLast(user.sessionId()), err)
	case auth.ErrorValidation,
		std.ErrorNotFound,
		auth.ErrorAccountAccessDenied,
		std.ErrorRuntime:
		return std.WrapErrorToRuntime(err, s, _m_)
	default:
		panic(err)
	}

	user.unAuthenticate()

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
