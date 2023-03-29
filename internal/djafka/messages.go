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
