package djafka

type ConnectionChangedMsg Connection
type ClientConnectedMsg struct{}

type TopicsSelectedMsg struct{}
type TopicsLoadedMsg []Topic
type TopicSelectedMsg Topic
type TopicSettingsLoadedMsg TopicConfig

type ConsumersSelectedMsg struct{}
type ConsumersLoadedMsg []Consumer
type ConsumerSelectedMsg Consumer

type ErrorMsg error
type ResetMsg struct{}
type InfoSelectedMsg struct{}
type TickMsg struct{}

type AddTopicSubmitMsg struct {
	name              string
	paritions         int
	replicationFactor int
}
type AddTopicCancel struct{}

type ResetOffsetMsg struct {
	consumerGroup string
	topicName     string
	offset        int64
}
type ResetOffsetCancel struct{}
