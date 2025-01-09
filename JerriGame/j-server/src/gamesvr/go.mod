module gamesvr

go 1.23.3

require (
	github.com/redis/go-redis/v9 v9.7.0
	github.com/walkon/wsgnet v1.0.0
	google.golang.org/protobuf v1.36.2
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.12
	jserver/src/common v0.0.0-00010101000000-000000000000
	jserver/src/protocol v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace jserver/src/common => ../common

replace jserver/src/protocol => ../protocol
