package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ReadFull until error occurs
func ReadFull(reader io.Reader, bs []byte) error {
	expect := len(bs)
	var nrRead, n int
	var err error

	for nrRead < expect {
		if n, err = reader.Read(bs[nrRead:expect]); err != nil {
			return err
		}
		nrRead += n
	}
	return nil
}

func WriteFull(writer io.Writer, bs []byte) error {
	expect := len(bs)
	var nrWrite, n int
	var err error

	for nrWrite < expect {
		if n, err = writer.Write(bs[nrWrite:expect]); err != nil {
			return err
		}
		nrWrite += n
	}
	return nil
}

func Command(path string, args ...string) (result string, err error) {
	cmd := exec.Command(path, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	return out.String(), err
}

// return 0 on invalid string
// for version "2.4.3", returns 2 << 16 + 4 << 8 + 3
func Version2Int(version string) int {
	vs := strings.Split(version, ".")
	// invalid
	if len(vs) == 0 {
		return 0
	}

	v1, err := strconv.Atoi(vs[0])
	if err != nil {
		return 0
	}
	v2, err := strconv.Atoi(vs[1])
	if err != nil {
		return 0
	}

	v3, err := strconv.Atoi(vs[2])
	if err != nil {
		return 0
	}

	return v1<<16 + v2<<8 + v3
}

func DownloadFile(filepath string, url string) error {
	rsp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, rsp.Body)
	return err
}

func Sha256OfFile(filepath string) (string, error) {
	hasher := sha256.New()
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func IPToUint(ip net.IP) uint {
	s := ip.To4()
	return uint(s[0])<<24 + uint(s[1])<<16 + uint(s[2])<<8 + uint(s[3])
}

func UintToIP(ip uint) net.IP {
	v3 := byte(ip & 0xFF)
	v2 := byte((ip >> 8) & 0xFF)
	v1 := byte((ip >> 16) & 0xFF)
	v0 := byte((ip >> 24) & 0xFF)
	return net.IPv4(v0, v1, v2, v3)
}

const (
	chanFullRation = 0.8
)

func ChannelAlmostFull(len, cap int) bool {
	return float32(len) >= float32(cap)*chanFullRation
}
