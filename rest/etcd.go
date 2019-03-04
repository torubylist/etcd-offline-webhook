package rest

import (
	"gopkg.in/resty.v1"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ETCDMember struct {
	Members []Member `json: "members"`
}

type Member struct {
	Id string `json: "id,omitempty"`
	Name string `json: "name, omitempty"`
}

func Get(url string) (ETCDMember, error) {
	logrus.Debugf("get members from url %s\n", url)
	resp, err := resty.R().Get(url)
	defer resp.RawBody().Close()
	var etcdmembers ETCDMember
	if resp.StatusCode() != http.StatusOK {
		logrus.Debugf("get %s failed, %s", url, err)
		return etcdmembers, err
	}
	err = json.Unmarshal(resp.Body(), &etcdmembers)
	if err != nil {
		logrus.Debugf("json unmarshal failed, %s", err)
		return etcdmembers, err
	}
	return etcdmembers, nil
}

func Delete(url string) error  {
	logrus.Debugf("delete member %s\n", url)
	resp, err := resty.R().Delete(url)
	defer resp.RawBody().Close()
	if resp.StatusCode() != http.StatusNoContent {
		logrus.Errorf("delete member %s failed: %s\n", url, err)
		return err
	}
	return nil
}