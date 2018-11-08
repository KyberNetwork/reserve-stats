
# Tradelogs 

Trade logs component crawls all tradelogs which go through Kyber contract in order to store and aggregation then provide appropriate stats.

## Get trade logs 



```shell
curl -X GET "http://localhost:8004/trade-logs?from="
```

> The above command returns JSON structured like this:

```json
{
  "data": "1517479497447",
  "success": true
}
```

Return list of trade logs **from** a point time and **to** another point of time

### HTTP Request
