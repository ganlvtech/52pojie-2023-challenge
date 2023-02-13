# 服务端

## 谜题页面

### 调试运行

```bash
go run . -listen :7999
```

然后访问 http://127.0.0.1:7999/

### 构建

```bash
go build .
```

### 部署

```bash
./challenge -listen :7999
# 或者使用 nohup 后台运行
nohup ./challenge -listen :7999 &
```

```nginx
server {
    listen       80;
    listen       443 ssl;
    server_name  2023challenge.52pojie.cn;
    server_name  2023challenge.ganlv.tech;
    index        index.html;

    location / {
        proxy_set_header  Host $host;
        proxy_pass        http://127.0.0.1:7999;
    }

    ssl_certificate            /.acme.sh/ganlv.tech/fullchain.cer;
    ssl_certificate_key        /.acme.sh/ganlv.tech/ganlv.tech.key;
    ssl_ciphers                ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
    ssl_protocols              TLSv1 TLSv1.1 TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers  on;
}
```

注意 flagB 需要手动修改 DNS

## 提交页面

### 调试运行

```bash
php -S 0.0.0.0:8000
```

然后访问 http://127.0.0.1:8000/

### 部署

实际部署时是嵌入到 DiscuzX 的代码中，这部分代码不公开。
