package gateway

import (
	"errors"
	"fmt"
	"strconv"
	"tg-mc/conf"
	"tg-mc/models"
	"time"

	"github.com/Tnze/go-mc/net"
	pk "github.com/Tnze/go-mc/net/packet"
	"github.com/Tnze/go-mc/offline"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type LoginRequest struct {
	ChatID  int64
	Resolve chan bool
}

func requestVerify(mcid string) (bool, error) {

	user, err := models.GetUserByMCName(mcid)
	if err != nil {
		return false, errors.New("mcid not found in settings")
	}

	// Send a message to the user to accept or reject using your bot logic
	logrus.Infof("Requesting verification for MCID: %s with ChatID: %d\n", mcid, user.TGID)

	// The channel to await a response from the verification process
	loginRequest := &LoginRequest{
		ChatID:  user.TGID,
		Resolve: make(chan bool),
	}

	GetAuthcator().RequestAuth(user, loginRequest)

	select {
	case result := <-loginRequest.Resolve:
		return result, nil
	case <-time.After(10 * time.Second): // Replace with your actual timeout handling logic
		return false, nil
	}
}

// This function should be called when a new client connection is accepted
func HandleClientConnection(clientConn net.Conn) {
	settings := conf.GetBotSettings().GatewaySettings

	targetConns, err := net.DialMC(settings.ServerHost + ":" + strconv.Itoa(settings.ServerPort))

	targetConn := *targetConns

	// targetConn, err := gonet.Dial("tcp", settings.ServerHost+":"+strconv.Itoa(settings.ServerPort))
	if err != nil {
		logrus.Errorf("Failed to connect to target: %s", err)
		return
	}
	protocol, intention, err := handshake(clientConn, targetConn)
	if err != nil {
		logrus.Errorf("Handshake error: %v", err)
		return
	}

	logrus.Infof("Protocol: %v, Intention: %v", protocol, intention)

	switch intention {
	default: // unknown error
		logrus.Errorf("Unknown handshake intention: %v", intention)
	case 1: // for status
		handleProxyConnection(clientConn, targetConn)
	case 2: // for login
		handlePlaying(clientConn, targetConn)
	}
}

type PlayerInfo struct {
	Name    string
	UUID    uuid.UUID
	OPLevel int
}

func handlePlaying(conn, target net.Conn) {
	// login, get player info
	info, err := acceptLogin(conn, target)
	if err != nil {
		logrus.Errorf("user [%s] Login failed", info.Name)
		return
	}

	logrus.Infof("Login successful: %s", info.Name)

	// Write LoginSuccess packet

	handleProxyConnection(conn, target)
}

// acceptLogin check player's account
func acceptLogin(conn, target net.Conn) (info PlayerInfo, err error) {
	// login start
	var p pk.Packet
	err = conn.ReadPacket(&p)
	if err != nil {
		return
	}

	err = p.Scan((*pk.String)(&info.Name)) // decode username as pk.String
	if err != nil {
		return
	}

	if ok, err := requestVerify(info.Name); !ok || err != nil {
		err = errors.New("verify failed")
		conn.Close()
		target.Close()
		return info, err
	}

	if err := target.WritePacket(p); err != nil {
		return info, err
	}

	// auth
	const OnlineMode = false
	if OnlineMode {
		info.UUID = offline.NameToUUID(info.Name)
	} else {
		// offline-mode UUID
		info.UUID = offline.NameToUUID(info.Name)
	}

	return
}

func handshake(conn, target net.Conn) (protocol, intention int32, err error) {
	var (
		p                   pk.Packet
		Protocol, Intention pk.VarInt
		ServerAddress       pk.String        // ignored
		ServerPort          pk.UnsignedShort // ignored
	)
	// receive handshake packet
	if err = conn.ReadPacket(&p); err != nil {
		return
	}
	err = p.Scan(&Protocol, &ServerAddress, &ServerPort, &Intention)
	if err != nil {
		return
	}
	err = target.WritePacket(p)

	return int32(Protocol), int32(Intention), err
}

func handleProxyConnection(clientConn net.Conn, targetConn net.Conn) {
	defer clientConn.Close()

	go func() {
		for {
			var p pk.Packet
			if err := clientConn.ReadPacket(&p); err != nil {
				logrus.Errorf("ReadPacket error: %v", err)
				break
			}
			if err := targetConn.WritePacket(p); err != nil {
				logrus.Errorf("WritePacket error: %v", err)
				break
			}
		}
	}()

	func() {
		for {
			var p pk.Packet
			if err := targetConn.ReadPacket(&p); err != nil {
				logrus.Errorf("ReadPacket error: %v", err)
				break
			}
			if err := clientConn.WritePacket(p); err != nil {
				logrus.Errorf("WritePacket error: %v", err)
				break
			}
		}
	}()
}

func startProxyServer(proxyHost string, proxyPort int) {

	settings := conf.GetBotSettings().GatewaySettings

	l, err := net.ListenMC(fmt.Sprintf(":%d", settings.ProxyPort))
	if err != nil {
		logrus.Fatalf("Error starting proxy server: %s", err)
	}
	defer l.Close()
	logrus.Infof("Proxy server started on %s:%d", proxyHost, proxyPort)

	for {
		clientConn, err := l.Accept()
		if err != nil {
			logrus.Errorf("Error accepting connection: %s", err)
			continue
		}
		go HandleClientConnection(clientConn)
	}
}

func StartGateway() {
	settings := conf.GetBotSettings().GatewaySettings
	startProxyServer(settings.ProxyHost, settings.ProxyPort)
}
