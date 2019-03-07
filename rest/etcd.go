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
	Name string `json: "name,omitempty"`
}

func Get(url string) (ETCDMember, error) {
	var etcdmembers ETCDMember

	logrus.Debugf("get members from url %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return etcdmembers, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logrus.Debugf("get %s failed, %s", url, err)
		return etcdmembers, err
	}
	err = json.NewDecoder(resp.Body).Decode(&etcdmembers)
	if err != nil {
		logrus.Debugf("json unmarshal failed, %s", err)
		return etcdmembers, err
	}
	logrus.Debugf("get %s", etcdmembers)
	return etcdmembers, nil
}

func Delete(url string) error  {
	logrus.Debugf("delete member %s\n", url)
	resp, err := resty.R().Delete(url)
	if err != nil {
		logrus.Debug("delete url failed! ", url)
		return err
	}
	defer resp.RawBody().Close()
	if resp.StatusCode() != http.StatusNoContent {
		logrus.Errorf("delete member %s failed: %s\n", url, err)
		return err
	}
	logrus.Debugf("delete %s success", url)
	return nil
}