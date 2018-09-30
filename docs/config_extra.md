## 上线之前

### 使用 iptables 针对每个用户限制流量

如果你选择了一个不提供流量限制的代理后端，你可能需要阅读这一段以添加相应的功能，其基本思想是：

> 在启动代理前使用 iptables 设置对应端口的策略，在退出代理时再进行清除。

在使用这个模板时，我们需要用户创建一个 python 脚本文件，并赋予执行权限（`chmod a+x`），内容如下：

请将对应的代理后端指令更换为你的，并按照实际情况灵活修改（如流量限制）。

> 请注意：该脚本可能需要以 root 权限运行，推荐使用将运行用户设为 sudoer，并在 iptables 命令前加 sudo 的方案来实现正常运行。

```python
#!/usr/bin/env python
from __future__ import print_function
import sys, os, subprocess

PROXY_CMD="/path/to/your/program"
QUOTA=1024*1024*1024 #1GBytes
port=sys.argv[2]

def run_cmds(*cmds):
	for cmd in cmds:
		subprocess.call(cmd.split(" "))

if sys.argv[1] == "run":
	run_cmds(
        "iptables -A OUTPUT -p tcp --dport %s -m quota --quota %s -j ACCEPT" % (port,QUOTA),
        "iptables -A OUTPUT -p tcp --dport %s -j DROP" % port,
		PROXY_CMD+" "+" ".join(sys.argv[3:]),
		)
elif sys.argv[1] == "exit":
	run_cmds(
		"iptables -D OUTPUT -p tcp --dport %s -m quota --quota %s -j ACCEPT" % (port,QUOTA),
        "iptables -D OUTPUT -p tcp --dport %s -j DROP" % port,
		)
```

在这个脚本中我们实现了启停的逻辑，并要求程序提供端口作为第二个参数、将从第三个开始的参数传递给后端代理程序，当然这对于 sShare 来说完全不是问题，所以只需要在配置中做出对应的配置即可：

```json
"run_command": { 
  "cmd": "/home/bobo/runproxy.py",
  "arg": "run {{.Port}} -p {{.Port}} -k {{.Pass pass_name}} -m aes-128-gcm", 
  "enabled": true
},
"exit_command": {
  "cmd": "/home/bobo/runproxy.py",
  "arg": "exit {{.Port}}",
  "enabled": true
},
```

### 使用 Nginx 反向代理 Web 服务

在上线时使用 Nginx 反代 sShare 的 Web 服务有助于更好地管理服务器上的 Web 服务，实现按域名提供内容，添加SSL等。

通常为了安全地进行这一步，我们需要：

- 将配置文件中的监听地址由":9527"修改为"127.0.0.1:9527"，以防止用户直接访问 sShare 本身。
- 在nginx中打开access log，或在nginx的配置中将客户端ip传递到sShare后端。

### 使用 iptables+tc 进行全局限速

用于对全局（全部代理端口）进行限速，参考内容来自网络。命令如下：

注意，请将下面的网卡名(eth0)，替换为自己的外网网卡名（通常来说，ovz是venet0:0，其他的是eth0）。

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

上面两篇文章也无法实现较高的安全性需求，如果要使用更加激进的策略，请：

```bash
# 添加已知用途的"安全"端口白名单，在下面的命令中选择适合你的一条运行（或者添加自己的端口白名单），将这里的2000:2200修改为你在config中设置的端口范围。
for i in 80 443; do iptables -A OUTPUT -p tcp --sport 2000:2200 --dport $i -j ACCEPT ; done # 如果仅允许访问网页
for i in 80 443 21 22 3389; do iptables -A OUTPUT -p tcp --sport 2000:2200 --dport $i -j ACCEPT ; done # 还允许ftp/ssh/远程桌面

# 禁止对其他的非白名单端口的访问
iptables -A OUTPUT -p tcp --sport 2000:2200 -j DROP
```
