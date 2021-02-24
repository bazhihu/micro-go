module micro-go

go 1.13

require (
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.15+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-kit/kit v0.9.0
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/golang/protobuf v1.4.1
	github.com/google/uuid v1.0.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/gorilla/mux v1.6.2
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.11.2 // indirect
	github.com/hashicorp/consul/api v1.1.0
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/openzipkin/zipkin-go v0.1.6
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/common v0.7.0 // indirect
	github.com/prometheus/procfs v0.0.5 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/smartystreets/assertions v0.0.0-20190116191733-b6c0e53d7304 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337 // indirect
	github.com/sony/gobreaker v0.4.1 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.4.0
	github.com/streadway/handy v0.0.0-20190108123426-d5acb3125c2a // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	golang.org/x/crypto v0.0.0-20191105034135-c7e5f84aec59 // indirect
	golang.org/x/net v0.0.0-20191105084925-a882066a44e0 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20191104094858-e8c54fb511f6 // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace (
	golang.org/x/crypto v0.0.0-20190313024323-a1f597ede03a => github.com/golang/crypto v0.0.0-20190313024323-a1f597ede03a
	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
