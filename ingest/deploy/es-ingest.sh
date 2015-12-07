#!/bin/bash

# ingest
INGEST_EXE="es-ingest.bin"

# hdfs
HDFS_HOST="uscc0-master0.uncharted.software"
HDFS_PORT="8020"

# es
ES_DEV_CLUSTER="http://192.168.0.41"
ES_PROD_CLUSTER="http://10.65.16.13"
ES_PORT="9200"

# time duration
#START_DATE="1433116800"
#END_DATE="1448928000"
DURATION="14515200" # 6 months in seconds

# isil
HDFS_ISIL_PATH="/xdata/data/twitter/isil-keywords"
HDFS_ISIL_COMPRESSION="gzip"
ES_ISIL_INDEX="isil_twitter_weekly"
ES_ISIL_DOC_TYPE="isil_twitter_deprecated"

# ingest into dev
./$INGEST_EXE \
    -es-host=$ES_DEV_CLUSTER \
    -es-port=$ES_PORT \
    -es-index=$ES_ISIL_INDEX \
    -es-doc-type=$ES_ISIL_DOC_TYPE \
    -hdfs-host=$HDFS_HOST \
    -hdfs-port=$HDFS_PORT \
    -hdfs-path=$HDFS_ISIL_PATH \
    -hdfs-compression=$HDFS_ISIL_COMPRESSION \
    -duration=$DURATION

# ingest into prod
./$INGEST_EXE \
    -es-host=$ES_PROD_CLUSTER \
    -es-port=$ES_PORT \
    -es-index=$ES_ISIL_INDEX \
    -es-doc-type=$ES_ISIL_DOC_TYPE \
    -hdfs-host=$HDFS_HOST \
    -hdfs-port=$HDFS_PORT \
    -hdfs-path=$HDFS_ISIL_PATH \
    -hdfs-compression=$HDFS_ISIL_COMPRESSION \
    -duration=$DURATION
