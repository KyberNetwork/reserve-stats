-- CONTINOUS QUERY 
CREATE CONTINUOUS QUERY "dst_volume" on trade_logs RESAMPLE EVERY 1h for 3h BEGIN SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades) GROUP BY "dst_addr", time(1h) END

CREATE CONTINUOUS QUERY "dst_volume_day" on trade_logs RESAMPLE EVERY 1h for 2d BEGIN SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO volume_day FROM volume_hour WHERE dst_addr!='' GROUP BY "dst_addr", time(1d) END  

CREATE CONTINUOUS QUERY "src_volume" on trade_logs RESAMPLE EVERY 1h for 3h BEGIN SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades) GROUP BY "dst_addr", time(1h) END

CREATE CONTINUOUS QUERY "src_volume_day" on trade_logs RESAMPLE EVERY 1h for 2d BEGIN SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO volume_day FROM volume_hour WHERE src_addr!=''GROUP BY "src_addr", time(1d) END  

-- WHEN Import new DB, historical data must be aggregate using these command manually :
SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades) GROUP BY "dst_addr", time(1h)
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO volume_day FROM volume_hour WHERE dst_addr!='' GROUP BY "dst_addr", time(1d)
SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades) GROUP BY "dst_addr", time(1h) 
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO volume_day FROM volume_hour WHERE src_addr!=''GROUP BY "src_addr", time(1d) 

-- Asset volume is queried as:
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume FROM volume_<freq> WHERE $timeFilter AND (dst_addr='<assetAddr>' OR src_addr='<assetAddr>') GROUP BY "dst_addr", time(1<freq>)
