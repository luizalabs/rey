# REY
![rey](https://pbs.twimg.com/profile_images/659732862514692096/XagnatJy_400x400.jpg)

A simple health checker.

## How it works?
Rey will read a json based file with a list of components (applications) to perform health checks.
On any status change, rey will expose a gauge metric to prometheus.

## How to run?
All rey configurations are read by environment variables (except for components list,
this configuration is read by a json based file on file system)

### Environment Variables
| Env | Description | Default |
|---|---|---|
| REY\_CHECKER\_TIMEOUT | http request timeout (in second) | 5 |
| REY\_CHECKER\_MAX\_RETRY | max retry to perform http request on check a component  | 3 |
| REY\_RUNNER\_CIRCLE\_INTERVAL | checker interval (in second) | 10 |
| REY\_COMPONENTS\_PATH | path of components list json based file | /etc/rey/components.json |
| REY\_METRICS\_SERVER\_PORT | Port to listen on Prometheus metrics server | 5000 |

### Deploying in a Kubernetes Cluster
We strongly recommend you to put all rey stuff in a new kubernetes namespace
```
$ kubectl create ns rey
```
For the component list you can create a simple json file with a list of components following those fields
```json
]
    {
        "name": "<component name>",
        "hc_endpoint": "<component healthcheck url>",
    }
]
```
And create a kubernetes configmap
```
$ kubectl create configmap rey-components --from-file=./components.json -n rey
```
> We'll use that configmap to mount a volume with componets file on rey container

Finally you can use the [rey-app.yaml](./rey-app.yaml) template to create a kubernetes deploy and
a kubernetes ClusterIP service (to expose prometheus metrics)

```
$ kubectl create -f rey-app.yaml -n rey
```

## Metrics
When you deploy rey on your Kubernetes cluster rey will expose simple metrics to [Prometheus](https://prometheus.io)
with components health status (same as HTTP status of each request and `500` for request errors [e.g. timeouts]).
With these metrics you can create a [Grafana](https://grafana.com) dashboard, just use queries like:
`rey{componente_name="<name of the component>"}`.

