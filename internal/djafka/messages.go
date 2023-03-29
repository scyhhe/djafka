package djafka

type ConnectionChangedMsg Connection
type ClientConnectedMsg struct{}
type TopicsSelectedMsg struct{}
type TopicsLoadedMsg []string

type ConsumersSelectedMsg struct{}
type ConsumersLoadedMsg []Consumer
type ConsumerSelectedMsg Consumer
type ErrorMsg error
type ResetMsg struct{}
