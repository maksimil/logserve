package main

import (
	"fmt"
	"os"
	"path"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

const HTMLTEMPLATE = `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8" /><meta http-equiv="X-UA-Compatible" content="IE=edge" /><meta name="viewport" content="width=device-width, initial-scale=1.0" /><title>Document</title></head><body></body><script>%s</script></html>`

const (
	MAINJS    = "web/main.js"
	BUILDHTML = "cmd/_gen/build.html"
)

func main() {
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints: []string{MAINJS},
		Write:       false,
		Outfile:     "bundle.js",

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

	outhtml := fmt.Sprintf(HTMLTEMPLATE, string(result.OutputFiles[0].Contents))

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

	fmt.Printf("build.html is %.2fkb", float32(len(outhtml))/1024)
}
