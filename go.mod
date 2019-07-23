module github.com/sagikazarmark/modern-go-application

require (
	contrib.go.opencensus.io/exporter/ocagent v0.5.0
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	contrib.go.opencensus.io/integrations/ocsql v0.1.4
	emperror.dev/emperror v0.21.2
	emperror.dev/errors v0.4.2
	emperror.dev/handler/logur v0.1.0
	github.com/99designs/gqlgen v0.9.1
	github.com/99designs/gqlgen-contrib v0.0.0-20190222015228-c654377d611c
	github.com/InVisionApp/go-health v2.1.0+incompatible
	github.com/ThreeDotsLabs/watermill v0.4.0
	github.com/cloudflare/tableflip v1.0.0
	github.com/go-kit/kit v0.9.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.1
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/goph/idgen v0.3.0
	github.com/goph/logur v0.11.2
	github.com/gorilla/mux v1.7.0
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mccutchen/go-httpbin v0.0.0-20190116014521-c5cb2f4802fa
	github.com/moogar0880/problems v0.1.1
	github.com/oklog/run v1.0.0
	github.com/olekukonko/tablewriter v0.0.1
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/sagikazarmark/ocmux v0.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0
	github.com/vektah/gqlparser v1.1.2
	go.opencensus.io v0.22.0
	google.golang.org/grpc v1.22.0
)

replace github.com/ThreeDotsLabs/watermill v0.4.0 => github.com/ThreeDotsLabs/watermill v0.4.1-0.20190614194842-8952d2e10035
