package kernel_ext

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"
)

func CalcFilesVersion(pathToDir string) string {
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
