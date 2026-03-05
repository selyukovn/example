package security

import (
	"context"
	"errors"
	"example/admin/bff/cmd/http/container"
	"example/admin/bff/cmd/http/kernel"
	"example/admin/bff/cmd/http/kernel_ext"
	"example/admin/bff/internal/infra/clients/auth"
	"fmt"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"net/http"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

const ctxKeyUser = "security.user"

type Security struct {
	ctr               *container.Container
	baseUrl           string
	sessionCookieName string
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func New(
	ctr *container.Container,
	baseUrl string,
	sessionCookieName string,
) *Security {
	assert.NotNilDeepMust(ctr)
	assert.Str().NotEmpty().Must(baseUrl)
	assert.Str().NotEmpty().Must(sessionCookieName)

	return &Security{
		ctr:               ctr,
		baseUrl:           baseUrl,
		sessionCookieName: sessionCookieName,
	}
}

// Сейчас id сессии передается через куки.
// Теоретически, может передаваться, например, в хедере Authorization -- зависит от клиента.

func (s *Security) getSessId(r *http.Request) string {
	c, err := r.Cookie(s.sessionCookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return ""
	}
	return c.Value
}

func (s *Security) setSessId(w http.ResponseWriter, sessId string, sessExpAt time.Time) {
	assert.Str().NotEmpty().Must(sessId)
	assert.Time().NotZero().Must(sessExpAt)

	isHttps, domain := func(baseUrl string) (bool, string) {
		parts := strings.Split(baseUrl, ":")
		return parts[0] == "https", strings.TrimPrefix(parts[1], "//")
	}(s.baseUrl)

	http.SetCookie(w, &http.Cookie{
		Name:     s.sessionCookieName,
		Value:    sessId,
		Expires:  sessExpAt,
		Domain:   domain,
		Path:     "/",
		HttpOnly: true,
		Secure:   isHttps,
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s *Security) allow(forGuests bool, forAuthorized bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessId := s.getSessId(r)

		// guest
		if sessId == "" && forGuests {
			// todo : некий guest-id ???
			user := newUserGuest()

			ctx := r.Context()
			ctx = context.WithValue(ctx, ctxKeyUser, user)
			r = r.WithContext(ctx)

			s.ctr.Logger.CtxInfoFf(ctx, "Позвольтевотель неидетинфицирован")

			next.ServeHTTP(w, r)
			return
		} else if sessId == "" && !forGuests {
			kernel.Error401(w)
			return
		}

		// authorized
		if !forAuthorized {
			kernel.Error403(w)
			return
		}

		ctx := r.Context()

		traceId := kernel_ext.TraceId(r)
		fromIp := kernel_ext.UserIp(r)
		fromUag := kernel_ext.UserAgent(r)

		// Внимание!
		// Может показаться, что логично было бы удалить куку в исключительных ситуациях
		// (например, при отсутствии связанного аккаунта) -- не стоит зависеть от предположений о работе клиента!
		// Пусть клиент сам разбирается с ситуацией, исходя из кода ответа.
		res, err := s.ctr.Services.Auth.CheckSession(ctx, traceId, fromIp, fromUag, sessId)
		switch vErr := err.(type) {
		case nil:
		case std.ErrorNotFound:
			s.ctr.Logger.CtxWarnFf(ctx, "связанный аккаунт не найден - %s: %q", vErr, sessId)
			kernel.Error401(w)
			return
		case auth.ErrorValidation:
			s.ctr.Logger.CtxWarnFf(ctx, "некорректный идентификатор сессии - %s: %q", vErr, sessId)
			kernel.Error400(w, "некорректный идентификатор сессии")
			return
		case auth.ErrorAccountAccessDenied:
			kernel.Error403(w)
			return
		case auth.ErrorSessionClosed:
			kernel.Error401(w)
			return
		case std.ErrorRuntime:
			s.ctr.Logger.CtxErrorFf(ctx, std.WrapErrorToRuntime(err, s, "allow", "CheckSession").Error())
			kernel.Error500(w)
			return
		default:
			panic(err)
		}

		user := newUserAuthorized(res.AccountId)
		ctx = context.WithValue(ctx, ctxKeyUser, user)

		// Имеет смысл поиск логов по аккаунту.
		// Ид сессии нельзя логировать в открытом виде, поэтому смысл его логировать теряется.
		ctx = s.ctr.Logger.AddExtraAttrToCtx(ctx, "account_id", res.AccountId)
		s.ctr.Logger.CtxInfoFf(
			ctx,
			"Позвольтевотель идетинфицирован. SessId: %q (exp: %s); AccId: %q",
			std.MaskStrNotFirstLast(sessId),
			res.SessionExpireAt.Format(time.RFC3339),
			res.AccountId,
		)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (s *Security) AllowOnlyGuests(next http.Handler) http.Handler {
	return s.allow(true, false, next)
}

func (s *Security) AllowOnlyAuthorized(next http.Handler) http.Handler {
	return s.allow(false, true, next)
}

func (s *Security) AllowAny(next http.Handler) http.Handler {
	return s.allow(true, true, next)
}

// ---------------------------------------------------------------------------------------------------------------------

// AssociatedUser
//
// Паникует при нулевых аргументах.
// Паникует, если обработчик запроса не был обернут в AllowOnlyGuests, AllowOnlyAuthorized или AllowAny.
func (s *Security) AssociatedUser(r *http.Request) *User {
	ctx := r.Context()

	vUser := ctx.Value(ctxKeyUser)
	if vUser == nil {
		panic(fmt.Errorf("%T.%s не нашел %s в контексте запроса", s, "AssociatedUser", ctxKeyUser))
	}

	user, ok := vUser.(*User)
	if !ok {
		panic(fmt.Errorf("%T.%s нашел %s, но это не %T: %#v", s, "AssociatedUser", ctxKeyUser, &User{}, vUser))
	}

	return user
}

// ---------------------------------------------------------------------------------------------------------------------

func (s *Security) AuthorizeClient(w http.ResponseWriter, sessId string, sessExpAt time.Time) {
	s.setSessId(w, sessId, sessExpAt)
}

func (s *Security) UnAuthorizeClient(w http.ResponseWriter) {
	s.setSessId(w, "-", time.Now().Add(-time.Hour))
}

// ---------------------------------------------------------------------------------------------------------------------
