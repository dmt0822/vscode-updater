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
	DebX64Url = "https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-x64"
)

func main() {
	res := download(DebX64Url)

	defer res.Body.Close()
	body := res.Body

	file := createNewFile("vscode_installer.deb")
	defer file.Close()

	copyToFile(file, body)

	cmd := newCmd("sudo", []string{"dpkg", "-i", "vscode_installer.deb"})

	if err := cmd.Run(); err != nil {
		cleanup()
		panic(err)
	}

	cleanup()
	log.Print("Update complete")
}

func download(url string) *http.Response {
	log.Print("Downloading vscode...")
	res, err := http.Get(DebX64Url)
	if err != nil {
		panic(err)
	}
	log.Print("Download finished")
	return res
}

func createNewFile(filename string) *os.File {
	log.Print("Creating installer file...")
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	return file
}

func copyToFile(file *os.File, data io.Reader) {
	log.Print("Copying to file...")
	bytesWritten, err := io.Copy(file, data)
	if err != nil {
		panic(err)
	}
	log.Print(fmt.Sprintf("Copy complete. %v written to disk.", bytesWritten))
}

func newCmd(executable string, args []string) *exec.Cmd {
	exe, err := exec.LookPath(executable)
	allArgs := []string{exe}
	allArgs = append(allArgs, args...)

	if err != nil {
		panic(err)
	}

	cmd := &exec.Cmd{
		Path:   exe,
		Args:   allArgs,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	return cmd
}

func cleanup() {
	log.Print("Removing installer...")
	os.Remove("vscode_installer.deb")
}
