#!/usr/bin/env bash

KUBECTL=`which kubectl`
if [ "$KUBECTL" == "" ];then
    echo "can't found kubectl command..."
    exit
fi


NAMESPACE="mlcloud"
kubectl -n $NAMESPACE delete configmap init-sql
sleep 5
kubectl -n $NAMESPACE create configmap init-sql --from-file=./mlcloud.sql
