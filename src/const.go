package main

const (
	CONFIG_PATH = "/data0/vad/ad_process/conf/ad_process.conf"

	LOG_SECTION       = "log"
	LOG_SECTION_LOG   = "log"
	LOG_SECTION_LEVEL = "level"

	REDIS_SECTION                 = "redis"
	REDIS_SECTION_HOST            = "host"
	REDIS_SECTION_PORT            = "port"
	REDIS_SECTION_CONNECT_TIMEOUT = "connect_timeout"
	REDIS_SECTION_READ_TIMEOUT    = "read_timeout"
	REDIS_SECTION_WRITE_TIMEOUT   = "write_timeout"

	MYSQL_FST_SECTION        = "mysql_fst"
	MYSQL_FSJJ_SECTION       = "mysql_fsjj"
	MYSQL_SECTION_USERNAME   = "username"
	MYSQL_SECTION_PASSWORD   = "password"
	MYSQL_SECTION_HOST       = "host"
	MYSQL_SECTION_PORT       = "port"
	MYSQL_SECTION_DBNAME     = "dbname"
	MYSQL_SECTION_TABLE_FEED = "table_feed"
	MYSQL_SECTION_TABLE_CUST = "table_cust"

	KAFKA_SECTION        = "kafka"
	KAFKA_SECTION_BROKER = "broker"
	KAFKA_SECTION_TOPIC  = "topic"

	KESTREL_SECTION         = "kestrel"
	KESTREL_SECTION_HOST    = "host"
	KESTREL_SECTION_PORT    = "port"
	KESTREL_SECTION_QUEUE   = "queue"
	KESTREL_SECTION_TIMEOUT = "timeout"

	HTTP_SECTION                  = "http"
	HTTP_SECTION_GET_FEEDTEXT     = "url_get_feedtext"
	HTTP_SECTION_GET_CUSTTYPE     = "url_get_custtype"
	HTTP_SECTION_GET_FEEDTAG      = "url_get_feedtag"
	HTTP_SECTION_GET_CREATIVETYPE = "url_get_creativetype"
	HTTP_SECTION_TIMEOUT          = "timeout"
)
