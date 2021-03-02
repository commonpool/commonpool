istioctl operator init

kubectl apply -f - <<EOF
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  namespace: istio-system
  name: example-istiocontrolplane
spec:
  profile: demo
  meshConfig:
    accessLogFile: /dev/stdout
  components:
    ingressGateways:
      - name: istio-ingressgateway
        k8s:
          service:
            ports:
              - name: status-port
                port: 15021
                targetPort: 15021
              - name: http2
                port: 80
                targetPort: 8080
              - name: https
                port: 443
                targetPort: 8443
              - name: tcp
                port: 31400
                targetPort: 31400
              - name: tls
                port: 15443
                targetPort: 15443
              - name: neo4j-http
                port: 7474
                targetPort: 7474
              - name: bolt-http
                port: 7687
                targetPort: 7687
              - name: bolt-tcp
                port: 7686
                targetPort: 7686
              - name: postgres
                port: 5432
                targetPort: 5432
              - name: amqp
                port: 5672
                targetPort: 5672
EOF

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

kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.2.0/cert-manager.yaml
kubectl apply -f https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml
kubectl apply -f cert.yaml
kubectl apply -f istio.yaml
kubectl apply -f neo4j.yaml
kubectl apply -f pg.yaml
kubectl apply -f rabbit.yaml

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

echo Commonpool Minikube platform started. Press CTRL+C to exit.

cleanup(){
  sudo sed -i '/^.* graphdb.commonpool.dev$/d' /etc/hosts
  sudo sed -i '/^.* 0.graphdb.commonpool.dev$/d' /etc/hosts
  sudo sed -i '/^.* 1.graphdb.commonpool.dev$/d' /etc/hosts
  sudo sed -i '/^.* 2.graphdb.commonpool.dev$/d' /etc/hosts
  sudo sed -i '/^.* db.commonpool.dev$/d' /etc/hosts
  sudo sed -i '/^.* rabbit.commonpool.dev$/d' /etc/hosts
  sudo sed -i '/^.* amqp.commonpool.dev$/d' /etc/hosts
}

trap cleanup SIGTERM
trap cleanup SIGINT

while true
do
  sleep 1
done

