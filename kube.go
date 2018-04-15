package main

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	apiaggclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
)

func getKubeClientset() (kubernetes.Interface, error) {
	apiCfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return nil, fmt.Errorf("error loading cluster config: %v", err)
	}

	cfg, err := clientcmd.NewDefaultClientConfig(*apiCfg, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error building client config: %v", err)
	}

	return kubernetes.NewForConfig(cfg)
}

func getKubeAggClientset() (apiaggclientset.Interface, error) {
	apiCfg, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return nil, fmt.Errorf("error loading cluster config: %s", err.Error())
	}

	cfg, err := clientcmd.NewDefaultClientConfig(*apiCfg, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error building client config: %v", err)
	}

	return apiaggclientset.NewForConfig(cfg)
}
