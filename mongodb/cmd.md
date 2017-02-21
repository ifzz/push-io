sudo su

mongod --master --dbpath /opt/mongodb/ --bind_ip 0.0.0.0 & 

mongod --slave  --source 10.2.68.215 --dbpath /opt/mongodb/ --bind_ip 0.0.0.0 &
