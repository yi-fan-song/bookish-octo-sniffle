#! /bin/bash

DB_USER=dbuser
DB_PASSWORD=Vj23urju

service postgresql start
su -c "psql -c \"CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';\"" postgres
su -c "psql -c \"CREATE DATABASE pomodoro\"" postgres

./build/octo --port 8080
