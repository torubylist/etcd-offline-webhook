package main

import (
	"github.com/torubylist/etcd-offline-webhook/server"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)


type Config struct {
	ListenOn string `default:"0.0.0.0:8080"`
	TlsCert  string `default:"/etc/webhook/certs/cert.pem"`
	TlsKey   string `default:"/etc/webhook/certs/key.pem"`
	Debug    bool   `default:"true"`
}

var (
	masterURL  string
	kubeconfig string
)

func main() {
	config := &Config{}
	envconfig.Process("", config)

	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.Infoln(config)
	nsac := &server.EtcdAdmission{}
	s := server.GetEtcdWehhookServer(nsac, config.TlsCert, config.TlsKey, config.ListenOn)
	s.ListenAndServeTLS("", "")
}
