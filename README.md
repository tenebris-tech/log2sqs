# log2sqs

This application facilitates distributed log collection in Graylog GELF via an AWS SQS queue. Communication with
SQS is over HTTPS, providing a relatively secure log collection mechanism.

For use on an AWS EC2 VM, providing access to SQS via an IAM role assigned to the EC2 instance is recommended.

For logging from non-AWS environments, an AWS IAM key must be added to the configuration file. The policy assigned
to the IAM user should only allow listing queues and write access to the single required SQS queue.

An open source application to read the SQS queue and send events to Graylog is available at
https://github.com/tenebris-tech/sqs2gl.

This application can:

- Read one or more log files in real time (like tail) and forward them in GELF to an AWS SQS queue.

- Receive RFC5424 and RFC3164 compliant syslog messages via UDP, parse them, and forward them to the 
AWS SQS queue in GELF. If the type of syslog message can not be identified, or if parsing fails, the
message is forwarded as-is. Parsing errors are reported via GELF to assist in debugging (this is a work in progress).

- Optionally add AWS EC2 instance metadata (instance ID, hostname, and tags) to each event.

The following log file formats are currently supported:

| Format Specifier | Description                               |
|------------------|-------------------------------------------|
| gelf             | Graylog GELF format messages (in JSON)    |
| error            | Apache2 error log                         |
| combined         | Apache2/NGINX combined log format         |
| combinedplus     | Apache2 log format with additional fields |
| text             | Plain text, not parsed                    |

For the combinedplus format, the following Apache definition is used
to add the time (in microseconds) required to process the request and
break the request into method, path, and query components.

```
LogFormat "%h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\" %D \"%m\" \"%U\" \"%q\"" combinedplus
```

### Development Status

This is a beta release.

### Linux Installation

1) Clone the repo and compile using "go build"
2) Copy the binary (log2sqs) and config file (log2sqs.conf) to /opt/log2sqs. If you put it elsewhere you will need to
   update the .service file.
3) Ensure that log2sqs has the execution bit set (i.e. chmod 700 or 755)
4) Copy log2sqs.service to /etc/systemd/system/
5) Update the User and Group in log2sqs.service if you do not wish to run as root
6) Update the configuration file (log2sqs.conf)
7) Run 'systemctl daemon-reload'
8) Run 'systemctl start log2sqs' to start the application
9) Configure logrotated to locate the log file specified in log2sqs.conf. A sample rotate file is contained in
   log2sqs.txt. It can be copied to /etc/logrotate.d/log2sqs. It should be owned by root and set to mode 0644.

### Windows Installation

The application has been tested on Windows, but there is currently no demand for it to run as a Windows service.

### Copyright

Copyright (c) 2021-2023 Tenebris Technologies Inc. All rights reserved.

Please see the LICENSE file for additional information.
