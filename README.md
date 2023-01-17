# crypto-server

## Step 1: Build the binary
* Uinx/Linux
```sh
$ mkdir -p ./bin > /dev/null
$ go build -mod=vendor -v -ldflags "-w -s" -o ./bin/crypto-server ./main.go
```

* Windows
```cmd
$ mkdir -p ./bin > /dev/null
$ go build -mod=vendor -v -ldflags "-w -s" -o ./bin/crypto-server.exe ./main.go
```
## Step 2: Run the binary 

* Uinx/Linux
```sh
$ ./bin/crypto-server
```

* Windows
```sh
$ ./bin/crypto-server.exe
```

## API DOC:

Server: `http://localhost:8080`

### 1. Get All Symbols 

>Method: GET  
Endpoint: /api/v1/symbol  

Response:
```json
{
    "error": "",
    "data": [
        "BTCUSD",
        "ETHBTC",
        "BTCETH",
        "ETHUSD"
    ]
}
```

### 2. Add new Symbol

>Method: POST  
Endpoint: /api/v1/symbol

Request: 
```json
{
    "base_currency":"ETH",
    "pair_currency":"USD"
}
```
Response:
```json

```
    
### 3. Get All Currencies

>Method: GET  
Endpoint: /api/v1/currency/all  

Response:
```json
{
    "error": "",
    "data": [
        {
            "id": "BTC",
            "fullName": "Bitcoin",
            "ask": "20900.67",
            "bid": "20896.99",
            "last": "20896.99",
            "open": "20841.90",
            "low": "20565.62",
            "high": "21076.73",
            "feeCurrency": "USD"
        },
        {
            "id": "ETH",
            "fullName": "Ethereum",
            "ask": "0.074308",
            "bid": "0.074286",
            "last": "0.074292",
            "open": "0.073662",
            "low": "0.073425",
            "high": "0.074541",
            "feeCurrency": "BTC"
        }
    ]
}
```
### 4. Get Currency by ID

>Method: GET  
Endpoint: /api/v1/currency/{symbol}

Response:
```json
{
            "id": "ETH",
            "fullName": "Ethereum",
            "ask": "1552.755",
            "bid": "1551.426",
            "last": "1553.000",
            "open": "1535.301",
            "low": "1517.172",
            "high": "1566.018",
            "feeCurrency": "USD"
        }
```
