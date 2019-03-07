
package server

import (
	"encoding/json"
	//"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	//"k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	//"k8s.io/api/admission/v1beta1"
)

var (
	AdmissionRequestStatefulset = v1beta1.AdmissionReview{
		TypeMeta: v1.TypeMeta{
			Kind: "AdmissionReview",
		},
		Request: &v1beta1.AdmissionRequest{
			UID: "  uid: 6c0eabcb-0e5a-11e9-8508-001e674fec5a",
			Kind: v1.GroupVersionKind{
				Kind: "StatefulSet",
			},
			Operation: "DELETE",
			Object: runtime.RawExtension{
				Raw: []byte(`{"metadata": {
        						"name": "etcd",
        						"uid": "6c0eabcb-0e5a-11e9-8508-001e674fec5a",
						        "creationTimestamp": "2018-12-28T12:40:19Z",
								"spec": {
                                    "serviceName": "etcd",
                                    "template": {
                                        "spec": {
 											"containers": [
												{"ports": [{"containerPort": 2379, "name": "client"}]}]
										}
									}
                                }
      						}}`),
			},
		},
	}
)

func decodeResponse(body io.ReadCloser) *v1beta1.AdmissionReview {
	response, _ := ioutil.ReadAll(body)
	review := &v1beta1.AdmissionReview{}
	codecs.UniversalDeserializer().Decode(response, nil, review)
	return review
}

func encodeRequest(review *v1beta1.AdmissionReview) []byte {
	ret, err := json.Marshal(review)
	if err != nil {
		logrus.Errorln(err)
	}
	return ret
}

func TestServeReturnsCorrectJson(t *testing.T) {
	nsc := &EtcdAdmission{}
	server := httptest.NewServer(GetEtcdServerNoSSL(nsc, ":8080").Handler)
	requestString := string(encodeRequest(&AdmissionRequestStatefulset))
	myr := strings.NewReader(requestString)
	r, _ := http.Post(server.URL, "application/json", myr)
	review := decodeResponse(r.Body)

	if review.Request.UID != AdmissionRequestStatefulset.Request.UID {
		t.Error("Request and response UID don't match")
	}
}
