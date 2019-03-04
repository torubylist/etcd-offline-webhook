package server

import (
	"github.com/sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//appsv1 "k8s.io/api/apps/v1"
	"strings"
	"github.com/torubylist/etcd-offline-webhook/rest"
	"fmt"

	"k8s.io/client-go/kubernetes"
)

type EtcdAdmission struct {
	Kclient *kubernetes.Clientset
}

func (ea *EtcdAdmission) HandleEtcdAdmission(ar *v1beta1.AdmissionReview) error {
	//raw := ar.Request.Object.Raw
	//statefulset := appsv1.StatefulSet{}
	//deserializer := codecs.UniversalDeserializer()
	//if _, _, err := deserializer.Decode(raw, nil, &statefulset); err != nil {
	//	logrus.Error(err)
	//	return err
	//}
	logrus.Debugln(ar.Request)

	ar.Response = &v1beta1.AdmissionResponse{}
	ar.Response.Allowed = false
	namespace := ar.Request.Namespace
	//assume servicename the same as statefulset name.
	serviceName := ar.Request.Name
	logrus.Debugln("statefulset servicename: ", serviceName)
	//namespace := statefulset.ObjectMeta.Namespace
	logrus.Debugln("statefulset namespace: ", namespace)
	logrus.Debugln("statefulset name: ", serviceName)
	service, err := ea.Kclient.CoreV1().Services(namespace).Get(serviceName)
	if err != nil {
		return err
	}
	var port int32
	for _,p := range service.Spec.Ports {
		if strings.Contains(p.Name, "client") {
			port = p.Port
		}
	}
	//etcd service url
	if  strings.Contains(ar.Request.Name, "etcd") {
		url := fmt.Sprintf("http://%s.%s.svc:%d/%s", "etcd", namespace, port, "v2/members")
		etcdmems, err := rest.Get(url)
		if err != nil {
			return err
		}
		for _, member := range etcdmems.Members {
			delUrl := fmt.Sprintf("%s/%s", url, member.Id)
			err := rest.Delete(delUrl)
			if err != nil {
				return err
			}
		}
		ar.Response.Result = &metav1.Status{
			Message: "etcd member has been deleted",
		}
	}else {
		ar.Response.Result = &metav1.Status{
			Message: "the statefulset isn't etcd",
		}
	}
	ar.Response.Allowed = true
	return nil
}