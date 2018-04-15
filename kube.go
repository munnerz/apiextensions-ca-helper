package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	apiaggclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
)

func getKubeClientset() (kubernetes.Interface, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
}

func getKubeAggClientset() (apiaggclientset.Interface, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return apiaggclientset.NewForConfig(cfg)
}
