# 受注登録(OK case2)

受注情報の登録処理（成功ケース）。

## テスト対象DBの初期化

```http
# @name sdequenceリセット
POST {{$dotenv DB_RESETER_API_HOST}}/schemas/orders/action-seq-reset HTTP/1.1
Content-Type: application/json

[
  "order_no_seed",
  "product_id_seed"
]
```

```http
# @name table初期化
POST {{$dotenv DB_RESETER_API_HOST}}/schemas/orders/action-truncate HTTP/1.1
Content-Type: application/json

[
  "products",
  "receivings"
]
```

## テストデータ登録

```http
# @name bulk insert(orders.products)
POST {{$dotenv POSTGREST_API_HOST}}/products HTTP/1.1
Content-Type: application/json
Content-Profile: orders

[
  {
    "cost_price": 20000,
    "product_name": "日本刀"
  },
  {
    "cost_price": 40000,
    "product_name": "火縄銃"
  },
  {
    "cost_price": 15000,
    "product_name": "弓"
  }
]
```

## テスト実行

```http
# @name execute API
POST {{$dotenv ORDERS_API_HOST}}/receivings HTTP/1.1
Content-Type: application/json
x-account-id: P0673822

{
  "customerName": "徳川商事株式会社",
  "details": [
    {
      "orderQuantity": 3,
      "productId": "P0001",
      "sellingPrice": 34800
    },
    {
      "orderQuantity": 1,
      "productId": "P0003",
      "sellingPrice": 106400
    }
  ],
  "operatorName": "織田信長"
}
```

## 検証対象

- API実行結果の http status が(201)であること
- テーブル (orders.receivings) のデータが正しいこと

```http
# @name get all data(orders.receivings)
GET {{$dotenv POSTGREST_API_HOST}}/receivings HTTP/1.1
Accept-Profile: orders
```

- テーブル (orders.receiving_details) のデータが正しいこと

```http
# @name get all data(orders.receiving_details)
GET {{$dotenv POSTGREST_API_HOST}}/receiving_details HTTP/1.1
Accept-Profile: orders
```

-----
api-test-specification.(v1.0.0)
