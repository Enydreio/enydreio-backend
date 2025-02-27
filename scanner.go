package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net"
	"os"
)

func BuildURL(host string, ip string, port int32, path string) string {
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
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err.Error())
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err.Error())
		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func ScanKubeApps(isExtern bool) {
	var config *rest.Config
	var err error
	if isExtern {
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	} else {
		config, err = rest.InClusterConfig()
	}

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

				serviceURL := BuildURL(host, ip, servicePort, path.Path)
				SaveApp(serviceName, serviceURL)
			}

		}
	}
}
func SaveApp(serviceName string, serviceURL string) {
	// First check if an app with this name already exists
	var existingAppByName Application
	resultByName := db.Where("name = ?", serviceName).First(&existingAppByName)

	if resultByName.Error == nil {
		// App with this name exists, update URL if changed
		if existingAppByName.Url != serviceURL {
			existingAppByName.Url = serviceURL
			db.Save(&existingAppByName)
			fmt.Printf("Service %s updated in the database successfully\n", serviceName)
		}
	} else {
		// App with this name doesn't exist, check if URL exists
		var existingAppByURL Application
		resultByURL := db.Where("url = ?", serviceURL).First(&existingAppByURL)

		if resultByURL.Error == nil {
			// App with this URL exists but name was changed, update the name
			existingAppByURL.Name = serviceName
			db.Save(&existingAppByURL)
			fmt.Printf("Service name updated from %s to %s in the database\n", existingAppByURL.Name, serviceName)
		} else {
			// No app with this name or URL exists, create new
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
func ScanDockerApps() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ctr := range containers {
		SaveApp(ctr.Names[0], BuildURL("", GetOutboundIP(), int32(ctr.Ports[0].PublicPort), "/"))
	}
}
