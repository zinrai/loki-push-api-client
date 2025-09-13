# Grafana Loki Push API Client

This is a simple client program written in Go that interacts with the Grafana Loki Push API.

It periodically sends log entries to a Grafana Loki instance using the Push API endpoint.

## Tested with Loki Version

This client has been tested with Grafana Loki version 2.9.5.

Test was performed using Loki in scalable mode installed with helm.

- https://grafana.com/docs/loki/latest/setup/install/helm/install-scalable/
- https://grafana.com/docs/loki/latest/get-started/deployment-modes/#simple-scalable

## Requirements

- Go programming language installed on your system
- Access to a Grafana Loki instance with Push API enabled
    - https://grafana.com/docs/loki/latest/reference/api/#push-log-entries-to-loki
- Access to a Grafana Loki instance with `auth_enabled: true`
    - https://grafana.com/docs/loki/latest/configure/#supported-contents-and-default-values-of-lokiyaml

## Configuration

The configuration for this client is done through a YAML file named `config.yaml`.

You need to provide the following configuration parameters:

- `labels`: An array of strings representing the labels to be associated with the log entries. These labels will be randomly selected for each log entry.
- `tenants`: An array of strings representing the X-Scope-OrgID values to be used in the HTTP request headers. One of these values will be randomly selected for each request.
- `endpoint`: The URL of the Grafana Loki Push API endpoint.
- `sleep_interval`: interval in seconds between POST to the Grafana Loki Push API endpoint.

Here is an example `config.yaml` file:

```yaml
labels:
  - val1
  - val2
  - val3
  - val4
  - val5
tenants:
  - tenant1
  - tenant2
  - tenant3
  - tenant4
  - tenant5
endpoint: "http://example.com/loki/api/v1/push"
sleep_interval: 30
```

## Usage

```
$ go run main.go
2024/03/22 00:35:12 Response Status: 204 No Content, X-Scope-OrgID: tenant3
2024/03/22 00:35:42 Response Status: 204 No Content, X-Scope-OrgID: tenant2
2024/03/22 00:36:12 Response Status: 204 No Content, X-Scope-OrgID: tenant1
2024/03/22 00:36:42 Response Status: 204 No Content, X-Scope-OrgID: tenant1
2024/03/22 00:37:12 Response Status: 204 No Content, X-Scope-OrgID: tenant2
2024/03/22 00:37:42 Response Status: 204 No Content, X-Scope-OrgID: tenant1
2024/03/22 00:38:12 Response Status: 204 No Content, X-Scope-OrgID: tenant1
2024/03/22 00:38:42 Response Status: 204 No Content, X-Scope-OrgID: tenant5
```

Here is an example JSON data:

```json
{
    "streams": [
        {
            "stream": {
                "label": "val5"
            },
            "values": [
                [
                    "1711062608762091837",
                    "8h9kg44ENf58CgVMaisAosF8uyhhBz"
                ],
                [
                    "1711062608762091837",
                    "I2vyItocBaq03NLCmLxybyA9tXSahu"
                ]
            ]
        }
    ]
}
```

## License

This project is licensed under the [MIT License](./LICENSE).
