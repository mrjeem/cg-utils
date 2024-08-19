package main

import (
	"os"
	"runtime"
)
import (
	"github.com/therecipe/qt/widgets"
)

func InstallerUI() *widgets.QMainWindow {

	var window = widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Installer")
	window.Resize2(750, 500)

	main_widget := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(main_widget)

	main_layout := widgets.NewQVBoxLayout()
	main_layout.AddStretch(1)
	main_widget.SetLayout(main_layout)

	dst_h_layout := widgets.NewQHBoxLayout()
	main_layout.AddLayout(dst_h_layout, 1)

	dst_label := widgets.NewQLabel2("Installation Directory", nil, 0)
	dst_h_layout.AddWidget(dst_label, 0, 0)

	dst_input := widgets.NewQLineEdit(nil)
	dst_h_layout.AddWidget(dst_input, 1, 0)

	browse_btn := widgets.NewQPushButton2("Browse", nil)
	browse_btn.ConnectClicked(func(bool) { _browse(main_widget, *dst_input) })
	dst_h_layout.AddWidget(browse_btn, 0, 0)

	main_layout.AddStretch(1)

	bottom_btn_layout := widgets.NewQHBoxLayout()
	bottom_btn_layout.AddStretch(1)
	main_layout.AddLayout(bottom_btn_layout, 0)

	cancel_btn := widgets.NewQPushButton2("Cancel", nil)
	cancel_btn.ConnectClicked(func(bool) { window.Close() })
	bottom_btn_layout.AddWidget(cancel_btn, 0, 0)

	install_btn := widgets.NewQPushButton2("Install", nil)
	install_btn.ConnectClicked(func(bool) { _install(dst_input.Text()) })
	bottom_btn_layout.AddWidget(install_btn, 0, 0)

	return window

}

func _browse(parent widgets.QWidget_ITF, line_edit widgets.QLineEdit) {

	dialog := widgets.NewQFileDialog(parent, 0)
	dialog.SetFileMode(widgets.QFileDialog__Directory)
	selected_dir := dialog.GetExistingDirectory(parent, "Select Instalation Directory...", "~", widgets.QFileDialog__ShowDirsOnly)
	line_edit.SetText(selected_dir)
}

func _install(dst string) {
	platform := runtime.GOOS
	// TODO: Reverse the if
	if platform != "windows" {
		_install_python(dst)
	}
}

func _install_python(dst string) {
	println("Instaling Python")
	println(dst)
}

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	var ui = InstallerUI()
	ui.Show()
	app.Exec()

}
