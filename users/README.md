# Kyber Reserve Stat User Component

# Build

```shell
    cd cmd
    go build - v .
```

# Run

```
    cd cmd
    ./cmd --host=127.0.01:5432 --user=admin --password="" --database=user_storage
```

## Available flags

**host**: where to connect to postgresql db, default is localhost: *127.0.0.1:5432*
**user**: postgresql db user, default is empty
**password**: postgresql db password, default is empty
**database**: postgresql db database to use, default is empty

## Available API

### Get user info by address

```http
GET http://localhost:9000/users/:userAddress
```

response: Return user info

**Sample response:**

```json
{
    "success": true,
    "data": {
        "kyc": false,
        "cap": 12313471038132,
        "rich": false,
    }
}
```

### Update user addresses

```http
POST http://localhost:9000/users
Content-Type: Application/x-www-form-urlencoded

user=halink0803@gmail.com&addresses=0x829bd824b016326a401d083b33d092293333a830-0xc499ae5806b7888aa3c539b3be7a691e83908a04&timestamps=1538020234000-1538020242000
```

Sample response:

```json
{
    "success": true
}
```