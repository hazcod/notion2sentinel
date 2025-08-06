# notion2sentinel

A Go program that exports Notion audit events to Microsoft Sentinel SIEM.

## Running

Get a [Notion API Key](https://www.notion.so/profile/integrations/form/new-integration). Also make note of your Notion Organization ID.

First create a yaml file, such as `config.yml`:
```yaml
log:
  level: INFO

microsoft:
  app_id: ""
  secret_key: ""
  tenant_id: ""
  subscription_id: ""
  resource_group: ""
  workspace_name: ""

  dcr:
    endpoint: ""
    rule_id: ""
    stream_name: ""

  expires_months: 6

notion:
  api_token: ""
  lookback: 168h
  organisation_id: ""
```

And now run the program from source code:
```shell
% make
go run ./cmd/... -config=dev.yml
INFO[0000] shipping logs                                 module=sentinel_logs table_name=NotionAuditLogs total=82
INFO[0002] shipped logs                                  module=sentinel_logs table_name=NotionAuditLogs
INFO[0002] successfully sent logs to sentinel            total=82
```

Or binary:
```shell
% notion2sen -config=config.yml
```

## Building

```shell
% make build
```
