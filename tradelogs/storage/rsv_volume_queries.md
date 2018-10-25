```SQL
-- Continous queries:
CREATE CONTINUOUS QUERY "rsv_volume_src_src_hr" ON trade_logs RESAMPLE Every 1h for 3h BEGIN 
SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO rsv_volume_hour FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") OR (src_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE') AND "src_rsv_addr"!='') GROUP BY "src_addr","src_rsv_addr",time(1h)
END

CREATE CONTINUOUS QUERY "rsv_volume_src_dst_hr" ON trade_logs RESAMPLE Every 1h for 3h BEGIN 
SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO rsv_volume_hour FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") OR (src_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE') AND "dst_rsv_addr"!='') GROUP BY "src_addr","dst_rsv_addr",time(1h)
END

CREATE CONTINUOUS QUERY "rsv_volume_dst_src_hr" ON trade_logs RESAMPLE Every 1h for 3h BEGIN 
SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO rsv_volume_hour FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") OR (src_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE') AND "src_rsv_addr"!='') GROUP BY "dst_addr","src_rsv_addr",time(1h)
END

CREATE CONTINUOUS QUERY "rsv_volume_dst_dst_hr" ON trade_logs RESAMPLE Every 1h for 3h BEGIN 
SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO rsv_volume_hour FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") OR (src_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE') AND "dst_rsv_addr"!='') GROUP BY "dst_addr","dst_rsv_addr",time(1h)
END

CREATE CONTINUOUS QUERY "rsv_volume_src_src_day" ON trade_logs RESAMPLE Every 1h for 2D BEGIN
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO rsv_volume_day FROM rsv_volume_hour WHERE src_addr!='' and src_rsv_addr!='' GROUP BY "src_addr","src_rsv_addr",time(1d)
END

CREATE CONTINUOUS QUERY "rsv_volume_src_dst_day" ON trade_logs RESAMPLE Every 1h for 2D BEGIN
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO rsv_volume_day FROM rsv_volume_hour WHERE src_addr!='' and dst_rsv_addr!='' GROUP BY "src_addr","dst_rsv_addr",time(1d)
END


CREATE CONTINUOUS QUERY "rsv_volume_dst_src_day" ON trade_logs RESAMPLE Every 1h for 2D BEGIN
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO rsv_volume_day FROM rsv_volume_hour WHERE dst_addr!='' and src_rsv_addr!='' GROUP BY "dst_addr","src_rsv_addr",time(1d)
END


CREATE CONTINUOUS QUERY "rsv_volume_dst_dst_day" ON trade_logs RESAMPLE Every 1h for 2D BEGIN
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO rsv_volume_day FROM rsv_volume_hour WHERE dst_addr!='' and dst_rsv_addr!='' GROUP BY "dst_addr","dst_rsv_addr",time(1d)
END

--- When import new DB, these queries has to be run to aggregate the volume. 
SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO rsv_volume_hour FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") OR (src_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE') AND "dst_rsv_addr"!='') GROUP BY "src_addr","dst_rsv_addr",time(1h)
SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO rsv_volume_hour FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") OR (src_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE') AND "src_rsv_addr"!='') GROUP BY "src_addr","src_rsv_addr",time(1h)
SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO rsv_volume_hour FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") OR (src_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE') AND "dst_rsv_addr"!='') GROUP BY "dst_addr","dst_rsv_addr",time(1h)
SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO rsv_volume_hour FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") OR (src_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE') AND "src_rsv_addr"!='') GROUP BY "dst_addr","src_rsv_addr",time(1h)

SELECT SUM(token_volume) as token_volume, SUM(eth_volume) AS eth_volume, SUM(usd_volume) AS usd_volume INTO rsv_volume_day FROM rsv_volume_hour WHERE src_addr!='' AND src_rsv_addr!='' GROUP BY "src_addr","src_rsv_addr",time(1d)
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO rsv_volume_day FROM rsv_volume_hour WHERE src_addr!='' and dst_rsv_addr!='' GROUP BY "src_addr","dst_rsv_addr",time(1d)
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO rsv_volume_day FROM rsv_volume_hour WHERE dst_addr!='' and src_rsv_addr!='' GROUP BY "dst_addr","src_rsv_addr",time(1d)
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO rsv_volume_day FROM rsv_volume_hour WHERE dst_addr!='' and dst_rsv_addr!='' GROUP BY "dst_addr","dst_rsv_addr",time(1d)

-- Reserver volume is queried as:
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) as usd_volume FROM rsv_volume_hour WHERE $timeFilter AND ((dst_addr='<assertAddr>' OR src_addr='<assertAddr>') AND (dst_rsv_addr='<rsvAddr>' OR src_rsv_addr='<rsvAddr>')) GROUP BY time(1<freq>)
```
