package util

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"

	"gitee.com/romeo_zpl/infra/logger"
)

var (
	ErrAuth = errors.New("authentication failed")
)

func HttpPostJson(url string, value interface{}, headers map[string]string, out interface{}) error {
	str, err := json.Marshal(value)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(str))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrAuth
	}

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("[Http] invalid status: %d", resp.StatusCode)
		return errors.New(fmt.Sprintf("invalid status: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	logger.Debugf("[Http], url: %s, value: %s, HttpGet: %s", url, value, string(body))

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			logger.Errorf("[Http], HttpPostJson: %s, Unmarshal error: %s", string(body), err)
			return err
		}
	}

	return nil
}

func HttpGet(url string, headers map[string]string, out interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrAuth
	}

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("[Http] invalid status: %d", resp.StatusCode)
		return errors.New(fmt.Sprintf("invalid status: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.Debugf("[Http], url: %s, HttpGet: %s", url, string(body))

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			logger.Errorf("[Http], HttpGet: %s, Unmarshal error: %s", string(body), err)
			return err
		}
	}

	return nil
}

type CheckSum string

const (
	MD5    CheckSum = "md5"
	Sha256          = "sha256"
)

type sumDownloader struct {
	method CheckSum
	f      *os.File
	hash   hash.Hash
}

func newSumDownloader(method CheckSum, out string) (*sumDownloader, error) {
	s := &sumDownloader{
		method: method,
	}

	var err error
	s.f, err = os.Create(out)
	if err != nil {
		return nil, err
	}

	if method == MD5 {
		s.hash = md5.New()
	} else if method == Sha256 {
		s.hash = sha256.New()
	} else {
		s.hash = nil
	}
	return s, nil
}

func (d *sumDownloader) Write(p []byte) (n int, err error) {
	if d.hash != nil {
		d.hash.Write(p)
	}
	return d.f.Write(p)
}

func (d *sumDownloader) Sum() string {
	bs := d.hash.Sum(nil)
	return hex.EncodeToString(bs)
}

func (d *sumDownloader) Close() {
	d.f.Close()
}

func HttpDownload(url string, headers map[string]string, file string, method CheckSum) (checksum string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("[Http] invalid status: %d", resp.StatusCode)
		return "", errors.New(fmt.Sprintf("invalid status: %d", resp.StatusCode))
	}

	d, err := newSumDownloader(method, file)
	if err != nil {
		return
	}
	defer d.Close()

	_, err = io.Copy(d, resp.Body)
	return d.Sum(), err
}
