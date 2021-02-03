module github.com/apache/skywalking-cli

go 1.13

replace skywalking/network v1.0.0 => ./gen-codes/skywalking/network

require (
	github.com/apache/skywalking-swck v0.0.0-20210107023854-d15ef19f8317
	github.com/gizak/termui/v3 v3.1.0
	github.com/machinebox/graphql v0.2.2
	github.com/mattn/go-runewidth v0.0.9
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/mum4k/termdash v0.12.1
	github.com/olekukonko/tablewriter v0.0.2
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.0
	github.com/urfave/cli v1.22.1
	google.golang.org/grpc v1.35.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.19.3
	sigs.k8s.io/controller-runtime v0.7.0-alpha.6
	skywalking/network v1.0.0
)
