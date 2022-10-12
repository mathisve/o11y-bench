### to get it back to a “good” state:

```
kubectl delete -f tracegen.yaml,jaeger-es.yaml,es.yaml
kubectl apply -f es.yaml
sed -i -E "s/password: .*/password: `kubectl get secret elasticsearch-sample-es-elastic-user -n default -o jsonpath="{.data.elastic}" | base64 --decode`/" jaeger-es.yaml
kubectl apply -f jaeger-es.yaml
```

Then you want to edit tracegen.yaml to either point to http://simple-prod-collector:14268/api/traces or 
http://tobs-opentelemetry-collector.tobs:14268/api/traces

Then apply that too

### port forwarding:
`kubectl port-forward svc/quickstart-kb-http 5601` will get you Kibana, good to see if data ingestion via Jaeger 
is even working
`kubectl port-forward svc/tobs-grafana 3000:80`  will get you grafana (password can be found via 
`kubectl get secret --namespace tobs tobs-grafana -o jsonpath="{.data.admin-password}" | base64 --decode`

Promscale logs will show you data flowing through there

That Grafana will automatically get both Jaeger and Promscale metrics - so you can graph the throughput

I would probably leave the tobs schema as is, and then create another collector/promscale/postgres (without 
prometheus or node-exporter etc…) to use for the Jaeger backend

I think you’ll also need to inject pod-anti-affinity rules (maybe just manually?) to make sure we get valid 
results

Any Q’s ask me early, always here to help (and love this stuff too so keen to help)


            - name: GRPC_STORAGE_SERVER
              value: tobs-promscale.tobs:9202