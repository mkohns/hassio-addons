package main

import "time"

type GroupInfo struct {
	GroupID   string `json:"groupId"`
	GroupName string `json:"groupName"`
	Revision  int    `json:"revision"`
	Type      string `json:"type"`
}

type Attachment struct {
	ContentType     string  `json:"contentType"`
	Filename        string  `json:"filename"`
	ID              string  `json:"id"`
	Size            int     `json:"size"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	Caption         *string `json:"caption,omitempty"`
	UploadTimestamp *int64  `json:"uploadTimestamp,omitempty"`
}

type DataMessage struct {
	Timestamp        int64        `json:"timestamp"`
	Message          string       `json:"message"`
	ExpiresInSeconds int          `json:"expiresInSeconds"`
	ViewOnce         bool         `json:"viewOnce"`
	Attachments      []Attachment `json:"attachments,omitempty"`
	GroupInfo        *GroupInfo   `json:"groupInfo,omitempty"`
}

type TypingMessage struct {
	Action    string `json:"action"`
	Timestamp int64  `json:"timestamp"`
	GroupID   string `json:"groupId"`
}

type Envelope struct {
	Source                   string         `json:"source"`
	SourceNumber             string         `json:"sourceNumber"`
	SourceUuid               string         `json:"sourceUuid"`
	SourceName               string         `json:"sourceName"`
	SourceDevice             int            `json:"sourceDevice"`
	Timestamp                int64          `json:"timestamp"`
	ServerReceivedTimestamp  int64          `json:"serverReceivedTimestamp"`
	ServerDeliveredTimestamp int64          `json:"serverDeliveredTimestamp"`
	DataMessage              *DataMessage   `json:"dataMessage,omitempty"`
	TypingMessage            *TypingMessage `json:"typingMessage,omitempty"`
}

type Message struct {
	Envelope Envelope `json:"envelope"`
	Account  string   `json:"account"`
}

type Slide struct {
	Filename    string
	ImageURL    string
	TumbnailURL string
	Message     string
	CreatedBy   string
	CreatedAt   time.Time
}
