#!/bin/bash

# hdfs
HDFS_HOST="uscc0-master0.uncharted.software"
HDFS_POST="8020"

# es
ES_HOST="http://10.64.16.120"
ES_PORT="9200"

# nyc
HDFS_NYC_PATH="/xdata/data/SummerCamp2015/JulyData-processed/nyc_twitter_merged"
HDFS_NYC_COMPRESSION=""
ES_NYC_INDEX="nyc_twitter"
ES_NYC_DOC_TYPE="nyc_twitter"

# isil
HDFS_ISIL_PATH="/xdata/data/twitter/isil-keywords"
HDFS_ISIL_COMPRESSION="gzip"
ES_ISIL_INDEX="isil_twitter"
ES_ISIL_DOC_TYPE="isil_twitter"

#go run main.go -es-host=$ES_HOST -es-port=$ES_PORT -es-index=$ES_NYC_INDEX -es-doc-type=$ES_NYC_DOC_TYPE -hdfs-host=$HDFS_HOST -hdfs-port=$HDFS_POST -hdfs-path=$HDFS_NYC_PATH -hdfs-compression=$HDFS_NYC_COMPRESSION
go run main.go -es-host=$ES_HOST -es-port=$ES_PORT -es-index=$ES_ISIL_INDEX -es-doc-type=$ES_ISIL_DOC_TYPE -hdfs-host=$HDFS_HOST -hdfs-port=$HDFS_POST -hdfs-path=$HDFS_ISIL_PATH -hdfs-compression=$HDFS_ISIL_COMPRESSION
