# log2sqs

This application reads one or more log files in real time (like tail) and forwards them to an AWS SQS queue.

It will optionally add AWS EC2 instance metadata (instance ID, hostname, and tags) to each log message.

Graylog GELF format messages (in JSON) and Apache2/NGINX combined log format messages are currently supported.

### Development Status

This is an alpha release.

### Installation

1) Clone the repo and compile using "go build"
2) Copy the binary (log2sqs) and config file (log2sqs.conf) to /opt/log2sqs. If you put it elsewhere you will need to update the .service file.
3) Ensure that log2sqs has the execution bit set (i.e. chmod 700 or 755)
4) Copy log2sqs.service to /etc/systemd/system/
5) Update the User and Group in log2sqs.service if you do not wish to run as root
6) Update the configuration file (log2sqs.conf)
7) Run 'systemctl daemon-reload'
8) Run 'systemctl start log2sqs' to start the application
9) Configure logrotated to locate the log file specified in log2sqs.conf. A sample rotate file is contained in log2sqs.txt. It can be copied to /etc/logrotate.d/log2sqs. It should be owned by root and set to mode 0644.

### Copyright

Copyright (c) 2021 Tenebris Technologies Inc. All rights reserved.

Please see the LICENSE file for additional information.
