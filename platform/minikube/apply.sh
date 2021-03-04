istioctl operator init

set +e
minikube tunnel --cleanup >/dev/null 2>&1 &
set -e

kubectl label namespace default istio-injection=enabled --overwrite

kubectl get secret -n cert-manager ca-key ||
  openssl req \
    -x509 \
    -new \
    -nodes \
    -subj "/C=CA/ST=QC/L=MTL/O=Commonpool/OU=DevOps/CN=www.commonpool.dev/emailAddress=dev@www.commonpool.dev" \
    -keyout myCA.key \
    -out myCA.pem &&
  key=$(cat myCA.key | base64 -w0) &&
  cert=$(cat myCA.pem | base64 -w0) &&
  cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: ca-key
  namespace: cert-manager
data:
  tls.crt: ${cert}
  tls.key: ${key}
EOF

kubectl apply -f istio.yaml
kubectl apply -f istio-controlplane.yaml
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.2.0/cert-manager.yaml
kubectl apply -f https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml
kubectl apply -f cert.yaml
kubectl apply -f neo4j.yaml
kubectl apply -f pg.yaml
kubectl apply -f rabbit.yaml
kubectl apply -f redis.yaml

loadBalancerIP=""
while [ -z $loadBalancerIP ]
do
  echo Waiting for load balancer IP
  loadBalancerIP=$(kubectl get svc -n istio-system istio-ingressgateway --template="{{range .status.loadBalancer.ingress}}{{.ip}}{{end}}")
  sleep 1
done

echo Load Balancer IP : $loadBalancerIP

replace(){
  grep -q ".* $1$" /etc/hosts && sudo sed -i "s/.* $1$/$2 $1/" /etc/hosts || echo "$2 $1" | sudo tee -a /etc/hosts
}

echo Fixing hosts file

replace graphdb.commonpool.dev ${loadBalancerIP}
replace 0.graphdb.commonpool.dev ${loadBalancerIP}
replace 1.graphdb.commonpool.dev ${loadBalancerIP}
replace 2.graphdb.commonpool.dev ${loadBalancerIP}
replace db.commonpool.dev ${loadBalancerIP}
replace rabbit.commonpool.dev ${loadBalancerIP}
replace amqp.commonpool.dev ${loadBalancerIP}
replace redis.commonpool.dev ${loadBalancerIP}

echo Commonpool Minikube platform started. Press CTRL+C to exit.

cleanup(){
  sudo sed -i '/^.* graphdb.commonpool.dev$/d' /etc/hosts
  sudo sed -i '/^.* 0.graphdb.commonpool.dev$/d' /etc/hosts
  sudo sed -i '/^.* 1.graphdb.commonpool.dev$/d' /etc/hosts
  sudo sed -i '/^.* 2.graphdb.commonpool.dev$/d' /etc/hosts
  sudo sed -i '/^.* db.commonpool.dev$/d' /etc/hosts
  sudo sed -i '/^.* rabbit.commonpool.dev$/d' /etc/hosts
  sudo sed -i '/^.* amqp.commonpool.dev$/d' /etc/hosts
  sudo sed -i '/^.* redis.commonpool.dev$/d' /etc/hosts
}

trap cleanup SIGTERM
trap cleanup SIGINT

while true
do
  sleep 1
done

