# go-jq

Go bindings for jq inspired by ruby-jq.

see [http://stedolan.github.io/jq/](http://stedolan.github.io/jq/).
see [https://bitbucket.org/winebarrel/ruby-jq](https://bitbucket.org/winebarrel/ruby-jq)

## Installation

First, please install libjq from HEAD of [git repository](https://github.com/stedolan/jq).

```sh
git clone https://github.com/stedolan/jq.git
cd jq
autoreconf -i
./configure --enable-shared
make
sudo make install
sudo ldconfig
```

## Usage

```go
package main

import (
   "github.com/bongole/go-jq"
)

func main(){
   src := "{\"foo\": 1}

   jq := jq.New(src)

   d, err := jq.Search(".foo")
   if err != nil {
      i := d.(float64)
      fmt.Printf("%d", i)
   }
}
```

## License
* MIT
