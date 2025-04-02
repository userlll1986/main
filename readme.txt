go初始安装指南
brew install go
配置go环境变量
vim ~/.bash_profile
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
查看环境变量 go env
初始化go.mod用来管理包的关键文件
go mod init + 模块名称
本地包的导入
例如在项目下新建目录 utils，创建一个tools.go文件
在根目录下的hello.go文件就可以 import “hello/utils” 引用utils
go常用的指令：
查看module下的所有依赖
go list -m all
清理无用的依赖
go mod tidy

gin项目配置
增加处理yaml文件的包
go get -u gopkg.in/yaml.v2
先添加gorm和mysql配置
go get -u github.com/gin-gonic/gin
go get -u github.com/go-sql-driver/mysql




