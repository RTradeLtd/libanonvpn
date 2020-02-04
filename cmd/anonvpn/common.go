package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

import (
	samforwardervpn "github.com/RTradeLtd/libanonvpn/client"
	samtunnelhandler "github.com/RTradeLtd/libanonvpn/ctrl"
	samforwardervpnserver "github.com/RTradeLtd/libanonvpn/server"
	"github.com/eyedeekay/canal/etc"
	checki2p "github.com/eyedeekay/checki2cp"
	i2pbrowserproxy "github.com/eyedeekay/httptunnel/multiproxy"
	"github.com/eyedeekay/outproxy"
	pc "github.com/eyedeekay/portcheck"
	i2ptunconf "github.com/eyedeekay/sam-forwarder/config"
	"github.com/justinas/nosurf"
	"github.com/zieckey/goini"
)

func lbMain(ctx context.Context) {
	if !*skipi2cp {
		home := os.Getenv("I2CP_HOME")
		if len(home) == 0 {
			os.Setenv("I2CP_HOME", "/etc/anonvpn")
		}
		conf := os.Getenv("GO_I2CP_CONF")
		if len(conf) == 0 {
			os.Setenv("GO_I2CP_CONF", "/i2cp.conf")
		}

		if here, err := checki2p.CheckI2PIsRunning(); here == false {
			if err != nil {
				log.Fatal(err)
			}
			log.Println("I2P is not running")
			if here, err := checki2p.CheckI2PIsInstalledDefaultLocation(); here == false {
				if err != nil {
					log.Fatal(err)
				}
				log.Println("I2P may not be not installed. Please run with the -update parameter or start your router.")
			}
			log.Fatal("Install or start i2p")
		}
	}
	if *canal {
		if *client {
			if err := fwscript.Setup(*tunName); err != nil {
				log.Fatal("Client firewall configuration error", err)
			}
		} else {
			if err := fwscript.ServerSetup(*tunName, *gateName); err != nil {
				log.Fatal("Server firewall configuration error", err)
			}
		}
	}

	config := &i2ptunconf.Conf{Config: goini.New()}
	if *iniFile != "none" && *iniFile != "" {
		config, err = i2ptunconf.NewI2PTunConf(*iniFile)
	} else if config == nil {
		config = i2ptunconf.NewI2PBlankTunConf()
		*startUp = true
	} else {
		config = i2ptunconf.NewI2PBlankTunConf()
		*startUp = true
	}
	log.Println(targetHost, config)
	config.TargetHost = config.GetHost(*targetHost, "10.0.0.2")
	config.SaveFile = config.GetSaveFile(*saveFile, true)
	config.SaveDirectory = config.GetDir(*targetDir, "../")
	config.SamHost = config.GetSAMHost(*samHost, "127.0.0.1")
	config.SamPort = config.GetSAMPort(*samPort, "7656")
	config.TunName = config.GetKeys(*tunName, "forwarder")
	config.SigType = config.GetSigType(*sigType, "EdDSA_SHA512_Ed25519")
	config.InLength = config.GetInLength(*inLength, 1)
	config.OutLength = config.GetOutLength(*outLength, 1)
	config.InVariance = config.GetInVariance(*inVariance, 0)
	config.OutVariance = config.GetOutVariance(*outVariance, 0)
	config.InQuantity = config.GetInQuantity(*inQuantity, 5)
	config.OutQuantity = config.GetOutQuantity(*outQuantity, 5)
	config.InBackupQuantity = config.GetInBackups(*inBackupQuantity, 3)
	config.OutBackupQuantity = config.GetOutBackups(*outBackupQuantity, 3)
	config.InAllowZeroHop = config.GetInAllowZeroHop(*inAllowZeroHop, false)
	config.OutAllowZeroHop = config.GetOutAllowZeroHop(*outAllowZeroHop, false)
	config.UseCompression = config.GetUseCompression(*useCompression, true)
	config.ReduceIdle = config.GetReduceOnIdle(*reduceIdle, true)
	config.ReduceIdleTime = config.GetReduceIdleTime(*reduceIdleTime, 600000)
	config.ReduceIdleQuantity = config.GetReduceIdleQuantity(*reduceIdleQuantity, 2)
	config.AccessListType = config.GetAccessListType(*accessListType, "none")
	config.CloseIdle = config.GetCloseOnIdle(*closeIdle, false)
	config.CloseIdleTime = config.GetCloseIdleTime(*closeIdleTime, 600000)
	config.ClientDest = config.GetClientDest(*targetDest, "", "")
	config.UserName = config.GetUserName(*userName, "anon")
	config.Password = config.GetPassword(*password, "")
	config.TunnelHost = config.GetEndpointHost(*tunnelHost, "10.0.0.1")

	s.eepProxy, err = i2pbrowserproxy.NewHttpProxy(
		i2pbrowserproxy.SetProxyAddr("127.0.0.1", "7980"),
		i2pbrowserproxy.SetControlAddr("127.0.0.1", "7981"),
		i2pbrowserproxy.SetReduceIdleTime(6000000),
	)
	go s.eepProxy.Serve()

	s.clientMux = samtunnelhandler.NewTunnelHandlerMux("localhost", *webPort, config.UserName, config.Password, *webCSS, *webJS)
	if err != nil {
		panic(err)
	}
	for _, label := range config.Labels {
		t, e := config.Get("type", label)
		log.Printf("%s |%s|\n", label, t)
		if e {
			if t == "vpnclient" {
				log.Println("Loading", t, "config from file")
				if f, e := samtunnelhandler.NewTunnelHandler(samforwardervpn.NewSAMVPNClientForwarderFromConfig(
					*iniFile,
					*samHost,
					*samPort,
					label,
				)); e == nil {
					log.Println("found vpnclient under", label)
					s.clientMux = s.clientMux.Append(f)
				} else {
					log.Println(e.Error())
					return
				}
			}
			if t == "vpnserver" {
				log.Println("Loading", t, "config from file")
				if f, e := samtunnelhandler.NewTunnelHandler(samforwardervpnserver.NewSAMVPNForwarderFromConfig(
					*iniFile,
					*samHost,
					*samPort,
					label,
				)); e == nil {
					log.Println("found vpnserver under", label)
					s.clientMux = s.clientMux.Append(f)
				} else {
					log.Println(e.Error())
					return
				}
				//TODO: Replace this with similar helper to above so it builds it out of the file
				if f, e := samtunnelhandler.NewTunnelHandler(outproxy.NewOutProxyFromOptions(
					outproxy.SetSAMHost(*samHost),
					outproxy.SetSAMPort(*samPort),
					outproxy.SetType("outproxy"),
					outproxy.SetSaveFile(true),
					outproxy.SetName("vpnsocks"),
					outproxy.SetInLength(*inLength),
					outproxy.SetOutLength(*outLength),
					outproxy.SetInVariance(*inVariance),
					outproxy.SetOutVariance(*outVariance),
					outproxy.SetInQuantity(*inQuantity),
					outproxy.SetOutQuantity(*outQuantity),
					outproxy.SetInBackups(*inBackupQuantity),
					outproxy.SetOutBackups(*outBackupQuantity),
					outproxy.SetCompress(*useCompression),
					outproxy.SetReduceIdle(*reduceIdle),
					outproxy.SetReduceIdleTimeMs(*reduceIdleTime),
					outproxy.SetReduceIdleQuantity(*reduceIdleQuantity),
				)); e == nil {
					log.Println("found vpnserver under", label)
					s.clientMux = s.clientMux.Append(f)
				} else {
					log.Println(e.Error())
					return
				}
			}
		}
	}
	if port, err := strconv.Atoi(*webPort); err != nil {
		log.Println("Error:", err)
		return
	} else {
		if pc.SCL(port) {
			log.Println("Service found, launching GUI")
			RunUI()
			Exit()
			return
		}
	}
	if len(s.clientMux.Tunnels()) < 1 {
		if !*client {
			f, e := samtunnelhandler.NewTunnelHandler(samforwardervpnserver.NewSAMClientServerVPNFromOptions(
				samforwardervpnserver.SetClientFilePath(*clientConfig),
				samforwardervpnserver.SetFilePath(*iniFile),
				samforwardervpnserver.SetEncrypt(*encryptLeaseSet),
				samforwardervpnserver.SetLeaseSetKey(*leaseSetKey),
				samforwardervpnserver.SetLeaseSetPrivateKey(*leaseSetPrivateKey),
				samforwardervpnserver.SetLeaseSetPrivateSigningKey(*leaseSetPrivateSigningKey),
				samforwardervpnserver.SetType("server"),
				samforwardervpnserver.SetPointHost(*targetHost),
				samforwardervpnserver.SetSAMHost(*samHost),
				samforwardervpnserver.SetSAMPort(*samPort),
				samforwardervpnserver.SetEndpointHost(*tunnelHost),
				samforwardervpnserver.SetName(*tunName),
				samforwardervpnserver.SetInLength(*inLength),
				samforwardervpnserver.SetOutLength(*outLength),
				samforwardervpnserver.SetInVariance(*inVariance),
				samforwardervpnserver.SetOutVariance(*outVariance),
				samforwardervpnserver.SetInQuantity(*inQuantity),
				samforwardervpnserver.SetOutQuantity(*outQuantity),
				samforwardervpnserver.SetInBackups(*inBackupQuantity),
				samforwardervpnserver.SetOutBackups(*outBackupQuantity),
				samforwardervpnserver.SetFilePath(*iniFile),
				samforwardervpnserver.SetAllowZeroIn(*inAllowZeroHop),
				samforwardervpnserver.SetAllowZeroOut(*outAllowZeroHop),
				samforwardervpnserver.SetCompress(*useCompression),
				samforwardervpnserver.SetReduceIdle(*reduceIdle),
				samforwardervpnserver.SetReduceIdleTimeMs(*reduceIdleTime),
				samforwardervpnserver.SetReduceIdleQuantity(*reduceIdleQuantity),
				samforwardervpnserver.SetCloseIdle(*closeIdle),
				samforwardervpnserver.SetCloseIdleTimeMs(*closeIdleTime),
				samforwardervpnserver.SetAccessListType(*accessListType),
				samforwardervpnserver.SetAccessList(accessList),
				// Wallet Options
			))
			if e != nil {
				return
			}
			g, e := samtunnelhandler.NewTunnelHandler(outproxy.NewOutProxyFromOptions(
				outproxy.SetType("outproxy"),
				outproxy.SetSaveFile(true),
				outproxy.SetName("vpnsocks"),
				outproxy.SetInLength(*inLength),
				outproxy.SetOutLength(*outLength),
				outproxy.SetInVariance(*inVariance),
				outproxy.SetOutVariance(*outVariance),
				outproxy.SetInQuantity(*inQuantity),
				outproxy.SetOutQuantity(*outQuantity),
				outproxy.SetInBackups(*inBackupQuantity),
				outproxy.SetOutBackups(*outBackupQuantity),
				outproxy.SetCompress(*useCompression),
				outproxy.SetReduceIdle(*reduceIdle),
				outproxy.SetReduceIdleTimeMs(*reduceIdleTime),
				outproxy.SetReduceIdleQuantity(*reduceIdleQuantity),
			))
			if e != nil {
				return
			}
			if e == nil {
				log.Println("setup vpnserver under", *tunName)
				s.clientMux = s.clientMux.Append(f)
				s.clientMux = s.clientMux.Append(g)
			} else {
				return
			}
		} else if *client {
			f, e := samtunnelhandler.NewTunnelHandler(samforwardervpn.NewSAMClientVPNFromOptions(
				samforwardervpn.SetClientFilePath(*iniFile),
				samforwardervpn.SetEncrypt(*encryptLeaseSet),
				samforwardervpn.SetLeaseSetKey(*leaseSetKey),
				samforwardervpn.SetLeaseSetPrivateKey(*leaseSetPrivateKey),
				samforwardervpn.SetLeaseSetPrivateSigningKey(*leaseSetPrivateSigningKey),
				samforwardervpn.SetType("client"),
				samforwardervpn.SetPointHost(*targetHost),
				samforwardervpn.SetSAMHost(*samHost),
				samforwardervpn.SetSAMPort(*samPort),
				samforwardervpn.SetEndpointHost(*tunnelHost),
				samforwardervpn.SetName(*tunName),
				samforwardervpn.SetInLength(*inLength),
				samforwardervpn.SetOutLength(*outLength),
				samforwardervpn.SetInVariance(*inVariance),
				samforwardervpn.SetOutVariance(*outVariance),
				samforwardervpn.SetInQuantity(*inQuantity),
				samforwardervpn.SetOutQuantity(*outQuantity),
				samforwardervpn.SetInBackups(*inBackupQuantity),
				samforwardervpn.SetOutBackups(*outBackupQuantity),
				samforwardervpn.SetClientFilePath(*iniFile),
				samforwardervpn.SetClientDest(*targetDest),
				samforwardervpn.SetAllowZeroIn(*inAllowZeroHop),
				samforwardervpn.SetAllowZeroOut(*outAllowZeroHop),
				samforwardervpn.SetCompress(*useCompression),
				samforwardervpn.SetReduceIdle(*reduceIdle),
				samforwardervpn.SetReduceIdleTimeMs(*reduceIdleTime),
				samforwardervpn.SetReduceIdleQuantity(*reduceIdleQuantity),
				samforwardervpn.SetCloseIdle(*closeIdle),
				samforwardervpn.SetCloseIdleTimeMs(*closeIdleTime),
				samforwardervpn.SetAccessListType(*accessListType),
				samforwardervpn.SetAccessList(accessList),
				// Wallet Options
			))
			if e == nil {
				log.Println("setup client under", *tunName)
				s.clientMux = s.clientMux.Append(f)
			} else {
				return
			}
		}
	}
	s.clientMux.Handler = nosurf.New(s.clientMux.Handler)
	log.Println("Launching service")
	if s.Serve() {
		log.Println("Ran successfully")
		return
	}
	log.Println("Server never ran")
	return
}

func ChownR(path string, uid, gid int) error {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chown(name, uid, gid)
		}
		return err
	})
}
