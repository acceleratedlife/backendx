package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
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
	lgr.Printf("[SSE] Broadcasting event - Original auctionID: %s", auctionID)

	// Try to parse and normalize the time format
	if t, err := time.Parse("2006-01-02 15:04:05.999 -0700 MST", auctionID); err == nil {
		auctionID = t.Format(time.RFC3339)
		lgr.Printf("[SSE] Normalized auctionID to: %s", auctionID)
	} else {
		lgr.Printf("[SSE] Could not parse time format: %v", err)
	}

	topic := "auction_" + auctionID
	lgr.Printf("[SSE] Broadcasting to topic: %s", topic)

	// If it's an auction, nil out the maxBid
	if auction, ok := data.(openapi.Auction); ok {
		auction.MaxBid = 0
		data = auction
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		lgr.Printf("[SSE] Error marshaling data: %v", err)
		return
	}
	lgr.Printf("[SSE] Marshaled data: %s", string(jsonData))

	msg := &sse.Message{}
	msg.Type = sse.Type(eventType)
	msg.AppendData(string(jsonData))

	lgr.Printf("[SSE] Publishing message - Type: %s, Topic: %s", eventType, topic)
	s.server.Publish(msg, topic)
	lgr.Printf("[SSE] Message published")
}

// HandleAuctionEventsSSE creates a new SSE subscription for specific auctions
func (s *SSEService) HandleAuctionEventsSSE(w http.ResponseWriter, r *http.Request) {
	lgr.Printf("[SSE] New connection attempt from %s", r.RemoteAddr)

	auctionIDs := r.URL.Query().Get("auction_ids")
	if auctionIDs == "" {
		lgr.Printf("[SSE] Error: Missing auction_ids parameter")
		http.Error(w, "auction_ids parameter is required", http.StatusBadRequest)
		return
	}

	// Split auction IDs and normalize them
	rawIds := strings.Split(auctionIDs, ",")
	topics := make([]string, 0, len(rawIds))
	for _, id := range rawIds {
		lgr.Printf("[SSE] Processing auction ID: %s", id)
		// Try to parse and normalize the time format
		if t, err := time.Parse(time.RFC3339, id); err == nil {
			topic := "auction_" + t.Format(time.RFC3339)
			topics = append(topics, topic)
			lgr.Printf("[SSE] Normalized topic: %s", topic)
		} else {
			topic := "auction_" + id
			topics = append(topics, topic)
			lgr.Printf("[SSE] Using raw topic: %s (parse error: %v)", topic, err)
		}
	}

	lgr.Printf("[SSE] Attempting to subscribe to topics: %v", topics)

	// Configure the server to handle this session
	s.server.OnSession = func(sess *sse.Session) (sse.Subscription, bool) {
		return sse.Subscription{
			Client: sess,
			Topics: topics,
		}, true
	}

	// Let the server handle the SSE connection
	s.server.ServeHTTP(w, r)
}
