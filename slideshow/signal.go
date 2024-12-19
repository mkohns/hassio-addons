package main

import (
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
)

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
					if downloadAttachment(attachment.ID, attachment.Filename, username, password) {
						// Append the slide to the list
						slides = append(slides, Slide{
							Filename:    attachment.Filename,
							ImageURL:    "images/" + attachment.Filename,
							TumbnailURL: "thumbnails/" + attachment.Filename,
							Message:     msg.Envelope.DataMessage.Message,
							CreatedBy:   msg.Envelope.SourceName,
							CreatedAt:   time.Now(),
						})
						saveSlides(slides)
					}
				}
			}

		}

		log.Println("Reconnecting to WebSocket...")
		time.Sleep(1 * time.Second) // Wait before reconnecting
	}
}

func downloadAttachment(attachmentID, filename, username, password string) bool {
	url := fmt.Sprintf("http://signal.kohns.eu:80/v1/attachments/%s", attachmentID)

	// Create the Basic Auth header
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return false
	}

	// Add the Authorization header
	req.Header.Add("Authorization", auth)

	// Perform the request with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second, // Set the connection timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error downloading attachment:", err)
		return false
	}
	defer resp.Body.Close()

	// Check if the content type is an image
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		log.Printf("Attachment %s is not an image, skipping", filename)
		return false
	}

	// check the filename
	if !strings.HasSuffix(filename, ".jpg") && !strings.HasSuffix(filename, ".jpeg") && !strings.HasSuffix(filename, ".png") {
		log.Println("Attachment is not a jpg, jpeg or png file, skipping")
		return false
	}

	// Save the downloaded image to a file
	file, err := os.Create(outputfolder + filename)
	if err != nil {
		log.Println("Error creating file:", err)
		return false
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println("Error saving attachment:", err)
		return false
	}

	file.Close()

	// create thumbnail
	err = createThumbnail(filename)
	if err != nil {
		log.Println("Error creating thumbnail:", err)
		os.Remove(outputfolder + filename)
		return false
	}

	log.Printf("Attachment %s downloaded and processed successfully", filename)
	return true
}

func createThumbnail(filename string) error {
	// Open the file
	file, err := os.Open(outputfolder + filename)
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
	out, err := os.Create(thumbnailfolder + filename)
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
