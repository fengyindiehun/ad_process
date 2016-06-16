package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-ini/ini"
	"github.com/golang/protobuf/proto"
	"strconv"
	"strings"
	"websidx_interface"
)

type kafkaProxy struct {
	broker            string
	topic             string
	partitionConsumer sarama.PartitionConsumer
	isConnected       bool
}

func (proxy *kafkaProxy) decodeMessage(message []byte) (feedId, custId string, err error) {
	documentMsg := &websidx_interface.DocumentMsg{}
	err = proto.Unmarshal(message, documentMsg)
	if err != nil {
		return
	}
	adBidFeedContentOps := &websidx_interface.AdBidFeedContentOps{}
	err = proto.Unmarshal(documentMsg.GetAdcontent(), adBidFeedContentOps)
	if err != nil {
		return
	}
	adfeedInfo := adBidFeedContentOps.GetAdfeedInfo()
	feedId = strconv.FormatUint(adfeedInfo.GetFeedId(), 10)
	custId = strconv.FormatUint(adfeedInfo.GetCustId(), 10)
	return
}

func newKafkaProxy(section *ini.Section, offset int64) (*kafkaProxy, error) {
	if offset <= 0 {
		offset = sarama.OffsetOldest
	}
	broker := section.Key(KAFKA_SECTION_BROKER).String()
	if broker == "" {
		return nil, fmt.Errorf("kafka broker parse error")
	}

	topic := section.Key(KAFKA_SECTION_TOPIC).String()
	if topic == "" {
		return nil, fmt.Errorf("kafka topic parse error")
	}

	consumer, err := sarama.NewConsumer(strings.Split(broker, ","), nil)
	if err != nil {
		return nil, err
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, offset)
	if err != nil {
		return nil, err
	}

	return &kafkaProxy{
		broker:            broker,
		topic:             topic,
		partitionConsumer: partitionConsumer,
		isConnected:       true,
	}, nil
}
