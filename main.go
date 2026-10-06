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
