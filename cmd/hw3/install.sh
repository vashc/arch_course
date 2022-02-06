# Install service chart
helm install hw3 ./resources/chart/

# Install Nginx ingress chart
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx \
  --set controller.metrics.enabled=true \
  --set-string controller.podAnnotations."prometheus\.io/scrape"="true" \
  --set-string controller.podAnnotations."prometheus\.io/port"="10254" \
  -f ./resources/nginx_controller.yaml

# Add repo for a prometheus stack
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Install Postgres exporter
helm install prom-postgres prometheus-community/prometheus-postgres-exporter -n hw3 -f ./resources/postgres_exporter.yaml

# Install prometheus stack
helm install prom prometheus-community/kube-prometheus-stack -n monitoring -f resources/prometheus.yaml
