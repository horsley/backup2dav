#!/bin/sh

echo "${CRON_SCHEDULE} /main > /proc/1/fd/1 2>&1" > /etc/crontabs/app
crond -f -l 2