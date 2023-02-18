package kubernetes

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	util "github.com/tonit/rekind/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getClient() (*k8s.Clientset, error) {
	var kubeconfig *string
	if given := os.Getenv("KUBECONFIG"); given != "" {
		kubeconfig = flag.String("kubeconfig", given, "(optional) absolute path to the kubeconfig file")
	} else if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	return k8s.NewForConfig(config)
}

func GetClusterVersion() string {
	client, err := getClient()
	if err != nil {
		panic(err)
	}
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}
	versionsFound := map[string]struct{}{}
	for _, d := range nodes.Items {
		versionsFound[d.Status.NodeInfo.KubeletVersion] = struct{}{}
	}
	if len(versionsFound) != 1 {
		panic("Not a unique version found..")
	}

	keys := make([]string, 0, len(versionsFound))
	for k := range versionsFound {
		keys = append(keys, k)
	}

	return util.NormalizeVersionToMinor(keys[0])
}
