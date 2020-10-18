
.PHONY: help
help:
	@echo 
	@echo 'Makefile to manage aws platform dependencies'
	@echo
	@echo -e "Usually just \e[36mmake install\033[0m will do the trick"
	@echo 
	@echo Usage:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36mmake %-30s\033[0m %s\n", $$1, $$2}'
	@echo
	@echo 'Prerequisites:'
	@echo '    helm'
	@echo '    aws-cli'
	@echo '    kubectl'
	@echo '    jq'
	@echo '    istioctl'
	@echo


master-dns-policy: ## Creates the master-dns-policy IAM policy that allows kops external-dns to update route53 records, load balancers, etc.
	-aws iam create-policy --policy-name master-dns-policy --policy-document file://policies/master-policy.json

cert-manager-policy: ## Creates the cert-manager-policy that allows cert-manager to update route53 records and resolve dns challenges
	-aws iam create-policy --policy-name cert-manager-policy --policy-document file://policies/cert-manager-policy.json

install: cert-manager-install istio-operator-install istio-controlplane-install cluster-issuer-install kops-install ## Installs platform dependencies

kops-install: master-dns-policy cert-manager-policy
	set -e; \
	P1=$$(aws iam list-policies --query 'Policies[?PolicyName==`cert-manager-policy`].Arn' --output text); \
	echo p1 = $$P1; \
	P2=$$(aws iam list-policies --query 'Policies[?PolicyName==`master-dns-policy`].Arn' --output text); \
	echo p2 = $$P2; \
	kops get cluster -o json \
	| jq ".spec.externalPolicies.nodes |= (. += [\"$$P1\"] | unique)" \
	| jq ".spec.externalPolicies.master |= (. += [\"$$P2\"] | unique)" \
	| kops replace -f -
	kops update cluster --yes 
	kops rolling-update cluster --yes

cert-manager-install: ## Installs cert-manager on a kubernetes cluster using helm
	-kubectl create namespace cert-manager
	helm repo add jetstack https://charts.jetstack.io
	helm repo update
	helm get values cert-manager --namespace cert-manager \
	|| helm install \
  		cert-manager jetstack/cert-manager \
  		--namespace cert-manager \
  		--version v1.0.3 \
  		--set installCRDs=true

istio-operator-install: ## Installs istio-operator on the k8s cluster
	istioctl operator init

istio-controlplane-install: ## Installs istio controlplane on the k8s cluster
	kubectl apply -f istio-controlplane.yaml

cluster-issuer-install: ## Installs cert-manager ClusterIssuer on the k8s cluster
	kubectl apply -f cert-issuer.yaml

.DEFAULT_GOAL := help