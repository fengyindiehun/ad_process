# ad_process
golang not support centos 5.x 
10.13.40.74

wget https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz --no-check-certificate
tar zxvf go1.6.2.linux-amd64.tar.gz -C /usr/local/
echo "export PATH=$PATH:/usr/local/go/bin:/root/go_workspace/bin" >> /etc/profile
source /etc/profile
mkdir -p /root/go_workspace
echo "export GOPATH=/root/go_workspace" >> /etc/profile
source /etc/profile
yum install git

go get github.com/go-ini/ini
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
go get github.com/Shopify/sarama
go get github.com/Sirupsen/logrus
go get github.com/bradfitz/gomemcache/memcache
go get github.com/go-sql-driver/mysql

wget https://github.com/google/protobuf/releases/download/v3.0.0-beta-3/protoc-3.0.0-beta-3-linux-x86_64.zip


protoc --go_out=. *.proto

modify gomemcache fronm gets to get
go build -o memcache.a memcache.go selector.go



golang has weak support of daemon so we just use nohup &
sarama not support consumer group

modify the gets to get in function getFromAddr, cause kestrel does not support gets then build go build -o memcache.a memcache.go 
