package moncfg

import (
	"os"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	monitoringv1alpha1 "github.com/platform9/mon-operator/pkg/apis/monitoring/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func setupPrometheus(c client.Client, moncfg *monitoringv1alpha1.MonCfg) error {

	var pm *monitoringv1.Prometheus
	pm, err := getPrometheus(c, moncfg.Spec.Prometheus.Name, moncfg.Spec.Global.NameSpace)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Prometheus object not found creating it")
			err = createPrometheus(c, moncfg)
			if err != nil {
				log.Error(err, "Failed to create prometheus object")
				return err
			}
		} else {
			log.Error(err, "Failed to get prometheus object")
			return err
		}
	} else {
		log.Info("Updating existing prometheus object: ", "Name", pm.Name)
		err = updatePrometheus(c, pm, moncfg)
		if err != nil {
			log.Error(err, "Failed to update prometheus object", "Name", pm.Name)
			return err
		}
	}

	return nil
}

func createPrometheus(c client.Client, moncfg *monitoringv1alpha1.MonCfg) error {
	replicas := moncfg.Spec.Prometheus.Replicas
	cpu, _ := resource.ParseQuantity(moncfg.Spec.Prometheus.Resources.Requests.CPU)
	mem, _ := resource.ParseQuantity(moncfg.Spec.Prometheus.Resources.Requests.Memory)
	name := moncfg.Spec.Prometheus.Name
	ns := moncfg.Spec.Global.NameSpace

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	mclient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		return err
	}

	pm := &monitoringv1.Prometheus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: monitoringv1.PrometheusSpec{
			Alerting: &monitoringv1.AlertingSpec{
				Alertmanagers: []monitoringv1.AlertmanagerEndpoints{
					monitoringv1.AlertmanagerEndpoints{
						Namespace: ns,
						Name:      name + "-alertmanager",
						Port:      intstr.FromString(moncfg.Spec.Prometheus.Port),
					},
				},
			},
			Retention: moncfg.Spec.Prometheus.Retention,
			Replicas:  &replicas,
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					"cpu":    cpu,
					"memory": mem,
				},
			},
			RuleSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"prometheus": name,
					"role":       "alert-rules",
				},
			},
			ServiceAccountName: moncfg.Spec.Global.SvcAccntName,
			ServiceMonitorSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"prometheus": name,
					"role":       "alert-rules",
				},
			},
		},
	}

	pm, err = mclient.MonitoringV1().Prometheuses(ns).Create(pm)
	if err != nil {
		log.Error(err, "Failed to create prometheus objects")
		return nil
	}

	return nil
}

func updatePrometheus(c client.Client, pm *monitoringv1.Prometheus, moncfg *monitoringv1alpha1.MonCfg) error {
	pm.Spec.Replicas = &moncfg.Spec.Prometheus.Replicas
	pm.Spec.Resources.Requests["cpu"], _ = resource.ParseQuantity(moncfg.Spec.Prometheus.Resources.Requests.CPU)
	pm.Spec.Resources.Requests["memory"], _ = resource.ParseQuantity(moncfg.Spec.Prometheus.Resources.Requests.Memory)
	pm.Spec.Retention = moncfg.Spec.Prometheus.Retention
	pm.Spec.Alerting.Alertmanagers[0].Port = intstr.FromString(moncfg.Spec.Prometheus.Port)
	pm.Spec.ServiceAccountName = moncfg.Spec.Global.SvcAccntName

	ns := moncfg.Spec.Global.NameSpace

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	mclient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		return err
	}

	pm, err = mclient.MonitoringV1().Prometheuses(ns).Update(pm)
	if err != nil {
		log.Error(err, "Failed to create prometheus objects")
		return nil
	}

	return nil
}

func getPrometheus(c client.Client, name, ns string) (*monitoringv1.Prometheus, error) {

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, os.ErrInvalid
	}

	mclient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		return nil, os.ErrInvalid
	}

	var options metav1.GetOptions
	var pm *monitoringv1.Prometheus
	pm, err = mclient.MonitoringV1().Prometheuses(ns).Get(name, options)
	if err != nil {
		return nil, err
	}
	return pm, nil
}
