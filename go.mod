module github.com/apache/skywalking-cli

go 1.16

replace golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0 => golang.org/x/crypto v0.0.0-20201216223049-8b5274cf687f

require (
	github.com/apache/skywalking-swck v0.2.0
	github.com/gizak/termui/v3 v3.1.0
	github.com/machinebox/graphql v0.2.2
	github.com/mattn/go-runewidth v0.0.9
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/mum4k/termdash v0.12.1
	github.com/olekukonko/tablewriter v0.0.2
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/crypto v0.0.0-20201216223049-8b5274cf687f // indirect
	google.golang.org/grpc v1.40.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.21.1
	sigs.k8s.io/controller-runtime v0.7.0
	skywalking.apache.org/repo/goapi v0.0.0-20220322033350-0661327d31e3
)
