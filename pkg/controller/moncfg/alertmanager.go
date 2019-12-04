package moncfg

import (
	"os"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	monitoringv1alpha1 "github.com/platform9/mon-operator/pkg/apis/monitoring/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func setupAlertmanager(c client.Client, moncfg *monitoringv1alpha1.MonCfg) error {

	var am *monitoringv1.Alertmanager
	am, err := getAlertmanager(c, moncfg.Spec.Prometheus.Name, moncfg.Spec.Global.NameSpace)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Alertmanager object not found creating it")
			err = createAlertmanager(c, moncfg)
			if err != nil {
				log.Error(err, "Failed to create Alertmanager object")
				return err
			}
		} else {
			log.Error(err, "Failed to get Alertmanager object")
			return err
		}
	} else {
		log.Info("Updating existing Alertmanager object: ", "Name", am.Name)
		err = updateAlertmanager(c, am, moncfg)
		if err != nil {
			log.Error(err, "Failed to update Alertmanager object", "Name", am.Name)
			return err
		}
	}

	return nil
}

func createAlertmanager(c client.Client, moncfg *monitoringv1alpha1.MonCfg) error {
	name := moncfg.Spec.Prometheus.Name
	ns := moncfg.Spec.Global.NameSpace
	replicas := moncfg.Spec.Alertmanager.Replicas
	cpu, _ := resource.ParseQuantity(moncfg.Spec.Alertmanager.Resources.Requests.CPU)
	mem, _ := resource.ParseQuantity(moncfg.Spec.Alertmanager.Resources.Requests.Memory)

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	mclient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		return err
	}

	sm := &monitoringv1.Alertmanager{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name + "-alertmanager",
			Namespace: ns,
		},
		Spec: monitoringv1.AlertmanagerSpec{
			ServiceAccountName: moncfg.Spec.Global.SvcAccntName,
			Replicas:           &replicas,
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					"cpu":    cpu,
					"memory": mem,
				},
			},
		},
	}

	sm, err = mclient.MonitoringV1().Alertmanagers(ns).Create(sm)
	if err != nil {
		log.Error(err, "Failed to create alertmanager objects")
		return nil
	}

	return nil
}

func updateAlertmanager(c client.Client, am *monitoringv1.Alertmanager, moncfg *monitoringv1alpha1.MonCfg) error {

	am.Spec.Replicas = &moncfg.Spec.Alertmanager.Replicas
	am.Spec.Resources.Requests["cpu"], _ = resource.ParseQuantity(moncfg.Spec.Alertmanager.Resources.Requests.CPU)
	am.Spec.Resources.Requests["memory"], _ = resource.ParseQuantity(moncfg.Spec.Alertmanager.Resources.Requests.Memory)
	am.Spec.ServiceAccountName = moncfg.Spec.Global.SvcAccntName

	ns := moncfg.Spec.Global.NameSpace

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	mclient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		return err
	}

	am, err = mclient.MonitoringV1().Alertmanagers(ns).Update(am)
	if err != nil {
		log.Error(err, "Failed to create alertmanager objects")
		return nil
	}

	return nil
}

func getAlertmanager(c client.Client, name, ns string) (*monitoringv1.Alertmanager, error) {

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, os.ErrInvalid
	}

	mclient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		return nil, os.ErrInvalid
	}

	var options metav1.GetOptions
	var sm *monitoringv1.Alertmanager
	sm, err = mclient.MonitoringV1().Alertmanagers(ns).Get(name+"-alertmanager", options)
	if err != nil {
		return nil, err
	}
	return sm, nil
}
