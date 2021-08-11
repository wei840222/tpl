package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

func main() {
	log.SetFlags(0)
	var valueFile, templateFile, outputDir string
	flag.StringVar(&valueFile, "v", "./values.csv", "File path for read value file used by render template file.")
	flag.StringVar(&templateFile, "t", "./template.txt", "File path for read template file.")
	flag.StringVar(&outputDir, "o", "./out", "Path for write template file rendered result.")
	flag.Parse()

	vf, err := os.Open(valueFile)
	if err != nil {
		panic(err)
	}
	defer vf.Close()

	cr := csv.NewReader(vf)

	rawCSV, err := cr.ReadAll()
	if err != nil {
		panic(err)
	}

	var header []string
	var values []map[string]string
	for idx, record := range rawCSV {
		if idx == 0 {
			for _, r := range record {
				header = append(header, strings.TrimSpace(r))
			}
			continue
		}
		line := map[string]string{}
		for i, r := range record {
			line[header[i]] = r
		}
		values = append(values, line)
	}

	tf, err := os.Open(templateFile)
	if err != nil {
		panic(err)
	}
	defer tf.Close()

	tfb, err := ioutil.ReadAll(tf)
	if err != nil {
		panic(err)
	}

	tpl, err := template.New("tpl").Funcs(sprig.TxtFuncMap()).Parse(string(tfb))
	if err != nil {
		panic(err)
	}

	os.RemoveAll(outputDir)
	os.Mkdir(outputDir, 0755)

	for _, v := range values {
		var buf bytes.Buffer
		if err := tpl.Execute(&buf, v); err != nil {
			log.Printf("Render file: %s, error: %s", v["file"], err)
			continue
		}

		fn := "nginx" + v["file"]

		f, err := os.OpenFile(path.Join(outputDir, fn), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Open file: %s, error: %s", fn, err)
		}
		if _, err := f.Write(buf.Bytes()); err != nil {
			log.Printf("Write file: %s, error: %s", fn, err)
		}
		if err := f.Close(); err != nil {
			log.Printf("Close file: %s, error: %s", fn, err)
		}
	}
}
