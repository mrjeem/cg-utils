package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// --------------------------------

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "CG Utils: ", log.Lshortfile)
)

var current_dir string

// --------------------------------

func newInstallerUI() *widgets.QMainWindow {

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowFlags(core.Qt__FramelessWindowHint)
	window.SetWindowTitle("Installer")
	window.Resize2(850, 550)

	stylesheet, err := os.ReadFile(filepath.Join(current_dir, "resources", "stylesheet.qss"))
	if err != nil {
		logger.Printf("[Warning] Could not read style sheet...")
	} else {
		window.SetStyleSheet(string(stylesheet))
	}

	mainWidget := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(mainWidget)

	mainLayout := widgets.NewQVBoxLayout()
	mainWidget.SetLayout(mainLayout)

	topBarLayout := widgets.NewQHBoxLayout()
	topBarLayout.AddStretch(1)
	mainLayout.AddLayout(topBarLayout, 0)

	closeBtn := widgets.NewQPushButton2("x", nil)
	closeBtn.SetStyleSheet("background-color: none; min-height: 40px; min-width: 40px;")
	closeBtn.ConnectClicked(func(bool) {
		window.Close()
	})
	topBarLayout.AddWidget(closeBtn, 0, 0)

	mainLayout.AddStretch(1)

	welcomeLabel := widgets.NewQLabel2("Welcome to CG Utils!", nil, 0)
	welcomeLabel.SetStyleSheet("font-size: 32px; min-height: 50px;")
	mainLayout.AddWidget(welcomeLabel, 0, core.Qt__AlignCenter)

	mainLayout.AddSpacing(48)

	dstLabel := widgets.NewQLabel2("Installation Directory:", nil, 0)
	dstLabel.SetStyleSheet("margin-bottom: 0px; margin-left: 48px;")
	mainLayout.AddWidget(dstLabel, 0, core.Qt__AlignLeft)

	dstHLayout := widgets.NewQHBoxLayout()
	mainLayout.AddLayout(dstHLayout, 1)

	dstInput := widgets.NewQLineEdit(nil)
	dstInput.SetStyleSheet("margin-top: 0px; margin-left: 48px;")
	dstHLayout.AddWidget(dstInput, 1, 0)

	var dst string
	browseBtn := widgets.NewQPushButton2("Browse", nil)
	browseBtn.SetStyleSheet("margin-top: 0px; margin-right: 48px;")
	browseBtn.ConnectClicked(func(bool) {
		browse(mainWidget, &dst)
		dstInput.SetText(dst)
	})

	dstHLayout.AddWidget(browseBtn, 0, 0)

	mainLayout.AddSpacing(48)
	mainLayout.AddStretch(1)

	bottomBtnLayout := widgets.NewQHBoxLayout()
	bottomBtnLayout.AddStretch(1)
	mainLayout.AddLayout(bottomBtnLayout, 0)

	cancelBtn := widgets.NewQPushButton2("Cancel", nil)
	cancelBtn.ConnectClicked(func(bool) { window.Close() })
	bottomBtnLayout.AddWidget(cancelBtn, 0, 0)

	installBtn := widgets.NewQPushButton2("Install", nil)
	installBtn.ConnectClicked(func(bool) { install(dstInput.Text()) })
	bottomBtnLayout.AddWidget(installBtn, 0, 0)

	bottomBtnLayout.AddStretch(1)

	mainLayout.AddStretch(1)
	return window

}

func browse(parent widgets.QWidget_ITF, dst *string) {

	dialog := widgets.NewQFileDialog(parent, 0)
	dialog.SetFileMode(widgets.QFileDialog__Directory)
	selectedDir := dialog.GetExistingDirectory(parent, "Select Instalation Directory...", "~", widgets.QFileDialog__ShowDirsOnly)
	*dst = selectedDir
}

func main() {
	logger.SetPrefix("[ui.go] ")
	logger.SetFlags(0)

	app := widgets.NewQApplication(len(os.Args), os.Args)

	current_dir, err := filepath.Abs("./")
	if err != nil {
		logger.Fatal(err)
	}

	fontPath := filepath.Join(current_dir, "resources", "fonts", "Anek_Devanagari")
	fontDb := gui.NewQFontDatabase()
	fontDb.AddApplicationFont(fontPath)

	font := gui.NewQFont2("Anek Devanagari", 10, 1, false)
	font.SetStyleName("Regular")
	app.SetFont(font, "")

	var ui = newInstallerUI()
	ui.Show()
	app.Exec()

}
