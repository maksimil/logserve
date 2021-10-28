package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

const HTMLTEMPLATE = `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8" /><meta http-equiv="X-UA-Compatible" content="IE=edge" /><meta name="viewport" content="width=device-width, initial-scale=1.0" /><title>Logserve</title></head><body></body><script>%s</script><style>%s</style></html>`

const (
	MAINCSS   = "web/main.css"
	MAINJS    = "web/main.ts"
	BUILDHTML = "cmd/_gen/build.html"
)

func main() {
	mainjs := make(chan string, 1)
	go func() {
		result := esbuild.Build(esbuild.BuildOptions{
			EntryPoints: []string{MAINJS},
			Write:       false,
			Outfile:     "bundle.js",

			Loader: map[string]esbuild.Loader{
				".ts": esbuild.LoaderTS,
			},

			MinifyWhitespace:  true,
			MinifyIdentifiers: true,
			MinifySyntax:      true,
			Bundle:            true,
		})

		if len(result.Errors) != 0 {
			for _, v := range result.Errors {
				fmt.Println(v.Text)
			}
			panic("Errored out because esbuild had errors")
		}

		for _, v := range result.Warnings {
			fmt.Println(v.Text)
		}

		mainjs <- string(result.OutputFiles[0].Contents)
	}()

	maincss := make(chan string)
	go func() {
		command := exec.Command("npm", "run", "build-css")
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr
		command.Dir = "web"

		err := command.Run()
		if err != nil {
			panic(err)
		}

		cssdata, err := os.ReadFile(MAINCSS)
		if err != nil {
			panic(err)
		}

		maincss <- string(cssdata)

		err = os.Remove(MAINCSS)
		if err != nil {
			panic(err)
		}
	}()

	outhtml := fmt.Sprintf(HTMLTEMPLATE, <-mainjs, <-maincss)

	if _, err := os.Stat(path.Dir(BUILDHTML)); os.IsNotExist(err) {
		err := os.MkdirAll(path.Dir(BUILDHTML), 0666)
		if err != nil {
			panic(err)
		}
	}

	err := os.WriteFile(BUILDHTML, []byte(outhtml), 0666)
	if err != nil {
		panic(err)
	}

	fmt.Printf("build.html is %.2fkb\n", float32(len(outhtml))/1024)
}
