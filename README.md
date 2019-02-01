# cni_plugin
```
[root@192 demo]# ip netns add ns1
[root@192 demo]# ip netns list
ns1 (id: 0)
[root@192 demo]# cat config 
{
	"name": "mynet",
	"BridgeName": "test",
	"IP": "192.0.2.1/24"
}
[root@192 demo]# CNI_COMMAND=ADD CNI_CONTAINERID=ns1 CNI_NETNS=/var/run/netns/ns1 CNI_IFNAME=eth10 CNI_PATH=/root/go/src/demo/ ./main  < config
{test 192.0.2.1/24}
{Name:vethe31ac0ed Mac: Sandbox:}
ipv4Addr:  192.0.2.1
ipv4Net:  192.0.2.1/24
addr:  192.0.2.1/24
[root@192 demo]# ifconfig
ens33: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 192.168.108.138  netmask 255.255.255.0  broadcast 192.168.108.255
        inet6 fe80::b588:5cc5:ac18:e0c8  prefixlen 64  scopeid 0x20<link>
        ether 00:0c:29:56:62:56  txqueuelen 1000  (Ethernet)
        RX packets 793036  bytes 1138033626 (1.0 GiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 123135  bytes 11802773 (11.2 MiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536
        inet 127.0.0.1  netmask 255.0.0.0
        inet6 ::1  prefixlen 128  scopeid 0x10<host>
        loop  txqueuelen 1  (Local Loopback)
        RX packets 4109  bytes 279988 (273.4 KiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 4109  bytes 279988 (273.4 KiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

test: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet6 fe80::8c9f:f7ff:fef0:edf0  prefixlen 64  scopeid 0x20<link>
        ether d2:93:fc:70:1d:0f  txqueuelen 1000  (Ethernet)
        RX packets 7  bytes 480 (480.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 7  bytes 578 (578.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

tun0: flags=4305<UP,POINTOPOINT,RUNNING,NOARP,MULTICAST>  mtu 1500
        inet 198.18.40.237  netmask 255.255.248.0  destination 198.18.40.237
        inet6 fe80::e6b6:8eb1:ba21:bcd0  prefixlen 64  scopeid 0x20<link>
        unspec 00-00-00-00-00-00-00-00-00-00-00-00-00-00-00-00  txqueuelen 100  (UNSPEC)
        RX packets 8985  bytes 8126832 (7.7 MiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 7262  bytes 1682751 (1.6 MiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

vethe31ac0ed: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet6 fe80::d093:fcff:fe70:1d0f  prefixlen 64  scopeid 0x20<link>
        ether d2:93:fc:70:1d:0f  txqueuelen 0  (Ethernet)
        RX packets 7  bytes 578 (578.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 13  bytes 1066 (1.0 KiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

virbr0: flags=4099<UP,BROADCAST,MULTICAST>  mtu 1500
        inet 192.168.122.1  netmask 255.255.255.0  broadcast 192.168.122.255
        ether 52:54:00:d4:07:2d  txqueuelen 1000  (Ethernet)
        RX packets 0  bytes 0 (0.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 0  bytes 0 (0.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

[root@192 demo]# ip netns exec ns1 ifconfig
eth10: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 192.0.2.1  netmask 255.255.255.0  broadcast 192.0.2.255
        inet6 fe80::c0ea:5aff:fe9e:328  prefixlen 64  scopeid 0x20<link>
        ether c2:ea:5a:9e:03:28  txqueuelen 0  (Ethernet)
        RX packets 15  bytes 1206 (1.1 KiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 8  bytes 648 (648.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

[root@192 demo]# ip netns exec ns1 ifconfig eth10 192.0.2.2
[root@192 demo]# ifconfig test 192.0.2.1
[root@192 demo]# ip netns exec ns1 ping 192.0.2.1
PING 192.0.2.1 (192.0.2.1) 56(84) bytes of data.
64 bytes from 192.0.2.1: icmp_seq=1 ttl=64 time=0.065 ms
64 bytes from 192.0.2.1: icmp_seq=2 ttl=64 time=0.063 ms
64 bytes from 192.0.2.1: icmp_seq=3 ttl=64 time=0.132 ms
64 bytes from 192.0.2.1: icmp_seq=4 ttl=64 time=0.087 ms
^C
--- 192.0.2.1 ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 2999ms
rtt min/avg/max/mdev = 0.063/0.086/0.132/0.030 ms
[root@192 demo]# 
[root@192 demo]# 
[root@192 demo]# CNI_COMMAND=DEL CNI_CONTAINERID=ns1 CNI_NETNS=/var/run/netns/ns1 CNI_IFNAME=eth10 CNI_PATH=/root/go/src/demo/ ./main  < config
[root@192 demo]# ip netns exec ns1 ifconfig
[root@192 demo]# ifconfig
ens33: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 192.168.108.138  netmask 255.255.255.0  broadcast 192.168.108.255
        inet6 fe80::b588:5cc5:ac18:e0c8  prefixlen 64  scopeid 0x20<link>
        ether 00:0c:29:56:62:56  txqueuelen 1000  (Ethernet)
        RX packets 793072  bytes 1138039188 (1.0 GiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 123157  bytes 11806316 (11.2 MiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536
        inet 127.0.0.1  netmask 255.0.0.0
        inet6 ::1  prefixlen 128  scopeid 0x10<host>
        loop  txqueuelen 1  (Local Loopback)
        RX packets 4131  bytes 281451 (274.8 KiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 4131  bytes 281451 (274.8 KiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

tun0: flags=4305<UP,POINTOPOINT,RUNNING,NOARP,MULTICAST>  mtu 1500
        inet 198.18.40.237  netmask 255.255.248.0  destination 198.18.40.237
        inet6 fe80::e6b6:8eb1:ba21:bcd0  prefixlen 64  scopeid 0x20<link>
        unspec 00-00-00-00-00-00-00-00-00-00-00-00-00-00-00-00  txqueuelen 100  (UNSPEC)
        RX packets 9003  bytes 8128928 (7.7 MiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 7279  bytes 1684295 (1.6 MiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

virbr0: flags=4099<UP,BROADCAST,MULTICAST>  mtu 1500
        inet 192.168.122.1  netmask 255.255.255.0  broadcast 192.168.122.255
        ether 52:54:00:d4:07:2d  txqueuelen 1000  (Ethernet)
        RX packets 0  bytes 0 (0.0 B)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 0  bytes 0 (0.0 B)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

[root@192 demo]#
```

