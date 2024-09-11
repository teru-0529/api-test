# 受注取得(OK case1)

受注情報の取得処理（成功ケース）。

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

```http
# @name bulk insert(orders.receivings)
POST {{$dotenv POSTGREST_API_HOST}}/receivings HTTP/1.1
Content-Type: application/json
Content-Profile: orders

[
  {
    "customer_name": "徳川商事株式会社",
    "operator_name": "織田信長"
  },
  {
    "customer_name": "株式会社武田物産",
    "operator_name": "豊臣秀吉"
  }
]
```

```http
# @name bulk insert(orders.receiving_details)
POST {{$dotenv POSTGREST_API_HOST}}/receiving_details HTTP/1.1
Content-Type: application/json
Content-Profile: orders

[
  {
    "order_no": "RO-0000001",
    "product_id": "P0002",
    "receiving_quantity": 5,
    "sellling_price": 58000
  },
  {
    "order_no": "RO-0000002",
    "product_id": "P0001",
    "receiving_quantity": 3,
    "sellling_price": 32000
  },
  {
    "order_no": "RO-0000002",
    "product_id": "P0003",
    "receiving_quantity": 1,
    "sellling_price": 19000
  }
]
```

## テスト実行

```http
# @name execute API
GET {{$dotenv ORDERS_API_HOST}}/receivings/RO-0000002 HTTP/1.1
Content-Type: application/json
x-account-id: P0673822
```

## 検証対象

- API実行結果の http status が(200)であること
- API実行結果が正しいこと

-----
api-test-specification.(v1.0.0)
