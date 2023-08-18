package authn

import (
	"testing"

	"github.com/alainQtec/FileServer/v3/pkg/cache"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	asserts := assert.New(t)
	cache.Set("setting_siteURL", "http://FileServer.org", 0)
	cache.Set("setting_siteName", "FileServer", 0)
	res, err := NewAuthnInstance()
	asserts.NotNil(res)
	asserts.NoError(err)
}
