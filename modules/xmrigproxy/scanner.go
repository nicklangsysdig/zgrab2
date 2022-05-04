// Package xmrigproxy contains the zgrab2 Module implementation for xmrig-proxy.
// The scan performs a banner grab
// The output is the banner
package xmrigproxy

import (
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/zmap/zgrab2"
)

// ScanResults is the output of the scan.
// Identical to the original from zgrab, with the addition of TLSLog.
type ScanResults struct {
	// Banner is the initial data banner sent by the server.
	Banner string `json:"banner,omitempty"`
}

type Flags struct {
	zgrab2.BaseFlags
}

// Module implements the zgrab2.Module interface.
type Module struct {
}

// Scanner implements the zgrab2.Scanner interface, and holds the state
// for a single scan.
type Scanner struct {
	config *Flags
}

// Connection holds the state for a single connection to the FTP server.
type Connection struct {
	// buffer is a temporary buffer for sending commands -- so, never interleave
	// sendCommand calls on a given connection
	buffer  [10000]byte
	config  *Flags
	results ScanResults
	conn    net.Conn
}

// RegisterModule registers the ftp zgrab2 module.
func RegisterModule() {
	var module Module
	_, err := zgrab2.AddCommand("xmrigproxy", "xmrig-proxy", module.Description(), 3333, &module)
	if err != nil {
		log.Fatal(err)
	}
}

// NewFlags returns the default flags object to be filled in with the
// command-line arguments.
func (m *Module) NewFlags() interface{} {
	return new(Flags)
}

// NewScanner returns a new Scanner instance.
func (m *Module) NewScanner() zgrab2.Scanner {
	return new(Scanner)
}

// Description returns an overview of this module.
func (m *Module) Description() string {
	return "Grab an XMRIG-PROXY banner"
}

// Validate flags
func (f *Flags) Validate(args []string) (err error) {
	return
}

// Help returns this module's help string.
func (f *Flags) Help() string {
	return "most xmrig-proxy instances run on 3333, 14444, 1444, or 1414!"
}

// Protocol returns the protocol identifer for the scanner.
func (s *Scanner) Protocol() string {
	return "xmrigproxy"
}

// Init initializes the Scanner instance with the flags from the command
// line.
func (s *Scanner) Init(flags zgrab2.ScanFlags) error {
	f, _ := flags.(*Flags)
	s.config = f
	return nil
}

// InitPerSender does nothing in this module.
func (s *Scanner) InitPerSender(senderID int) error {
	return nil
}

// GetName returns the configured name for the Scanner.
func (s *Scanner) GetName() string {
	return s.config.Name
}

// GetTrigger returns the Trigger defined in the Flags.
func (scanner *Scanner) GetTrigger() string {
	return scanner.config.Trigger
}

// isOKResponse returns true iff and only if the given response code indicates
// success (e.g. 2XX)
func (xmrigp *Connection) isOKResponse(retCode string) bool {
	return strings.HasPrefix(retCode, "{\"jsonrpc\":\"2.0\",\"method\":\"job\",")
}

// readResponse reads an FTP response chunk from the server.
// It returns the full response, as well as the status code alone.
func (xmrigp *Connection) readResponse() (string, string, error) {
	var buf []byte
	buf, err := zgrab2.ReadAvailable(xmrigp.conn)
	log.Warn(buf)
	if err != nil {
		return "", "", err
	}
	for i := range buf {
		xmrigp.buffer[i] = buf[i]
	}
	ret := string(xmrigp.buffer[0:31])
	retCode := ret
	return ret, retCode, nil
}

// GetFTPBanner reads the data sent by the server immediately after connecting.
// Returns true if and only if the server returns a success status code.
// Taken over from the original zgrab.
func (xmrigp *Connection) GetXmrigProxyBanner() (bool, error) {
	banner, retCode, err := xmrigp.readResponse()
	log.Info("HELP", banner, retCode, err)
	if err != nil {
		return false, err
	}
	xmrigp.results.Banner = banner
	return xmrigp.isOKResponse(retCode), nil
}

// Scan performs the configured scan on the FTP server, as follows:
// * Read the banner into results.Banner (if it is not a 2XX response, bail)
// * If the FTPAuthTLS flag is not set, finish.
// * Send the AUTH TLS command to the server. If the response is not 2XX, then
//   send the AUTH SSL command. If the response is not 2XX, then finish.
// * Perform ths TLS handshake / any configured TLS scans, populating
//   results.TLSLog.
// * Return SCAN_SUCCESS, &results, nil
func (s *Scanner) Scan(t zgrab2.ScanTarget) (status zgrab2.ScanStatus, result interface{}, thrown error) {
	exampleConfigSend := `{"id":1,"jsonrpc":"2.0","method":"login","params":{"login":"x","pass":"x","agent":"XMRig/6.17.0 (Macintosh; macOS; x86_64) libuv/1.41.0 clang/10.0.0","algo":["cn/1","cn/2","cn/r","cn/fast","cn/half","cn/xao","cn/rto","cn/rwz","cn/zls","cn/double","cn/ccx","cn-lite/1","cn-heavy/0","cn-heavy/tube","cn-heavy/xhv","cn-pico","cn-pico/tlo","cn/upx2","rx/0","rx/wow","rx/arq","rx/graft","rx/sfx","rx/keva","argon2/chukwa","argon2/chukwav2","argon2/ninja","astrobwt","astrobwt/v2","ghostrider"]}}` + "\x0a"
	var err error
	conn, err := t.Open(&s.config.BaseFlags)
	if err != nil {
		return zgrab2.TryGetScanStatus(err), nil, err
	}
	cn := conn
	defer func() {
		cn.Close()
	}()
	conn.Write([]byte(exampleConfigSend))
	results := ScanResults{}

	xmrigp := Connection{conn: cn, config: s.config, results: results}
	_, err = xmrigp.GetXmrigProxyBanner()
	if err != nil {
		return zgrab2.TryGetScanStatus(err), &xmrigp.results, err
	}

	return zgrab2.SCAN_SUCCESS, &xmrigp.results, nil
}
