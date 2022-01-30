# Setup
```helm install hw2 ./resources/chart/```

All resources are automatically installed within Kubernetes namespace named `hw2`.

# Test cases

All test scenarios could be acquired via postman collection which is located in the `/tests` folder.

# Models
`user`:
```json
{
  "id":         integer(int64),
  "username":   string,
  "firstName":  string,
  "lastName":   string,
  "email":      string,
  "phone":      string
}
```

`response`:
```json
{
  "code":    int,
  "message": string
}
```

# Application description
There are five endpoints that application can handle:
* GET `/health` is used for healthcheck and returns `{"status": "OK"}`
* POST `/user` is used for user creation and returns response with "`user created`" message and status 200
* GET `/user/<user_id>` is used for user information acquiring and returns response with user data in JSON format
* PUT `/user/<user_id>` is used for user information updating and returns response with "`user updated`" message and status 200
* DELETE `/user/<user_id>` is used for user deletion and returns response with "`user deleted`" message and status 200

All endpoints return response with "`internal server error`" message and status 500 in case there were any errors during request processing.

# Resources description
Application resources include the following kinds:
* `Namespace`, which is applied first;
* `ConfigMap` + `Secret`, that are used for application configuration and storing DB credentials;
* `Deployment` + `Service` for application;
* `Ingress`, which can be turned off;
* `Service` + `StatefulSet` for Postgres;
* `Job` for migration, which creates `user` table in case it doesn't exist already and populates it with one record (`id=1`).

All resources are parameterized and should be installed via Helm.
