package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/logrusorgru/aurora"
)

func walk(cwd string, hideClean bool) {
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
			inspect(cwd, hideClean)
			continue
		}
		walk(fullpath, hideClean)
	}
}

func inspect(dir string, hideClean bool) {
	status := getGitStatus(dir)
	if !hideClean || status.Reset().String() != "clean" {
		fmt.Printf("%s: %s\n", status.String(), dir)
	}
}

func getGitStatus(dir string) aurora.Value {
	cmd := exec.Command("git", "-C", dir, "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	if len(out) != 0 {
		return aurora.Red("dirty")
	}

	cmd = exec.Command("git", "-C", dir, "log", "--branches", "--not", "--remotes")
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	if len(out) != 0 {
		return aurora.Yellow("not-pushed")
	}

	return aurora.Green("clean")
}

func main() {
	hideCleanPtr := flag.Bool("hide-clean", false, "Hide clean repository")
	flag.Parse()

	walk(".", *hideCleanPtr)
}
