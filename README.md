# Benchmarking

Install the operators (jaeger, ES, etc.).

## Elastic Search

### Install ES
```
kubectl apply -f es.yaml
sed -i -E "s/password: .*/password: `kubectl get secret elasticsearch-sample-es-elastic-user -n default -o jsonpath="{.data.elastic}" | base64 --decode`/" jaeger-es.yaml
kubectl apply -f jaeger-es.yaml
```

### Install Jaeger for ES
```
kubectl apply -f jaeger-es.yaml
```

After Jaeger created all the pods, set the replicas to 0 for the `jaeger-operator`.

Manually these args to the `jaeger-collector` pod in the `simple-prod-collector` deployment:
```
- '--collector.num-workers=1000'
- '--collector.queue-size=100000'
```

### Start Tracegen
```
kubectl delete -f tracegen/synthetic-load-gen.yaml
```

### Get service hostname
```
kubectl get svc simple-prod-query -o json | jq -r '.status.loadBalancer.ingress[0].hostname'
```

### Run query benchmark
Update the URL in `main.go` with the response from the query above
```
cd go-get-latency
go run .
```

## Promscale

### Install tobs
```
kubectl create ns tobs
helm upgrade --install tobs timescale/tobs -f helm/values.yaml -n tobs
```

### Install Jaeger for Promscale
```
kubectl apply -f jaeger-promscale.yaml
```

After Jaeger created all the pods, set the replicas to 0 for the `jaeger-operator`.

Manually env vars to the `jaeger-collector` and `jaeger-query` pod in the `simple-prod-collector` deployment and `simple-prod-query` deployment respectively:
```
- name: GRPC_STORAGE_SERVER
  value: tobs-promscale.tobs:9202
```

Manually these args to the `jaeger-collector` pod in the `simple-prod-collector` deployment:
```
- '--collector.num-workers=1000'
- '--collector.queue-size=100000'
```

### Start Tracegen
```
kubectl delete -f tracegen/synthetic-load-gen.yaml
```

### Get service hostname
```
kubectl get svc simple-prod-query -o json | jq -r '.status.loadBalancer.ingress[0].hostname'
```

### Run query benchmark
Update the URL in `main.go` with the response from the query above
```
cd go-get-latency
go run .
```