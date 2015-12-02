#!/bin/bash

# hdfs
HDFS_HOST="uscc0-master0.uncharted.software"
HDFS_POST="8020"

# es
ES_DEV_CLUSTER="http://192.168.0.41"
ES_PROD_CLUSTER="http://10.65.16.13"
ES_OPENSTACK="http://10.64.16.120"
ES_HOST=$ES_DEV_CLUSTER
ES_PORT="9200"

# time duration
#START_DATE="1433116800"
#END_DATE="1448928000"
DURATION="14515200" # 6 months in seconds

# nyc
HDFS_NYC_PATH="/xdata/data/SummerCamp2015/JulyData-processed/nyc_twitter_merged"
HDFS_NYC_COMPRESSION=""
ES_NYC_INDEX="nyc_twitter_test"
ES_NYC_DOC_TYPE="nyc_twitter"

# go run main.go \
#     -es-host=$ES_HOST \
#     -es-port=$ES_PORT \
#     -es-index=$ES_NYC_INDEX \
#     -es-doc-type=$ES_NYC_DOC_TYPE \
#     -hdfs-host=$HDFS_HOST \
#     -hdfs-port=$HDFS_POST \
#     -hdfs-path=$HDFS_NYC_PATH \
#     -hdfs-compression=$HDFS_NYC_COMPRESSION

# isil
HDFS_ISIL_PATH="/xdata/data/twitter/isil-keywords"
HDFS_ISIL_COMPRESSION="gzip"
ES_ISIL_INDEX="isil_twitter_weekly"
ES_ISIL_DOC_TYPE="isil_twitter_deprecated"

# go run main.go \
#     -es-host=$ES_DEV_CLUSTER \
#     -es-port=$ES_PORT \
#     -es-index=$ES_ISIL_INDEX \
#     -es-doc-type=$ES_ISIL_DOC_TYPE \
#     -hdfs-host=$HDFS_HOST \
#     -hdfs-port=$HDFS_POST \
#     -hdfs-path=$HDFS_ISIL_PATH \
#     -hdfs-compression=$HDFS_ISIL_COMPRESSION \
#     -duration=$DURATION

go run main.go \
    -es-host=$ES_PROD_CLUSTER \
    -es-port=$ES_PORT \
    -es-index=$ES_ISIL_INDEX \
    -es-doc-type=$ES_ISIL_DOC_TYPE \
    -hdfs-host=$HDFS_HOST \
    -hdfs-port=$HDFS_POST \
    -hdfs-path=$HDFS_ISIL_PATH \
    -hdfs-compression=$HDFS_ISIL_COMPRESSION \
    -duration=$DURATION
