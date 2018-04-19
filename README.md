# REY
![rey](https://pbs.twimg.com/profile_images/659732862514692096/XagnatJy_400x400.jpg)

A simple health checker with [status.io](https://status.io) integration.

## How it works?
Rey will read a json based file with a list of components (applications) to perform health checks.
On any status change, rey will consume status.io api to notify this change.

## How to run?
All rey configurations are read by environment variables (except for components list,
this configuration is read by a json based file on file system)

### Environment Variables
| Env  | Description  | Default  |
|---|---|---|
| REY\_AGGREGATOR\_API\_ID | status.io api id  |   |
| REY\_AGGREGATOR\_API\_KEY | status.io api key  |   |
| REY\_CHECKER\_TIMEOUT | http request timeout (in second) | 5 |
| REY\_CHECKER\_MAX\_RETRY | max retry to perform http request on check a component  | 3 |
| REY\_RUNNER\_CIRCLE\_INTERVAL | checker interval (in second) | 10 |
| REY\_COMPONENTS\_PATH | path of components list json based file | /etc/rey/components.json |

### Deploying in a Kubernetes Cluster
We strongly recommend you to put all rey stuff in a new kubernetes namespace
```
$ kubectl create ns rey
```
For the component list you can create a simple json file with a list of components following those fields
```json
]
    {
        "id": "<status.io component id>",
        "name": "<component name>",
        "container_id": "<status.io component id>",
        "hc_endpoint": "<component healthcheck url>",
        "status_page_id": "<status.io status page id>"
    }
]
```
And create a kubernetes configmap
```
$ kubectl create configmap rey-components --from-file=./components.json -n rey
```
> We'll use that confimap to mount a volume with componets file on rey container

For the status.io credentials we strongly recommend you to use kubernetes secrets and mount
that secret as an enviroment variable

```
$ echo -n "<status.io api id>" > ./id
$ echo -n "<status.io api key>" > ./key
$ kubectl create secret generic rey-aggregator --from-file=./id --from-file=./key -n rey
```
Finally you can use the [rey-app.yaml](./rey-app.yaml) template to create a kubernetes deploy

```
$ kubectl create -f rey-app.yaml -n rey
```
