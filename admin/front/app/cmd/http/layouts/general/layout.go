package general

import (
	"example/admin/front/cmd/http/layouts/general/handlers"
	"example/admin/front/internal/infra/clients/gateway"
	assert "github.com/selyukovn/go-wm-assert"
	"html/template"
	"net/http"
)

const pathToTemplate = "static/layouts/general/layout.html"
const urlSignOut = "/layouts/general/sign-out/"

var config = struct {
	appName     string
	menuNameUrl map[string]string
	isSet       bool
}{}

func Register(
	apiClient *gateway.ApiClient,
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
	config.isSet = true
}

// ---------------------------------------------------------------------------------------------------------------------

func MakeView(
	pathToPageTemplate string,
	pathToPageTemplates ...string,
) *View {
	assert.TrueMust(config.isSet)

	paths := make([]string, 0, len(pathToPageTemplates)+2)
	paths = append(paths, pathToTemplate)
	paths = append(paths, pathToPageTemplate)
	paths = append(paths, pathToPageTemplates...)

	tpl := template.Must(template.ParseFiles(paths...))

	return &View{
		tpl: tpl,
	}
}

type View struct {
	tpl *template.Template
}

func (v *View) Render(w http.ResponseWriter, requestUrlPath string, pageData any) error {
	type MenuItem = struct {
		IsActive bool
		Name     string
		Url      string
	}

	// --

	// todo : брать из апи
	userPic := "/static/favicon.ico"
	userName := "username@example.com"

	menuItems := make([]MenuItem, 0)
	for name, url := range config.menuNameUrl {
		menuItems = append(menuItems, MenuItem{
			IsActive: requestUrlPath == url,
			Name:     name,
			Url:      url,
		})
	}

	// --

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
