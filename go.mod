module github.com/apache/skywalking-cli

go 1.16

replace golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0 => golang.org/x/crypto v0.0.0-20201216223049-8b5274cf687f

require (
	github.com/apache/skywalking-swck/operator v0.0.0-20220720103355-e822fca26c0f
	github.com/gizak/termui/v3 v3.1.0
	github.com/google/uuid v1.3.0
	github.com/machinebox/graphql v0.2.2
	github.com/matryer/is v1.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.9
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/mum4k/termdash v0.12.1
	github.com/olekukonko/tablewriter v0.0.2
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20220421235706-1d1ef9303861 // indirect
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150 // indirect
	golang.org/x/term v0.0.0-20220411215600-e5f449aeb171 // indirect
	golang.org/x/text v0.3.7
	google.golang.org/grpc v1.40.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.22.1
	sigs.k8s.io/controller-runtime v0.10.0
	skywalking.apache.org/repo/goapi v0.0.0-20220714130828-0d56d1f4c592
)
