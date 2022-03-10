/**
* @license
* Copyright 2020 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// By default container metrics for apps running
// in namespaces `kube-system` and `dynatrace`
// will get ignored
var excludedNamespaces = []string{
	"dynatrace",
	"kube-system",
}

func metrics(w http.ResponseWriter, req *http.Request) {
	if len(nodeName) == 0 {
		return
	}
	data, err := clientset.CoreV1().RESTClient().Get().Resource("nodes").SubResource(nodeName+":10250", "proxy", "metrics", "cadvisor").DoRaw(context.TODO())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	data, _ = filter(data)
	w.Write(data)
}

var clientset *kubernetes.Clientset
var nodeName string

func filter(data []byte) ([]byte, int) {
	numLines := 0
	buffer := new(bytes.Buffer)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		txt := scanner.Text()
		// We don't need to pass on comment lines
		if strings.HasPrefix(txt, "#") {
			continue
		}
		// Metrics for `kube-system` and `dynatrace` are
		// by default getting ignored
		for _, namespace := range excludedNamespaces {
			if strings.Contains(txt, fmt.Sprintf("namespace=\"%s\"", namespace)) {
				continue
			}
		}
		numLines++
		fmt.Fprintln(buffer, txt)
	}
	return buffer.Bytes(), numLines
}

func main() {
	var err error
	var config *rest.Config

	// Create a REST client for accessing K8s API Server
	config, err = rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Discovering Node Name of this specific POD
	podName, _ := os.LookupEnv("HOSTNAME")
	// TODO: Not a good idea to hard code the namespace `dynatrace` here
	pod, err := clientset.CoreV1().Pods("dynatrace").Get(context.TODO(), podName, metav1.GetOptions{})
	nodeName = pod.Spec.NodeName
	if err != nil {
		fmt.Println("Error resolving node name", err.Error())
	}

	// Launch HTTP Server - acting as Prometheus Exporter
	http.HandleFunc("/metrics", metrics)
	http.ListenAndServe(":9001", nil)
}
