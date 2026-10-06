## 1. Что должно получиться

В конце у тебя будет проект:

```text
hello-gui/
├── .github/
│   └── workflows/
│       └── ci.yml
├── Dockerfile.test
├── logic.go
├── logic_test.go
├── main.go
├── go.mod
├── go.sum
└── .gitignore
```

GitHub Actions будет:

1. проверять форматирование;
2. запускать `go vet`;
3. запускать тесты;
4. проверять сборку;
5. при создании тега `v0.1.0`, `v0.2.0` и т. д. собирать приложение под **Linux, macOS и Windows**;
6. автоматически создавать GitHub Release и прикреплять 3 бинарника. CD_Go_GUI (1)

---

# 2. Что установить заранее

Если работаешь в **Windows + VS Code**, тебе понадобятся:

- VS Code;
- Git;
- Docker Desktop;
- аккаунт GitHub.

Go локально устанавливать необязательно, потому что тестирование и сборку можно выполнять через Docker.

Открой **VS Code → Terminal → New Terminal**.

Если у тебя Git Bash, используй команды из раздела Git Bash.

---

# 3. Создаём проект

В терминале:

```bash
cd ~
```

Создаём папку:

```bash
mkdir -p hello-gui/.github/workflows
cd hello-gui
```

Проверь:

```bash
pwd
```

Должно быть что-то вроде:

```text
/c/Users/ТВОЁ_ИМЯ/hello-gui
```

---

# 4. Создаём `go.mod`

Создай файл:

```text
go.mod
```

В него:

```go
module hello-gui

go 1.23
```

На этом этапе Fyne в `go.mod` ещё нет — это нормально. Она добавится после `go mod tidy`. CD_Go_GUI (1)

---

# 5. Создаём `logic.go`

Создай:

```text
logic.go
```

Содержимое:

```go
package main

import "fmt"

// Greeting возвращает приветствие для указанного имени.
func Greeting(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// SumRange считает сумму чисел от from до to включительно.
func SumRange(from, to int) int {
	sum := 0

	for i := from; i <= to; i++ {
		sum += i
	}

	return sum
}
```

Здесь находятся функции, которые потом будут тестироваться.

---

# 6. Создаём тесты

Создай:

```text
logic_test.go
```

Вставь:

```go
package main

import "testing"

func TestGreeting(t *testing.T) {
	got := Greeting("Fyne")
	want := "Hello, Fyne!"

	if got != want {
		t.Errorf("Greeting() = %q, want %q", got, want)
	}
}

func TestSumRange(t *testing.T) {
	got := SumRange(1, 10)
	want := 55

	if got != want {
		t.Errorf("SumRange(1, 10) = %d, want %d", got, want)
	}
}
```

Здесь проверяются две функции:

```text
Greeting()
SumRange()
```

---

# 7. Создаём GUI

Создай:

```text
main.go
```

Вставь:

```go
package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// version перезаписывается через -ldflags "-X main.version=..." во время сборки.
var version = "dev"

func main() {
	a := app.New()
	w := a.NewWindow("Hello GUI " + version)

	output := widget.NewLabel("Нажмите кнопку ниже")

	greetBtn := widget.NewButton("Поздороваться", func() {
		output.SetText(Greeting("GitHub"))
	})

	quitBtn := widget.NewButton("Выход", func() {
		a.Quit()
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Hello from Go GUI! 🎨🐹"),
		widget.NewSeparator(),
		greetBtn,
		output,
		widget.NewSeparator(),
		widget.NewLabel(fmt.Sprintf("Version: %s", version)),
		widget.NewLabel(fmt.Sprintf("Sum 1..10 = %d", SumRange(1, 10))),
		widget.NewSeparator(),
		quitBtn,
	))

	w.ShowAndRun()
}
```

Это основное GUI-приложение на Fyne. В нём есть кнопки, вывод версии и результат суммы `1..10`. CD_Go_GUI (1)

---

# 8. Создаём Dockerfile для тестов

Создай:

```text
Dockerfile.test
```

Вставь:

```dockerfile
# Тестовый образ для Go GUI (Fyne)

FROM golang:1.23

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        libgl1-mesa-dev \
        libegl1-mesa-dev \
        xorg-dev \
        libwayland-dev \
        libxkbcommon-dev \
        wayland-protocols && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
```

Здесь устанавливаются зависимости Fyne/GLFW для Linux: OpenGL, EGL, X11 и Wayland. CD_Go_GUI (1)

---

# 9. Собираем Docker-образ

В терминале:

```bash
docker build -f Dockerfile.test -t hello-gui-test .
```

Если всё нормально, в конце будет:

```text
FINISHED
naming to docker.io/library/hello-gui-test
```

Проверить наличие образа:

```bash
docker images
```

Должен появиться:

```text
hello-gui-test
```

---

# 10. Создаём `go.sum`

Это **очень важный этап**.

Выполни:

```bash
mkdir -p ~/.go-docker-cache
```

Затем:

```bash
docker run --rm \
  -u "$(id -u):$(id -g)" \
  -e HOME=/tmp \
  -e GOPATH=/tmp/go \
  -e GOCACHE=/tmp/go-cache \
  -v "$(pwd)":/app \
  -v ~/.go-docker-cache:/tmp/go \
  -w /app \
  hello-gui-test \
  go mod tidy
```

<img width="1922" height="1044" alt="image" src="https://github.com/user-attachments/assets/9c039bef-e7cb-42ff-b8b6-fc4e72df6d72" />

Команда скачает Fyne и остальные зависимости.

После этого:

```bash
ls
```

Должны появиться:

```text
go.mod
go.sum
```

А в `go.mod` появится зависимость:

```text
fyne.io/fyne/v2
```

`go.sum` обязательно нужно будет отправить на GitHub. CD_Go_GUI (1)

---

# 11. Проверяем `go.mod`

В VS Code открой `go.mod`.

Он уже будет выглядеть примерно так:

```go
module hello-gui

go 1.23

require fyne.io/fyne/v2 v2.8.1
```

Версия Fyne у тебя может отличаться — это нормально.

---

# 12. Запускаем тесты через Docker

Для Git Bash:

```bash
docker run --rm \
  -u "$(id -u):$(id -g)" \
  -e HOME=/tmp \
  -e GOPATH=/tmp/go \
  -e GOCACHE=/tmp/go-cache \
  -v "$(pwd)":/app \
  -v ~/.go-docker-cache:/tmp/go \
  -w /app \
  hello-gui-test \
  go test ./... -v
```

В конце должно быть:

```text
=== RUN   TestGreeting
--- PASS: TestGreeting

=== RUN   TestSumRange
--- PASS: TestSumRange

PASS
ok      hello-gui
```

Это означает, что тесты прошли. CD_Go_GUI (1)

---

# 13. Проверяем сборку

Теперь проверим, что GUI действительно собирается:

```bash
docker run --rm \
  -u "$(id -u):$(id -g)" \
  -e HOME=/tmp \
  -e GOPATH=/tmp/go \
  -e GOCACHE=/tmp/go-cache \
  -e CGO_ENABLED=1 \
  -v "$(pwd)":/app \
  -v ~/.go-docker-cache:/tmp/go \
  -w /app \
  hello-gui-test \
  sh -c "go build -o /tmp/hello-gui . && echo 'BUILD OK'"
```

В конце:

```text
BUILD OK
```

Если увидел это — приложение компилируется.

---

# 14. Создаём `.gitignore`

Создай:

```text
.gitignore
```

Вставь:

```gitignore
/hello-gui
/hello-gui-*
.env
.idea/
.vscode/
*.iml
```

Он нужен, чтобы случайные бинарники и настройки IDE не попали в Git. CD_Go_GUI (1)

---

# 15. Создаём GitHub Actions

Теперь самое главное — CI/CD.

Открой:

```text
.github/workflows/ci.yml
```

И вставь:

```yaml
name: Go GUI CI/CD

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'
          cache: true

      - name: Install GUI dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y \
            libgl1-mesa-dev \
            libegl1-mesa-dev \
            xorg-dev \
            libwayland-dev \
            libxkbcommon-dev \
            wayland-protocols

      - name: Format check
        run: |
          UNFORMATTED=$(gofmt -l .)

          if [ -n "$UNFORMATTED" ]; then
            echo "Не отформатировано:"
            echo "$UNFORMATTED"
            echo "Запустите локально: gofmt -w ."
            exit 1
          fi

      - name: Lint with go vet
        run: go vet ./...

      - name: Run tests
        run: go test ./... -v

      - name: Build
        run: go build -o /tmp/hello-gui .

  release:
    needs: test

    if: startsWith(github.ref, 'refs/tags/v')

    permissions:
      contents: write

    strategy:
      matrix:
        include:
          - os: ubuntu-latest
            artifact: hello-gui-linux-x64
            extra_ldflags: ""

          - os: macos-14
            artifact: hello-gui-macos-arm64
            extra_ldflags: ""

          - os: windows-latest
            artifact: hello-gui-windows-x64.exe
            extra_ldflags: "-H windowsgui"

    runs-on: ${{ matrix.os }}

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'
          cache: true

      - name: Install GUI dependencies
        if: runner.os == 'Linux'
        run: |
          sudo apt-get update
          sudo apt-get install -y \
            libgl1-mesa-dev \
            libegl1-mesa-dev \
            xorg-dev \
            libwayland-dev \
            libxkbcommon-dev \
            wayland-protocols

      - name: Setup MSYS2
        if: runner.os == 'Windows'
        uses: msys2/setup-msys2@v2
        with:
          msystem: MINGW64
          update: true
          install: >-
            mingw-w64-x86_64-gcc
            mingw-w64-x86_64-pkg-config

      - name: Set CC for Windows
        if: runner.os == 'Windows'
        shell: pwsh
        run: |
          $gccPath = Join-Path $env:RUNNER_TEMP "msys64\mingw64\bin\gcc.exe"

          if (-not (Test-Path $gccPath)) {
            $gccPath = "C:\msys64\mingw64\bin\gcc.exe"
          }

          if (-not (Test-Path $gccPath)) {
            $found = Get-ChildItem `
              -Path $env:RUNNER_TEMP `
              -Filter "gcc.exe" `
              -Recurse `
              -ErrorAction SilentlyContinue |
              Where-Object { $_.FullName -like "*mingw64*" } |
              Select-Object -First 1

            if ($found) {
              $gccPath = $found.FullName
            }
          }

          if (-not (Test-Path $gccPath)) {
            throw "gcc.exe not found"
          }

          Write-Host "Found gcc at: $gccPath"

          $gccUnix = $gccPath -replace '\\', '/'
          $gccDir = (Split-Path $gccPath -Parent) -replace '\\', '/'

          echo "CC=$gccUnix" >> $env:GITHUB_ENV
          echo "$gccDir" >> $env:GITHUB_PATH

      - name: Build GUI binary
        shell: bash
        env:
          CGO_ENABLED: 1
        run: |
          go build \
            -ldflags="-s -w ${{ matrix.extra_ldflags }} -X main.version=${{ github.ref_name }}" \
            -o ${{ matrix.artifact }} .

      - name: Upload to Release
        uses: softprops/action-gh-release@v2
        with:
          files: ${{ matrix.artifact }}
          generate_release_notes: true
```

В этом workflow принципиально важно, что GUI собирается на **трёх разных runner'ах**, поскольку Fyne использует CGO и нативные GUI-зависимости. CD_Go_GUI (1)

---

# 16. Проверяем проект

Выполни:

```bash
git status
```

Перед GitHub у тебя примерно:

```text
hello-gui/
├── .github/
│   └── workflows/
│       └── ci.yml
├── Dockerfile.test
├── logic.go
├── logic_test.go
├── main.go
├── go.mod
├── go.sum
└── .gitignore
```

Особенно проверь:

```bash
git status
```

`go.sum` должен присутствовать.

---

# 17. Создаём GitHub repository

На GitHub создай новый репозиторий:

```text
hello-gui
```
<img width="667" height="281" alt="image" src="https://github.com/user-attachments/assets/7fd27d90-aee1-4e75-a91d-aa8de2006336" />

**Важно:** репозиторий должен быть пустым.

Не ставь галочки:

```text
☐ Add a README file
☐ Add .gitignore
☐ Choose a license
```

В задании отдельно указано не создавать эти файлы на GitHub перед первым push. CD_Go_GUI (1)

---

# 18. Загружаем проект на GitHub

В терминале:

```bash
git init
```

Добавляем файлы:

```bash
git add .
```

Проверяем:

```bash
git status
```

Затем:

```bash
git commit -m "Initial commit: Go GUI with Fyne and CI/CD to Releases"
```
<img width="874" height="233" alt="image" src="https://github.com/user-attachments/assets/ff7f983b-6951-47ac-b9e6-ace30d3f56cb" />

Переименовываем ветку:

```bash
git branch -M main
```

Теперь добавляем GitHub:

```bash
git remote add origin https://github.com/ТВОЙ_USERNAME/hello-gui.git
```

Например:

```bash
git remote add origin https://github.com/SherKron/hello-gui.git
```

Проверить:

```bash
git remote -v
```
<img width="769" height="79" alt="image" src="https://github.com/user-attachments/assets/2f4772b5-fcf4-4f7f-ae22-0c23c3108437" />

Должно быть:

```text
origin  https://github.com/ТВОЙ_USERNAME/hello-gui.git
```

Отправляем:

```bash
git push -u origin main
```

---

# 19. Проверяем GitHub Actions

Открой репозиторий на GitHub.

Перейди:

```text
Actions
```

Там должен появиться:

```text
Go GUI CI/CD
```

Нажми на него.

Запустится:

```text
test
```

Внутри будут:

```text
✓ Setup Go
✓ Install GUI dependencies
✓ Format check
✓ Lint with go vet
✓ Run tests
✓ Build
```

При обычном `push` в `main` задача `release` запускаться не должна. Это нормально. CD_Go_GUI (1)

---

# 20. Если Actions зелёный

Если напротив workflow появилась:

```text
✓
```

значит CI работает.

Теперь можно делать первый Release.

---

# 21. Создаём тег `v0.1.0`

В терминале:

```bash
git tag v0.1.0
```

Проверить:

```bash
git tag
```

Должно быть:

```text
v0.1.0
```

Отправляем тег:

```bash
git push origin v0.1.0
```
<img width="695" height="226" alt="image" src="https://github.com/user-attachments/assets/392169de-d90d-405d-a48f-e3b3b547431a" />

---

# 22. Что произойдёт после тега

GitHub Actions снова запустит:

```text
test
```

После успешного `test` запустится:

```text
release
```

И GitHub создаст 3 сборки:

```text
Linux
macOS
Windows
```

Именно так устроена matrix-сборка в задании. CD_Go_GUI (1)

---

# 23. Проверяем Release

На GitHub открой:

```text
Releases
```

Там должен появиться:

```text
v0.1.0
```

Внутри:

```text
hello-gui-linux-x64
hello-gui-macos-arm64
hello-gui-windows-x64.exe
```
<img width="1230" height="1180" alt="image" src="https://github.com/user-attachments/assets/38553e44-4b27-414f-8604-74b68e04ebb6" />

Это и есть главный результат лабораторной. CD_Go_GUI (1)

---

# 24. Проверяем Windows-бинарник

Скачай:

```text
hello-gui-windows-x64.exe
```

Запусти.

<img width="945" height="886" alt="image" src="https://github.com/user-attachments/assets/e38fbd97-e9b3-43ac-87ae-983b2c58d7a5" />

Должно открыться GUI-приложение.

На Windows используется:

```text
-H windowsgui
```

поэтому консольное окно при запуске GUI скрывается. CD_Go_GUI (1)

---

# 25. Делаем версию `0.2.0`

Теперь по заданию нужно продемонстрировать обновление приложения.

Открой:

```text
main.go
```

Замени содержимое на новую версию из задания:

```go
package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// version перезаписывается через -ldflags "-X main.version=..." во время сборки.
var version = "dev"

func main() {
	a := app.New()
	w := a.NewWindow("Приветствие " + version)

	greeting := widget.NewLabel("Привет, мир! 👋")
	greeting.TextStyle = fyne.TextStyle{Bold: true}
	greeting.Alignment = fyne.TextAlignCenter

	versionLabel := widget.NewLabel(fmt.Sprintf("Version: %s", version))
	versionLabel.Alignment = fyne.TextAlignCenter

	closeBtn := widget.NewButton("Закрыть", func() {
		w.Close()
	})

	content := container.NewVBox(
		container.NewPadded(greeting),
		versionLabel,
		container.NewCenter(closeBtn),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(320, 200))
	w.CenterOnScreen()
	w.ShowAndRun()
}
```

Этот вариант предусмотрен в исходном задании для версии `0.2.0`. CD_Go_GUI (1)

---

# 26. Форматируем код

Перед отправкой:

```bash
docker run --rm \
  -u "$(id -u):$(id -g)" \
  -e HOME=/tmp \
  -e GOPATH=/tmp/go \
  -e GOCACHE=/tmp/go-cache \
  -v "$(pwd)":/app \
  -v ~/.go-docker-cache:/tmp/go \
  -w /app \
  hello-gui-test \
  sh -c "gofmt -w . && gofmt -l ."
```

Если команда ничего не вывела после выполнения — форматирование корректное. CD_Go_GUI (1)

---

# 27. Запускаем тесты ещё раз

```bash
docker run --rm \
  -u "$(id -u):$(id -g)" \
  -e HOME=/tmp \
  -e GOPATH=/tmp/go \
  -e GOCACHE=/tmp/go-cache \
  -v "$(pwd)":/app \
  -v ~/.go-docker-cache:/tmp/go \
  -w /app \
  hello-gui-test \
  go test ./... -v
```

В конце:

```text
PASS
```

---

# 28. Коммитим `0.2.0`

```bash
git add .
```

Затем:

```bash
git commit -m "feat: redesign UI with centered greeting version to 0.2.0"
```

Отправляем:

```bash
git push origin main
```
<img width="1176" height="1123" alt="image" src="https://github.com/user-attachments/assets/670fc72c-1008-4456-99d2-ba5b5a459917" />

После этого GitHub Actions снова запустит `test`.

---

# 29. Создаём тег `v0.2.0`

Когда `test` завершился успешно:

```bash
git tag v0.2.0
```

И:

```bash
git push origin v0.2.0
```

---

# 30. Проверяем второй Release

Снова:

```text
GitHub
→ Releases
```

Теперь должны быть:

```text
v0.1.0
v0.2.0
```

Старый `v0.1.0` остаётся без изменений, а `v0.2.0` содержит новые бинарники. CD_Go_GUI (1)

---
