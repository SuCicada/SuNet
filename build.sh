#!/bin/bash
# https://www.cnblogs.com/cangqinglang/p/12101493.html

os_all='aix darwin dragonfly freebsd linux netbsd openbsd solaris windows'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle'

bin_dir=$(dirname $0)/bin
go_src="port_forward process"

if [ ! -f $bin_dir ]; then
  mkdir -p $bin_dir
  echo mkdir -p $bin_dir
fi

#并发数
threadTask=100
#创建fifo管道
fifoFile="/tmp/test_fifo"
rm ${fifoFile}
mkfifo ${fifoFile}
# 建立文件描述符关联
exec 9<>${fifoFile}
rm -f ${fifoFile}
# 预先向管道写入数据
for ((i = 0; i < ${threadTask}; i++)); do
  echo "" >&9
done
echo "wait all task finish,then exit!!!"

for src in $go_src; do
  for os in $os_all; do
    for arch in $arch_all; do
      read -u9
      {
        exe=''
        if [ $os == "windows" ]; then
          exe=".exe"
        fi
        file=$bin_dir/$src/$src-$os-$arch$exe
        srcfile=cmd/${src}.go
        echo build $src with $os-$arch in $file
        CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -o $file $srcfile
        echo "" >&9
      } &
    done
  done
done
wait
# 关闭管道
exec 9>&-
echo
echo "success"
