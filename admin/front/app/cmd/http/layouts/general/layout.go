package general

import (
	"example/admin/front/cmd/http/kernel"
	"example/admin/front/cmd/http/layouts/general/handlers"
	"example/admin/front/internal/infra/clients/gateway"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"html/template"
	"net/http"
)

const pathToTemplate = "static/layouts/general/layout.html"
const urlSignOut = "/layouts/general/sign-out/"

var config = struct {
	appName              string
	menuNameUrl          map[string]string
	redirectUrlForGuests string
	isSet                bool
}{}

func Register(
	apiClient gateway.ApiClient,
	mux *http.ServeMux,
	appName string,
	redirectUrlForGuests string,
	menuNameUrl map[string]string,
) {
	assert.Str().NotEmpty().Must(appName)
	assert.Str().NotEmpty().Must(redirectUrlForGuests)

	mux.Handle("DELETE "+urlSignOut, handlers.NewSignOut(apiClient, redirectUrlForGuests))

	staticDir := "/layouts/general/"
	mux.Handle("GET "+staticDir, http.StripPrefix(staticDir, http.FileServer(http.Dir("./static"+staticDir))))

	config.appName = appName
	config.menuNameUrl = menuNameUrl
	config.redirectUrlForGuests = redirectUrlForGuests
	config.isSet = true
}

// ---------------------------------------------------------------------------------------------------------------------

func MakeView(
	apiClient gateway.ApiClient,
	pathToPageTemplate string,
	pathToPageTemplates ...string,
) View {
	assert.TrueMust(config.isSet)

	paths := make([]string, 0, len(pathToPageTemplates)+2)
	paths = append(paths, pathToTemplate)
	paths = append(paths, pathToPageTemplate)
	paths = append(paths, pathToPageTemplates...)

	tpl := template.Must(template.ParseFiles(paths...))

	return View{
		apiClient: apiClient,
		tpl:       tpl,
	}
}

type View struct {
	apiClient gateway.ApiClient
	tpl       *template.Template
}

func (v View) Render(w http.ResponseWriter, r *http.Request, pageData any) error {
	sessId := kernel.CookieGetSessId(r)
	if sessId == "" {
		kernel.Redirect307(w, r, config.redirectUrlForGuests)
		return nil
	}

	requestUrlPath := r.URL.Path

	fromIp := kernel.ClientIp(r)
	fromUag := kernel.ClientUag(r)

	// User info
	// ----------------------------------------------------------------

	uResp, err := v.apiClient.LayoutUserInfo(fromIp, fromUag, sessId)
	if err != nil {
		kernel.Error500(w)
		return nil
	} else if uResp.JSON422 != nil && uResp.JSON422.Code == 400 {
		kernel.Error400(w, uResp.JSON422.Message)
		return nil
	} else if uResp.JSON422 != nil && uResp.JSON422.Code == 401 {
		kernel.CookieUnsetSessId(w)
		kernel.Redirect307(w, r, config.redirectUrlForGuests)
		return nil
	} else if uResp.JSON422 != nil && uResp.JSON422.Code == 403 {
		kernel.Error403(w, uResp.JSON422.Message)
		return nil
	} else if uResp.JSON422 != nil && uResp.JSON422.Code == 404 {
		kernel.Error404(w, uResp.JSON422.Message)
		return nil
	} else if uResp.JSON422 != nil && uResp.JSON422.Code == 422 {
		kernel.Error422(w, uResp.JSON422.Message)
		return nil
	} else {
		assert.NotNilDeepMust(uResp.JSON200)
	}

	// todo : рефакторинг : дефолтный аватар-урл понадобится не только тут
	userPic := std.Ternary(*uResp.JSON200.AvatarUrl == "", "/static/favicon.ico", *uResp.JSON200.AvatarUrl)
	userName := *uResp.JSON200.Username

	// Menu
	// ----------------------------------------------------------------

	type MenuItem = struct {
		IsActive bool
		Name     string
		Url      string
	}

	menuItems := make([]MenuItem, 0)
	for name, url := range config.menuNameUrl {
		menuItems = append(menuItems, MenuItem{
			IsActive: requestUrlPath == url,
			Name:     name,
			Url:      url,
		})
	}

	// ----------------------------------------------------------------

	return v.tpl.Execute(w, struct {
		AppName    string
		MenuItems  []MenuItem
		UserPic    string
		UserName   string
		UrlSignOut string
		Page       any
	}{
		AppName:    config.appName,
		MenuItems:  menuItems,
		UserPic:    userPic,
		UserName:   userName,
		UrlSignOut: urlSignOut,
		Page:       pageData,
	})
}

// ---------------------------------------------------------------------------------------------------------------------
