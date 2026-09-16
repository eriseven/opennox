package localization

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
	"unicode/utf16"

	"github.com/opennox/libs/strman"
	"github.com/stretchr/testify/require"
)

func writeResource(t *testing.T, dir, name string, data []byte) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, data, 0600))
	return path
}

func TestDirectory(t *testing.T) {
	base := t.TempDir()
	want := filepath.Join(base, "locales", "zh-CN")
	require.NoError(t, os.MkdirAll(want, 0700))
	dir, err := Directory(base, "zh-CN")
	require.NoError(t, err)
	require.Equal(t, want, dir)
	dir, err = Directory(base, "")
	require.NoError(t, err)
	require.Empty(t, dir)
	for _, language := range []string{".", "..", "../zh-CN", `..\zh-CN`, "/zh-CN", `C:\zh-CN`, "missing"} {
		t.Run(language, func(t *testing.T) {
			_, err := Directory(base, language)
			require.Error(t, err)
		})
	}
	writeResource(t, filepath.Join(base, "locales"), "file", nil)
	_, err = Directory(base, "file")
	require.Error(t, err)
}

func TestResolve(t *testing.T) {
	base, locale := t.TempDir(), t.TempDir()
	baseFont := writeResource(t, base, "default.ttf", nil)
	localFont := writeResource(t, locale, "default.fnt", nil)
	number := writeResource(t, base, "number.fnt", nil)
	path, err := Resolve(base, locale, "default.ttf", "default.otf", "default.fnt")
	require.NoError(t, err)
	require.Equal(t, localFont, path, "locale directory takes precedence over font format")
	path, err = Resolve(base, "", "default.ttf", "default.otf", "default.fnt")
	require.NoError(t, err)
	require.Equal(t, baseFont, path, "no language uses original resources")
	path, err = Resolve(base, locale, "number.ttf", "number.otf", "number.fnt")
	require.NoError(t, err)
	require.Equal(t, number, path, "missing localized font uses original font")
	_, err = Resolve(base, locale, "missing.fnt")
	require.ErrorIs(t, err, os.ErrNotExist)
	require.NoError(t, os.Mkdir(filepath.Join(locale, "number.ttf"), 0700))
	_, err = Resolve(base, locale, "number.ttf", "number.fnt")
	require.Error(t, err, "an invalid override must not silently fall back")
}

func TestLoadStrings(t *testing.T) {
	base, locale := t.TempDir(), t.TempDir()
	writeResource(t, base, "nox.csf.json", []byte(`{"lang":0,"entries":[{"id":"menu:quit","vals":[{"str":"Quit"}]}]}`))
	writeResource(t, locale, "nox.csf", testCSF(t, 8, "menu:quit", "退出"))
	sm := strman.New()
	require.NoError(t, LoadStrings(sm, base, locale))
	require.Equal(t, "退出", sm.GetString("menu:quit"))
	require.Equal(t, 8, sm.Lang())

	// JSON has priority within the selected language directory.
	jsonPath := writeResource(t, locale, "nox.csf.json", []byte(`{"lang":8,"entries":[{"id":"menu:quit","vals":[{"str":"返回"}]}]}`))
	require.NoError(t, LoadStrings(sm, base, locale))
	require.Equal(t, "返回", sm.GetString("menu:quit"))
	require.NoError(t, os.Remove(jsonPath))
	require.NoError(t, os.Remove(filepath.Join(locale, "nox.csf")))
	require.NoError(t, LoadStrings(sm, base, locale))
	require.Equal(t, "Quit", sm.GetString("menu:quit"))
	require.Equal(t, 0, sm.Lang())

	writeResource(t, locale, "nox.csf", []byte("broken CSF"))
	require.Error(t, LoadStrings(sm, base, locale), "corrupt translations must be reported")
	writeResource(t, locale, "nox.csf.json", []byte("broken JSON"))
	require.Error(t, LoadStrings(sm, base, locale))
}

// A synthetic CSF fixture avoids depending on copyrighted game resources.
func testCSF(t *testing.T, language uint32, id, value string) []byte {
	t.Helper()
	var buf bytes.Buffer
	buf.WriteString(" FSC")
	require.NoError(t, binary.Write(&buf, binary.LittleEndian, []uint32{3, 1, 1, 0, language}))
	buf.WriteString(" LBL")
	require.NoError(t, binary.Write(&buf, binary.LittleEndian, []uint32{1, uint32(len(id))}))
	buf.WriteString(id)
	buf.WriteString(" RTS")
	chars := utf16.Encode([]rune(value))
	require.NoError(t, binary.Write(&buf, binary.LittleEndian, uint32(len(chars))))
	for _, ch := range chars {
		require.NoError(t, binary.Write(&buf, binary.LittleEndian, ^ch))
	}
	return buf.Bytes()
}
