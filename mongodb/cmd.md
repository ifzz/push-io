sudo su

mongod --master --dbpath /opt/mongodb/ &

mongod --slave  --source 10.2.68.215:57017 --dbpath /opt/mongodb/ &
