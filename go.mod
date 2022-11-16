module github.com/jackzing/gosdk
require (
	github.com/buger/jsonparser v1.1.1
	github.com/davecgh/go-spew v1.1.1
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/google/gofuzz v1.2.0
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/go-uuid v1.0.2
	github.com/hyperchain/go-crypto-gm v0.2.25
	github.com/hyperchain/go-crypto-standard v0.2.10
	github.com/hyperchain/go-hpc-common v0.3.9
	github.com/hyperchain/go-hpc-msp v0.2.14
	github.com/json-iterator/go v1.1.12
	github.com/magiconair/properties v1.8.6
	github.com/meshplus/crypto v0.0.15
	github.com/mholt/archiver/v3 v3.5.1
	github.com/mitchellh/mapstructure v1.5.0
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/segmentio/kafka-go v0.4.35
	github.com/spf13/viper v1.13.0
	github.com/streadway/amqp v0.0.0-20180528204448-e5adc2ada8b8
	github.com/stretchr/testify v1.8.0
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	google.golang.org/grpc v1.46.2
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
)

require (
	github.com/andybalholm/brotli v1.0.1 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/fatih/set v0.2.1 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/klauspost/compress v1.15.7 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.5 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pingcap/errors v0.11.4 // indirect
	github.com/pingcap/failpoint v0.0.0-20220423142525-ae43b7f4e5c3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.6.1 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	github.com/ulikunitz/xz v0.5.9 // indirect
	github.com/willf/bitset v1.1.11 // indirect
	github.com/willf/bloom v2.0.3+incompatible // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/net v0.0.0-20220921203646-d300de134e69 // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
	golang.org/x/text v0.3.8 // indirect
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/hyperchain/go-crypto-standard => git.hyperchain.cn/hyperchain/go-crypto-standard.git v0.2.10

replace github.com/hyperchain/go-crypto-gm => git.hyperchain.cn/hyperchain/go-crypto-gm.git v0.2.25

replace github.com/hyperchain/go-hpc-msp => git.hyperchain.cn/hyperchain/go-hpc-msp.git v0.3.0

replace github.com/hyperchain/go-hpc-common => git.hyperchain.cn/hyperchain/go-hpc-common.git v0.3.12

go 1.17
