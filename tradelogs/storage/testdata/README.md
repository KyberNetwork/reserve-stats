# How to generate export.dat test data

```sh
# start with a fresh database
# in root directory of this project
docker-compose down
sudo rm -rf ./data
docker-compose up

# fetch some sample trade logs to database
./trade-logs-crawler \
   --core-url xxx \
   --core-signing-key yyy \
   --from-block=6494037 \
   --to-block=6494137
   
# inside influxdb docker container
docker-compose exec influxdb /bin/bash
influx_inspect export -datadir /var/lib/influxdb/data \
    -waldir /var/lib/influxdb/wal -out /var/lib/influxdb/data/export.dat \
    -database trade_logs -retention autogen
exit

# exit to host shell
cp data/influxdb/data/export.dat ./tradelogs/storage/testdata/

# remove this comment section in export.dat
# INFLUXDB EXPORT:
# ...
# writing wal data
```
