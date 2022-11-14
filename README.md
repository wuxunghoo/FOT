[![pipeline status](https://api.travis-ci.org/33cn/plugin.svg?branch=master)](https://travis-ci.org/33cn/plugin/)
[![Go Report Card](https://goreportcard.com/badge/github.com/33cn/plugin?branch=master)](https://goreportcard.com/report/github.com/33cn/plugin)


# 基于 chain33 区块链开发 框架 开发的 FOT公有链系统


### 编译

```
git clone https://github.com/wuxunghoo/FOT $GOPATH/src/github.com/wuxunghoo/FOT
cd $GOPATH/src/github.com/wuxunghoo/FOT
go build -i -o fot
go build -i -o fot-cli github.com/wuxunghoo/FOT/cli
```

### 运行
拷贝编译好的fot, fot-cli, fot.toml这三个文件置于同一个文件夹下，执行：
```
./fot -f fot.toml
```
