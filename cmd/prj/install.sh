# Install service chart
helm install prj ./resources/chart/

# Add repos for Helm charts
helm repo add stable https://charts.helm.sh/stable
#helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Install nginx ingress controller
helm install ingress-nginx bitnami/nginx-ingress-controller \
  --namespace ingress-nginx \
  --create-namespace \
  -f ./resources/nginx_controller.yaml

# Install RabbitMQ
helm install rabbitmq bitnami/rabbitmq \
  --namespace prj \
  -f ./resources/rabbitmq.yaml

# Install Postgres exporter
#helm install prom-postgres prometheus-community/prometheus-postgres-exporter \
#  --namespace hw3 \
#  -f ./resources/postgres_exporter.yaml

# Install prometheus stack
#helm install prom prometheus-community/kube-prometheus-stack \
#  --namespace monitoring \
#  -f ./resources/prometheus.yaml

# Import Grafana dashboard from ConfigMap
#kubectl apply -f ./resources/grafana_dashboard.yaml
