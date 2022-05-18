# Encrypt Secret

## Install
Go v1.17 or above
```
    go install github.com/wijhuang/go-secure
```
below
```
    go get github.com/wijhuang/go-secure
```


## Usage
```
usage of go-secure
  -k string
        key to decrypt/encrypt in base64 if empty random string will be generated
  -m int
        mode: 0: encrypt(DEFAULT), 1: decrypt
  -o string
        output file path (default "output.txt")
```

## Encrypt
```
    go-secure example-secret.json
```

## Decrypt
```
    go-secure -m=1 -k=ry3M2iGzPO9jkG2TZYx2bStxvILZt79QdDvZeatEiAU= example-output.txt
```

## Removing
```
    rm $GOPATH/bin/go-secure
```