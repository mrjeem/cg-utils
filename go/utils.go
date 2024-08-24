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
)

func install(dst string) {
	dst = filepath.Join(dst, "installer")
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		os.MkdirAll(dst, 0775)

	}

	go func() {
		platform := runtime.GOOS
		if platform == "windows" {
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

	// FIXME: You gotta install after downloading

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
	// FIXME: For windows, use downloaded python
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
