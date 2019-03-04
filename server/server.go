package server

import (
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/json"

	"net/http"
)

var (
	scheme          = runtime.NewScheme()
	codecs          = serializer.NewCodecFactory(scheme)
	tlscert, tlskey string
)

type EtcdWebHookController interface {
	HandleEtcdAdmission(review *v1beta1.AdmissionReview) error
}

type EtcdWebHook struct {
	EtcdWebHookController EtcdWebHookController
	Decoder             runtime.Decoder
}

func (ewh *EtcdWebHook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body []byte
		if r.Body != nil {
			if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	logrus.Debugln(string(body))
	review := &v1beta1.AdmissionReview{}
	_, _, err := ewh.Decoder.Decode(body, nil, review)
	if err != nil {
		logrus.Errorln("Can't decode request", err)
	}
	ewh.EtcdWebHookController.HandleEtcdAdmission(review)
	responseInBytes, err := json.Marshal(review)
	logrus.Debugln(string(responseInBytes))

	if _, err := w.Write(responseInBytes); err != nil {
		logrus.Errorln(err)
	}
}

func GetEtcdServerNoSSL(etcdwhc EtcdWebHookController, listenOn string) *http.Server {
	server := &http.Server{
		Handler: &EtcdWebHook{
			EtcdWebHookController: etcdwhc,
			Decoder:             codecs.UniversalDeserializer(),
		},
		Addr: listenOn,
	}

	return server
}

func GetEtcdWehhookServer(etcdwhc EtcdWebHookController, tlsCert, tlsKey, listenOn string) *http.Server {
	sCert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	server := GetEtcdServerNoSSL(etcdwhc, listenOn)
	server.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{sCert},
	}
	if err != nil {
		logrus.Error(err)
	}
	return server
}

