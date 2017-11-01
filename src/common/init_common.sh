#!/bin/sh

# 安装common所有库的依赖
# @Author: Golion
# @Date: 2017.5

go get github.com/theckman/go-flock
go get github.com/bitly/go-simplejson
go get github.com/Shopify/sarama
go get github.com/wvanbergen/kazoo-go
go get github.com/wvanbergen/kafka/consumergroup
go get github.com/bradfitz/gomemcache/memcache
go get github.com/go-sql-driver/mysql
go get github.com/go-redis/redis
go get github.com/yunge/sphinx
go get github.com/hongst/rend
go get github.com/spaolacci/murmur3

echo "Init Common Succeed!"
