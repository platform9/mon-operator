package moncfg

import (
	"os"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	monitoringv1alpha1 "github.com/platform9/mon-operator/pkg/apis/monitoring/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func setupServiceMonitor(c client.Client, moncfg *monitoringv1alpha1.MonCfg) error {

	var pm *monitoringv1.ServiceMonitor
	pm, err := getServiceMonitor(c, moncfg.Spec.Prometheus.Name, moncfg.Spec.Global.NameSpace)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("ServiceMonitor object not found creating it")
			err = createServiceMonitor(c, moncfg)
			if err != nil {
				log.Error(err, "Failed to create ServiceMonitor object")
				return err
			}
		} else {
			log.Error(err, "Failed to get ServiceMonitor object")
			return err
		}
	} else {
		log.Info("Updating existing ServiceMonitor object: ", "Name", pm.Name)
		err = updateServiceMonitor(c, pm, moncfg)
		if err != nil {
			log.Error(err, "Failed to update ServiceMonitor object", "Name", pm.Name)
			return err
		}
	}

	return nil
}

func createServiceMonitor(c client.Client, moncfg *monitoringv1alpha1.MonCfg) error {
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

	sm := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name + "-svcmon",
			Namespace: ns,
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			Endpoints: []monitoringv1.Endpoint{
				monitoringv1.Endpoint{
					Port: moncfg.Spec.Prometheus.Port,
				},
			},
			Selector: metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					metav1.LabelSelectorRequirement{
						Key:      moncfg.Spec.Prometheus.Selector[0].Key,
						Operator: metav1.LabelSelectorOpIn,
						Values:   moncfg.Spec.Prometheus.Selector[0].Values,
					},
				},
			},
			NamespaceSelector: monitoringv1.NamespaceSelector{
				MatchNames: moncfg.Spec.Prometheus.NameSpaceSelector,
			},
		},
	}

	sm, err = mclient.MonitoringV1().ServiceMonitors(ns).Create(sm)
	if err != nil {
		log.Error(err, "Failed to create service monitor objects")
		return nil
	}

	return nil
}

func updateServiceMonitor(c client.Client, sm *monitoringv1.ServiceMonitor, moncfg *monitoringv1alpha1.MonCfg) error {
	sm.Spec.Endpoints[0].Port = moncfg.Spec.Prometheus.Port
	sm.Spec.Selector.MatchExpressions[0].Key = moncfg.Spec.Prometheus.Selector[0].Key
	sm.Spec.Selector.MatchExpressions[0].Values = moncfg.Spec.Prometheus.Selector[0].Values
	sm.Spec.NamespaceSelector.MatchNames = moncfg.Spec.Prometheus.NameSpaceSelector

	ns := moncfg.Spec.Global.NameSpace

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	mclient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		return err
	}

	sm, err = mclient.MonitoringV1().ServiceMonitors(ns).Update(sm)
	if err != nil {
		log.Error(err, "Failed to create service monitor objects")
		return nil
	}

	return nil
}

func getServiceMonitor(c client.Client, name, ns string) (*monitoringv1.ServiceMonitor, error) {

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, os.ErrInvalid
	}

	mclient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		return nil, os.ErrInvalid
	}

	var options metav1.GetOptions
	var sm *monitoringv1.ServiceMonitor
	sm, err = mclient.MonitoringV1().ServiceMonitors(ns).Get(name+"-svcmon", options)
	if err != nil {
		return nil, err
	}
	return sm, nil
}
