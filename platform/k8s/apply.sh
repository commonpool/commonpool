istioctl operator init

kubectl label namespace default istio-injection=enabled --overwrite
kubectl apply -f istio.yaml
kubectl apply -f istio-controlplane.yaml
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.2.0/cert-manager.yaml
kubectl apply -f https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml
kubectl apply -f cert.yaml
kubectl apply -f pg.yaml
kubectl apply -f rabbit.yaml
kubectl apply -f redis.yaml
kubectl apply -f frontend.yaml
kubectl apply -f backend.yaml
