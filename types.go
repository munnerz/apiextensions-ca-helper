package main

type Config struct {
	APIServices                     []APIService                     `json:"apiServices"`
	ValidatingWebhookConfigurations []ValidatingWebhookConfiguration `json:"validatingWebhookConfigurations"`
	MutatingWebhookConfigurations   []MutatingWebhookConfiguration   `json:"mutatingWebhookConfigurations"`
}

type Source struct {
	Secret *Secret `json:"secret"`
	File   *File   `json:"file"`
}

type Secret struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	Namespace string `json:"namespace"`
}

type File struct {
	Path string `json:"path"`
}

type APIService struct {
	Source `json:",inline"`
	Name   string `json:"name"`
}

type ValidatingWebhookConfiguration struct {
	Source `json:",inline"`
	Name   string `json:"name"`
}

type MutatingWebhookConfiguration struct {
	Source `json:",inline"`
	Name   string `json:"name"`
}
