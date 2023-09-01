package main

import (
	"time"
	"crypto/tls"
	"net/http"
	"github.com/rs/zerolog/log"

	"github.com/unitoftime/flow/net"
)

type Hello struct {
	Msg string
}

type S struct {
	union *net.UnionBuilder
}
func (s *S) Marshal(v any) ([]byte, error) {
	return s.union.Serialize(v)
}
func (s *S) Unmarshal(dat []byte) (any, error) {
		return s.union.Deserialize(dat)
}

func main() {
	certPem := []byte(`-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIQIRi6zePL6mKjOipn+dNuaTAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMB4XDTE3MTAyMDE5NDMwNloXDTE4MTAyMDE5NDMwNlow
EjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0d
7VNhbWvZLWPuj/RtHFjvtJBEwOkhbN/BnnE8rnZR8+sbwnc/KhCk3FhnpHZnQz7B
5aETbbIgmuvewdjvSBSjYzBhMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggr
BgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdEQQiMCCCDmxvY2FsaG9zdDo1
NDUzgg4xMjcuMC4wLjE6NTQ1MzAKBggqhkjOPQQDAgNIADBFAiEA2zpJEPQyz6/l
Wf86aX6PepsntZv2GYlA5UpabfT2EZICICpJ5h/iI+i341gBmLiAFQOyTDT+/wQc
6MF9+Yw1Yy0t
-----END CERTIFICATE-----`)
	keyPem := []byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIrYSSNQFaA2Hwf1duRSxKtLYX5CB04fSeQ6tF1aY/PuoAoGCCqGSM49
AwEHoUQDQgAEPR3tU2Fta9ktY+6P9G0cWO+0kETA6SFs38GecTyudlHz6xvCdz8q
EKTcWGekdmdDPsHloRNtsiCa697B2O9IFA==
-----END EC PRIVATE KEY-----`)
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		panic(err)
	}
	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	cfg := net.Config{
		Url: "webrtc://localhost:8080",
		Serdes: &S{net.NewUnion(Hello{})},
		TlsConfig: tlsCfg,
		HttpServer: &http.Server{
			TLSConfig: tlsCfg,
			// ReadTimeout: 10 * time.Second,
			// WriteTimeout: 10 * time.Second,
		},
	}

	listener, err := cfg.Listen()
	if err != nil {
		panic(err)
	}

	for {
		log.Print("Accepting")
		sock, err := listener.Accept()
		if err != nil {
			log.Print("Error", err)
			continue
		}

		log.Print("Accepted")
		go func() {
			for {
				msg, err := sock.Recv()
				if err != nil {
					log.Print("Server RecvError: ", err)
					return
				}
				log.Print("ServerRecv: ", msg)
			}
		}()

		go func() {
			for {
				log.Print("ServerSend")
				err := sock.Send(Hello{"Hi From Server"})
				if err != nil {
					log.Print("Server Send Error: ", err)
					// TODO - Close?
					return
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}
}
