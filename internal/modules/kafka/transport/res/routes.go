package kafkarest

import "net/http"

func (handler *kafkaHandlers) RegisterRouter() {
	s := handler.router.PathPrefix("/kafka").Subrouter()
	// Topic
	s.HandleFunc("/topics", handler.ListTopicHandler()).Methods(http.MethodGet)
	// Publish
	s.HandleFunc("/publish", handler.SendMessageHandler()).Methods(http.MethodPost)
	// Subscribe
	s.HandleFunc("/subscribe/topic/{topic_name}", handler.SubscriberHandler()).Methods(http.MethodGet)
	// Request
	s.HandleFunc("/requests", handler.ListRequestHandler()).Methods(http.MethodGet)
	s.HandleFunc("/requests", handler.CreateRequestHandler()).Methods(http.MethodPost)
	s.HandleFunc("/requests/{request_id}", handler.UpdateRequestHandler()).Methods(http.MethodPut)
}
