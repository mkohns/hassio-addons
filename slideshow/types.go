package main

import "time"

type SlideInfo struct {
	SlidesCount   int    `json:"slidesCount"`
	FavoriteCount int    `json:"favoriteCount"`
	ActiveCount   int    `json:"activeCount"`
	RemoteIP      string `json:"remoteIP"`
	SlidesSize    int    `json:"slidesSize"`
	ThumbnailSize int    `json:"thumbnailSize"`
	Version       string `json:"version"`
	GitCommit     string `json:"gitCommit"`
}

type HAConfig struct {
	SignalUsername        string `json:"SIGNAL_USERNAME"`
	SignalPassword        string `json:"SIGNAL_PASSWORD"`
	SignalAccountNo       string `json:"SIGNAL_ACCOUNTNO"`
	SignalOutputFolder    string `json:"SIGNAL_OUTPUTFOLDER"`
	SignalThumbnailFolder string `json:"SIGNAL_THUMBNAILFOLDER"`
	SignalGroupID         string `json:"SIGNAL_GROUPID"`
	SignalGroupIDReal     string `json:"SIGNAL_GROUPID_REAL"`
	SignalSignalWS        string `json:"SIGNAL_SIGNALWS"`
	SignalSignalAPI       string `json:"SIGNAL_SIGNALAPI"`
	SlideshowPort         string `json:"SLIDESHOW_PORT"`
	SlideshowFrontendDist string `json:"SLIDESHOW_FRONTEND_DIST"`
	SlideshowConfigDir    string `json:"SLIDESHOW_CONFIGDIR"`
}

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
	RemoteDelete     *Timestamp   `json:"remoteDelete,omitempty"`
}
type Timestamp struct {
	Timestamp int64 `json:"timestamp"`
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
	MsgTimestamp int64
	AttachmentID string
	ImageURL     string
	TumbnailURL  string
	Message      string
	CreatedBy    string
	CreatedAt    time.Time
	Enabled      bool
	Favorite     bool
}

type SlidePatchBody struct {
	Enabled  *bool
	Favorite *bool
}

type SlideShowConfig struct {
	SessionID                 string
	ShowOnlyFavorites         bool
	ShowOnlyActive            bool
	ShowOnlyInTimeFrame       bool
	ShowNewImagesWithPriority bool
	PortraitMode              bool
	ModeRandom                bool
	ModeChronological         bool
	ModeReverseChronological  bool
	StartDate                 *time.Time
	EndDate                   *time.Time
}

const (
	ModeRandom               string = "random"
	ModeChronological        string = "chronological"
	ModeReverseChronological string = "reverseChronological"
)

type PortraitPatchBody struct {
	PortraitMode bool `json:"PortraitMode"`
}

type SlideSessionInfo struct {
	LastConfig        *SlideShowConfig
	LastSlideIndex    int
	NewSlidesPriority *[]int
	PrioNewSlides     bool
}
