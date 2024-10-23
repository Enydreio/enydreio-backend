package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func BuildIngressURL(host string, ip string, port int32, path string) string {
	url := "https://" + ip
	if host != "" {
		url = "https://" + host
	}
	url += fmt.Sprintf(":%d%s", port, path)
	return url
}
func GetLoadBalancerIP(ingressStatus []v1.IngressLoadBalancerIngress) string {
	for _, loadBalancer := range ingressStatus {
		if loadBalancer.IP != "" {
			return loadBalancer.IP
		}
	}
	return ""
}
func ScanApps() {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ingressList, err := clientset.NetworkingV1().Ingresses("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, ingress := range ingressList.Items {
		ip := GetLoadBalancerIP(ingress.Status.LoadBalancer.Ingress)
		for _, rule := range ingress.Spec.Rules {
			host := rule.Host

			for _, path := range rule.HTTP.Paths {
				backend := path.Backend.Service
				if backend == nil {
					continue
				}

				serviceName := backend.Name
				servicePort := backend.Port.Number

				serviceURL := BuildIngressURL(host, ip, servicePort, path.Path)
				var existingApp Application
				result := db.Where("name = ?", serviceName).First(&existingApp)
				if result.Error == nil {
					if existingApp.Url != serviceURL {
						existingApp.Url = serviceURL
						db.Save(&existingApp)
						fmt.Printf("Service %s updated in the database successfully\n", serviceName)
					}
				} else {
					app := Application{
						Name: serviceName,
						Url:  serviceURL,
					}
					result := db.Create(&app)
					if result.Error != nil {
						fmt.Printf("Failed to save service %s to the database: %v\n", serviceName, result.Error)
					} else {
						fmt.Printf("Service %s saved to the database successfully\n", serviceName)
					}
				}
			}

		}
	}
}