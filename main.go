package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "content-loader"
	app.Usage = "cli for managing serialized json content in the database"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "hooks",
			Usage: "hook names for sending log in addition to stderr",
			Value: &cli.StringSlice{},
		},
		cli.StringFlag{
			Name:  "log-level",
			Usage: "log level for the application",
			Value: "error",
		},
		cli.StringFlag{
			Name:  "log-format",
			Usage: "format of the logging out, either of json or text",
			Value: "json",
		},
		cli.StringFlag{
			Name:   "slack-channel",
			EnvVar: "SLACK_CHANNEL",
			Usage:  "Slack channel where the log will be posted",
		},
		cli.StringFlag{
			Name:   "slack-url",
			EnvVar: "SLACK_URL",
			Usage:  "Slack webhook url[required if slack channel is provided]",
		},
		cli.StringFlag{
			Name:   "content-api-host",
			EnvVar: "CONTENT_API_SERVICE_HOST",
			Usage:  "content api host address",
			Value:  "content-api",
		},
		cli.StringFlag{
			Name:   "content-api-port",
			EnvVar: "CONTENT_API_SERVICE_PORT",
			Usage:  "content api port",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "import",
			Usage:  "import serialized json content in upsert mode",
			Action: importAction,
			Before: validateImport,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "s3-host",
					Usage:  "S3 server host",
					EnvVar: "MINIO_SERVICE_HOST",
					Value:  "minio",
				},
				cli.StringFlag{
					Name:   "s3-port",
					Usage:  "S3 server port",
					EnvVar: "MINIO_SERVICE_PORT",
				},
				cli.StringFlag{
					Name:  "s3-bucket",
					Usage: "S3 bucket where the import data is kept",
					Value: "content",
				},
				cli.StringFlag{
					Name:   "access-key, akey",
					EnvVar: "S3_ACCESS_KEY",
					Usage:  "access key for S3 server, required based on command run",
				},
				cli.StringFlag{
					Name:   "secret-key, skey",
					EnvVar: "S3_SECRET_KEY",
					Usage:  "secret key for S3 server, required based on command run",
				},
				cli.StringFlag{
					Name:  "remote-path, rp",
					Usage: "full path(relative to the bucket) of s3 folder which will be used for import[required]",
				},
				cli.StringFlag{
					Name:  "extension, e",
					Usage: "Extension of the file(s) that will be matched in the remote folder",
					Value: "json",
				},
				cli.StringFlag{
					Name:  "namespace, n",
					Usage: "namespace under which the content will be imported",
				},
				cli.StringFlag{
					Name:  "email",
					Usage: "Email of user who will either create or update the content",
				},
				cli.StringFlag{
					Name:   "user-api-host",
					EnvVar: "USER_API_SERVICE_HOST",
					Usage:  "user api host address",
					Value:  "user-api",
				},
				cli.StringFlag{
					Name:   "user-api-port",
					EnvVar: "USER_API_SERVICE_PORT",
					Usage:  "user api port",
				},
			},
		},
	}
	app.Run(os.Args)
}
