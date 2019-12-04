package moncfg

import (
	"context"
	"io/ioutil"
	"os"

	monitoringv1alpha1 "github.com/platform9/mon-operator/pkg/apis/monitoring/v1alpha1"
	yaml "gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	configDir = "/etc/alertmgrcfg"
)

func setupSecret(c client.Client, moncfg *monitoringv1alpha1.MonCfg) error {
	file, err := os.Open(configDir + "/alertmanager.yaml")
	if err != nil {
		log.Error(err, "Failed to open alert manager config file")
		return os.ErrInvalid
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err, "Failed to read alert manager config file")
		return os.ErrInvalid
	}

	var acfg alertConfig
	yaml.Unmarshal(data, &acfg)

	for _, recv := range moncfg.Spec.Alertmanager.Receivers {
		log.Info("Listing receiver: ", "Type", recv.Type)
		err = formatReceiver(moncfg, &acfg)
		if err != nil {
			log.Error(err, "Failed to format receiver for ", "Type", recv.Type)
			return err
		}
	}

	data, err = yaml.Marshal(&acfg)
	if err != nil {
		log.Error(err, "Failed to marshal alert mgr secret ")
		return err
	}

	secretName := "alertmanager-" + moncfg.Spec.Prometheus.Name + "-alertmanager"

	_, err = deleteSecret(c, moncfg.Spec.Global.NameSpace, secretName)
	if err != nil {
		log.Error(err, "Failed to delete secret", "secretname", secretName)
		return err
	}

	err = createSecret(c, moncfg.Spec.Global.NameSpace, secretName, data)
	if err != nil {
		log.Error(err, "Failed to create secret", "secretname", secretName)
		return err
	}
	log.Info("Created secret: ", "secretname", secretName)

	return nil
}

func deleteSecret(c client.Client, ns string, secretName string) (bool, error) {
	sec := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: ns,
		},
	}

	err := c.Delete(context.TODO(), &sec)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

/*func checkSecretExists(c client.Client, ns string, secretName string) (bool, error) {
	key := types.NamespacedName{Name: secretName, Namespace: ns}
	var sec v1.Secret

	err := c.Get(context.TODO(), key, &sec)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}*/

func createSecret(c client.Client /*obj *metav1.ObjectMeta,*/, ns string, secretName string /*kind string,*/, data []byte) error {

	//trueVar := true
	cfg := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: ns,
			Annotations: map[string]string{
				"created_by": "alertmgr-controller",
			},
		},
		Data: map[string][]byte{
			"alertmanager.yaml": data,
		},
	}

	/*if obj != nil {
		cfg.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
			metav1.OwnerReference{
				APIVersion: monitoringv1.SchemeGroupVersion.String(),
				Name:       obj.GetName(),
				Kind:       kind,
				UID:        obj.GetUID(),
				Controller: &trueVar,
			},
		}
	}*/

	if err := c.Create(context.TODO(), cfg); err != nil {
		return err
	}

	return nil
}
