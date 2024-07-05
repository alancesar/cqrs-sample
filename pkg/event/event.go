package event

import (
	"encoding/json"
	"errors"
)

const (
	ArtistSubscribedEvent Event = "ARTIST_SUBSCRIBED"
	AlbumPublishedEvent   Event = "ALBUM_PUBLISHED"
	SongPublishedEvent    Event = "SONG_PUBLISHED"
	SongPlayedEvent       Event = "SONG_PLAYED"
)

var (
	InvalidPayloadErr = errors.New("invalid payload")
)

type (
	Event string

	Message struct {
		Body    []byte
		Headers map[string]interface{}
	}
)

func NewMessage(payload any) Message {
	return NewMessageWithHeaders(payload, nil)
}

func NewMessageWithHeaders(payload any, headers map[string]interface{}) Message {
	body, _ := json.Marshal(payload)

	return Message{
		Body:    body,
		Headers: headers,
	}
}
