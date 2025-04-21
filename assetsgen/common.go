package assetsgen

import (
	"errors"
)

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
)

type platformType string

const (
	platformTypeAndroid platformType = "android"
	platformTypeIos     platformType = "ios"
)

type intention string

const (
	intentAppIcon          intention = "app_icon"
	intentNotificationIcon intention = "notification_icon"
	intentAsset            intention = "asset"
)

type Asset interface {
	Name() string
	CalcSize(w, h int) (int, int)
	CalcPadding(w, h int) int
}
