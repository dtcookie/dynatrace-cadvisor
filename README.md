# dynatrace-cadvisor
The cAdvisor Extension enables [Dynatrace](https://www.dynatrace.com) to consume [cAdvisor](https://github.com/google/cadvisor) container metrics that are getting published by the [kublet](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/) node agents of a [Kubernetes](https://kubernetes.io/docs/home/) Cluster.

<p>
  <a href="https://github.com/dtcookie/dynatrace-cadvisor/actions/workflows/release.yml"><img alt="GitHub Workflow Status (main branch)" src="https://img.shields.io/github/workflow/status/dynatrace-oss/dt-cli/Build%20Test%20Release/main?logo=github"></a>
</p>

## Quick Start
> Prerequisite for this Extension is the installation of the automated [Dynatrace Operator](https://www.dynatrace.com/support/help/setup-and-configuration/setup-on-container-platforms/kubernetes/set-up-k8s-monitoring#automated).
> The Dynatrace Operator contributes the Service Account `dynatrace-kubernetes-monitoring` within the namespace `dynatrace`.

### Installation
```sh
  kubectl apply -f https://github.com/dtcookie/dynatrace-cadvisor/releases/latest/download/dynatrace-cadvisor.yaml
```
This configures a [DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/) within your Kubernetes Cluster. It utilizes the prebuilt Docker Image [dtcookie/dynatrace-cadvisor](https://hub.docker.com/repository/docker/dtcookie/dynatrace-cadvisor).

> The DaemonSet ensures that a [pod](https://kubernetes.io/docs/concepts/workloads/pods/) is getting scheduled to run on every node the Kubernetes Cluster consists of. These pods are now offering the cAdvisor metrics by acting as [Prometheus](https://prometheus.io) Exporters. In addition to that the pods are properly annotated in order for [Dynatrace to scrape these metrics](https://www.dynatrace.com/support/help/how-to-use-dynatrace/infrastructure-monitoring/container-platform-monitoring/kubernetes-monitoring/monitor-prometheus-metrics) automatically.

It will take a minute or two until these metrics are available within the Dynatrace WebUI.
![Metrics](https://raw.githubusercontent.com/dtcookie/dynatrace-cadvisor/main/docs/img/metrics.png)

They are also available for charting via [Data Explorer](https://www.dynatrace.com/support/help/how-to-use-dynatrace/dashboards-and-charts/explorer) and can get split and filtered by dimensions published by cAdvisor aswell as by dimensions discovered by Dynatrace.
![Data Explorer](https://raw.githubusercontent.com/dtcookie/dynatrace-cadvisor/main/docs/img/data-explorer.png)

### Configuration
By default [dynatrace-cadvisor.yaml](https://github.com/dtcookie/dynatrace-cadvisor/releases/latest/download/dynatrace-cadvisor.yaml) restricts the collected metrics to those prefixed with `container_` (which covers almost all metrics offered by cAdvisor).

```yaml
metrics.dynatrace.com/filter: | 
  {
    "mode" : "include",
    "names" : [
        "container_*"
    ]
  }   
```

Please check out the [Dynatrace Documentation](https://www.dynatrace.com/support/help/how-to-use-dynatrace/infrastructure-monitoring/container-platform-monitoring/kubernetes-monitoring/monitor-prometheus-metrics) in case you would like adjust which metrics should get filtered.

### Uninstall dynatrace-cadvisor
```sh
  kubectl delete -f https://github.com/dtcookie/dynatrace-cadvisor/releases/latest/download/dynatrace-cadvisor.yaml
```

## Limitations
* This integration supports up to 200 metric data points each per minute.
* This integration supports only the counter and gauge [Prometheus metric types](https://prometheus.io/docs/concepts/metric_types/).

## Monitoring consumption
Prometheus metrics in Kubernetes environments are subject to DDU consumption. Metrics are first deducted from your quota of included metrics per host unit. Once this quota is exceeded, the remaining metrics consume DDUs.

## License

`dynatrace-cadvisor` is an Open Source Project. Please see [LICENSE](https://github.com/dtcookie/dynatrace-cadvisor/blob/main/LICENSE) for more information.
