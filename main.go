package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	apiaggclientset "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
)

// This is a quick tool that should be run as a CronJob to automatically sync
// the contents of Secret resources with APIServices, ValidatingWebhookConfiguration
// and MutatingWebhookConfiguration resources.

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "Path to config file (required)")
}

func main() {
	flag.Parse()
	if configPath == "" {
		log.Fatalf("-config must be specified")
	}

	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	kubeClient, err := getKubeClientset()
	if err != nil {
		log.Fatalf("error building kubeclient: %v", err)
	}

	apiAggClient, err := getKubeAggClientset()
	if err != nil {
		log.Fatalf("error building kube aggregator clientset: %v", err)
	}

	p := &processor{kubeClient: kubeClient, apiAggClient: apiAggClient}
	err = p.processConfig(cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

type processor struct {
	kubeClient   kubernetes.Interface
	apiAggClient apiaggclientset.Interface
}

func (p *processor) processConfig(cfg *Config) error {
	var errs []error
	for _, apisvc := range cfg.APIServices {
		errs = append(errs, p.processAPIService(apisvc))
	}
	for _, vwc := range cfg.ValidatingWebhookConfigurations {
		errs = append(errs, p.processValidatingWebhookConfiguration(vwc))
	}
	for _, mwc := range cfg.MutatingWebhookConfigurations {
		errs = append(errs, p.processMutatingWebhookConfiguration(mwc))
	}
	return utilerrors.NewAggregate(errs)
}

func (p *processor) processAPIService(a APIService) error {
	if a.Name == "" {
		return fmt.Errorf("APIService name not set")
	}
	data, err := p.loadSource(a.Source)
	if err != nil {
		return err
	}

	apiService, err := p.apiAggClient.ApiregistrationV1beta1().APIServices().Get(a.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	apiService.Spec.CABundle = data
	_, err = p.apiAggClient.ApiregistrationV1beta1().APIServices().Update(apiService)
	if err != nil {
		return err
	}

	return nil
}

func (p *processor) processValidatingWebhookConfiguration(a ValidatingWebhookConfiguration) error {
	if a.Name == "" {
		return fmt.Errorf("ValidatingWebhookConfiguration name not set")
	}
	data, err := p.loadSource(a.Source)
	if err != nil {
		return err
	}

	vwc, err := p.kubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Get(a.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for i := range vwc.Webhooks {
		vwc.Webhooks[i].ClientConfig.CABundle = data
	}

	_, err = p.kubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Update(vwc)
	if err != nil {
		return err
	}

	return nil
}

func (p *processor) processMutatingWebhookConfiguration(a MutatingWebhookConfiguration) error {
	if a.Name == "" {
		return fmt.Errorf("MutatingWebhookConfiguration name not set")
	}
	data, err := p.loadSource(a.Source)
	if err != nil {
		return err
	}

	mwc, err := p.kubeClient.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Get(a.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for i := range mwc.Webhooks {
		mwc.Webhooks[i].ClientConfig.CABundle = data
	}

	_, err = p.kubeClient.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Update(mwc)
	if err != nil {
		return err
	}

	return nil
}

func (p *processor) loadSource(s Source) ([]byte, error) {
	var data []byte
	var err error
	switch {
	case s.File != nil:
		path := s.File.Path
		data, err = ioutil.ReadFile(path)
	case s.Secret != nil:
		data, err = p.loadAPISecret(*s.Secret)
	default:
		return nil, fmt.Errorf("certificate source not specified")
	}
	if err != nil {
		return nil, err
	}
	return data, err
}
func (p *processor) loadAPISecret(s Secret) ([]byte, error) {
	secret, err := p.kubeClient.CoreV1().Secrets(s.Namespace).Get(s.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	data, ok := secret.Data[s.Key]
	if !ok {
		return nil, fmt.Errorf("key %q not found in secret '%s/%s'", s.Key, s.Namespace, s.Name)
	}
	return data, nil
}

func getConfig() (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(file)
	out := &Config{}
	err = d.Decode(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
