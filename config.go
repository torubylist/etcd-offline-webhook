package main

import (
	"crypto/tls"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/sirupsen/logrus"
)

// Get a clientset with in-cluster config.
func getClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		logrus.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatal(err)
	}
	return clientset
}

func configTLS(config Config, clientset *kubernetes.Clientset) *tls.Config {
	sCert, err := tls.LoadX509KeyPair(config.TlsCert, config.TlsKey)
	if err != nil {
		logrus.Fatal(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
		// TODO: uses mutual tls after we agree on what cert the apiserver should use.
		// ClientAuth:   tls.RequireAndVerifyClientCert,
	}
}