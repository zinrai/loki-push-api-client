# Grafana Loki Push API Client

This is a simple client program written in Go that interacts with the Grafana Loki Push API.

It periodically sends log entries to a Grafana Loki instance using the Push API endpoint.

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

Here is an example `config.yaml` file:

```yaml
labels:
  - val1
  - val2
  - val3
tenants:
  - tenant1
  - tenant2
  - tenant3
endpoint: "http://example.com/loki/api/v1/push"
```

## Usage

```
$ go run main.go
2024/03/22 08:10:05 Response Status: 200 OK, X-Scope-OrgID: tenant3
2024/03/22 08:10:08 Response Status: 200 OK, X-Scope-OrgID: tenant3
2024/03/22 08:10:11 Response Status: 200 OK, X-Scope-OrgID: tenant5
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

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
