package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const (
	WinX64Url = "https://go.microsoft.com/fwlink/?Linkid=852157"
	DebX64Url = "https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-x64"
)

func main() {
	log.Print("Downloading vscode...")

	res, err := http.Get(DebX64Url)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Download finished")

	defer res.Body.Close()
	body := res.Body

	log.Print("Creating installer file...")

	file, err := os.Create("vscode_installer.deb")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	bytesWritten, err := io.Copy(file, body)

	if err != nil {
		log.Fatal(err)
	}

	log.Print(fmt.Sprintf("Installer file created. %v written to disk.", bytesWritten))

	sudo, err := exec.LookPath("sudo")

	if err != nil {
		log.Fatal(err)
	}

	cmd := &exec.Cmd{
		Path:   sudo,
		Args:   []string{sudo, "dpkg", "-i", "vscode_installer.deb"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	log.Print("You're all set and ready to code!")
}
