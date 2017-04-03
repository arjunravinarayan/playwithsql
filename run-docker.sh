#!/bin/bash

DB=$1
HOST=$2
LOOPS=$3
PAUSETIME=$4
MAXCONNS=$5


docker service rm pws_cmd-$DB

docker service create \
    --name pws_cmd-$DB  \
    --restart-condition none \
    --network pws_default \
    vincentserpoul/playwithsql-cmd \
    -db=$DB -host=$HOST -loops=$LOOPS -pausetime=$PAUSETIME -maxconns=$MAXCONNS 

WAITTILFINISH=$(($LOOPS/10));
sleep $WAITTILFINISH;

docker service logs pws_cmd-$DB | awk '{ print $3."," }' >> results.log;
docker service rm pws_cmd-$DB;
docker service rm pws_$DB;