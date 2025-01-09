#!/bin/bash
mysqlhost="192.168.31.97"
mysqlport="3306"
mysqlusr="root"
mysqlpasswd="123456"

databasename="j_server_game_db"

mysql -h${mysqlhost} -P${mysqlport} -u${mysqlusr} -p${mysqlpasswd}  -e \
"CREATE DATABASE IF NOT EXISTS ${databasename}; USE ${databasename}; source sql/game_db.sql;"