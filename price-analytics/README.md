# Price analytic component
Component to store and response price analytic data for stat purpose

## Build

```shell
    cd price-analytics/cmd
    go build -v .
```

## Run
in the cmd directory, run shell command:

```shell
    ./cmd --postgres_host=127.0.0.1:5432 --postgres_user=<postgres_username> --postgres_password=<postgres_password> --postgres_database=<postgres_database>
```

**postgres_host**: host to connect to postgres, default value is *127.0.0.1:5432* 
**postgres_user**: user to connect to postgres, default value is empty
**postgres_password**: password to connect to postgres with provided user, default value is empty
**postgres_database**: database to store data in postgres,default value is empty

## APIs

### Update Price Analytic Data - (signing required) set a record marking the condition because of which the set price is called.

```
<host>:8000/update-price-analytic-data
POST request
params:
 - timestamp - the timestamp of the action (real time ) in millisecond
 - value - the json enconded object to save. 

Note: the data sent over must be encoded in Json in order to make it valid for output operation
  In Python, the data would be encoded as:
   data = {"timestamp": timestamp, "value": json.dumps(analytic_data)} 

```

**response:**

on success:
```json
{
    "success": true
}
```
on failure:

```json
{
    "success":false,
    "reason": <error>
}
```

### Get Price Analytic Data - (signing required) list of price analytic data, sorted by timestamp

```
<host>:8000/get-get-price-analytic-data
GET request
params:
 - fromTime (integer) - from timestamp (millisecond)
 - toTime (integer) - to timestamp (millisecond)
example:

curl -x GET \
  http://localhost:8000/get-price-analytic-data?fromTime=1522753160000&toTime=1522755792000
```

**response:**

```json
{
  "data": [
    {
      "Timestamp": 1522755271000,
      "Data": {
        "block_expiration": false,
        "trigger_price_update": true,
        "triggering_tokens_list": [
          {
            "ask_price": 0.002,
            "bid_price": 0.003,
            "mid afp_old_price": 0.34555,
            "mid_afp_price": 0.6555,
            "min_spread": 0.233,
            "token": "OMG"
          },
          {
            "ask_price": 0.004,
            "bid_price": 0.005,
            "mid afp_old_price": 0.21555,
            "mid_afp_price": 0.4355,
            "min_spread": 0.133,
            "token": "KNC"
          }
        ]
      }
    }
  ],
  "success": true
}
```