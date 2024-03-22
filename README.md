# 一、使用
1. 先调用ip_create接口创建防火墙规则
2. 再调用ip_renewal接口给防火墙规则续期
3. 调用ip_monitor接口开启监控，如果key到期了，那就从redis中删除这个规则


# 二、编译
`SET CGO_ENABLED=0`

`SET GOARCH=amd64`

`SET GOOS=linux`

`go build iptables_manage`


# 三、能力
1. 支持部署到多个机器上，部署到哪一台机器上，就可以通过接口操纵哪一台机器上的iptables

# 四、部署
要使用 iptables 实现白名单功能，只允许通过的 IP 地址访问，而禁止其他 IP 地址访问，需要采取以下步骤：

1. 清除现有的 iptables 规则（可选）：如果你已经有现有的规则，可能需要清除它们以确保从一个干净的状态开始。可以使用以下命令清除所有规则：

   ```bash
   iptables -F
   iptables -X
   iptables -Z
   iptables -t nat -F
   iptables -t nat -X
   iptables -t nat -Z
   ```

2. 添加允许规则：通过SSH远程操作的情况，需要把操作机器的 IP 地址添加允许规则，以便后续能正常访问机器。

   ```bash
   iptables -A INPUT -s 172.19.128.1 -j ACCEPT
   ```

   将 `<允许的IP地址>` 替换为实际允许访问的 IP 地址。你可以根据需要添加多个规则，每个规则对应一个允许的 IP 地址。


3. 设置默认策略：将默认策略设置为拒绝（DROP）所有数据包。这将确保除了白名单中的 IP 地址之外的所有来源 IP 地址都被拒绝。

   ```bash
   iptables -P INPUT DROP
   iptables -P FORWARD DROP
   iptables -P OUTPUT ACCEPT
   ```

4. 保存规则（可选）：如果你希望在系统重启后保留 iptables 规则，可以将规则保存到适当的防火墙规则配置文件中。具体的保存方法可能因你使用的 Linux 发行版和防火墙软件而异。

   在大多数情况下，可以使用以下命令将规则保存到 `/etc/iptables/rules.v4` 文件中：

   ```bash
   iptables-save > /etc/iptables/rules.v4
   ```

   这样，规则将在系统重启后自动加载。 以上步骤将根据白名单中的 IP 地址允许或拒绝访问。只有在白名单中的 IP 地址才能访问系统，其他 IP 地址将被拒绝。

请注意，使用 iptables 配置防火墙规则需要以 root 用户或具有适当权限的用户身份运行命令，所以启动程序需要使用root用户来操作。
