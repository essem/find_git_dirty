package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

func walk(cwd string) {
	files, err := ioutil.ReadDir(cwd)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fullpath := path.Join(cwd, file.Name())
		fi, err := os.Lstat(fullpath)
		if err != nil {
			log.Fatal(err)
		}
		if !fi.IsDir() {
			continue
		}
		if file.Name() == ".git" {
			inspect(cwd)
			continue
		}
		walk(fullpath)
	}
}

func inspect(dir string) {
	fmt.Print(dir + ": ")

	cmd := exec.Command("git", "-C", dir, "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	if len(out) != 0 {
		fmt.Println("dirty")
		return
	}

	cmd = exec.Command("git", "-C", dir, "log", "--branches", "--not", "--remotes")
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	if len(out) != 0 {
		fmt.Println("not-pushed")
		return
	}

	fmt.Println("clean")
}

func main() {
	walk(".")
}
