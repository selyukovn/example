package static

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// RegisterFileHandler
//
// Возвращает префикс пути и префикс урла к статике бандла.
func RegisterFileHandler(mux *http.ServeMux, bundleName string) (string, string) {
	// см. bff/build/http/Dockerfile
	dirPath := "./static/" + bundleName
	// см. rproxy/build/server/nginx.conf
	urlPrefix := "/static/" + bundleName + "/" + calcFilesVersion(dirPath)

	// `rproxy` кеширует файлы, полученные из `bff` -- cм. README.md.
	mux.Handle(
		// GET разрешает и HEAD
		"GET "+urlPrefix+"/",
		http.StripPrefix(urlPrefix, http.FileServer(http.Dir(dirPath))),
	)

	return dirPath, urlPrefix
}

func calcFilesVersion(pathToDir string) string {
	maxModTime := time.Time{}

	if err := filepath.Walk(pathToDir, func(oath string, info os.FileInfo, err error) error {
		// никаких ошибок (права, например) быть не может -- иначе что-то очень не так.
		if err != nil {
			panic(err)
		}

		// проход уже рекурсивный, поэтому папки пропускаем
		if info.IsDir() {
			return nil
		}

		if info.ModTime().After(maxModTime) {
			maxModTime = info.ModTime()
		}

		return nil
	}); err != nil {
		panic(err)
	}

	hasher := md5.New()
	hasher.Write([]byte(maxModTime.String()))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
