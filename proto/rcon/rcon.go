package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	InputPrefix  = "# "
	OutputPrefix = ">>> "
	_Expr        = "(?i)§[0-9a-fk-or]"
)

var (
	flagPwd, flagAddr           string
	flagVerbose, flagChatColors bool
)

var (
	ErrNoPwd = errors.New("You must first provide a password.")
	ErrRid   = errors.New("RequestIDs don't match.")
)

var re = regexp.MustCompile(_Expr)

func init() {
	flag.StringVar(&flagPwd, "p", "", "RCON password.")
	flag.StringVar(&flagAddr, "addr", "127.0.0.1:25575", "RCON service address.")
	flag.BoolVar(&flagVerbose, "v", false, "Enable verbosity.")
	flag.BoolVar(&flagChatColors, "cc", true, "Strip out chat colors.")

	flag.Parse()

	if !flagVerbose {
		log.SetOutput(ioutil.Discard)
	}

	// Check if address contains also port
	log.Println("Checking address...")
	if strings.Index(flagAddr, ":") == -1 {
		fmt.Println("Provided address has no port. Using default one (25575).")

		flagAddr += ":25575"
	}
	log.Println("Address seems OK.")
}

func main() {
	var err error

	fmt.Printf("Connecting to %q...\n", flagAddr)

	rcon, err := Dial("tcp", flagAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer rcon.Close()

	fmt.Println("Connected!")
	defer fmt.Println("Disconnecting.")

	// Authentication
	if flagPwd != "" {
		// Via flags
		_, err = rcon.Login(flagPwd)
		log.Printf("User provided password (%q) via flags.\n", flagPwd)
	} else {
		// Via console
		pwd := Repl("Please enter your RCON password: ")
		_, err = rcon.Login(pwd)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	// Command REPL
	for {
		cmd := Repl(InputPrefix)

		// Closing commands
		if cmd == "Q" || cmd == "exit" {
			fmt.Println("Caught exit message.")

			break
		}

		// REPL commands
		r, err := rcon.Cmd(cmd)
		if r == nil || err != nil {
			fmt.Println(err)

			break
		}

		if r.Payload != "" {
			// Ignore chat colors
			if !flagChatColors {
				temp := strings.Replace(r.Payload, "\xa7", "§", -1)
				r.Payload = re.ReplaceAllString(temp, "")
			}

			fmt.Printf("%s%q\n", OutputPrefix, r.Payload)
		} else {
			fmt.Println("<No response>")
		}

		// Exceptional case
		if cmd == "stop" {
			break
		}
	}
}

func Repl(prefix string) string {
	var line string
	buf := bufio.NewReader(os.Stdin)

	fmt.Print(prefix)

	for line == "" || line == "\n" {
		line, _ = buf.ReadString('\n')
	}

	// Remove space chars
	line = strings.TrimSpace(line)

	return line
}

type RCON struct {
	conn net.Conn
}

func Dial(network, address string) (*RCON, error) {
	// TODO(toqueteos): Customize timeout via flag
	conn, err := net.DialTimeout(network, address, 1*time.Second)
	if err != nil {
		return nil, err
	}

	r := &RCON{conn}

	return r, nil
}

func (r *RCON) Close() error {
	return r.conn.Close()
}

func (r *RCON) Addr() string {
	return r.conn.RemoteAddr().String()
}

type Header struct {
	Length    int32       // Length of remainder of packet
	RequestID int32       // Client-generated ID
	Type      MessageType // 3 for login, 2 for command
}

func (h Header) String() string {
	return fmt.Sprintf("Header{%d, %d, %d}", h.Length, h.RequestID, h.Type)
}

type Response struct {
	Header
	Payload string
}

func (r *RCON) Send(typ MessageType, payload string) (*Response, error) {
	const (
		// byte-length of RequestID + Type + Pad = [4] + [4] + [2]
		remSize = 10
		padSize = 2
	)

	var (
		err     error
		pad     [2]byte
		plBytes = []byte(payload)
		n       = len(plBytes)
	)

	psend := Header{
		Length:    int32(remSize + n),
		RequestID: 0,
		Type:      typ,
	}

	// Outgoing buffer (request packet)
	var buf bytes.Buffer

	// Write request packet to a buffer
	err = binary.Write(&buf, binary.LittleEndian, &psend)
	if err != nil {
		return nil, fmt.Errorf("binary.Write: %s", err)
	}
	log.Println("'packet header' has been written.")

	buf.Write(plBytes)
	log.Println("'payload' has been written.")

	buf.Write(pad[:])
	log.Println("'pad' has been written.")

	log.Printf("outgoing buffer: % x\n", buf.Bytes())

	// Write buffer to conn
	_, err = r.conn.Write(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("r.conn.Write: %s", err)
	}

	// Discard bytes of outgoing packet
	buf.Reset()

	// Read: [4]Length + [4]RequestID + [4]Type
	var pread Header
	err = binary.Read(r.conn, binary.LittleEndian, &pread)
	if err != nil {
		return nil, fmt.Errorf("binary.Read: %s", err)
	}
	log.Printf("Header (read): %s\n", pread)

	// Read: [n]Payload + [2]Pad
	var endSize = int64(pread.Length - remSize + padSize)
	log.Printf("Payload + Pad is %d bytes long.", endSize)

	var nn int64
	nn, err = io.CopyN(&buf, r.conn, endSize)
	if nn != endSize {
		return nil, fmt.Errorf("io.CopyN: read %d bytes, expected %d.", nn, endSize)
	}
	if err != nil {
		return nil, fmt.Errorf("io.CopyN: %s", err)
	}

	switch endSize {
	case 2:
		log.Println("'payload' not present, nothing to read.")
	default:
		log.Println("'payload' has been read.")
	}

	switch pread.RequestID {
	// Success!
	case psend.RequestID:
	// No password provided
	case int32(AuthFailed):
		return nil, ErrNoPwd
	}

	resp := &Response{
		Header:  pread,
		Payload: strings.TrimRight(buf.String(), "\x00"),
	}

	return resp, nil
}

type MessageType int32

const (
	AuthFailed MessageType = -1
	Command                = 2
	Login                  = 3
)

func (r *RCON) Login(pwd string) (*Response, error) {
	return r.Send(Login, pwd)
}

func (r *RCON) Cmd(cmd string) (*Response, error) {
	return r.Send(Command, cmd)
}
