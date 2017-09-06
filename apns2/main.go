package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/antlinker/apns2"
	"github.com/antlinker/apns2/certificate"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	certificatePath = kingpin.Flag("certificate-path", "Path to certificate file.").Required().Short('c').String()
	topic           = kingpin.Flag("topic", "The topic of the remote notification, which is typically the bundle ID for your app").Required().Short('t').String()
	mode            = kingpin.Flag("mode", "APNS server to send notifications to. `production` or `development`. Defaults to `production`").Default("production").Short('m').String()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.1").Author("Alisson Sales")
	kingpin.CommandLine.Help = `Listens to STDIN to send notifications and writes APNS response code and reason to STDOUT.
	The expected format is: <DeviceToken> <APNS Payload>
	Example: test c9b0a6ec0903b4654b1dce1acdf80b1d39ff3f93c38ded7c0b1c34683ccafda0 {"aps": {"alert": "hi"}}`
	kingpin.Parse()

	cert, pemErr := certificate.FromPemFile(*certificatePath, "123456")

	if pemErr != nil {
		log.Fatalf("Error retrieving certificate `%v`: %v", certificatePath, pemErr)
	}

	client := apns2.NewClient(cert)

	if *mode == "development" {
		client.Development()
	} else {
		client.Production()
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		in := scanner.Text()
		notificationArgs := strings.SplitN(in, " ", 3)
		CollapseID := notificationArgs[0]
		token := notificationArgs[1]
		payload := notificationArgs[2]
		notification := &apns2.Notification{
			DeviceToken: token,
			Topic:       *topic,
			Payload:     payload,
			CollapseID:  CollapseID,
		}

		res, err := client.Push(notification)

		if err != nil {
			log.Fatal("Error: ", err)
		} else {
			fmt.Printf("%v: '%v'\n", res.StatusCode, res.Reason)
		}
	}
}
