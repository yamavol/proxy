package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const privateKey = `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDJol2kYsFJeGQK
tNEZhpjwMPjBGxSg/UNbX2F6JNoDI/lUgrYdZpva828lH60KtIHNjrGAZafY/K1t
vN7lXel4YaBrhcPQIMJ6cYHMq85aqyt81LNhJu2NKQIZvuTphiTlNCdWPz9Bjf+0
q4EDIQsim4Dsc3byCJRAO1fuSPVlcxRAY/cqQMJQEzSr81lt6pAErnVxVe/4NetY
XHCkPUp+RUKoFSvyzBRj7gjIrjBxKAek5TrumnB0sTUXdq9oBi5PcQ3Ql1H5G3R9
5zgKl9VGHkKxKTcj+c/wcn7o3gdxaOHLRSk8/pzzSSed6EccegWBn32Yh56bsLBC
9CYwb/JfAgMBAAECggEAD/PRlszdX/Ovbo1psaxNc0tckuKSmj4PUy5Tpvc9bFwv
QLlqsR7KG+OAmp5L8XngfyPX9UGVqvwquHDl7Z2leAm5SGh32oKNAGT0kP3SwKek
NCcb3gbXaoChEupgb/1V8/BRYGh2l7glT+T1uwqlN+K3q31jHrkBCafoSAjrqU/y
fTJxyYg41jTywPUVvdVck4W4/pNCr+nB0nfDXzRZPwPH+34yyAWb2wNh9NJ38DML
agiaDab49cP31EnRqA6gK6J4Em7i3+5JRG7w+qbV0/qbw2xBcTTbwqI0y/kKf0Mv
ybHa2pQwk+EALu7Ao2jQLNx+oN1639H6fuY4sktLgQKBgQDvffBubrxaSkVNXJKy
XwzbCyGfrhvIVTgixkxsmQ2HEBmSY2l24Tfj9dl/M3PpHM3wDtW30feawXai4YaZ
yUKC61tSgfNu+4xjfCV5n8GOGwNjE3CKUY1Rpna1VVcY+uslfytzqrzhiyDpHhS+
TZgV6UMdq9vYAjvVU97oG8a+QQKBgQDXiGRHUG0iI3BX4rcSDcZVS5f3xBpxvL4B
ooY5HxkHjCOKx8wIf5zJQjmeQpr2kkRJy3tKH2ZSh0pKtKWVQYioRID8FvMERuti
b+Xp/A2AvydUCVqgNYA+hh8qrzvDJcgWzd7tlqYZm4axJayM/zJWAuAHYW7pV+Vp
Fdranh3InwKBgBQ5/7dj9NZvVWEOQ3l7G5vYWdOhockOoXoWY1f8qS7SBkbdzId0
yAKhvefHUa/Ldf0jU5t9yTqxwjJJd9O/MrXZ6NGUFho2donkb0nRW0iEMYoJl0Sn
VJcjxvzTo1KBxqBZGDNhpSgrVvE5UCkuZnzbQYbc/+lDbwg6WCYkSmnBAoGAXjtO
hHNgU3WlD3eazLTjCrWzKms9mI6JkBNrlZvICKm3fFygEvMgLEndARljwPvwCUeC
jsStqtVloMXcQyZUxiS1NAIgm7UaAn6jyaoeiTSJ0E8KpVLez/c5tyLIASkKkxXN
Kpkb48RAnkC3cSm96yb0paVupWx9a3VXqw9IPEcCgYAcywJucHFBX1Le00EvusNq
wkjj2Y9bW+ZkBNCugcoV4E1yAd6VVrs36ACAySrg7qF0Nq+eX6Mzr2uNILxVpijv
uFRbeV805P0+EnBJGpeCBhE8cwezg41DqDuG4JttM/cD7lGGUXm1Xutn8XN5Gx65
yvALzH8PbqY9PQtU7RfIXg==
-----END PRIVATE KEY-----
`

const certificate = `-----BEGIN CERTIFICATE-----
MIIDejCCAmKgAwIBAgIEYrGB5DANBgkqhkiG9w0BAQsFADBWMQswCQYDVQQGEwJK
UDESMBAGA1UEAwwJbG9jYWxob3N0MRUwEwYDVQQHDAxEZWZhdWx0IENpdHkxHDAa
BgNVBAoME0RlZmF1bHQgQ29tcGFueSBMdGQwHhcNMjQwMzEyMDgyOTEwWhcNMzQw
MzEwMDgyOTEwWjBWMQswCQYDVQQGEwJKUDESMBAGA1UEAwwJbG9jYWxob3N0MRUw
EwYDVQQHDAxEZWZhdWx0IENpdHkxHDAaBgNVBAoME0RlZmF1bHQgQ29tcGFueSBM
dGQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDJol2kYsFJeGQKtNEZ
hpjwMPjBGxSg/UNbX2F6JNoDI/lUgrYdZpva828lH60KtIHNjrGAZafY/K1tvN7l
Xel4YaBrhcPQIMJ6cYHMq85aqyt81LNhJu2NKQIZvuTphiTlNCdWPz9Bjf+0q4ED
IQsim4Dsc3byCJRAO1fuSPVlcxRAY/cqQMJQEzSr81lt6pAErnVxVe/4NetYXHCk
PUp+RUKoFSvyzBRj7gjIrjBxKAek5TrumnB0sTUXdq9oBi5PcQ3Ql1H5G3R95zgK
l9VGHkKxKTcj+c/wcn7o3gdxaOHLRSk8/pzzSSed6EccegWBn32Yh56bsLBC9CYw
b/JfAgMBAAGjUDBOMB0GA1UdDgQWBBQ4fSlYIqbxVcGUiGCd1G4H7a3OszAfBgNV
HSMEGDAWgBQ4fSlYIqbxVcGUiGCd1G4H7a3OszAMBgNVHRMEBTADAQH/MA0GCSqG
SIb3DQEBCwUAA4IBAQBkwJQSJaNsyy6GJHJJ2vSQy8onrUznjx585VJWBNc9CsaH
VJF2CDzUS+TZMRq4Ua5QomHMJUgbnDlbdFa7uPogyEtMvaXKoThzFb6CWjwacVvu
Jtf30Flr4XPY4epEi9x0jgU72TGXTuGEmLMDK0TZbEO2lJVCPta4Fg8f7eSXGl15
jPf18vmo36gSx6cDFBtVuhcBvOIYDv501dyU1xeBYhNy+DkkwtrNIpY9/sv/sSOc
oEOQdMcdZRPDaAZ/ZbxmTSIWTIUq8D1veprYYjyLIikfhmn/+JgN45CmeR8TQRoN
JTdumKGtMxLZIdQX8rZgSDM3KmS0+QEwdyMcYctJ
-----END CERTIFICATE-----`

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!"))
}

func main() {

	cert, err := tls.X509KeyPair([]byte(certificate), []byte(privateKey))

	if err != nil {
		panic(err)
	}

	serverPort := 443
	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%v", serverPort),
		Handler: http.HandlerFunc(helloHandler),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	fmt.Println("https server is running on port", serverPort)
	if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
		fmt.Println("Error:", err)
	}
}
