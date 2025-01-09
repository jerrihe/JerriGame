module robot

go 1.23.3

require (
	github.com/walkon/wsgnet v1.0.0
	google.golang.org/protobuf v1.36.2
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	jserver/src/common v0.0.0-00010101000000-000000000000
	jserver/src/protocol v0.0.0-00010101000000-000000000000
)

require (
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/sys v0.0.0-20220224120231-95c6836cb0e7 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace jserver/src/common => ../../src/common

replace jserver/src/protocol => ../../src/protocol
