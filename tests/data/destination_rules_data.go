package data

import (
	"github.com/kiali/kiali/kubernetes"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateEmptyDestinationRule(namespace string, name string, host string) kubernetes.IstioObject {
	destinationRule := (&kubernetes.GenericIstioObject{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			ClusterName: "svc.cluster.local",
		},
		Spec: map[string]interface{}{
			"host": host,
		},
	}).DeepCopyIstioObject()
	return destinationRule
}

func CreateTestDestinationRule(namespace string, name string, host string) kubernetes.IstioObject {
	destinationRule := AddSubsetToDestinationRule(CreateSubset("v1", "v1"),
		AddSubsetToDestinationRule(CreateSubset("v2", "v2"), CreateEmptyDestinationRule(namespace, name, host)))
	return destinationRule
}

func CreateSubset(name string, versionLabel string) map[string]interface{} {
	return map[string]interface{}{
		"name": name,
		"labels": map[string]interface{}{
			"version": versionLabel,
		},
	}
}

func AddSubsetToDestinationRule(subset map[string]interface{}, dr kubernetes.IstioObject) kubernetes.IstioObject {
	if subsetTypeExists, found := dr.GetSpec()["subsets"]; found {
		if subsetTypeCasted, ok := subsetTypeExists.([]interface{}); ok {
			subsetTypeCasted = append(subsetTypeCasted, subset)
			dr.GetSpec()["subsets"] = subsetTypeCasted
		}
	} else {
		dr.GetSpec()["subsets"] = []interface{}{subset}
	}
	return dr
}
