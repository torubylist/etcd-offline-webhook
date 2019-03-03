package server

import (
	"github.com/sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"strings"
	"github.com/torubylist/etcd-offline-webhook/rest"
	"fmt"

)

type EtcdAdmission struct {
}

func (ea *EtcdAdmission) HandleEtcdAdmission(ar *v1beta1.AdmissionReview) error {
	raw := ar.Request.Object.Raw
	statefulset := appsv1.StatefulSet{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &statefulset); err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Debugln(statefulset)

	ar.Response = &v1beta1.AdmissionResponse{}
	ar.Response.Allowed = false
	serviceName := statefulset.Spec.ServiceName
	logrus.Debugln("statefulset servicename: ", serviceName)
	namespace := statefulset.ObjectMeta.Namespace
	logrus.Debugln("statefulset namespace: ", namespace)
	logrus.Debugln("statefulset name: ", statefulset.ObjectMeta.Name)
	if  strings.Contains(statefulset.ObjectMeta.Name, "etcd") {
		containers := statefulset.Spec.Template.Spec.Containers
		ports := make(map[string]int32)
		var portName string
		//get etcd  client port
		for _, container := range containers {
			for _, port := range container.Ports {
				if strings.Contains(port.Name, "client") {
					portName = port.Name
					ports[portName]= port.ContainerPort
					break
				}
			}
		}
		//etcd service url
		url := fmt.Sprintf("http://%s.%s.svc:%d/%s", serviceName, namespace, ports[portName], "v2/members")
		etcd := &rest.ETCD{}
		etcdmems, err := etcd.Get(url)
		if err != nil {
			return err
		}
		for _, member := range etcdmems.Members {
			delUrl := fmt.Sprintf("%s/%s", url, member.Id)
			err := etcd.Delete(delUrl)
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