package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	//"strings"
	"github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
)

var (
	log              *logrus.Logger
	redisHandler     *redisProxy
	mysqlFstHandler  *mysqlProxy
	mysqlFsjjHandler *mysqlProxy
	kafkaHandler     *kafkaProxy
	kestrelHandler   *kestrelProxy
	httpHandler      *httpProxy
)

func initLog(config string) error {
	cfg, err := ini.Load(config)
	if err != nil {
		return err
	}

	section, err := cfg.GetSection(LOG_SECTION)
	if err != nil {
		return err
	}

	logFile := section.Key(LOG_SECTION_LOG).String()
	if logFile == "" {
		return fmt.Errorf("log section logFile parse error")
	}

	log = logrus.New()
	switch level := section.Key(LOG_SECTION_LEVEL).String(); level {
	case "DEBUG":
		log.Level = logrus.DebugLevel
	case "INFO":
		log.Level = logrus.InfoLevel
	case "WARN":
		log.Level = logrus.WarnLevel
	case "ERROR":
		log.Level = logrus.ErrorLevel
	case "FATAL":
		log.Level = logrus.FatalLevel
	default:
		log.Level = logrus.ErrorLevel
	}
	log.Formatter = new(logrus.JSONFormatter)
	os.Create(logFile)
	fd, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	log.Out = fd
	return nil
}

func loadConfig(config string) error {
	cfg, err := ini.Load(config)
	if err != nil {
		return err
	}
	log.Debug("config load success")

	section, err := cfg.GetSection(REDIS_SECTION)
	if err != nil {
		return err
	}
	redisHandler, err = newRedisProxy(section)
	if err != nil {
		return err
	}
	log.Debug("redis init success")

	section, err = cfg.GetSection(MYSQL_FST_SECTION)
	if err != nil {
		return err
	}
	mysqlFstHandler, err = newMysqlProxy(section)
	if err != nil {
		return err
	}
	log.Debug("mysql_fst init success")

	section, err = cfg.GetSection(MYSQL_FSJJ_SECTION)
	if err != nil {
		return err
	}
	mysqlFsjjHandler, err = newMysqlProxy(section)
	if err != nil {
		return err
	}
	log.Debug("mysql_fsjj init success")

	//get offset of kafka from redis
	var off string
	for i := 0; i < 3; i++ {
		off, err := redisHandler.get("objectinfo_kafka_offset")
		if err == nil {
			log.Debug("get kafka offset from redis success, offset:", off)
			break
		}
		log.Error("get kafka offset from redis failed, try again")
		time.Sleep(time.Second)
	}
	offset, err := strconv.ParseInt(off, 10, 64)

	section, err = cfg.GetSection(KAFKA_SECTION)
	if err != nil {
		return err
	}
	kafkaHandler, err = newKafkaProxy(section, offset)
	if err != nil {
		return err
	}
	log.Debug("kafka init success")

	section, err = cfg.GetSection(KESTREL_SECTION)
	if err != nil {
		return err
	}
	kestrelHandler, err = newKestrelProxy(section)
	if err != nil {
		return err
	}
	log.Debug("kestrel init success")

	section, err = cfg.GetSection(HTTP_SECTION)
	if err != nil {
		return err
	}
	httpHandler, err = newHttpProxy(section)
	if err != nil {
		return err
	}
	log.Debug("ad_process init success")

	return nil
}

func process(feedId, custId string, messageType string) {
	text, err := httpHandler.getFeedText(feedId)
	if err != nil {
		log.Error(err)
	}

	custType, err := httpHandler.getCustType(custId)
	if err != nil {
		log.Error(err)
	}

	feedTag, err := httpHandler.getFeedTag(feedId, custId, text)
	if err != nil {
		log.Error(err)
	}

	creativeType, err := httpHandler.getCreativeType(feedId)
	if err != nil {
		log.Error(err)
	}

	var args []interface{} = []interface{}{
		"objectinfo_mid_" + feedId,
		creativeType + ":" + custId + ":" + feedTag + ":" + text,
		"objectinfo_custid_" + custId,
		custType}
	if err := redisHandler.mSet(args); err != nil {
		log.Error(err)
	}

	if messageType == "KAFKA" {
		err := mysqlFstHandler.updateFeedInfo(feedId, custId, text, creativeType, feedTag)
		if err != nil {
			log.Error(err)
		}
		err = mysqlFstHandler.updateCustInfo(custId, custType)
		if err != nil {
			log.Error(err)
		}
	} else if messageType == "KESTREL" {
		err := mysqlFsjjHandler.updateFeedInfo(feedId, custId, text, creativeType, feedTag)
		if err != nil {
			log.Error(err)
		}
		err = mysqlFsjjHandler.updateCustInfo(custId, custType)
		if err != nil {
			log.Error(err)
		}
	} else {
		log.Error("wrong messageType")
	}
}

func main() {
	err := initLog(CONFIG_PATH)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = loadConfig(CONFIG_PATH)
	if err != nil {
		log.Error(err)
		return
	}

	go kestrelHandler.receiveMessage()

	for {
		select {
		case msg := <-kafkaHandler.partitionConsumer.Messages():
			offset := msg.Offset
			feedId, custId, err := kafkaHandler.decodeMessage(msg.Value)
			log.Debug("kafka message,feedId:", feedId, " custId:", custId, " offset:", offset)
			if err != nil {
				log.Error(err)
			} else {
				process(feedId, custId, "KAFKA")
				fmt.Println(strconv.FormatInt(offset, 10))
				redisHandler.set("objectinfo_kafka_offset", strconv.FormatInt(offset, 10))
			}
		case msg := <-kestrelHandler.messages():
			feedId, custId, err := kestrelHandler.decodeMessage(msg)
			log.Debug("kestrel message,feedId:", feedId, " custId:", custId)
			if err != nil {
				log.Error(err)
			} else {
				process(feedId, custId, "KESTREL")
			}
		}
	}
}
