package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Config struct {
	ProxyPort           string `json:"proxy_port"`
	ServerTarget        string `json:"server_target"`
	TargetVersionString string `json:"target_version_string"`
}

func main() {
	var cfg Config
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		cfg = Config{
			ProxyPort:           "7778",
			ServerTarget:        "127.0.0.1:7777",
			TargetVersionString: "Terraria318",
		}
		data, _ := json.MarshalIndent(cfg, "", "  ")
		os.WriteFile("config.json", data, 0644)
		log.Println("已生成默认 config.json，请检查端口配置。")
	} else {
		if err := json.Unmarshal(configFile, &cfg); err != nil {
			log.Fatalf("配置解析失败: %v", err)
		}
	}

	listener, err := net.Listen("tcp", ":"+cfg.ProxyPort)
    if err != nil {
        log.Fatalf("监听失败: %v", err)
    }
    defer listener.Close()

    log.Printf("代理启动: [::]:%s -> %s (伪装版本: %s)", cfg.ProxyPort, cfg.ServerTarget, cfg.TargetVersionString)

    for {
        clientConn, err := listener.Accept()
        if err != nil {
            log.Printf("接受连接失败: %v", err) 
            continue
        }
        
        if tcpConn, ok := clientConn.(*net.TCPConn); ok {
            tcpConn.SetKeepAlive(true)
            tcpConn.SetKeepAlivePeriod(30 * time.Second)
        }

        go handleClient(clientConn, cfg)
    }
}

func handleClient(clientConn net.Conn, cfg Config) {
    
    
    defer clientConn.Close()

    remoteAddr := clientConn.RemoteAddr().String()

    
    serverConn, err := net.DialTimeout("tcp", cfg.ServerTarget, 5*time.Second)
    if err != nil {
        log.Printf("[%s] 无法连接后端服务器: %v", remoteAddr, err)
        return
    }
    defer serverConn.Close()

    
    if tcpConn, ok := serverConn.(*net.TCPConn); ok {
        tcpConn.SetKeepAlive(true)
        tcpConn.SetKeepAlivePeriod(30 * time.Second)
    }

    
    clientConn.SetReadDeadline(time.Now().Add(5 * time.Second))

    header := make([]byte, 2) 
    if _, err := io.ReadFull(clientConn, header); err != nil {
        return 
    }

    packetLen := binary.LittleEndian.Uint16(header)
    
    if packetLen > 4096 {
        log.Printf("[%s] 异常包长度: %d", remoteAddr, packetLen)
        return
    }

    body := make([]byte, int(packetLen)-2)
    if _, err := io.ReadFull(clientConn, body); err != nil {
        return
    }

    
    clientConn.SetReadDeadline(time.Time{})

    
    if len(body) > 0 && body[0] == 1 {
        clientVer := parseVersionString(body)
        log.Printf("[%s] 握手: %s -> %s", remoteAddr, clientVer, cfg.TargetVersionString)

        newPacket := buildFakeConnectPacket(cfg.TargetVersionString)
        if _, err := serverConn.Write(newPacket); err != nil {
            return
        }
    } else {
        
        if _, err := serverConn.Write(header); err != nil { return }
        if _, err := serverConn.Write(body); err != nil { return }
    }

    
    
    
    var once sync.Once
    closeLog := func(err error, src string) {
        once.Do(func() {
            
            if err != io.EOF {
                
            }
        })
    }


    pipe := func(dst, src net.Conn, tag string) {

        buf := make([]byte, 32*1024) 
        _, err := io.CopyBuffer(dst, src, buf)
        
        closeLog(err, tag)
        

        src.Close() 
        dst.Close()
    }

    go pipe(serverConn, clientConn, "client->server")
    pipe(clientConn, serverConn, "server->client")
}


func parseVersionString(body []byte) string {
	if len(body) < 2 { return "Unknown" }
	slice := body[1:]
	strLen, bytesRead := read7BitEncodedInt(slice)
	if strLen > 0 && strLen <= len(slice)-bytesRead {
		return string(slice[bytesRead : bytesRead+strLen])
	}
	return "Unknown"
}

func buildFakeConnectPacket(versionStr string) []byte {
	buf := new(bytes.Buffer)
	buf.Write([]byte{0, 0, 1}) 
	write7BitEncodedInt(buf, len(versionStr))
	buf.WriteString(versionStr)
	data := buf.Bytes()
	binary.LittleEndian.PutUint16(data[0:2], uint16(len(data)))
	return data
}

func write7BitEncodedInt(w *bytes.Buffer, value int) {
	v := uint32(value)
	for v >= 0x80 {
		w.WriteByte(byte(v | 0x80))
		v >>= 7
	}
	w.WriteByte(byte(v))
}

func read7BitEncodedInt(data []byte) (int, int) {
	count, shift := 0, 0
	for i, b := range data {
		count |= (int(b) & 0x7F) << shift
		shift += 7
		if (b & 0x80) == 0 { return count, i + 1 }
	}
	return 0, 0
}
