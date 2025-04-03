package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/tmaxmax/go-sse"
)

type SSEServiceInterface interface {
	BroadcastAuctionEvent(auctionID string, eventType string, data interface{})
	HandleAuctionEventsSSE(w http.ResponseWriter, r *http.Request)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type SSEService struct {
	server *sse.Server
}

func NewSSEService() *SSEService {
	return &SSEService{
		server: &sse.Server{}, // zero value is ready to use
	}
}

// ServeHTTP implements http.Handler
func (s *SSEService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
}

// BroadcastAuctionEvent sends an event to all clients subscribed to a specific auction
func (s *SSEService) BroadcastAuctionEvent(auctionID string, eventType string, data interface{}) {
	// Try to parse and normalize the time format
	if t, err := time.Parse("2006-01-02 15:04:05.999 -0700 MST", auctionID); err == nil {
		auctionID = t.Format(time.RFC3339Nano)
	}

	topic := "auction_" + auctionID

	// If it's an auction, nil out the maxBid
	if auction, ok := data.(openapi.Auction); ok {
		auction.MaxBid = 0
		data = auction
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	msg := &sse.Message{}
	msg.Type = sse.Type(eventType)
	msg.AppendData(string(jsonData))

	s.server.Publish(msg, topic)
}

// HandleAuctionEventsSSE creates a new SSE subscription for specific auctions
func (s *SSEService) HandleAuctionEventsSSE(w http.ResponseWriter, r *http.Request) {
	// Clear existing headers
	for k := range w.Header() {
		w.Header().Del(k)
	}

	// Set SSE headers only
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	auctionIDs := r.URL.Query().Get("auction_ids")
	if auctionIDs == "" {
		http.Error(w, "auction_ids parameter is required", http.StatusBadRequest)
		return
	}

	// Split auction IDs and normalize them
	rawIds := strings.Split(auctionIDs, ",")
	topics := make([]string, 0, len(rawIds))

	for _, id := range rawIds {
		// Try to parse and normalize the time format
		if t, err := time.Parse(time.RFC3339Nano, id); err == nil {
			topic := "auction_" + t.Format(time.RFC3339Nano)
			topics = append(topics, topic)
		} else {
			topic := "auction_" + id
			topics = append(topics, topic)
		}
	}

	// Configure the server to handle this session
	s.server.OnSession = func(sess *sse.Session) (sse.Subscription, bool) {
		return sse.Subscription{
			Client: sess,
			Topics: topics,
		}, true
	}

	s.server.ServeHTTP(w, r)
}
