package main

import (
	"crypto/tls"
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"net"
	"os"
	"strings"
	"time"
)

var port = "443"

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	domainsPointer := flag.String("domains", "", "-domains=domain1.com,domain2.com")
	flag.Parse()
	if *domainsPointer == "" {
		log.Fatal().Msg("No params")
	}
	for _, v := range strings.Split(*domainsPointer, ",") {
		scan(v)
	}
}

func scan(host string) {
	var sans []string
	host = strings.ToLower(host)
	conn, err := tls.DialWithDialer(&net.Dialer{
		Timeout: 5 * time.Second,
	}, "tcp", host+":"+port, &tls.Config{
		InsecureSkipVerify: true,
		MaxVersion:         0,
	})
	if err != nil {
		return
	}
	defer conn.Close()

	for _, v := range conn.ConnectionState().PeerCertificates {
		sans = append(sans, v.DNSNames...)
	}
	sans = lo.Uniq(sans)
	log.Info().Str("Domain", host).Strs("SANs", sans).Msg("Found SANs")
}
