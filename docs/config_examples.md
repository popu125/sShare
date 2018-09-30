## 配置示例

考虑到在真正执行代理程序前，可能需要执行其他命令，我们建议用户使用 bash/python 脚本作为启动命令，我们在下面的示例中（如果需要）将使用脚本。

### 基础模板

基础的配置模板，其他的模板都将以修改本模板和创建本地脚本的方式来创建。本文件位于 Repo 或者 Release 包解压后的根目录中，文件名为 `config.example.json`。

### 验证码

验证码是防止 sShare 惨遭爬虫光顾或者前期维持分享热情的重要手段，我们为您提供下面四种验证码接口，如需要其他请发 issue：

**注意**：在您使用验证码前，您可能需要访问对应的主页以获取验证码接口所需要的参数。

- base - 根本没有验证码，参数可以随便设置，**如果您未正确配置验证码的名字，则默认为base**。
- coinbase - 最著名的挖矿验证码，site_id 参数代表 secret_id，extra 参数为您想要用户为您挖的 hash 数（hash 的含义不再做科普，请自行搜索）。要使用该验证码，请访问 https://coinbase.com 注册并创建一个 site。
- recaptcha - 由谷歌提供的验证码服务，被认为是目前较为安全的验证码。在此处 site_id 输入 Site key，extra 设置为 Secret key。要使用 recaptcha 验证码，请先访问 https://www.google.com/recaptcha/admin 创建一个 site。

### ShadowsocksR

SSR 本身的管理功能已经较为完备，我们不再在模板中提供用户限制的额外脚本，故而本模板不需要创建额外的文件。

如果需要限速功能，请在阅读完本段后下翻至“基于 iptables 的端口连接限制和流量限制”段落。

在使用 SSR 作为 sShare 的后端时，您需要在命令行中将配置项作为参数写入，我们推荐的参数配置为：`-p {{.Port}} -k {{.Pass pass_name}} -m aes-128-cfb -o tls1.2_ticket_auth [-g 可选的混淆参数] -O auth_aes128_md5 [-G 可选的用户连接限制，推荐为1]` 

将其填入`run_command`配置项的`arg`条目，`cmd`中填写 SSR 的“单体版”脚本 server.py 的绝对路径（如我的就是`/home/bobo/ssr/shadowsocks/server.py`并赋予对应文件以执行权限(`chmod a+x`)，配置就宣告完成。

在此配置下 sShare 会为每一个用户自动启动一个对应的 ShadowsocksR 服务端进程，并在时间限制（配置详解中有说明）到期后结束该进程。

示例配置如下（仅包括 run_command 段，请禁用 exit_command）：

```json
"run_command": { 
  "cmd": "/home/bobo/ssr/shadowsocks/server.py",
  "arg": "-p {{.Port}} -k {{.Pass pass_name}} -m aes-128-cfb -O auth_aes128_md5 -o tls1.2_ticket_auth -G 1", 
  "enabled": true
},
```

### ShadowsocksR mujson mode

SSR 在后来的 mu（manyuser）版本中添加了用于用户管理的脚本`mujson_mgr.py`，这也使得我们对于用户的管理更为便捷。额外的功能包括限速、用户限制和流量限制。

在使用 mujson mode（以下简记作mujson）时，因为 mujson 已经存在非阻塞的管理脚本，所以我们最好打开“不检查存活”(nca)，同时使用`run_command`运行添加用户，使用`exit_command`运行清除用户。

请注意，SSR 的 mujson mode 要求您在使用之前初始化并修改 mu API 配置，且服务端需事先运行，相关方法请参考网络。

示例配置如下（仅包括两个 command 段）：

```json
"run_command": { 
  "cmd": "/home/bobo/ssr/mujson_mgr.py",
  "arg": "-a -p {{.Port}} -k {{.Pass pass_name}} -m aes-128-cfb -O auth_aes128_md5 -o tls1.2_ticket_auth -G 1 -t 1", 
  "enabled": true
},
"exit_command": {
  "cmd": "/home/bobo/ssr/mujson_mgr.py",
  "arg": "-d -p {{.Port}}",
  "enabled": true
},
```

但切记还需要配置一个nca：

```json
"no_check_alive": true
```

### Shadowsocks

Shadowsocks 原分支并没有“贴心”的用户限制/限速服务，如果需要，请在阅读完本段后下翻至“使用 iptables 针对每个用户限制流量”段落。

Shadowsocks 的配置更少，也更容易完成（通常我们也不建议直接使用Shadowsocks，至少要加个流量限制吧），此处不再给出示例，请参照 SSR（非 mujson）的配置和注释自行修改。

### Brook

因为 Brook 的启动与 ss 一样简单，故不再给出示例。

### V2Ray

V2Ray 同样没有提供“贴心”的用户限制/限速服务，如果需要，请在阅读完本段后下翻至“使用 iptables 针对每个用户限制流量”段落。

V2Ray 要实现多用户多端口（同时还要考虑到不能重启现有进程以免影响）就需要多个配置文件，这时我们可以通过一个脚本来实现生成配置+启动的过程：

```python
#!/usr/bin/env python3
import sys, os, subprocess, json, shlex

PROXY_CMD="/path/to/your/program"
CONF_PATH="/tmp/sshare/"
CONF_TPL="/tmp/sshare/tpl.conf"

port, pw = sys.argv[1:3]

def run_p(cmd,data):
    p = subprocess.Popen(shlex.split(cmd), stdin=subprocess.PIPE, encoding="u8")
    json.dump(data, p.stdin)
    p.stdin.close()
    p.wait()

with open(CONF_TPL) as tpl:
    c = json.load(tpl)
c["inbound"]["settings"]["clients"][0]["id"] = pw
c["inbound"]["port"] = port
run_p(PROXY_CMD+"-config stdin", c)
```

配置项可参考下面“[使用 iptables 针对每个用户限制流量](#使用 iptables 针对每个用户限制流量)”段落。由于 V2Ray 的特别机制，请记得启用`gen_uuid`选项。

同时同端口多用户也是可行的，方案给出思路如下：

- start脚本中插入client项，stop删除，同时重启服务器。
- 多个服务器后端使用ws，设置不同的path，然后使用Nginx反代。
