# ZEX APP Helm Chart

## Local deployment - prerequisites
* [docker](https://www.docker.com/) v25.0.3+
* [helm](https://helm.sh/) v.3.8.1+
* [kind](https://kind.sigs.k8s.io/) v0.15.0+ or any other K8s local cluster.
* [ingress-nginx](https://github.com/kubernetes/ingress-nginx) ingress controller.
* [metrics-server](https://github.com/kubernetes-sigs/metrics-server)
* [istio](https://istio.io/latest/docs/setup/getting-started/#download) service mesh.

### Local setup scripts
Check out the `scripts` folder for easy deployment of the prerequisites.
Docker, helm and kind cli must be installed before running the scripts.

Enter the `scripts` folder, and run the following scripts in order:
* `kind-cluster.sh` - bash script deploys a kind cluster.
* `ingress-controller.sh` - bash script for deploying the `ingress-nginx` ingress controller.
* `metrics-server.sh` - bash script which installs metrics server needed for autoscaling.
Feel free to tinker with the configuration files if needed.

After these scripts, only `istio` service mesh needs to be installed manually via the `istio` cli tool `istioctl install --set profile=demo -y`

## Chart install
Use helm cli to install chart
    `helm upgrade --install -n <your_namespace> --create-namespace <your_chart_name> <zex-app folder location>`
For example:
`helm upgrade --install -n zex --create-namespace zex-app .`

Add `-f <your_yaml_variables_file>` to the `helm` command to customize charts variables.

## Local DNS
An A record needs to be added you your local DNS server od `hosts` file.
By default, this will be `zex-app.local` pointing to `127.0.0.1`, but you can set your own domain with `ingress.hosts`.

## Istio dashboard
Install `addons` to be able to use istio dashboard. The `samples/addons` folder has been [downloaded](https://istio.io/latest/docs/setup/getting-started/#download)
together with the `istio` cli tool.
Run `kubectl apply -f <istio_folder/samples/addons>` and open the dashboard with `istioctl dashboard kiali`.

## Running stress tests
* Backend - `kubectl run -n zex-backend -i --tty load-generator --rm --image=busybox:1.28 --restart=Never -- /bin/sh -c "while sleep 0.01; do wget -q -O- http://zex-app-backend.zex-backend:8080/api/v1/title; done"`
* Frontend - `kubectl run -n zex -i --tty load-generator --rm --image=busybox:1.28 --restart=Never -- /bin/sh -c "while sleep 0.01; do wget -q -O- http://zex-app-frontend:8080; done"`

The `frontend` test will stress all components, as it needs the `backend` which needs the `database`.

## Cleanup
Just destroy your kind cluster with `kind destroy cluster` command.

## Miscellaneous
### Core App
The `core-app` folder contains the small and very primitive `go` code base used for `frontend` and `backend`.
It has been dockerized and published to docker hub where it is publicly available.
The code is here just for reference.

### Scripts
As described above, `scripts` folder contains utility deployment scripts.

### Helm Chart repo folders
These extra folders `core-app` and `scripts` would not normally be a part of the Helm Chart repository.
In this exceptional case, they are placed here just for convenience,
but in reality they would probably have their own dedicated repositories.


## Variables

## General parameters

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| global.additionalLabels | string | `""` | Additional labels for all resources |
| kubeVersionOverride | string | `""` | Override kubernetes version |

## Frontend

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| frontend.affinity | object | `{}` | Frontend component affinity |
| frontend.autoscaling.enabled | bool | `true` | Enable frontend component autoscaling |
| frontend.autoscaling.maxReplicas | int | `10` | Maximum number of frontend component replicas |
| frontend.autoscaling.minReplicas | int | `1` | Minimum number of frontend component replicas |
| frontend.autoscaling.targetCPUUtilizationPercentage | int | `80` | Frontend component target CPU percentage |
| frontend.autoscaling.targetMemoryUtilizationPercentage | int | `80` | Frontend component target Memory utilization percentage |
| frontend.bootCommands | list | `[]` | Additional boot commands |
| frontend.image.pullPolicy | string | `"IfNotPresent"` | Frontend component pull policy |
| frontend.image.repository | string | `"zeljkoiphouse/zex-app"` | Frontend component container repository |
| frontend.image.tag | string | `"latest"` | Frontend component image tag |
| frontend.imagePullSecrets | object | `{}` | Frontend component image pull secrets |
| frontend.ingress.annotations | object | `{"nginx.ingress.kubernetes.io/force-ssl-redirect":"true","nginx.ingress.kubernetes.io/proxy-body-size":"100m","nginx.ingress.kubernetes.io/proxy-connect-timeout":"300","nginx.ingress.kubernetes.io/proxy-read-timeout":"300","nginx.ingress.kubernetes.io/proxy-send-timeout":"300","nginx.ingress.kubernetes.io/ssl-passthrough":"true"}` | The annotations for frontend component |
| frontend.ingress.enabled | bool | `true` | Frontend component enable ingress |
| frontend.ingress.extraPaths | list | `[]` | Extra ingress paths |
| frontend.ingress.hosts | list | `["zex-app.local"]` | A list of ingress hosts |
| frontend.ingress.ingressClassName | string | `"nginx"` | The ingress class name |
| frontend.ingress.labels | object | `{}` | The labels for frontend component |
| frontend.ingress.pathType | string | `"Prefix"` | The ingress path type |
| frontend.ingress.paths | list | `["/"]` | A list of ingress paths |
| frontend.ingress.tls | list | `[]` | A list of ingress TLS configuration |
| frontend.livenessProbe.failureThreshold | int | `3` | Minimum consecutive failures for the [probe] to be considered failed after having succeeded |
| frontend.livenessProbe.initialDelaySeconds | int | `10` | Number of seconds after the container has started before [probe] is initiated |
| frontend.livenessProbe.periodSeconds | int | `10` | How often (in seconds) to perform the [probe] |
| frontend.livenessProbe.successThreshold | int | `1` | Minimum consecutive successes for the [probe] to be considered successful after having failed |
| frontend.livenessProbe.timeoutSeconds | int | `1` | Number of seconds after which the [probe] times out |
| frontend.name | string | `"frontend"` | Frontend component name |
| frontend.nodeSelector | object | `{}` | Frontend component node selector |
| frontend.podAnnotations | object | `{}` | Frontend component pod annotations |
| frontend.podLabels | object | `{}` | Frontend component pod labels |
| frontend.podSecurityContext | object | `{}` | Frontend component pod security context |
| frontend.readinessProbe.failureThreshold | int | `3` | Minimum consecutive failures for the [probe] to be considered failed after having succeeded |
| frontend.readinessProbe.initialDelaySeconds | int | `10` | Number of seconds after the container has started before [probe] is initiated |
| frontend.readinessProbe.periodSeconds | int | `10` | How often (in seconds) to perform the [probe] |
| frontend.readinessProbe.successThreshold | int | `1` | Minimum consecutive successes for the [probe] to be considered successful after having failed |
| frontend.readinessProbe.timeoutSeconds | int | `1` | Number of seconds after which the [probe] times out |
| frontend.replicaCount | int | `1` | Frontend component number of replicas if autoscaling disabled |
| frontend.resources | object | `{"limits":{"cpu":"300m","memory":"128Mi"},"requests":{"cpu":"100m","memory":"64Mi"}}` | Frontend component resources and limits |
| frontend.securityContext | object | `{}` | Frontend component security context |
| frontend.service.annotations | object | `{}` | Frontend component service annotations |
| frontend.service.externalIPs | list | `[]` | The service external IPs |
| frontend.service.externalTrafficPolicy | string | `""` | The service external traffic policy |
| frontend.service.http | object | `{"name":"http","nodePort":8080,"port":8080}` | Frontend component service ports |
| frontend.service.http.name | string | `"http"` | Frontend http service name |
| frontend.service.http.nodePort | int | `8080` | Frontend http service nodePort |
| frontend.service.http.port | int | `8080` | Frontend http service port |
| frontend.service.https.name | string | `"https"` | Frontend https service name |
| frontend.service.https.nodePort | int | `443` | Frontend https service nodePort |
| frontend.service.https.port | int | `443` | Frontend https service port |
| frontend.service.labels | object | `{}` | Frontend component service labels |
| frontend.service.loadBalancerIP | string | `""` | The load balancer IP to use, if service type is `LoadBalancer` |
| frontend.service.loadBalancerSourceRanges | list | `[]` | The IP CIDR to whitelist, if service type is `LoadBalancer` |
| frontend.service.namedTargetPort | bool | `true` | Uses the target port name alias instead of the explicit port |
| frontend.service.sessionAffinity | string | `""` | The session affinity. either `ClientIP` or `None` |
| frontend.service.type | string | `"NodePort"` | Frontend component service type |
| frontend.serviceAccount.annotations | object | `{}` | The annotations for the created service account |
| frontend.serviceAccount.automountServiceAccountToken | bool | `true` | Automount the service account API credentials |
| frontend.serviceAccount.create | bool | `false` | Create a service account for the Frontend component |
| frontend.serviceAccount.serviceAccountName | string | `""` | A service account name override to use |
| frontend.tolerations | object | `{}` | Frontend component tolerations |
| frontend.updateStrategy.maxSurge | int | `1` | Max surge |
| frontend.updateStrategy.maxUnavailable | int | `1` | Max unavailable |
| frontend.updateStrategy.type | string | `"RollingUpdate"` | Strategy type |
| frontend.volumeMounts | object | `{}` | Frontend component volume mounts |
| frontend.volumes | object | `{}` | Frontend component volumes |

## Backend

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| backend.affinity | object | `{}` | Backend component node affinity |
| backend.autoscaling.enabled | bool | `true` | Enable backend autoscaling |
| backend.autoscaling.maxReplicas | int | `10` | Maximum number of replicas |
| backend.autoscaling.minReplicas | int | `1` | Minimum number of replicas |
| backend.autoscaling.targetCPUUtilizationPercentage | int | `80` | Target CPU percentage |
| backend.autoscaling.targetMemoryUtilizationPercentage | int | `80` | Target Memory utilization percentage |
| backend.bootCommands | list | `[]` | Additional boot commands |
| backend.image.pullPolicy | string | `"IfNotPresent"` | Container pull policy |
| backend.image.repository | string | `"zeljkoiphouse/zex-app"` | Container repo and name |
| backend.image.tag | string | `"latest"` | Container version |
| backend.imagePullSecrets | object | `{}` | Backend component image pull secrets |
| backend.livenessProbe.failureThreshold | int | `3` | Minimum consecutive failures for the [probe] to be considered failed after having succeeded |
| backend.livenessProbe.initialDelaySeconds | int | `10` | Number of seconds after the container has started before [probe] is initiated |
| backend.livenessProbe.periodSeconds | int | `10` | How often (in seconds) to perform the [probe] |
| backend.livenessProbe.successThreshold | int | `1` | Minimum consecutive successes for the [probe] to be considered successful after having failed |
| backend.livenessProbe.timeoutSeconds | int | `1` | Number of seconds after which the [probe] times out |
| backend.name | string | `"backend"` | The name of the backend component |
| backend.nodeSelector | object | `{}` | Backend component node selector |
| backend.podAnnotations | object | `{}` | Backend component pod annotations |
| backend.podLabels | object | `{}` | Backend component pod labels |
| backend.podSecurityContext | object | `{}` | Backend pod security context |
| backend.readinessProbe.failureThreshold | int | `3` | Minimum consecutive failures for the [probe] to be considered failed after having succeeded |
| backend.readinessProbe.initialDelaySeconds | int | `10` | Number of seconds after the container has started before [probe] is initiated |
| backend.readinessProbe.periodSeconds | int | `10` | How often (in seconds) to perform the [probe] |
| backend.readinessProbe.successThreshold | int | `1` | Minimum consecutive successes for the [probe] to be considered successful after having failed |
| backend.readinessProbe.timeoutSeconds | int | `1` | Number of seconds after which the [probe] times out |
| backend.replicaCount | int | `1` | Backend component replica count when autoscaling is not enabled |
| backend.resources | object | `{"limits":{"cpu":"500m","memory":"128Mi"},"requests":{"cpu":"250m","memory":"64Mi"}}` | Backend component resource requests and limits |
| backend.securityContext | object | `{}` | Backend component container security context |
| backend.service.annotations | object | `{}` | Backend service annotations |
| backend.service.externalIPs | list | `[]` | The service external IPs |
| backend.service.externalTrafficPolicy | string | `""` | The service external traffic policy |
| backend.service.http.name | string | `"http"` | Backend http service name |
| backend.service.http.port | int | `8080` | Backend http service port |
| backend.service.https.name | string | `"https"` | Backend https service name |
| backend.service.https.port | int | `8081` | Backend https service port |
| backend.service.labels | object | `{}` | Backend component service labels |
| backend.service.loadBalancerIP | string | `""` | The load balancer IP to use, if service type is `LoadBalancer` |
| backend.service.loadBalancerSourceRanges | list | `[]` | The IP CIDR to whitelist, if service type is `LoadBalancer` |
| backend.service.namedTargetPort | bool | `true` | Uses the target port name alias instead of the explicit port |
| backend.service.sessionAffinity | string | `""` | The session affinity. either `ClientIP` or `None` |
| backend.service.type | string | `"ClusterIP"` | Backend component service type |
| backend.serviceAccount.annotations | object | `{}` | The annotations for the created service account |
| backend.serviceAccount.automountServiceAccountToken | bool | `true` | Automount the service account API credentials |
| backend.serviceAccount.serviceAccountName | string | `""` | A service account name override to use |
| backend.tolerations | object | `{}` | Backend node tolerations |
| backend.updateStrategy.maxSurge | int | `1` | Max surge |
| backend.updateStrategy.maxUnavailable | int | `1` | Max unavailable |
| backend.updateStrategy.type | string | `"RollingUpdate"` | Update strategy type |
| backend.volumeMounts | object | `{}` | Backend component volume mounts |
| backend.volumes | object | `{}` | Backend component volumes |

## Database

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| database.affinity | object | `{}` | Database component affinity |
| database.autoscaling.enabled | bool | `false` | Enable autoscaling |
| database.autoscaling.maxReplicas | int | `10` | Minimum number of replicas |
| database.autoscaling.minReplicas | int | `1` | Maximum number of replicas |
| database.autoscaling.targetCPUUtilizationPercentage | int | `90` | Target CPU percentage |
| database.autoscaling.targetMemoryUtilizationPercentage | int | `90` | Target Memory utilization percentage |
| database.customConfig | string | `"[mysqld] \nskip-name-resolve\nexplicit_defaults_for_timestamp\nport={{ .Values.database.service.mysql.port }}\nmax_allowed_packet=16M\nbind-address=*\nslow_query_log=0\nlong_query_time=10.0\n\n[client]\nport={{ .Values.database.service.mysql.port }}\n\n[manager]\nport={{ .Values.database.service.mysql.port }}"` | Custom config file |
| database.databaseName | string | `"zex-app"` | Name of the database that will be created |
| database.databaseUser | string | `"zex-app"` | New database user |
| database.enabled | bool | `true` | Enable local database |
| database.image.pullPolicy | string | `"IfNotPresent"` | Container pull policy |
| database.image.repository | string | `"mariadb"` | Container repo and name |
| database.image.tag | string | `"latest"` | Container version |
| database.imagePullSecrets | object | `{}` | Database component image pull secrets |
| database.livenessProbe.failureThreshold | int | `3` | Minimum consecutive failures for the [probe] to be considered failed after having succeeded |
| database.livenessProbe.initialDelaySeconds | int | `30` | Number of seconds after the container has started before [probe] is initiated |
| database.livenessProbe.periodSeconds | int | `10` | How often (in seconds) to perform the [probe] |
| database.livenessProbe.successThreshold | int | `1` | Minimum consecutive successes for the [probe] to be considered successful after having failed |
| database.livenessProbe.timeoutSeconds | int | `1` | Number of seconds after which the [probe] times out |
| database.name | string | `"database"` | The name of the database component |
| database.nodeSelector | object | `{}` | Database component node selector |
| database.persistence.accessModes | list | `["ReadWriteOnce"]` | Persistence access modes |
| database.persistence.enabled | bool | `true` | Enable database persistence |
| database.persistence.persistentVolumeClaimRetentionPolicy | object | `{}` | Persistent volume claim retention policy |
| database.persistence.selector | object | `{}` | Persistence selector |
| database.persistence.storageClassName | string | `""` | Persistence storage class |
| database.persistence.storageSize | string | `"10Gi"` | Persistence disk size |
| database.podAnnotations | object | `{}` | Database component pod annotations |
| database.podSecurityContext | object | `{}` | Database component pod security context |
| database.readinessProbe.failureThreshold | int | `3` | Minimum consecutive failures for the [probe] to be considered failed after having succeeded |
| database.readinessProbe.initialDelaySeconds | int | `30` | Number of seconds after the container has started before [probe] is initiated |
| database.readinessProbe.periodSeconds | int | `10` | How often (in seconds) to perform the [probe] |
| database.readinessProbe.successThreshold | int | `1` | Minimum consecutive successes for the [probe] to be considered successful after having failed |
| database.readinessProbe.timeoutSeconds | int | `1` | Number of seconds after which the [probe] times out |
| database.replicaCount | int | `1` | Database component number of replicas |
| database.resources | object | `{"limits":{"cpu":"1000m","memory":"256Mi"},"requests":{"cpu":"500m","memory":"128Mi"}}` | Database component resources and requests |
| database.securityContext | object | `{}` | Database component security context |
| database.service.annotations | object | `{}` | Database component service annotations |
| database.service.externalIPs | list | `[]` | The service external IPs |
| database.service.externalTrafficPolicy | string | `""` | The service external traffic policy |
| database.service.labels | object | `{}` | Database component service labels |
| database.service.loadBalancerIP | string | `""` | The load balancer IP to use, if service type is `LoadBalancer` |
| database.service.loadBalancerSourceRanges | list | `[]` | The IP CIDR to whitelist, if service type is `LoadBalancer` |
| database.service.mysql.name | string | `"mysql"` | Database component port name |
| database.service.mysql.nodePort | int | `3306` | Database component nodePort |
| database.service.mysql.port | int | `3306` | Database component port |
| database.service.namedTargetPort | bool | `true` | Uses the target port name alias instead of the explicit port |
| database.service.sessionAffinity | string | `""` | The session affinity. either `ClientIP` or `None` |
| database.service.type | string | `"ClusterIP"` | Database component service type |
| database.serviceAccount.annotations | object | `{}` | The annotations for the created service account |
| database.serviceAccount.automountServiceAccountToken | bool | `true` | Automount the service account API credentials |
| database.serviceAccount.create | bool | `false` | Create a service account for the database component |
| database.serviceAccount.serviceAccountName | string | `""` | A service account name override to use |
| database.tolerations | object | `{}` | Database component tolerations |
| database.volumeMounts | object | `{}` | Database component volume mounts |
| database.volumes | object | `{}` | Database component volumes |

### Bitnami mysql

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| bitnamiMysql.enabled | bool | `false` | Enable bitnami/mysql helm chart dependency as database backend |
| mysql.auth.database | string | `"zex-app"` | Bitnami mysql database name |
| mysql.auth.username | string | `"zex-app"` | Bitnami mysql username |

### Remotely hosted database

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| remoteDb.dbName | string | `""` | Remote database name |
| remoteDb.enabled | bool | `false` | Enable remotely hosted database |
| remoteDb.host | string | `""` | Remote database host |
| remoteDb.password | string | `""` | Remote database password |
| remoteDb.port | string | `""` | Remote database port |
| remoteDb.username | string | `""` | Remote database username |
