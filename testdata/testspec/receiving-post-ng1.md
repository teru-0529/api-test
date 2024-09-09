# 受注登録(NG case1)

受注情報の登録処理（失敗ケース）、受注明細で指定する商品の受注数量が正しくない(負数)場合。登録に失敗し400が返る。

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
      "orderQuantity": -1,
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

- API実行結果の http status が(400)であること
-----
api-test-specification.(v1.0.0)
