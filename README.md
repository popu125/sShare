# sShare

> 生活不止眼前的苟，还有身后的苟。

sShare 是一个致力于用户身心健康发展的智障项目，其主要目的是帮助用户与“关系不太亲近、但可以还算可以的人群”分享自己的代理工具（如同事、同学等）。

sShare 在运行时将提供一个 web 页面，用于为用户分发代理账户。本 Repo 为其管理后端，但是Release包中包含了一个简单的 web 服务页面。

**推荐**与 sShare 搭配使用的代理工具为：ShadowsocksR，因为其内置用户数限定功能而不需配置额外工具。但是 sShare **完全支持其他代理工具**（如 SS 原分支、Brook 等），并且配备有**高自由度**的启停命令执行配置选项，同时为了方便普通用户，为常见工具提供了简单的配置模板。

sShare 不是多服务器（机场）出售/管理工具，sShare 是为**单服务器（本机）**准备的，**不记名**的“慈善型”服务工具，其最初设计原因是作者想要向学校里的陌生人提供代理服务。

sShare 在一定程度上考虑到了“回血”的需求，并且配备了两条路径：增加广告、或者使用挖矿验证码。通常我们建议您使用后者，尽管它收入较广告低得多得多得多得多，但是其准入门槛更低。

## 工作原理

sShare 的最初设计思想是“一人一进程，互相不干扰，来了开进程，到期就杀掉（考虑到进程人道主义的影响势大，我们会考虑使用SIGTERM）”。

且看图（图像使用 Microsoft Visio 创建）：

![](https://user-images.githubusercontent.com/7552030/34298930-1e712056-e75b-11e7-9979-e678db5f888a.png)

## 配置简解

sShare 的配置方式为单 json 文件，文件名为`config.json`。下面是示例配置文件及简解（请不要直接复制，json不支持注释的）：

```json
{
  "run_command": {   // 启动代理程序的配置，sShare会为每一个用户执行一次run_command
    "cmd": "ssserver", // 命令
    "arg": "-p {{port}} -k {{pass}}", // 参数，其中的{{port}}会被替换为端口，{{pass}}替换为密码，这两个参数都是随机生成的
    "enabled": true // 是否启用，run_command必须启用
  },
  "exit_command": {  // 当程序退出时，sShare应执行的任务，通常为清理临时文件、临时防火墙规则等
    "cmd": "exit", // 命令
    "arg": "{{port}}", // 参数，exit_command的参数中只有{{port}}会被替换
    "enabled": false // 是否启用，如果你不需要，你可以选择不启用exit_command
  },
  "captcha": { // 验证码选项
    "name": "base", // 验证码接口的名称，目前已实现的验证码见下方“配置示例”中“验证码”部分
    "site_id": "23333", // Site ID，该属性的含义因接口而异，详见示例
    "extra": "66666" // 额外数据，该属性的含义因接口而异，详见示例
  },
  "ttl": "20m", // 一个用户在获取到账号后可以使用的时间，超时后对应的进程将被Kill，用户需要重新在web界面获取
  "limit": 20, // 限制的最大用户数量
  "web_addr": ":9527", // Web界面监听地址，使用"ip:port"可以指定监听ip，使用":port"监听所有ip，建议监听本地并使用Nginx反代
  "port_start": 2000, // 分配端口起始值
  "port_range": 200, // 分配端口范围，最终用户得到的端口将在[port_start, port_start+port_range]之间，请务必保证该范围内端口没有被占用
  "rand_seed": 23343, // 随机种子，用于生成端口和密码，可以不设置
  "no_check_alive": false // 运行代理程序时不检查是否存活，具体用途参考下面的“配置实例”
}
```



## 配置示例

考虑到在真正执行代理程序前，可能需要执行其他命令，我们建议用户使用 bash/python 脚本作为启动命令，我们在下面的示例中（如果需要）将使用脚本。

### 基础模板

基础的配置模板，其他的模板都将以修改本模板和创建本地脚本的方式来创建。本文件位于 Repo 或者 Release 包解压后的根目录中，文件名为 `config.example.json`。

### 验证码

验证码是防止 sShare 惨遭爬虫光顾或者前期维持分享热情的重要手段，我们为您提供下面四种验证码接口，如需要其他请发 issue：

**注意**：在您使用验证码前，您可能需要访问对应的主页以获取验证码接口所需要的参数。

- base - 根本没有验证码，参数可以随便设置。
- coinbase - 最著名的挖矿验证码，site_id 参数代表 secret_id，extra 参数为您想要用户为您挖的 hash 数（hash 的含义不再做科普，请自行搜索）。要使用该验证码，请访问 https://coinbase.com 注册并创建一个 site。
- ppoi - 全称 ProjectPoi，coinbase 的仿制品，相比而言ppoi收取的手续费更低，ppoi 是目前第三大 xmr 验证码平台，由国人创建并运营，官网也是中文版。ppoi 的参数含义与coinbase是一致的。要使用该验证码，请访问 https://ppoi.org 注册并创建一个 site。
- recaptcha.v2 - 由谷歌提供的验证码服务，被认为是目前较为安全的验证码。在此处 site_id 输入 Site key，extra 设置为 Secret key。要使用 recaptcha 验证码，请先访问 https://www.google.com/recaptcha/admin 创建一个 site，目前 sShare 仅提供 v2 版本支持。

### ShadowsocksR

SSR 本身的管理功能已经较为完备，我们不再在模板中提供用户限制的额外脚本，故而本模板不需要创建额外的文件。

如果需要限速功能，请在阅读完本段后下翻至“基于 iptables 的端口连接限制和流量限制”段落。

在使用 SSR 作为 sShare 的后端时，您需要在命令行中将配置项作为参数写入，我们推荐的参数配置为：`-p {{port}} -k {{pass}} -m aes-128-cfb -o tls1.2_ticket_auth [-g 可选的混淆参数] -O auth_aes128_md5 [-G 可选的用户连接限制，推荐为1]` 

将其填入`run_command`配置项的`arg`条目，`cmd`中填写 SSR 的“单体版”脚本 server.py 的绝对路径（如我的就是`/home/bobo/ssr/shadowsocks/server.py`并赋予对应文件以执行权限(`chmod a+x`)，配置就宣告完成。

在此配置下 sShare 会为每一个用户自动启动一个对应的 ShadowsocksR 服务端进程，并在时间限制（配置详解中有说明）到期后结束该进程。

示例配置如下（仅包括 run_command 段，请禁用 exit_command）：

```json
"run_command": { 
  "cmd": "/home/bobo/ssr/shadowsocks/server.py",
  "arg": "-p {{port}} -k {{pass}} -m aes-128-cfb -O auth_aes128_md5 -o tls1.2_ticket_auth -G 1", 
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
  "arg": "-a -p {{port}} -k {{pass}} -m aes-128-cfb -O auth_aes128_md5 -o tls1.2_ticket_auth -G 1 -t 1", 
  "enabled": true
},
"exit_command": {
  "cmd": "/home/bobo/ssr/mujson_mgr.py",
  "arg": "-d -p {{port}}",
  "enabled": true
},
```

### Shadowsocks

Shadowsocks 原分支并没有“贴心”的用户限制/限速服务，如果需要，请在阅读完本段后下翻至“基于 iptables 的端口连接限制和流量限制”段落。

Shadowsocks 的配置更少，也更容易完成，此处不再给出示例，请参照 SSR（非 mujson）的配置和注释自行修改。

### Brook

考虑到用户少没人权，作者懒得写了，请参考上面ss的配置自行编写，谢谢❤

### Brook

同上desu

### 使用 iptables 进行（每用户）流量限制

用于为不提供流量限制的代理后端添加相应的功能，基本思想是：

> 在启动代理前使用 iptables 设置对应端口的流量策略，在退出代理时再进行清除。

在使用这个模板时，我们需要用户创建一个 python 脚本文件，并赋予执行权限（`chmod a+x`），内容如下：

```python
#!/usr/bin/env python
from __future__ import print_function
import sys, os, subprocess

PROXY_CMD="/path/to/your/program"
port=sys.argv[2]

def run_cmds(*cmds):
	for cmd in cmds:
		subprocess.call(cmd)

if sys.argv[1] == "run":
	run_cmds(
        "iptables -A FORWARD -s 208.8.14.36 -m limit --limit 700/s -j ACCEPT",
		PROXY_CMD+" "+" ".join(sys.argv[3:]),
		)
elif sys.argv[1] == "exit":
	run_cmds(
		"echo b", #todo
		)
```

在这个脚本中我们实现了启停的逻辑，并要求程序提供端口作为第二个参数、将从第三个开始的参数传递给后端代理程序，当然这对于 sShare 来说完全不是问题，所以只需要在配置中做出对应的配置即可：

```json
"run_command": { 
  "cmd": "/home/bobo/runproxy.py",
  "arg": "run {{port}} -p {{port}} -k {{pass}}", 
  "enabled": true
},
"exit_command": {
  "cmd": "/home/bobo/runproxy.py",
  "arg": "stop {{port}}",
  "enabled": true
},
```

## Web 配置

sShare 建议用户使用一个基于 ajax 技术实现的页面提供 web 服务，作者提供一个简单的使用 jq+bootstrap 实现的页面（vue真的好难，看了俩小时还是没搞懂啊）

关于这个简单页面的配置说明，请点击项目主页查看 readme：[sShare-simple](https://github.com/popu125/sShare-simple)

只需要在 sShare 所在目录创建一个 static 目录，并将静态页面文件放入，他们就会正常工作在 sShare 的 web 服务端口。

通常 sShare 的 release 包中已包括了这个简单的页面。

## 上线之前

### 使用 Nginx 反向代理 Web 服务

在上线时使用 Nginx 反代 sShare 的 Web 服务有助于更好地管理服务器上的 Web 服务，实现按域名提供内容，添加SSL等。

### 使用 iptables+tc 进行限速

用于对全局（全部代理端口）进行限速，参考内容来自网络。命令如下：

注意，请将下面的网卡名(eth0)，替换为自己的外网网卡名（通常来说，ovz是venet0，其他的是eth0）。

该规则会在 OS/iptables 重启后失效，如需要保持，可以通过`/etc/rc.local`做一个简单的开机自启。

```bash
tc qdisc add dev eth0 root handle 1:0 htb default 123
tc class add dev eth0 parent 1:0 classid 1:1 htb rate 100Mbit ceil 100Mbit prio 0
tc class add dev eth0 parent 1:1 classid 1:2 htb rate 10Mbit ceil 10Mbit prio 1 burst 96kbit # 在这里设置rate后面的值为要限制的速度，ceil后面为突发速度（临时可以达到的最大速度），celi >= rate
tc qdisc add dev eth0 parent 1:2 handle 111:0 sfq perturb 10
tc filter add dev eth0 parent 1:0 protocol ip prio 1 handle 9527 fw classid 1:2

iptables -A OUTPUT -t mangle -p tcp --sport 2000:2200 -j MARK --set-mark 9527 # 将这里的2000:2200修改为你在config中设置的端口范围。
```

### 防止 BT 下载

最好实行，以保护自身vps的安全，参考 https://www.dwhd.org/20150915_162703.html 和 https://dreamcreator108.com/dreams/iptables-ban-bt/index.html 。

## API 说明

目前 sShare 有两个公开api接口，分别用于查询当前客户端数和获取账号。文档如下：

### /api/count

请求类型：`POST`

直接请求即可得到目前已分配的客户端数量。

### /api/new

请求类型：`POST`

请求数据：

| Key   | Value    |
| ----- | -------- |
| token | 验证码服务返回值 |

如果验证通过则创建服务并返回一个 JSON 类型的返回值：

| Key  | Value |
| ---- | ----- |
| Port | 端口    |
| Pass | 密码    |

## 日志说明

sShare 有一个极其简陋的日志系统，在未指定日志保存文件时，日志会被输出到标准错误管道，而当你使用`-l [文件名]`参数指定保存文件后，日志会被保存在文件中。

日志格式为：`[模块名] 时间 日志内容`

各模块的日志内容为下：

### [MAIN]

这个模块几乎啥都不会输出，只会输出你设置的 Web 服务地址。

### [POOL]

这个模块用于记录程序池的进程创建与销毁，其中以`PROC_SPAWN`起始的是进程创建，以`CLEANUP`起始的是进程销毁。

其后面的两个数据分别为：端口和密码。

### [WEB]

这个模块用于记录 api 的访问，其中：

-  `ERR_NO_CAPTCHA` 为验证码验证出错，后面的数据为访问者的地址。
-  `ERR_FULL` 为用户池已满，后面的数据同样为访问者的地址。
-  `ACCEPT` 为成功分配代理，后面的数据分别是：访问者地址、端口、密码。

## LICENSE

AGPL 3.0，我们要求您保留原作者的LICENSE，同时您对该程序做出的任何修改（无论是用于分发还是仅提供服务）都必须开源。

额外条款：禁止商用。（尽管这破程序也不太适合商用的样子就是了）
