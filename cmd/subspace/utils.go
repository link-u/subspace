package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"
	"strconv"
)

func RandomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)[:n]
}

func Overwrite(filename string, data []byte, perm os.FileMode) error {
	f, err := ioutil.TempFile(filepath.Dir(filename), filepath.Base(filename)+".tmp")
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Chmod(f.Name(), perm); err != nil {
		return err
	}
	return os.Rename(f.Name(), filename)
}

func bash(tmpl string, params interface{}) (string, error) {
	preamble := `
set -o nounset
set -o errexit
set -o pipefail
set -o xtrace
`
	t, err := template.New("template").Parse(preamble + tmpl)
	if err != nil {
		return "", err
	}
	var script bytes.Buffer
	err = t.Execute(&script, params)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	output, err := exec.CommandContext(ctx, "/bin/bash", "-c", string(script.Bytes())).CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %s\n%s", err, string(output))
	}
	return string(output), nil
}

func ipv4host(Number int, IPV4_CIDR string) string {
	var cidr int
	cidr, _ = strconv.Atoi(IPV4_CIDR)
	var IPV4_HOST string = ""
	if cidr <= 0 {
		IPV4_HOST += fmt.Sprintf("%d.", Number/256/256/256%256)
	}
	if cidr <= 8 {
		IPV4_HOST += fmt.Sprintf("%d.", Number/256/256%256)
	}
	if cidr <= 16 {
		IPV4_HOST += fmt.Sprintf("%d.", Number/256%256)
	}
	if cidr <= 24 {
		IPV4_HOST += fmt.Sprintf("%d", Number%256) //上３つと違いピリオドがないのは最後の数だから
	}
	return IPV4_HOST
}
