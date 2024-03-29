# log2sqs

**WARNING: The configuration file has transitioned to YAML. Please see the example configuration file and the notes below.**

This application facilitates distributed log collection and parsing into Graylog GELF via an AWS SQS queue. Communication with SQS is over HTTPS, providing a relatively secure log collection mechanism.

If you are installing log2sqs on an AWS EC2 VM, providing access to SQS via an IAM role assigned to the EC2 instance is recommended.

For logging from non-AWS environments, an AWS IAM key must be added to the configuration file. The policy assigned
to the IAM user should only allow listing queues (ListQueues) and sending messages (SendMessage) to the required SQS
queue.

The following policy example provides the minimum permissions required to locate and add messages to the SQS Queue.
(Replace the region, account number, and queue name with appropriate values.)

```
{
   "Version": "2012-10-17",
   "Statement": [
      {
         "Sid": "VisualEditor0",
         "Effect": "Allow",
         "Action": "sqs:SendMessage",
         "Resource": "arn:aws:sqs:ca-central-1:88888888:GELF"
      },
      {
         "Sid": "VisualEditor1",
         "Effect": "Allow",
         "Action": "sqs:ListQueues",
         "Resource": "*"
      }
   ]
}
```

An open-source companion application to read the SQS queue and send events to Graylog is available at 

https://github.com/tenebris-tech/sqs2gl

log2sqs can:

- Read one or more log files in real-time (like tail) and forward them in GELF to an AWS SQS queue.

- Receive RFC5424 and RFC3164 compliant syslog messages via UDP, parse them, and forward them to the
  AWS SQS queue in GELF. If the type of syslog message can not be identified, the entire message is sent as text.
  If a received syslog message contains a valid GELF message, the GELF message is extracted and the syslog header
  discarded. This allows sending GELF messages by leveraging standard syslog mechanisms.

- Optionally add AWS EC2 instance metadata (instance ID, hostname, and tags) to each event.

- Optionally override the host name and/or add a site name to each log entry (see log2sqs.conf).

The following log file formats are currently supported:

| Format Specifier     | Description                                                     |
|----------------------|-----------------------------------------------------------------|
| gelf                 | Graylog GELF format messages (in JSON)                          |
| error                | Apache2 error log                                               |
| combined             | Apache2/NGINX combined log format                               |
| combinedplus         | Apache2 log format with additional fields                       |
| combinedplusvhost    | Apache2 log format with vhost information and additional fields |
| combinedloadbalancer | Apache2 log format with load balancer info, etc.                |
| text                 | Plain text, not parsed                                          |

Log file format specifiers are case-insensitive.

For the combinedplus format, the following Apache definition is used to add the time (in microseconds) required to process the request and break the request into method, path, and query components:

```
LogFormat "%h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\" %D \"%m\" \"%U\" \"%q\"" combinedplus
```

For the combinedplusvhost format, the following Apache definition is used to add vhost information:

```
LogFormat "%v:%p %h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\" %D \"%m\" \"%U\" \"%q\"" combinedplusvhost
```

The combinedloadbalancer format includes most of the fields above, but also includes the X-Forwarded information so that the original client IP address and protocol can be determined when a load balancer is used:

```
LogFormat "%{X-Forwarded-Proto}i %{Host}i:%{X-Forwarded-Port}i %v:%p %{X-Forwarded-For}i %h %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\" %D \"%m\" \"%U\" \"%q\"" combinedloadbalancer
```

**As of version 0.7.0, the preferred configuration file format is YAML and the default configuration filename is log2sqs.yaml. If a filename ending in .conf is specified, it will be parsed using the legacy format for backward compatibility. However, this will be deprecated in a future version.**

**User-defined regex-based parsing formats can be added to the YAML-format configuration file.**

### Command Line Arguments

log2sqs now supports the following command line arguments:

​	`-config <configuration file path and name>`

​	`-ingest <file>,<format>`

​	`-dryrun`

The ingest argument is designed for testing and for edge cases in which an existing log file must be ingested in its entirety (unlike the default behaviour of starting a tail at the end of an existing file).

Dryrun will stop anything (the ingested file and any other files specified in the config file) from being sent to SQS and will turn on a JSON pretty-print of the GELF message that would have otherwise been sent to SQS. Note that this feature does not currently change syslog processing. This is intended for interactive testing with log files only.

### Development Status

This is a beta release and should be thoroughly tested prior to use in production environments.

### Linux Installation

1) Clone the repo and compile using "go build"
2) Copy the binary (log2sqs) and config file (log2sqs.yaml) to /opt/log2sqs. If you put it elsewhere, you will need to
   update the .service file.
3) Ensure that log2sqs has the execute bit set (i.e. chmod 700 or 755)
4) Set restrictive permissions on log2sqs.yaml, especially if it contains an AWS key (i.e. chmod 600)
5) Copy log2sqs.service to /etc/systemd/system/
6) Update the User and Group in log2sqs.service if you do not wish to run as root
7) Update the configuration file (log2sqs.conf)
8) Run `systemctl daemon-reload`
9) Run `systemctl enable log2sqs` to configure automatic start at boot
10) Run `systemctl start log2sqs` to start the application
11) Configure logrotated to locate the log file specified in log2sqs.yaml. A sample rotate file is contained in
    log2sqs.txt. It can be copied to /etc/logrotate.d/log2sqs. It should be owned by root and set to mode 0644.

### Windows Installation

The application has been tested on Windows, but there is currently no demand for it to run as a Windows service. If
you have a requirement, please open an issue to discuss.

### Copyright

Copyright (c) 2021-2023 Tenebris Technologies Inc. All rights reserved.

Please see the LICENSE file for additional information.
