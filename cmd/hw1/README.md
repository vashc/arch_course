# Setup
```kubectl apply -f ./resources```

# Test cases

* `curl arch.homework/health`
* `curl arch.homework/health/`
* `curl arch.homework/otusapp/such_a_student/whatever`
* `curl arch.homework/otusapp/alp4num3r1c/`

# Application description
There are two endpoints that application can handle:
* GET `/health` returns `{"status": "OK"}`
* GET `/student/<name>` returns greeting `Hello, <name>`

# Ingress description
Internally ingress uses configuration-snippet annotation with Nginx rewrites.\
It is possible to make two ingress resources with and without rewrite annotation though.

The forwarding is as follows:\
`^/health/?$` -> `/health`\
`^/otusapp/(\w+)/.*` -> `/student/<name>`


The rest of requests return HTTP 404.