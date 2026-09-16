package noxrender

import (
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/opennox/libs/noxfont"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/font/gofont/goregular"
)

func TestLoadLocaleFont(t *testing.T) {
	base, locale := t.TempDir(), t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(base, "default.ttf"), goregular.TTF, 0600))
	fnt := &noxfont.Font{Ranges: []noxfont.Range{{
		StartChar: '中', EndChar: '中',
		Glyphs: []*noxfont.Bitmap{{Pix: []byte{0x80}, Stride: 1, Rect: image.Rect(0, 0, 1, 1)}},
	}}}
	data, err := fnt.Encode()
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(locale, "default.fnt"), data, 0600))
	face, err := loadLocaleFont(base, locale, "default", 10)
	require.NoError(t, err)
	require.IsType(t, &noxfont.Font{}, face, "localized bitmap font overrides the base TrueType font")
	require.NotNil(t, face.(*noxfont.Font).Char16('中'))
	require.NoError(t, face.Close())

	// The supplied translation need not duplicate the original numeric font.
	require.NoError(t, os.WriteFile(filepath.Join(base, "number.fnt"), data, 0600))
	face, err = loadLocaleFont(base, locale, "number", 7)
	require.NoError(t, err)
	require.NoError(t, face.Close())

	face, err = loadLocaleFont(base, "", "default", 10)
	require.NoError(t, err)
	_, ok := face.(*noxfont.Font)
	require.False(t, ok, "no locale keeps the original font")
	require.NoError(t, face.Close())

	// A malformed higher-priority font must not silently load a different one.
	require.NoError(t, os.WriteFile(filepath.Join(locale, "default.ttf"), []byte("broken font"), 0600))
	_, err = loadLocaleFont(base, locale, "default", 10)
	require.Error(t, err)
}
