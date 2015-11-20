#!/bin/bash

HDFS_HOST="uscc0-master0.uncharted.software"
HDFS_POST="8020"
HDFS_NYC_PATH="/xdata/data/SummerCamp2015/JulyData-processed/nyc_twitter_new"

ES_HOST="http://10.64.16.120"
ES_PORT="9200"
ES_NYC_INDEX="nyc_twitter"

go run main.go -es-host=$ES_HOST -es-port=$ES_PORT -es-index=$ES_NYC_INDEX -hdfs-host=$HDFS_HOST -hdfs-port=$HDFS_POST -hdfs-path=$HDFS_NYC_PATH
