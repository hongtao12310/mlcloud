#!/bin/bash

# create docker bridge for mlcloud
function create_network
{
    name="mlcloud"
    echo "name: $1"
    if [ "$1" != "" ]; then
        name=$1
    fi
    docker network create $name >/dev/null 2&>1
}

function run_mysql
{
    docker stop mysql >/dev/null 2&>1
    docker rm mysql >/dev/null 2&>1
    docker run -d --name mysql --net mlcloud \
    -e MYSQL_ROOT_PASSWORD=root \
    -e MYSQL_USER=mlcloud \
    -e MYSQL_PASSWORD=mlcloud \
    -e MYSQL_DATABASE=mlcloud -p 3306:3306 \
    -v $HOME/mlcloud/data/mysql:/var/lib/mysql \
    -v $HOME/mlcloud/data/init_sql:/docker-entrypoint-initdb.d mysql:5.7
}

create_network mlcloud
run_mysql