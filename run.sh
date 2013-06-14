#!/bin/bash
BASE_DIR=`dirname $0`
cd $BASE_DIR
pkill kindle-delivery
#nohup $BASE_DIR/bin/comics 2>&1 >> ~/work/comics/logs/comics.log &
