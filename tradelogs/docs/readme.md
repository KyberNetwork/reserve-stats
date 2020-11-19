# Documentation for Katalyst upgrade


## Database design

Object entity: 
    1. tradelog
    2. reserve
    3. user
    4. platform
    5. token

### Tables

### tradelog

column_name | data_type | is_nullable
----------- | --------- | -----------
timestamp | timestamp | false

### reserve

column_name | data_type | is_nullable | description
----------- | --------- | ----------- | -----------
id  | sequence | false | id of reserve
address | text (hex value) | false | address of reserve
reserve_id | text | true  | reserve id of the reserve which is set on contract
block_number | integer | true | block number where addReserveToStorage update
rebate_wallet_address | text | true | wallet to receive rebate

### user

column_name | data_type | is_nullable | description 
----------- | --------- | ----------- | -----------
id | sequence (integer) | false |
address | text | false | user address
timestamp | timestamp | false 

### platform

column_name | data_type | is_nullable | description
----------- | --------- | ----------- | -----------
id | sequence (integer) | false |
name | text | true |
address | text (hex address) | true | platform address to receive platform fee

### token

column_name | data_type | is_nullable | description
----------- | --------- | ----------- | -----------
id | sequence (integer) | false |
address | text | false |
name | text | false |


## Entity Relationship 

tradelog