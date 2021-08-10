# log2sqs

This application reads Graylog GELF-format messages from a text file and forwards them to an AWS SQS queue.

It will optionally add AWS EC2 instance metadata (instance ID, hostname, and tags) to each log message.

### Development Status

This is an alpha release.

### Installation

1) Clone the repo and compile using "go build"
2) Copy the binary (log2sqs) and config file (log2sqs.conf) to /opt/log2sqs. If you put it elsewhere you will need to
   update the .service file.
3) Ensure that log2sqs has the execution bit set (i.e. chmod 700 or 755)
4) Copy log2sqs.service to /etc/systemd/system/
5) Update the User and Group in log2sqs.service if you do not wish to run as root
6) Update the configuration file (log2sqs.conf)
7) Run 'systemctl daemon-reload'
8) Run 'systemctl start log2sqs' to start the application
9) Configure logrotated to locate the log file specified in log2sqs.conf

### Apache2 Configuration

To configure Apache2 to log in GELF format, add the following line to apache2.conf after other LogFormat lines:

```
LogFormat "{ \"version\": \"1.1\", \"host\": \"%V\", \"short_message\": \"%r\", \"timestamp\": %{%s}t, \"level\": 6, \"_user_agent\": \"%{User-Agent}i\", \"_src_ip\": \"%a\", \"_duration_usec\": %D, \"_duration_sec\": %T, \"_request_size_byte\": %O, \"_http_status_orig\": %s, \"_http_status\": %>s, \"_http_request_path\": \"%U\", \"_http_request\": \"%U%q\", \"_http_method\": \"%m\", \"_http_referer\": \"%{Referer}i\", \"_from_apache\": \"true\" }" GELF
```
Next, add or replace existing logging (typically within each vhost) using the GELF format defined above:

```
CustomLog ${APACHE_LOG_DIR}/gelf.log GELF
```

### Copyright

Copyright (c) 2021 Tenebris Technologies Inc. All rights reserved.

Please see the LICENSE file for additional information.
