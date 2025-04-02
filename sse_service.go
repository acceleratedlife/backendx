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
	lgr.Printf("[SSE] Broadcasting event - Type: %s, Auction ID: %s", eventType, auctionID)

	// Try to parse and normalize the time format
	if t, err := time.Parse("2006-01-02 15:04:05.999 -0700 MST", auctionID); err == nil {
		auctionID = t.Format(time.RFC3339Nano)
		lgr.Printf("[SSE] Normalized auction ID timestamp to: %s", auctionID)
	} else {
		lgr.Printf("[SSE] Using raw auction ID (parse error: %v)", err)
	}

	topic := "auction_" + auctionID
	lgr.Printf("[SSE] Broadcasting to topic: %s", topic)

	// If it's an auction, nil out the maxBid
	if auction, ok := data.(openapi.Auction); ok {
		auction.MaxBid = 0
		data = auction
		lgr.Printf("[SSE] Modified auction data (cleared MaxBid)")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		lgr.Printf("[SSE] ERROR marshaling data: %v", err)
		return
	}

	lgr.Printf("[SSE] Raw message data being sent: %s", string(jsonData))

	msg := &sse.Message{}
	msg.Type = sse.Type(eventType)
	msg.AppendData(string(jsonData))

	lgr.Printf("[SSE] Publishing message to server")
	s.server.Publish(msg, topic)
	lgr.Printf("[SSE] Message published successfully")
}

// HandleAuctionEventsSSE creates a new SSE subscription for specific auctions
func (s *SSEService) HandleAuctionEventsSSE(w http.ResponseWriter, r *http.Request) {
	// Clear existing headers
	for k := range w.Header() {
		w.Header().Del(k)
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Handle allowed origins
	origin := r.Header.Get("Origin")
	allowedOrigins := map[string]bool{
		"https://test.schoolbucks.net": true,
		"https://www.schoolbucks.net":  true,
		"https://schoolbucks.net":      true,
		"http://localhost:3000":        true,
	}

	if origin != "" && allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		lgr.Printf("[SSE] Allowed origin: %s", origin)
	} else {
		lgr.Printf("[SSE] Warning: Unknown origin: %s", origin)
	}

	lgr.Printf("[SSE] New SSE connection request from %s", r.RemoteAddr)

	auctionIDs := r.URL.Query().Get("auction_ids")
	if auctionIDs == "" {
		lgr.Printf("[SSE] ERROR: Missing auction_ids parameter from %s", r.RemoteAddr)
		http.Error(w, "auction_ids parameter is required", http.StatusBadRequest)
		return
	}

	lgr.Printf("[SSE] Requested auction IDs: %s", auctionIDs)

	// Split auction IDs and normalize them
	rawIds := strings.Split(auctionIDs, ",")
	topics := make([]string, 0, len(rawIds))
	lgr.Printf("[SSE] Processing %d auction IDs", len(rawIds))

	for _, id := range rawIds {
		// Try to parse and normalize the time format
		if t, err := time.Parse(time.RFC3339Nano, id); err == nil {
			topic := "auction_" + t.Format(time.RFC3339Nano)
			topics = append(topics, topic)
			lgr.Printf("[SSE] Normalized topic: %s", topic)
		} else {
			topic := "auction_" + id
			topics = append(topics, topic)
			lgr.Printf("[SSE] Using raw topic: %s (parse error: %v)", topic, err)
		}
	}

	// Configure the server to handle this session
	s.server.OnSession = func(sess *sse.Session) (sse.Subscription, bool) {
		lgr.Printf("[SSE] Creating new session for client %s", r.RemoteAddr)
		lgr.Printf("[SSE] Subscribing to topics: %v", topics)
		return sse.Subscription{
			Client: sess,
			Topics: topics,
		}, true
	}

	lgr.Printf("[SSE] Starting SSE stream for client %s", r.RemoteAddr)
	s.server.ServeHTTP(w, r)
	lgr.Printf("[SSE] SSE stream ended for client %s", r.RemoteAddr)
}
