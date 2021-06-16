module github.com/codenotary/immudb

go 1.13

require (
	github.com/fatih/color v1.9.0
	github.com/gizak/termui/v3 v3.1.0
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/jaswdr/faker v1.0.2
	github.com/o1egl/paseto v1.0.0
	github.com/peterh/liner v1.2.0
	github.com/prometheus/client_golang v1.5.1
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.9.1
	github.com/rogpeppe/go-internal v1.6.2
	github.com/rs/xid v1.2.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.8.0
	github.com/stretchr/testify v1.7.0
	github.com/takama/daemon v0.12.0
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/takama/daemon v0.12.0 => github.com/codenotary/daemon v0.0.0-20200507161650-3d4bcb5230f4
