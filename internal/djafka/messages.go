package djafka

type ConnectionChangedMsg Connection
type ClientConnectedMsg struct{}
type TopicsSelectedMsg struct{}
type TopicsLoadedMsg []string
