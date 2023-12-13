# Debug: HTTPS 跳转

## 问题描述

对于如下请求：

`POST https://192.168.0.200/api/2/list?token=dr95wfdg9023qaz0`

它能够正常返回数据，但若将 `https` 修改为 `http`，就会返回如下内容：

```json
{
  "code": 3030,
  "data": null,
  "message": "认证失败"
}
```

## 分析与解答

问题表象：通过返回体结构可以判断出是下游服务返回的，但问题不是这里导致的。

分析：注意到 token 是以参数的形式加在 URI 里面，应该是 token 丢失导致的认证失败。

结论：排查出来是网关配置中做 HTTPS 转发的地方导致的。`ngx.var.uri` 应该改成 `ngx.var.request_uri`，才会保留参数。

此外，**默认的重定向码 302，按照标准，客户端可能会将 POST 请求转为 GET**，要用 307 Temporary Redirect 或者 308 Permanent Redirect 返回码，才会保留 POST 请求。

修改之后如下：

```
rewrite_by_lua_block {
    if (os.getenv('NO_AUTO_HTTPS') == nil or os.getenv('NO_AUTO_HTTPS') == '') and (not string.find(ngx.var.http_host, ':')) then
        return ngx.redirect('https://' .. ngx.var.http_host .. ngx.var.request_uri, 307)
    end
}
```
