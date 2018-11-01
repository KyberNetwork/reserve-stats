CREATE CONTINUOUS QUERY "burn_fee_1h" ON "trade_logs"
RESAMPLE EVERY 1h FOR 1d
BEGIN
    SELECT SUM("amount") as "sum_amount" INTO "hourly_burn_fees" FROM "burn_fees" GROUP BY "reserve_addr", time(1h)
END

CREATE CONTINUOUS QUERY "burn_fee_1d" ON "trade_logs"
RESAMPLE EVERY 1d FOR 3d
BEGIN
    SELECT SUM("amount") as "sum_amount" INTO "daily_burn_fees" FROM "burn_fees" GROUP BY "reserve_addr", time(1d)
END