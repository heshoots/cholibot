module github.com/quorauk/cholibot

go 1.12

require (
	github.com/bwmarrin/discordgo v0.19.0
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/mux v1.7.3
	github.com/jinzhu/configor v1.1.1
	github.com/quorauk/dmux v0.0.6
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80
	google.golang.org/grpc v1.22.1
)

replace github.com/quorauk/dmux => ../dmux
