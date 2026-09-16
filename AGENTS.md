# Repository Guidelines

## Project Structure & Module Organization

OpenNox extends the Nox engine using Go and legacy C. The Go module lives in `src/`, not the repository root.

- `src/client/`, `src/server/`, and `src/common/`: client, server, and shared functionality.
- `src/legacy/`: legacy implementation and C interoperability.
- `src/cmd/`: executable entry points; `src/internal/`: build, refactoring, and support tools.
- Tests sit beside implementation files as `*_test.go`.
- `docs/`: build and refactoring guides; `res/`: project icons; `app/`, `docker/`, `snap/`, and `scripts/`: packaging and utilities.

Original game assets are not included; use a legally owned Nox installation for runtime testing.

## Build, Test, and Development Commands

Use Go 1.25 or newer, as specified in `src/go.mod`, plus a 32-bit C toolchain, SDL2, and OpenAL. Follow `docs/build-linux.md` or `docs/build-windows.md`; Windows builds use MSYS2 MINGW32.

From `src/`:

- `go run ./internal/noxbuild`: build default client, HD client, and server targets, plus Windows debug variants.
- `go run ./internal/noxbuild client-hd`: build only the HD client.
- `go run ./internal/noxbuild server`: build the dedicated server.
- `./opennox-hd` (Windows: `./opennox-hd.exe`): launch the built client with game data configured and required runtime libraries installed.

For the full test suite, use Bash/MSYS2 from `src/`:

```bash
GOARCH=386 CGO_ENABLED=1 \
CGO_CFLAGS_ALLOW='(-fshort-wchar)|(-fno-strict-aliasing)|(-fno-strict-overflow)' \
go test ./...
```

## Coding Style & Naming Conventions

Implement new features in Go and format changes with `gofmt`. Keep legacy refactoring in plain C; avoid new C dependencies. Follow `.editorconfig`: UTF-8, LF, and tabs for C with width 4. C formatting uses `.clang-format` with a 120-column limit; run `make format` from the root and inspect the diff.

Preserve legacy procedure names using `<module>_<component>_<action>_<original_address>`. Read `docs/refactoring.md` before changing legacy types or logic. Preserve Nox 1.2b file and network compatibility.

## Testing Guidelines

Use Go's `testing` framework and existing assertion helpers. Name tests `TestXxx` in `*_test.go`; add regression tests for behavior changes. No numeric coverage threshold is configured. Build and smoke-test affected gameplay; some tests need external game data.

## Commit & Pull Request Guidelines

Use short, imperative subjects matching history, such as `Fix potential buffer overflow.` Separate renames, type changes, and logic restoration into reviewable commits. Discuss significant changes in an issue first. Explain what changed and why, link related issues, and report validation. Complete both required sign-offs in `.github/pull_request_template.md`. Read `CONTRIBUTING.md` before contributing.
