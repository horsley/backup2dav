#!/bin/sh

echo "${CRON_SCHEDULE} /main" > /etc/crontabs/app
crond -f -L 2