# repository-with-tx

Repository Pattern での Transaction 管理

基本的に Interface では ctx と model もしくは id を受け取る

Repository を実装する構造体の生成時に、DB の Connection を渡し、それを内部で利用する形式

Transaction の共通処理については、 [transaction-common-func](https://github.com/trewanek/transaction-common-func) より流用

DB の Connection 呼び出し側には、Transaction 実行メソッドと DB 切断メソッドを公開し、

他の Query 実行メソッドについては、同一パッケージ内での利用を想定し、 private にしている。
