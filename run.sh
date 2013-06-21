#!/bin/bash
BASE_DIR=`dirname $0`
cd $BASE_DIR
pkill kindle-delivery
nohup $BASE_DIR/bin/kindle-delivery 2>&1 >> ~/.logs/kindle-delivery/logs/kindle-delivery.log &
