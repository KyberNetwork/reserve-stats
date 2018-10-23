```SQL
-- CONTINOUS QUERY for volume aggregation
CREATE CONTINUOUS QUERY "dst_volume_hour" on trade_logs RESAMPLE EVERY 1h for 3h BEGIN SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") OR (src_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE')) GROUP BY "dst_addr", time(1h) END

CREATE CONTINUOUS QUERY "dst_volume" on trade_logs RESAMPLE EVERY 1m for 3m BEGIN SELECT SUM(dst_amount) AS tok_vol, SUM(fiat_amount) as fiat_vol,SUM(eth_receival_amount) AS eth_vol  INTO volume FROM trades GROUP BY "dst_addr",time(1m) END
CREATE CONTINUOUS QUERY "src_volume" on trade_logs RESAMPLE EVERY 1m for 3m BEGIN SELECT SUM(src_amount) AS tok_vol,SUM(fiat_amount) as fiat_vol, SUM(eth_receival_amount) AS eth_vol  INTO volume FROM trades GROUP BY "src_addr",time(1m) END

-- The four following cq can be removed and added to previous query if we modify eth_receival_amount from each eth-tok trade and tok-eth to the eth value of that trade.
CREATE CONTINUOUS QUERY "src_eth_volume_tok" on trade_logs RESAMPLE EVERY 1m for 3m BEGIN SELECT SUM(eth_receival_amount) AS eth_vol INTO volume FROM trades WHERE eth_receival_amount >0 GROUP BY "src_addr", time(1m) END
CREATE CONTINUOUS QUERY "dst_eth_volume_tok" on trade_logs RESAMPLE EVERY 1m for 3m BEGIN SELECT SUM(eth_receival_amount) AS eth_vol INTO volume FROM trades WHERE eth_receival_amount >0 GROUP BY "dst_addr", time(1m) END
CREATE CONTINUOUS QUERY "src_eth_volume_eth" on trade_logs RESAMPLE EVERY 1m for 3m BEGIN SELECT SUM(src_amount) AS eth_vol INTO volume FROM trades WHERE "src_addr"='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' GROUP BY "dst_addr",time(1m) END
CREATE CONTINUOUS QUERY "dst_eth_volume_eth" on trade_logs RESAMPLE EVERY 1m for 3m BEGIN SELECT SUM(dst_amount) AS eth_vol INTO volume FROM trades WHERE "dst_addr"='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' GROUP BY "src_addr",time(1m) END

-- Sum these for volume : 
SELECT eth_vol_eth_tok, eth_vol_tok_tok_dst, dst_vol, dst_fiat_vol where $time_filter and "dst_addr" = Token_addr
SELECT eth_vol_tok_eth, eth_vol_tok_tok_src, src_vol, src_fiat_vol where $time_filter and "src_addr" = Token_addr

-- ETH_VOL= eth_vol_eth_tok + eth_vol_tok_tok_dst + eth_vol_tok_eth + eth_vol_tok_tok_src
-- FIAT_VOL = dst_fiat_vol +src_fiat_vol 
-- cTOKEN_VOL  = dst_vol + src_vol
-- WHEN Import new DB, historical data must be aggregate using these command manually :
SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") OR (src_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE')) GROUP BY "dst_addr", time(1h)
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO volume_day FROM volume_hour WHERE dst_addr!='' GROUP BY "dst_addr", time(1d)
SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2") OR (src_addr!="0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE')) GROUP BY "src_addr", time(1h) 
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO volume_day FROM volume_hour WHERE src_addr!=''GROUP BY "src_addr", time(1d) 

-- Asset volume is queried as:
SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume FROM volume_<freq> WHERE $timeFilter AND (dst_addr='<assetAddr>' OR src_addr='<assetAddr>') GROUP BY "dst_addr", time(1<freq>)
```
