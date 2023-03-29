package djafka

type ConnectionChangedMsg Connection
type ClientConnectedMsg struct{}
type TopicsSelectedMsg struct{}
type TopicsLoadedMsg []string
type ErrorMsg error
type ResetMsg struct{}
type AddTopicSubmitMsg struct {
	name              string
	paritions         int
	replicationFactor int
}
type AddTopicCancel struct{}
