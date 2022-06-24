package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"wasmhook/pkg/wasmhelper/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	cms, err := clientset.CoreV1().ConfigMaps("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Preloading WASM modules from configmaps")

	modules := map[string]LoadedModule{}

	for _, cm := range cms.Items {
		if strings.HasSuffix(cm.ObjectMeta.Name, "-wasm") {
			fmt.Printf("↳ ConfigMap %s found to be a WASM module\n", cm.ObjectMeta.Name)
			modules[cm.ObjectMeta.Name] = NewLoadedModule(cm.BinaryData["module.wasm"])
			defer modules[cm.ObjectMeta.Name].Close()
		}
	}

	fmt.Println("Done! ☑️\n")

	input := types.HookRequest{
		UID:  k8stypes.UID(123),
		Name: "Hello, Go",
	}
	fmt.Printf("Initial struct => %v\n\n", prettyPrint(input))

	for name, module := range modules {
		marshaled, _ := json.Marshal(input)
		fmt.Printf("Passing to %s WASM module...\n", name)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		errChan := make(chan error, 1)
		resultChan := make(chan []interface{}, 1)
		go func() {
			result, err := module.Run(string(marshaled))
			resultChan <- result
			errChan <- err
		}()
		select {
		case err := <-errChan:
			log.Fatal(err)
		case <-ctx.Done():
			log.Fatal(ctx.Err())
		case results := <-resultChan:
			json.Unmarshal([]byte(results[0].(string)), &input)
			fmt.Printf("Returned struct => %v\n\n", prettyPrint(input))
		}

	}
}

func prettyPrint(b interface{}) string {
	s, _ := json.MarshalIndent(b, "", "\t")
	return fmt.Sprintf(string(s))
}

func loadModule(path string) []byte {
	wasmBytes, _ := ioutil.ReadFile(path)
	return wasmBytes
}
