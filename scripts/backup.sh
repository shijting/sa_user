#!/bin/bash

# 数据库名称
dbname=$1
# 备份路径（包含文件名）
backup_file=$2
# 数据库端口
port=$3
/data/pg14sql/bin/pg_dump -h 127.0.0.1 -U postgres -p ${port} -Fc -f ${backup_file} ${dbname}
