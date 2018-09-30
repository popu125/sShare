## API 说明

目前 sShare 有两个公开api接口，分别用于查询当前客户端数和获取账号。文档如下：

### /api/count

请求类型：`POST`

直接请求即可得到目前已分配的客户端数量。当开启反CC时，该接口也可能会收到`{"Status":"ERR_ANTICC"}`，应在收到后刷新页面。

### /api/new

请求类型：`POST`

请求数据：

| Key   | Value    |
| ----- | -------- |
| token | 验证码服务返回值 |

该请求将返回一个 JSON 类型的返回值：

| Key    | Value |
| ------ | ----- |
| Status | 状态    |
| Info   | JSON字典形式的返回，包含`pass_map`和`uuid_map`两个hashmap，和一个Port |

其中状态码含义如下：

 - `ERR_NO_CAPTCHA`: 验证码未通过
 - `ERR_FULL`: 用户池已满
 - `ERR_ANTICC`: 触发反CC，应在收到后刷新页面
 - `ACCEPT`: 成功分配
