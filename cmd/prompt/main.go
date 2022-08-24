// Copyright (c) 2022 The Swedish Internet Foundation
//
// Distributed under the MIT License. (See accompanying LICENSE file or copy at
// <https://opensource.org/licenses/MIT>.)

package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"log"
	"net"
	"os"

	"github.com/dotse/epp-client/pkg"
	"github.com/dotse/epp-client/pkg/prompt"
	"github.com/dotse/epp-client/pkg/validation"

	"github.com/alecthomas/kingpin"
)

func main() {
	var (
		port = kingpin.
			Flag("port", "the port to send requests to").
			Short('p').
			Default("7000").
			String()

		host = kingpin.
			Flag("host", "the host to send requests to").
			Short('h').
			Default("127.0.0.0").
			String()

		cert = kingpin.
			Flag("cert", "path to the cert to use for tls").
			Short('c').
			Default("some-cert-path.cert").
			String()

		key = kingpin.
			Flag("key", "path to the key to use for tls").
			Short('k').
			Default("some-key-path.key").
			String()

		keepAlive = kingpin.
				Flag("keep-alive", "keep connection to the epp server alive").
				Short('a').
				Default("false").
				Bool()

		validateResponses = kingpin.
					Flag("validate-responses", "validate responses from epp server").
					Short('v').
					Default("true").
					Bool()
	)

	kingpin.Parse()

	logger := log.New(os.Stdout, "", 0)
	ctx := context.Background()

	cl := connect(host, port, cert, key)
	logger.Println(cl.Greeting)

	if *keepAlive {
		cl.KeepAlive(ctx)
	}

	scnr := bufio.NewScanner(os.Stdin)
	scnr.Split(func(data []byte, _ bool) (int, []byte, error) {
		if i := bytes.Index(data, []byte{'\n', '\n'}); i >= 0 {
			return i + 2, data[0:i], nil
		}

		return 0, nil, nil
	})

	p := prompt.Prompt{
		Client:           cl,
		MultilineScanner: scnr,
		XMLValidator: &validation.XMLValidator{
			XSDIndexFile: "pkg/validation/xsd/index.xsd",
		},
		Cli: logger,
	}

	p.Run(ctx, *validateResponses)
}

func connect(host, port, cert, key *string) *pkg.Client {
	tlsCert, err := tls.LoadX509KeyPair(
		os.ExpandEnv(*cert),
		os.ExpandEnv(*key),
	)
	if err != nil {
		panic(err)
	}

	cl, err := pkg.Connect(net.JoinHostPort(*host, *port), &tls.Config{
		InsecureSkipVerify: true, //nolint:gosec // should only be used for testing..
		Certificates:       []tls.Certificate{tlsCert},
	})
	if err != nil {
		panic(err)
	}

	return cl
}
