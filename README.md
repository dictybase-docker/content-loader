# content-loader

```
NAME:
   content-loader - cli for managing serialized json content in the database

USAGE:
   go-content-loader [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     import   import serialized json content in upsert mode
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --hooks value             hook names for sending log in addition to stderr
   --log-level value         log level for the application (default: "error")
   --log-format value        format of the logging out, either of json or text (default: "json")
   --slack-channel value     Slack channel where the log will be posted [$SLACK_CHANNEL]
   --slack-url value         Slack webhook url[required if slack channel is provided] [$SLACK_URL]
   --content-api-host value  content api host address (default: "content-api") [$CONTENT_API_SERVICE_HOST]
   --content-api-port value  content api port [$CONTENT_API_SERVICE_PORT]
   --help, -h                show help
   --version, -v             print the version
```
## Subcommand
```
NAME:
   go-content-loader import - import serialized json content in upsert mode

USAGE:
   go-content-loader import [command options] [arguments...]

OPTIONS:
   --s3-host value                   S3 server host (default: "minio") [$MINIO_SERVICE_HOST]
   --s3-port value                   S3 server port [$MINIO_SERVICE_PORT]
   --s3-bucket value                 S3 bucket where the import data is kept (default: "content")
   --access-key value, --akey value  access key for S3 server, required based on command run [$S3_ACCESS_KEY]
   --secret-key value, --skey value  secret key for S3 server, required based on command run [$S3_SECRET_KEY]
   --remote-path value, --rp value   full path(relative to the bucket) of s3 folder which will be used for import[required]
   --extension value, -e value       Extension of the file(s) that will be matched in the remote folder (default: "json")
   --namespace value, -n value       namespace under which the content will be imported
   --email value                     Email of user who will either create or update the content
   --user-api-host value             user api host address (default: "user-api") [$USER_API_SERVICE_HOST]
   --user-api-port value             user api port [$USER_API_SERVICE_PORT]
```
