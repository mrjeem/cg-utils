package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/go-git/go-git/v5"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func newInstallerUI() *widgets.QMainWindow {
	dst := ""

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowFlags(core.Qt__FramelessWindowHint)
	window.SetWindowTitle("Installer")
	window.Resize2(750, 500)

	// FIXME: Relative path
	stylesheet, err := os.ReadFile("/home/mrjeem/dev/installer/stylesheet.qss")
	if err != nil {
		log.Printf("[Warning] Could not read style sheet...")
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

func install(dst string) {
	dst = filepath.Join(dst, "installer")
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		os.MkdirAll(dst, 0775)

	}

	go func() {
		platform := runtime.GOOS
		// FIXME: Reverse the if
		if platform != "windows" {
			log.Printf("[Info] Installing Python...")
			installPython(dst)
		}

		log.Printf("[Info] Installing Rez...")
		installRez(dst)
	}()
}

func installPython(dst string) {
	url := "https://www.python.org/ftp/python/3.12.5/python-3.12.5.exe"

	log.Printf("[Info] Downloading python from %s to %s", url, dst)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("\t[Fatal] Failed to download the python installer. Status code: %d", resp.StatusCode)
	}

	filename := "python-3.12.5.exe"
	pathname := filepath.Join(dst, filename)
	if _, err := os.Stat(pathname); err == nil {
		log.Printf("\t[Info] Python install exe already exists, removing first...")
		os.RemoveAll(pathname)
	}
	file, err := os.Create(pathname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[Info] Python downloaded successfully...")

}

func installRez(dst string) {

	url := "https://github.com/AcademySoftwareFoundation/rez.git"
	rez_repo_dir := filepath.Join(dst, "rez_src")
	log.Printf("[Info] git clone %s %s", url, rez_repo_dir)

	if _, err := os.Stat(rez_repo_dir); err == nil {
		log.Printf("[Info] Rez repo already exists, removing first...")
		os.RemoveAll(rez_repo_dir)
	}
	_, err := git.PlainClone(rez_repo_dir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Fatal(err)
	} else {

		log.Printf("\t[Info] Repo cloned successfully...")
	}

	setupPath := filepath.Join(dst, "rez_src", "install.py")
	installDir := filepath.Join(dst, "rez")
	// FIXME: Proper command for windows
	installCmd := fmt.Sprintf("python3 %s -v %s", setupPath, installDir)
	log.Printf("[Info] %s", installCmd)
	cmd := exec.Command("sh", "-c", installCmd)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("[Fatal] Error executing command: %s\n", err)
		log.Printf("[Fatal] Output: %s\n", string(output))
	} else {

		log.Printf("[Info] Rez installed successfully...")
	}

}

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// FIXME: Relative path
	fontPath := filepath.Join("/home/mrjeem/dev/installer", "fonts", "Anek_Devanagari")
	fontDb := gui.NewQFontDatabase()
	fontDb.AddApplicationFont(fontPath)

	font := gui.NewQFont2("Anek Devanagari", 10, 1, false)
	font.SetStyleName("Regular")
	app.SetFont(font, "")

	var ui = newInstallerUI()
	ui.Show()
	app.Exec()

}
