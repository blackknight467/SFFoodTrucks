# SFFoodTrucks
requires go version greater than or equal to 1.11 because of go mod dependency manager

### Install Go
If you have a mac and use homebrew you can use

```
$ brew install go
```

otherwise you can follow the install instruction on the go website
https://golang.org/doc/install

### Build program
From the directory of this file/project run
```
go build -o SFFoodTrucks
```

this will generate an executable called "SFFoodTrucks"

### Run the program
You can now execute the program by running
```
./SFFoodTrucks  openNow
```

there is also a verbose flag if you'd like that gives more details during output

```
./SFFoodTrucks  openNow -v
```