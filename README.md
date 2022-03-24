# jwt-cracker

HS256 JWT brute-force secret cracker in go (inspired by [lmammino/go-jwt-cracker](https://github.com/alexrsagen/go-jwt-cracker)). I made some improvements to the concurrent brute force logic.

## Dependencies

> This project uses cgo with these libraries to deal with multiple-precision number calculations.

- [gmp](https://gmplib.org)
- [mpfr](https://www.mpfr.org)

## Build

```bash
$ go build

$ # If pkg-config files are in other directories
$ PKG_CONFIG_PATH=/path/to/lib/pkgconfig go build
```

## Usage

```
Usage of ./jwt-cracker:
  -char string
        Characters to generate secret texts during brute-force process. (default "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
  -jobs int
        Number of concurrent goroutines to crack jwt (default ?(value depends on number of cores or threads))
  -maxlen int
        Max length of the secret text during brute-force process (default 12)
  -report-interval int
        Running status report interval in seconds (default 10)
  -token string
        HS256 JWT token string to crack
```

## Example

Example from [lmammino/go-jwt-cracker](https://github.com/alexrsagen/go-jwt-cracker)

```bash
$ ./jwt-cracker -token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o" -char "abcdefghijklmnopqrstuwxyz" -maxlen 6
```

### Output

```
Parsed JWT:
- Algorithm: HS256
- Type: JWT
- Payload: {"sub":"1234567890","name":"John Doe","iat":1516239022}
- Signature (hex): 5db3df6c81cc23a6ab67763ddb60618d6810cd65dc5cdaf3d2882d5617c4776a

partition: 0 // 1 startSecret: a // 25431315 endSecret: bobobo 
partition: 1 // 25431316 startSecret: bobobp // 50862630 endSecret: eeeeee 
partition: 2 // 50862631 startSecret: eeeeef // 76293945 endSecret: gtgtgt 
partition: 3 // 76293946 startSecret: gtgtgu // 101725260 endSecret: jjjjjj 
partition: 4 // 101725261 startSecret: jjjjjk // 127156575 endSecret: lzlzlz 
partition: 5 // 127156576 startSecret: lzlzma // 152587890 endSecret: oooooo 
partition: 6 // 152587891 startSecret: ooooop // 178019205 endSecret: rerere 
partition: 7 // 178019206 startSecret: rererf // 203450520 endSecret: tttttt 
partition: 8 // 203450521 startSecret: tttttu // 228881835 endSecret: xjxjxj 
partition: 9 // 228881836 startSecret: xjxjxk // 254313150 endSecret: zzzzzz 

(Partition: 0) Running: oqzsq / bobobo
(Partition: 2) Running: etokqw / gtgtgt
(Partition: 5) Running: mozici / oooooo
(Partition: 3) Running: hjdgej / jjjjjj
(Partition: 9) Running: yaikfu / zzzzzz
(Partition: 7) Running: ruducm / tttttt
(Partition: 1) Running: celzng / eeeeee
(Partition: 8) Running: ukfidm / xjxjxj
(Partition: 4) Running: jzuwee / lzlzlz
(Partition: 6) Running: peglke / rerere

Found Secret (in 16.015405 seconds): secret
```

### Time Spent

- `Apple M1 Max (2 E Cores + 8 P Cores) @ 3.22 GHz`
  - jobs: 10
  - time taken: **15 ~ 20 seconds**
- `Intel Xeon Gold 6142 (2 CPU, Total: 32 Cores 64 Threads) @ 3.7 GHz`
  - jobs: 64
  - time taken: **40 ~ 45 seconds**
- `Intel Core i9-9900K (8 Cores 16 Threads) @ 5.0 GHz`
  - jobs: 16
  - time taken: **35 ~ 40 seconds**
- `Intel Core i5-4690 (4 Corese 4 Threads) @ 3.9 GHz`
  - jobs: 4
  - time taken **215 ~ 220 seconds**