# log2sqs

This application reads Graylog GELF-format messages from a text file and forwards them to an AWS SQS queue.

### Development Status

This is an alpha release.

### Installation

1) Clone the repo and compile using "go build"
2) Copy the binary (log2sqs) and config file (log2sqs.conf) to /opt/log2sqs. If you put it elsewhere you will need to
   update the .service file.
3) Copy log2sqs.service to /etc/systemd/system/
4) Update the User and Group in log2sqs.service if you do not wish to run as root
5) Update the configuration file (log2sqs.conf)
6) Run 'systemctl daemon-reload'
7) Run 'systemctl start log2sqs' to start the application

### Copyright

Copyright (c) 2021 Tenebris Technologies Inc. All rights reserved.

Please see the LICENSE file for additional information.
