
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
