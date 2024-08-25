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

var resources_dir, _ = filepath.Abs("./resources/")

// --------------------------------

func newInstallerUI() *widgets.QMainWindow {

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowFlags(core.Qt__FramelessWindowHint)
	window.SetWindowTitle("Installer")
	window.Resize2(850, 550)

	stylesheet, err := os.ReadFile(filepath.Join(resources_dir, "stylesheet.qss"))
	if err != nil {
		logger.Printf("[Warning] Could not read style sheet...")
	} else {
		window.SetStyleSheet(string(stylesheet))
	}

	widget := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(widget)

	mainLayout := widgets.NewQVBoxLayout()
	widget.SetLayout(mainLayout)

	// Top Bar
	topBarLayout := widgets.NewQHBoxLayout()
	topBarLayout.AddStretch(1)
	mainLayout.AddLayout(topBarLayout, 0)

	// Close button
	closeBtn := widgets.NewQPushButton2("x", nil)
	closeBtn.SetStyleSheet("background-color: none; min-height: 40px; min-width: 40px;")
	closeBtn.ConnectClicked(func(bool) {
		window.Close()
	})
	topBarLayout.AddWidget(closeBtn, 0, 0)

	// Cancel Button
	cancelBtn := widgets.NewQPushButton2("Cancel", nil)
	cancelBtn.ConnectClicked(func(bool) { window.Close() })

	// Welcome screen
	welcomeWidget(*mainLayout, cancelBtn)

	return window
}

func welcomeWidget(mainLayout widgets.QVBoxLayout, cancelBtn widgets.QPushButton_ITF) {
	// This widget
	widget := widgets.NewQWidget(nil, 0)
	mainLayout.AddWidget(widget, 0, 0)

	// This widget's layout
	layout := widgets.NewQVBoxLayout()
	widget.SetLayout(layout)

	// Welcome
	label := widgets.NewQLabel2("Welcome to CG Utils!", nil, 0)
	label.SetStyleSheet("font-size: 40px; min-height: 56px; margin-top: 64px;")
	layout.AddWidget(label, 0, core.Qt__AlignCenter)

	// Bottom bar
	bottomBtnLayout := widgets.NewQHBoxLayout()
	bottomBtnLayout.AddStretch(1)
	layout.AddLayout(bottomBtnLayout, 0)

	bottomBtnLayout.AddWidget(cancelBtn, 0, core.Qt__AlignCenter)

	// Next button
	next_btn := widgets.NewQPushButton2("Next", nil)
	next_btn.ConnectClicked(func(bool) { widget.Hide(); baseInstallationWidget(mainLayout, cancelBtn) })
	bottomBtnLayout.AddWidget(next_btn, 0, core.Qt__AlignCenter)
	bottomBtnLayout.AddStretch(1)

}

func baseInstallationWidget(mainLayout widgets.QVBoxLayout, cancelBtn widgets.QPushButton_ITF) {

	widget := widgets.NewQWidget(nil, 0)

	layout := widgets.NewQVBoxLayout()
	widget.SetLayout(layout)

	baseDirLabel := widgets.NewQLabel2("Install Base Application", nil, 0)
	baseDirLabel.SetStyleSheet("font-size: 32px; min-height: 56px; margin-bottom: 0px;")
	layout.AddWidget(baseDirLabel, 0, core.Qt__AlignCenter)

	info := widgets.NewQLabel2("This directory will be used to install the launcher and all the tools you select next/later.", nil, 0)
	info.SetStyleSheet("margin-top: 0px; color: #afafaf;")
	layout.AddWidget(info, 0, core.Qt__AlignCenter)

	layout.AddStretch(1)

	// Destination
	dstLabel := widgets.NewQLabel2("Installation Directory:", nil, 0)
	dstLabel.SetStyleSheet("margin-bottom: 0px; margin-left: 48px;")
	layout.AddWidget(dstLabel, 0, core.Qt__AlignLeft)

	dstHLayout := widgets.NewQHBoxLayout()
	layout.AddLayout(dstHLayout, 1)

	layout.AddStretch(1)

	// Input Line Edit
	dstInput := widgets.NewQLineEdit(nil)
	dstInput.SetStyleSheet("margin-top: 0px; margin-left: 48px;")
	dstHLayout.AddWidget(dstInput, 1, 0)

	var dst string

	// Browse Button
	browseBtn := widgets.NewQPushButton2("Browse", nil)
	browseBtn.SetStyleSheet("margin-top: 0px; margin-right: 48px;")
	browseBtn.ConnectClicked(func(bool) {
		browse(widget, &dst)
		dstInput.SetText(dst)
	})

	dstHLayout.AddWidget(browseBtn, 0, 0)

	layout.AddStretch(1)

	// Bottom bar
	bottomBtnLayout := widgets.NewQHBoxLayout()
	bottomBtnLayout.AddStretch(1)
	layout.AddLayout(bottomBtnLayout, 0)

	// Cancel Button
	bottomBtnLayout.AddWidget(cancelBtn, 0, core.Qt__AlignCenter)

	// Back button
	backBtn := widgets.NewQPushButton2("Back", nil)
	backBtn.ConnectClicked(func(bool) { widget.Hide(); welcomeWidget(mainLayout, cancelBtn) })
	bottomBtnLayout.AddWidget(backBtn, 0, core.Qt__AlignCenter)

	// Install Button
	installBtn := widgets.NewQPushButton2("Install", nil)
	installBtn.ConnectClicked(func(bool) { install(dstInput.Text()) })
	bottomBtnLayout.AddWidget(installBtn, 0, core.Qt__AlignCenter)

	bottomBtnLayout.AddStretch(1)

	layout.AddStretch(1)

	mainLayout.AddWidget(widget, 0, 0)

}

func browse(parent widgets.QWidget_ITF, dst *string) {

	dialog := widgets.NewQFileDialog(parent, 0)
	dialog.SetFileMode(widgets.QFileDialog__Directory)
	selectedDir := dialog.GetExistingDirectory(parent, "Select Installation Directory...", "~", widgets.QFileDialog__ShowDirsOnly)
	*dst = selectedDir
}

func main() {
	logger.SetPrefix("[ui.go] ")
	logger.SetFlags(0)

	app := widgets.NewQApplication(len(os.Args), os.Args)

	// Adding font to db
	fontDb := gui.NewQFontDatabase()
	fontDir := filepath.Join(resources_dir, "fonts", "Anek_Devanagari")
	fontFiles, _ := filepath.Glob(filepath.Join(fontDir, "*.ttf"))
	for _, file := range fontFiles {
		fontDb.AddApplicationFont(file)
	}

	// Applying font
	font := gui.NewQFont2("Anek Devanagari", 10, 1, false)
	font.SetStyleName("Regular")
	app.SetFont(font, "")

	var ui = newInstallerUI()
	ui.Show()
	app.Exec()

}
