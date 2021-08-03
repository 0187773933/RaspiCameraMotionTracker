#!/bin/bash
# APP_NAME="rpmt-motion-tracker"
APP_NAME="xp6qhg9fmuolztbd2ixwdbtd1/raspi-motion-tracker:arm32test"
# sudo docker rm $APP_NAME -f || echo "failed to remove existing ssh server"
sudo docker build --squash -t $APP_NAME .