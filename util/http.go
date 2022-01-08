package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitee.com/romeo_zpl/infra/logger"
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

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("[Http] invalid status: %d", resp.StatusCode)
		return errors.New(fmt.Sprintf("invalid status: %d", resp.StatusCode))
	}

	body, _ := ioutil.ReadAll(resp.Body)
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

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("[Http] invalid status: %d", resp.StatusCode)
		return errors.New(fmt.Sprintf("invalid status: %d", resp.StatusCode))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	logger.Debugf("[Http], url: %s, HttpGet: %s", url, string(body))

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			logger.Errorf("[Http], HttpGet: %s, Unmarshal error: %s", string(body), err)
			return err
		}
	}

	return nil
}
