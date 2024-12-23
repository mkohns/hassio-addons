package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nfnt/resize"
	"kohns.eu/signal2S3/signal"
)

// SWAGGER: https://bbernhard.github.io/signal-cli-rest-api/

var httpClient *http.Client = nil
var signalClient *signal.ClientWithResponses = nil

func initClient(baseurl string) {
	httpClient = &http.Client{
		Timeout: 10 * time.Second, // Set the connection timeout
	}
	var err error

	signalClient, err = signal.NewClientWithResponses(baseurl, signal.WithHTTPClient(httpClient))
	if err != nil {
		log.Panicln("Error creating signal client:", err)
	}
}

func connectToWebSocket(socketURL, username, password string) {
	for {
		u, err := url.Parse(socketURL)
		if err != nil {
			log.Fatal("Error parsing URL:", err)
		}

		// Create the Basic Auth header
		auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

		// Create a custom dialer with the Basic Auth header and a connection timeout
		dialer := websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 10 * time.Second, // Set the connection timeout
		}

		headers := http.Header{}
		headers.Add("Authorization", auth)

		conn, _, err := dialer.Dial(u.String(), headers)
		if err != nil {
			log.Println("Error connecting to WebSocket:", err)
			time.Sleep(1 * time.Second) // Wait before retrying
			continue
		}
		log.Println("Connected to WebSocket:", socketURL)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				conn.Close()
				break
			}

			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Println("Error unmarshalling message:", err)
				continue
			}

			log.Printf("Received: %s", message)

			if msg.Envelope.DataMessage != nil &&
				msg.Envelope.DataMessage.GroupInfo != nil &&
				msg.Envelope.DataMessage.GroupInfo.GroupID == groupId &&
				msg.Envelope.DataMessage.Attachments != nil &&
				len(msg.Envelope.DataMessage.Attachments) > 0 {
				log.Println("Got Message with Attachments!")

				for _, attachment := range msg.Envelope.DataMessage.Attachments {
					if downloadAttachment(attachment.ID) {
						// Append the slide to the list
						newIndex := addSlide(Slide{
							MsgTimestamp: msg.Envelope.DataMessage.Timestamp,
							AttachmentID: attachment.ID,
							ImageURL:     "images/" + attachment.ID,
							TumbnailURL:  "thumbnails/" + attachment.ID,
							Message:      msg.Envelope.DataMessage.Message,
							CreatedBy:    msg.Envelope.SourceName,
							CreatedAt:    time.Now(),
							Enabled:      true,
							Favorite:     false,
						})

						// check image orientation
						err, imageOrientation := isImagePortraitInOrientation(outputfolder + attachment.ID)
						if err != nil {
							log.Println("Error checking image orientation:", err)
						}

						// get the default session for orientation handling
						session := getSession("")

						if err == nil && session.LastConfig != nil && imageOrientation != session.LastConfig.PortraitMode {
							currentMode := "portrait"
							if !session.LastConfig.PortraitMode {
								currentMode = "landscape"
							}
							sendMessage("The image is in wrong orientation 👎. Image should be " + currentMode)
							sendReaction(msg, "👎")
						} else {
							sendReaction(msg, "🚀")
						}

						handleNewPrioritySession(newIndex)
					}
					// remove the attachment from the server
					err := removeAttachment(attachment.ID)
					if err != nil {
						log.Println("Error removing attachment:", err)
					}
				}
			} else if msg.Envelope.DataMessage != nil &&
				msg.Envelope.DataMessage.RemoteDelete != nil {
				log.Printf("Got RemoteDelete Message for timestamp: %d \n", msg.Envelope.DataMessage.RemoteDelete.Timestamp)
				removeSlideByTimestamp(msg.Envelope.DataMessage.RemoteDelete.Timestamp)
			}

		}

		log.Println("Reconnecting to WebSocket...")
		time.Sleep(1 * time.Second) // Wait before reconnecting
	}
}

func isImagePortraitInOrientation(filename string) (error, bool) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening file:", err)
		return err, false
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Println("Error decoding image:", err)
		return err, false
	}

	if img.Width > img.Height {
		return nil, false
	} else {
		return nil, true
	}
}

func doNothing(ctx context.Context, req *http.Request) error {
	return nil
}

func sendReaction(msg Message, emoji string) {
	intPtr := func(i int) *int { return &i }
	reaction := emoji
	receipient := msg.Envelope.Source
	timestamp := msg.Envelope.DataMessage.Timestamp
	body := signal.PostV1ReactionsNumberJSONRequestBody{
		Reaction:     &reaction,
		Recipient:    &groupIdReal,
		TargetAuthor: &receipient,
		Timestamp:    intPtr(int(timestamp)),
	}
	_, err := signalClient.PostV1ReactionsNumber(context.Background(), accountNo, body, doNothing)
	if err != nil {
		log.Println("Error sending reaction:", err)
		return
	}
	log.Println("Reaction sent successfully")
}

func sendMessage(message string) {
	recipients := []string{groupIdReal}
	textMode := "normal"

	body := signal.PostV2SendJSONRequestBody{
		Message:    &message,
		Number:     &accountNo,
		Recipients: &recipients,
		TextMode:   (*signal.ApiSendMessageV2TextMode)(&textMode),
	}
	_, err := signalClient.PostV2Send(context.Background(), body, doNothing)
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}
	log.Println("Message sent successfully")
}

func removeAttachment(attachmentID string) error {
	res, err := signalClient.DeleteV1AttachmentsAttachmentWithResponse(context.Background(), attachmentID)
	if err != nil {
		log.Panicln("Error getting signal about:", err)
	}

	if res.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("error removing attachment: %s", res.Status())
	}

	log.Printf("Attachment %s removed successfully", attachmentID)
	return nil
}

func downloadAttachment(attachmentID string) bool {
	// check the suffix
	if !strings.HasSuffix(attachmentID, ".jpg") && !strings.HasSuffix(attachmentID, ".jpeg") && !strings.HasSuffix(attachmentID, ".png") {
		log.Println("Attachment is not a jpg, jpeg or png file, skipping")
		return false
	}

	res, err := signalClient.GetV1AttachmentsAttachment(context.Background(), attachmentID)
	if err != nil {
		log.Println("Error getting attachment:", err)
		return false
	}

	if res.StatusCode != http.StatusOK {
		log.Println("Error getting attachment:", res.StatusCode)
		return false
	}

	// Check if the content type is an image
	contentType := res.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		log.Printf("Attachment %s is not an image, skipping .. (real type: %s)", attachmentID, contentType)
		return false
	}

	// Save the downloaded image to a file
	file, err := os.Create(outputfolder + attachmentID)
	if err != nil {
		log.Println("Error creating file:", err)
		return false
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		log.Println("Error saving attachment:", err)
		return false
	}

	file.Close()

	// create thumbnail
	err = createThumbnail(attachmentID)
	if err != nil {
		log.Println("Error creating thumbnail:", err)
		os.Remove(outputfolder + attachmentID)
		return false
	}

	log.Printf("Attachment %s downloaded and processed successfully", attachmentID)
	return true

}

func createThumbnail(attachmentID string) error {
	// Open the file
	file, err := os.Open(outputfolder + attachmentID)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Create a thumbnail while preserving the aspect ratio
	thumbnail := resize.Thumbnail(100, 100, img, resize.Lanczos3)

	// Save the thumbnail to a new file
	out, err := os.Create(thumbnailfolder + attachmentID)
	if err != nil {
		return err
	}
	defer out.Close()

	// Encode the thumbnail as a JPEG
	err = jpeg.Encode(out, thumbnail, nil)
	if err != nil {
		return err
	}

	return nil
}
