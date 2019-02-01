package main

import (
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/containernetworking/plugins/pkg/ns"
	"fmt"
	"encoding/json"
	"github.com/vishvananda/netlink"
	"syscall"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/plugins/pkg/ip"
	"net"
	"runtime"
)

type DemoBridge struct {
	BridgeName string `json:"bridgeName"`
	IP string `json:"ip"`
}

func init()  {
	runtime.LockOSThread()
}

func createBridge(name string) (*netlink.Bridge, error) {

	// prepare bridge object
	br := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
			MTU: 1500,
			TxQLen: -1,
		},
	}

	// create a bridge
	err := netlink.LinkAdd(br)
	if err != nil && err != syscall.EEXIST {
		return nil, err
	}

	// up the linux bridge
	if err :=  netlink.LinkSetUp(br); err != nil {
		return nil, err
	}

	// get the bridge object from the Bridge we created before
	l, err := netlink.LinkByName(name)
	if err != nil {
		return nil, fmt.Errorf("could not find %q: %v", name, err)
	}

	newBr, ok := l.(*netlink.Bridge)
	if !ok {
		return nil, fmt.Errorf("%q already exists but is not a bridge", name)
	}

	return newBr, nil
}


// create veth pair and connect host veth to the bridge
func setupVeth(netns ns.NetNS, br *netlink.Bridge,ifName string) error {

	guestIface := &current.Interface{}
	hostIface := &current.Interface{}
	mtu := 1500

	// create veth pair
	err := netns.Do(
		func(hostNS ns.NetNS) error {
			hostVeth, containerVeth, err := ip.SetupVeth(ifName, mtu, hostNS)
			if err != nil {
				return err
			}

			guestIface.Name = containerVeth.Name
			guestIface.Mac = containerVeth.HardwareAddr.String()
			guestIface.Sandbox = netns.Path()
			hostIface.Name = hostVeth.Name
			return nil
		})

	if err != nil {
		return err
	}

	fmt.Println(hostIface)

	// get host veth object by name
	hostVeth, err := netlink.LinkByName(hostIface.Name)
	if err != nil {
		return fmt.Errorf("fail to lookup %q: %v", hostIface.Name, err)
	}

	hostIface.Mac = hostVeth.Attrs().HardwareAddr.String()

	// connect host veth to the  bridge
	if err := netlink.LinkSetMaster(hostVeth, br); err != nil {
		return fmt.Errorf("failed to connect %q to the bridge %v: %v", hostVeth.Attrs().Name, br.Attrs().Name, err)
	}

	return nil
}


func cmdAdd(args *skel.CmdArgs) error {

	db := DemoBridge{}
	if err := json.Unmarshal(args.StdinData, &db); err != nil {
		return err
	}

	// print DemoBridge configuration
	fmt.Println(db)

	// create a Linux bridge
	br, err := createBridge(db.BridgeName)
	if err != nil {
		return err
	}

	// get the namespace of the container
	netns, err := ns.GetNS(args.Netns)
	if err != nil {
		return err
	}

	// create veth pair and connect host veth to the bridge
	if err := setupVeth(netns, br, args.IfName); err != nil {
		return err
	}

	err = netns.Do(func(hostNS ns.NetNS) error {
		link, err := netlink.LinkByName(args.IfName)
		if err != nil {
			return err
		}
		// parse CIDR
		ipv4Addr, ipv4Net, err := net.ParseCIDR(db.IP)
		addr := &netlink.Addr{IPNet: ipv4Net, Label:""}
		ipv4Net.IP = ipv4Addr

		fmt.Println("ipv4Addr: ", ipv4Addr)
		fmt.Println("ipv4Net: ", ipv4Net)
		fmt.Println("addr: ", addr)

		// ifconfig eth10 x.y.z.w
		if err := netlink.AddrAdd(link, addr); err != nil {
			return fmt.Errorf("failed to add IP address %v to %q : %v", ipv4Net, args.IfName, err)
		}
		return nil
	})

	return err
}

func cmdDel(args *skel.CmdArgs) error {

	db := DemoBridge{}
	// get configuration
	if err := json.Unmarshal(args.StdinData, &db); err != nil {
		return err
	}

	// get link object(bridge)
	l, err := netlink.LinkByName(db.BridgeName)
	if err != nil {
		return fmt.Errorf("could not find %q: %v", db.BridgeName, err)
	}

	// delete the bridge
	if err := netlink.LinkDel(l); err != nil {
		return fmt.Errorf("could not delete %q: %v", db.BridgeName, err)
	}

	// delete links(veth pair) in ns
	ns.WithNetNSPath(args.Netns, func(_ ns.NetNS) error {
		var err error
		_, err = ip.DelLinkByNameAddr(args.IfName)
		if err != nil && err == ip.ErrLinkNotFound {
			return nil
		}
		return err
	})

	return nil
}

func cmdCheck(args *skel.CmdArgs) error {
	// TODO
	return fmt.Errorf("not implemented")
}

func main() {

	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, "CNI Demo v0.1")

}
