package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	pb "github.com/clawfinger/ratelimiter/api/generated"
	"github.com/clawfinger/ratelimiter/cli"
)

const (
	addBlacklist    = 1
	removeBlacklist = 2
	addWhitelist    = 3
	removeWhitelist = 4
	dropIP          = 5
	dropLogin       = 6
	dropPassword    = 7
	exit            = 0
)

func printActions() {
	fmt.Println("1 Add subnet to blacklist")
	fmt.Println("2 Remove subnet from blacklist")
	fmt.Println("3 Add subnet to whitelist")
	fmt.Println("4 Remove subnet from whitelist")
	fmt.Println("5 Drop ip restrictions")
	fmt.Println("6 Drop login restrictions")
	fmt.Println("7 Drop password restrictions")
	fmt.Println("0 Exit")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	fmt.Println("Limiter service control interface")
	reader := bufio.NewReader(os.Stdin)

	addrValidator := regexp.MustCompile(`\d+.\d+.\d+:\d+`)

	var client *cli.CtlAgent
Label:
	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("Please enter ip and port. e.g: 127.0.0.1:50051")

			text, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			ipOK := addrValidator.MatchString(text)
			if !ipOK {
				fmt.Println("address string is invalid")
			} else {
				client, err = cli.NewClient(strings.TrimSuffix(text, "\n"))
				if err != nil {
					fmt.Printf("Failed to connect to service. Reason: %s\n", err.Error())
				} else {
					break Label
				}
			}
		}
	}

	for {
		fmt.Println("\nPlease enter number correspoding to action")
		printActions()
		text, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		text = strings.TrimSuffix(text, "\n")
		selected, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Wrong action")
		}
		switch selected {
		case addBlacklist:
			handleAddBlacklist(ctx, client, reader)
		case removeBlacklist:
			handleRemoveBlacklist(ctx, client, reader)
		case addWhitelist:
			handleAddWhitelist(ctx, client, reader)
		case removeWhitelist:
			handleRemoveWhitelist(ctx, client, reader)
		case dropIP:
			handleIPDropStats(ctx, client, reader)
		case dropLogin:
			handleLoginDropStats(ctx, client, reader)
		case dropPassword:
			handlePasswordDropStats(ctx, client, reader)
		case exit:
			return
		default:
			fmt.Println("Wrong action")
		}
	}
}

func handleAddBlacklist(ctx context.Context, client *cli.CtlAgent, reader *bufio.Reader) {
	fmt.Println("Please enter the subnet to add to blacklist")
	text, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	subnet := &pb.Subnet{
		IPWithMask: text,
	}
	result, err := client.AddBlacklist(ctx, subnet)
	if err != nil {
		fmt.Printf("Error handling add blacklist %s", err.Error())
		return
	}
	fmt.Printf("Finished adding subnet to blacklist. Status %s. Reason: %s\n", result.Status.String(), result.Reason)
}

func handleRemoveBlacklist(ctx context.Context, client *cli.CtlAgent, reader *bufio.Reader) {
	fmt.Println("Please enter the subnet to remove from blacklist")
	text, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	subnet := &pb.Subnet{
		IPWithMask: text,
	}
	result, err := client.RemoveBlacklist(ctx, subnet)
	if err != nil {
		fmt.Printf("Error handling remove blacklist %s", err.Error())
		return
	}
	fmt.Printf("Finished removing subnet from blacklist. Status %s. Reason: %s\n", result.Status.String(), result.Reason)
}

func handleAddWhitelist(ctx context.Context, client *cli.CtlAgent, reader *bufio.Reader) {
	fmt.Println("Please enter the subnet to add to whitelist")
	text, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	subnet := &pb.Subnet{
		IPWithMask: text,
	}
	result, err := client.AddWhitelist(ctx, subnet)
	if err != nil {
		fmt.Printf("Error handling add whitelist %s", err.Error())
		return
	}
	fmt.Printf("Finished adding subnet to whitelist. Status %s. Reason: %s\n", result.Status.String(), result.Reason)
}

func handleRemoveWhitelist(ctx context.Context, client *cli.CtlAgent, reader *bufio.Reader) {
	fmt.Println("Please enter the subnet to remove from whitelist")
	text, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	subnet := &pb.Subnet{
		IPWithMask: text,
	}
	result, err := client.RemoveBlacklist(ctx, subnet)
	if err != nil {
		fmt.Printf("Error handling remove whitelist %s", err.Error())
		return
	}
	fmt.Printf("Finished removing subnet from whitelist. Status %s. Reason: %s\n", result.Status.String(), result.Reason)
}

func handleIPDropStats(ctx context.Context, client *cli.CtlAgent, reader *bufio.Reader) {
	fmt.Println("Please enter the IP to drop the restrictions")
	text, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	ip := &pb.Stats{
		Data: text,
	}
	result, err := client.DropIPStats(ctx, ip)
	if err != nil {
		fmt.Printf("Error handling drop ip restrictions %s", err.Error())
		return
	}
	fmt.Printf("Finished drop ip restrictions. Status %s. Reason: %s\n", result.Status.String(), result.Reason)
}

func handleLoginDropStats(ctx context.Context, client *cli.CtlAgent, reader *bufio.Reader) {
	fmt.Println("Please enter the login to drop the restrictions")
	text, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	login := &pb.Stats{
		Data: text,
	}
	result, err := client.DropLoginStats(ctx, login)
	if err != nil {
		fmt.Printf("Error handling drop login restrictions %s", err.Error())
		return
	}
	fmt.Printf("Finished drop login restrictions. Status %s. Reason: %s\n", result.Status.String(), result.Reason)
}

func handlePasswordDropStats(ctx context.Context, client *cli.CtlAgent, reader *bufio.Reader) {
	fmt.Println("Please enter the password to drop the restrictions")
	text, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	pass := &pb.Stats{
		Data: text,
	}
	result, err := client.DropPasswordStats(ctx, pass)
	if err != nil {
		fmt.Printf("Error handling drop password restrictions %s", err.Error())
		return
	}
	fmt.Printf("Finished drop password restrictions. Status %s. Reason: %s\n", result.Status.String(), result.Reason)
}
