# Localized text and fonts

OpenNox can use a language pack without duplicating the original game data.
Place the additional files in `locales/<language>/` inside your Nox installation:

```text
Nox/
  nox.csf
  default.fnt
  large.fnt
  small.fnt
  number.fnt
  maps/
  video.bag
  video.idx
  ...
  locales/
    zh-CN/
      nox.csf
      default.fnt
      large.fnt
      small.fnt
```

In an MSYS2 MINGW32 terminal, start the client with:

```bash
./opennox-hd.exe -data "E:/GOG_Galaxy/Games/Nox" -lang zh-CN -window
```

Omit `-lang` to use the original installation's language. The value is an existing
folder name under `locales/`, not a path or a CSF language number. Other names,
such as `fr` or `zh-TW`, work the same way. Restart the game to change languages;
the command-line selection is not saved in configuration.

## Resource lookup and fallback

- Text: look for `nox.csf.json`, then `nox.csf`, in the selected language folder.
  If neither exists, load the original string table using the same priority.
- Fonts: for each required font, try `.ttf`, `.otf`, then `.fnt` in the language
  folder before trying those formats in the original data directory. For example,
  a localized `default.fnt` overrides an original `default.ttf`.
- Missing files fall back individually. A pack with only the four files shown
  above uses the original `number` font. String tables are selected as whole files;
  missing translation entries are not merged from the original table.
- A missing language directory, unreadable file, or malformed selected resource
  is an error. Such errors are not silently replaced with another language.

The language ID inside the string table still controls the existing font and
keyboard adjustments. Use text and fonts from a compatible language pack.
Maps, images, interface layouts, music, and voice continue to use the original
data directory. This option changes text and fonts; it does not translate content
or add input-method support. No rebuild is needed when replacing language files.

Language packs and original game assets are user-provided data; do not commit
them to this repository.
