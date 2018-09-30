## 配置简解

sShare 的配置方式为单 json 文件，文件名为`config.json`。下面是示例配置文件及简解：

```json5
{
  "run_command": {   // 启动代理程序的配置，sShare会为每一个用户执行一次run_command
    "cmd": "ssserver", // 命令
    "arg": "-p {{.Port}} -k {{.Pass pass_name}}", // 参数，其中的{{.Port}}会被替换为端口，{{.Pass pass_name}}替换为密码，这两个参数都是随机生成的
    "enabled": true // 是否启用，run_command必须启用
  },
  "exit_command": {  // 当程序退出时，sShare应执行的任务，通常为清理临时文件、临时防火墙规则等
    "cmd": "exit", // 命令
    "arg": "{{.Port}}", // 参数，exit_command的参数中只有{{.Port}}会被替换
    "enabled": false // 是否启用，如果你不需要，可以选择不启用exit_command
  },
  "captcha": { // 验证码选项
    "name": "base", // 验证码接口的名称，目前已实现的验证码见下方“配置示例”中“验证码”部分
    "site_id": "23333", // Site ID，该属性的含义因接口而异，详见示例
    "extra": "66666" // 额外数据，该属性的含义因接口而异，详见示例
  },
  "ttl": "20m", // 一个用户在获取到账号后可以使用的时间，超时后对应的进程将被Kill，用户需要重新在web界面获取，单位可以为s（秒），m（分），h（时）
  "limit": 20, // 限制的最大用户数量
  "web_addr": ":9527", // Web界面监听地址，使用"ip:port"可以指定监听ip，使用":port"监听所有ip，建议监听本地(127.0.0.1)并使用Nginx反代
  "port_start": 2000, // 分配端口起始值
  "port_range": 200, // 分配端口范围，最终用户得到的端口将在[port_start, port_start+port_range]之间，请务必保证该范围内端口没有被占用
  "rand_seed": 23343, // 随机种子，可以不设置
  "no_check_alive": false, // 运行代理程序时不检查是否存活，具体用途参考下面的“配置示例”中ssr mujson部分
  "safe": { // 安全特性们
      "anti_cc": true,  // 反CC（大概有用？）
      "city_check": true,  // 检测访问的源ip所在地级市
      "city_name": "Beijing", // 检测的城市名称
      "city_file":"/path/to/your/file.datx", // IPIP的DATX库文件路径
      "cdn_enabled": false  // 是否开启了CDN或使用Nginx反代（如果开启则需要提取xff）
  }
}
```

需要注意的是 sShare 目前使用`text/template`包实现可替换的参数模板，并支持多密码和多uuid同时生成，使用`{{.Port}}`替换端口，使用`{{.Pass \"pass_name\"}}`和`{{.UUID \"uuid_name\"}}`来生成一个密码，所有生成的密码将被返回，所有同名的密码都将是一致的，详情参阅api部分。
