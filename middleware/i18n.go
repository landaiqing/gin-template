package middleware

import (
	"github.com/BurntSushi/toml"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func I18n() gin.HandlerFunc {
	return ginI18n.Localize(
		ginI18n.WithBundle(&ginI18n.BundleCfg{
			RootPath:         "i18n/language/",
			AcceptLanguage:   []language.Tag{language.Chinese, language.English},
			DefaultLanguage:  language.Chinese,
			UnmarshalFunc:    toml.Unmarshal,
			FormatBundleFile: "toml",
		}),
		ginI18n.WithGetLngHandle(
			func(context *gin.Context, defaultLng string) string {
				lang := context.Query("lang")
				if lang == "" {
					return defaultLng
				}
				return lang
			},
		),
	)
}
