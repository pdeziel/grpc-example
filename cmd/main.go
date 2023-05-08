package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	hello "github.com/pdeziel/grpc-example"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	app := cli.NewApp()
	app.Name = "grpc-example"
	app.Usage = "A simple gRPC client"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "endpoint",
			Aliases: []string{"e"},
			Usage:   "gRPC server endpoint",
			Value:   "localhost:443",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:     "serve",
			Usage:    "Start the gRPC server",
			Category: "server",
			Action:   serve,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "bindaddr",
					Aliases: []string{"a"},
					Usage:   "Address to bind the server to",
					Value:   ":443",
				},
			},
		},
		{
			Name:   "hello",
			Usage:  "Say hello in a given language",
			Before: initClient,
			Action: getHello,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "lang",
					Aliases: []string{"l"},
					Usage:   "Language code",
					Value:   "en",
				},
			},
		},
		{
			Name:   "hello:many",
			Usage:  "Say hello in many languages",
			Before: initClient,
			Action: getManyHellos,
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:    "langs",
					Aliases: []string{"l"},
					Usage:   "Language codes",
					Value:   cli.NewStringSlice("en", "fr", "es"),
				},
			},
		},
		{
			Name:   "hello:stream",
			Usage:  "Say hello in real time",
			Before: initClient,
			Action: getStreamHellos,
			Flags:  []cli.Flag{},
		},
	}

	app.Run(os.Args)
}

// Start the gRPC server and block until stopped
func serve(c *cli.Context) (err error) {
	addr := c.String("bindaddr")

	var server *hello.Server
	if server, err = hello.NewServer(); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Println("Starting the server on", addr)
	if err = server.Serve(addr); err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

var client *hello.Client

// Initialize a gRPC client
func initClient(c *cli.Context) (err error) {
	endpoint := c.String("endpoint")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if client, err = hello.NewClient(endpoint, opts...); err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

// Get a hello message in a given language
func getHello(c *cli.Context) (err error) {
	ctx := context.Background()
	lang := c.String("lang")

	var msg string
	if msg, err = client.SayHello(ctx, lang); err != nil {
		fmt.Println("Error:", err)
		return cli.Exit(err, 1)
	}

	fmt.Println(msg)
	return nil
}

// Get hello messages in many languages
func getManyHellos(c *cli.Context) (err error) {
	ctx := context.Background()
	langs := c.StringSlice("langs")

	var msgs []string
	if msgs, err = client.SayServerStream(ctx, langs); err != nil {
		fmt.Println("Error:", err)
		return cli.Exit(err, 1)
	}

	for _, msg := range msgs {
		fmt.Println(msg)
	}
	return nil
}

// Stream hello messages in real time
func getStreamHellos(c *cli.Context) (err error) {
	ctx := context.Background()

	langs := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter a language code: ")
			lang, _ := reader.ReadString('\n')
			if lang == "" {
				break
			}
			langs <- strings.TrimSpace(lang)
		}
		close(langs)
	}()

	var msgs []string
	if msgs, err = client.SayClientStream(ctx, langs); err != nil {
		fmt.Println("Error:", err)
		return cli.Exit(err, 1)
	}

	for _, msg := range msgs {
		fmt.Println(msg)
	}

	wg.Wait()
	return nil
}
