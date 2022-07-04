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
	"strings"
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

func ipv4address(Number int,IPv4_POOL string) (string,error) {
    var err error = nil
    
    var hostcidrS[] string = strings.Split(IPv4_POOL,"/")
    cidr, _ := strconv.Atoi(hostcidrS[1])
    subnet := ((1<<cidr)-1)
    subnetmask := ((1<<32)-1)^subnet
    
    var hostOctetS[] string = strings.Split(hostcidrS[0],".")
    var hostOctet[4] int
    hostOctet[0], _=strconv.Atoi(hostOctetS[0])
    hostOctet[1], _=strconv.Atoi(hostOctetS[1])
    hostOctet[2], _=strconv.Atoi(hostOctetS[2])
    hostOctet[3], _=strconv.Atoi(hostOctetS[3])
    host := (hostOctet[0]<<24)|(hostOctet[1]<<16)|(hostOctet[2]<<8)|(hostOctet[3]<<0)
    
    ip := host|Number
    ipOctet := [4]int {(ip>>24)&255, (ip>>16)&255, (ip>>8)&255, (ip>>0)&255}
    
    var ipS string = fmt.Sprintf("%d.%d.%d.%d",ipOctet[0],ipOctet[1],ipOctet[2],ipOctet[3])
    
    if (host&subnet) != 0 {
        err = fmt.Errorf("IPv4_POOL=\"%s\"",IPv4_POOL)
    }else if (Number&subnetmask) != 0 {
        err = fmt.Errorf("Number does not fit in the POOL. Number=%d, POOL=\"/%d\"",Number,cidr)
    }else if Number == subnet {
        err = fmt.Errorf("Broadcast address. Number=%d, IPv4_POOL=\"%s\", address=\"%s\"",Number,IPv4_POOL,ipS)
    }else if Number == 0 {
        err = fmt.Errorf("Network address. Number=%d, IPv4_POOL=\"%s\", address=\"%s\"",Number,IPv4_POOL,ipS)
    }
    
    return ipS,err
}


