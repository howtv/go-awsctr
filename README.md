# go-awsctr

go-awsctr is the package which is intended to make daily aws operation easy

## install

`go get github.com/howtv/go-awsctr/cmd/awsctr`


## `awsctr logs watch`

watch cloudwatch logs like tail -f

```
awsctr logs watch logGroupName \
       -format=ecs \
       -filter=someTag \
       -interval=3 \
       -region=ap-northeast-1
```