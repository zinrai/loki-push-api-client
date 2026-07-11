# Grafana Loki Push API Client

This is a simple client program written in Go that interacts with the Grafana Loki Push API.

It periodically sends log entries to a Grafana Loki instance using the Push API endpoint.

## Scope

This client sends log entries directly to Loki, without promtail or Grafana Alloy in the path. Its purpose is to verify that Loki's own multi-tenant separation works, so the client sets the `X-Scope-OrgID` header itself and does not exercise how an agent assigns tenants. For verifying the agent to Loki pipeline, use a promtail or Alloy based setup instead.

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

## Usage

Point the client at a Loki endpoint and choose how many tenants and labels to spread across:

```
$ loki-push-api-client -endpoint http://localhost:3100/loki/api/v1/push -tenants 5 -labels 5 -interval 30s
tenants: tenant1 tenant2 tenant3 tenant4 tenant5
labels: label1 label2 label3 label4 label5
2026/07/12 08:26:16 Response Status: 204 No Content, X-Scope-OrgID: tenant2
2026/07/12 08:26:17 Response Status: 204 No Content, X-Scope-OrgID: tenant1
```

Tenant and label names are generated from the counts and printed at startup, so you know which `X-Scope-OrgID` values to query with logcli. Run with `-help` for the full list of options.

Each pushed log line embeds the tenant it was sent to and a run global sequence number, so a per tenant query can confirm that only that tenant's lines are visible under it:

```json
{
    "streams": [
        {
            "stream": {
                "label": "label2"
            },
            "values": [
                [
                    "1783812376808899720",
                    "tenant=tenant2 seq=1 MyestiPglkU2hjCkxhf2WZkjwhESdd"
                ],
                [
                    "1783812376808904254",
                    "tenant=tenant2 seq=2 7jnBLlviQze8x7VyEpPnA4p9iIvuBW"
                ]
            ]
        }
    ]
}
```

## License

This project is licensed under the [MIT License](./LICENSE).
