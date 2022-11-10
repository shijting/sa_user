#!/bin/bash

#备份路径

dbname=$1
backup_file=$2
port=$3
/data/pg14sql/bin/pg_dump -h 127.0.0.1 -U postgres -p ${port} -Fc -f ${backup_file} ${dbname}
