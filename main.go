package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"runtime"
	"strings"
	"sync"

	"github.com/zhixian001/jwt-cracker/counter"
	"github.com/zhixian001/jwt-cracker/jwt"
)

type SecretPartition struct {
	startSecret string
	endSecret   string
	alphabets   string
	maxLength   int
}

func generateSecretPartitions(alphabet string, maxLength int, partitionCount int) []*SecretPartition {
	cntr := counter.MakeCounter(alphabet, maxLength)
	cntr.ToMaxValue()

	numberOfCombinations := cntr.ToBigInt()

	calc := big.NewInt(0)

	// cntr.ToMinValue()
	// initialSecret := cntr.ToString()

	jobsPerEachPartition := calc.Div(numberOfCombinations, big.NewInt(int64(partitionCount)))

	partitions := make([]*SecretPartition, partitionCount)

	// Calculate Partition
	cntr.ToMinValue()

	accumStart := cntr.ToBigInt()
	accumEnd := big.NewInt(0).Set(jobsPerEachPartition)

	for i := 0; i < partitionCount; i++ {
		fmt.Printf("partition: %d // ", i)

		fmt.Print(accumStart)
		// cntr.LoadBigIntSlow(accumStart)
		cntr.LoadBigInt(accumStart)
		startSecret := cntr.ToString()
		fmt.Printf(" startSecret: %s // ", startSecret)

		// to end
		fmt.Print(accumEnd)
		if i == partitionCount-1 {
			accumEnd.Set(numberOfCombinations)
		}
		// cntr.LoadBigIntSlow(accumEnd)
		cntr.LoadBigInt(accumEnd)

		endSecret := cntr.ToString()
		fmt.Printf(" endSecret: %s \n", endSecret)

		// update accum
		accumStart.Add(accumStart, jobsPerEachPartition)
		accumEnd.Add(accumEnd, jobsPerEachPartition)

		// add partition

		toAdd := &SecretPartition{
			startSecret: startSecret,
			endSecret:   endSecret,
			alphabets:   alphabet,
			maxLength:   maxLength,
		}
		partitions = append(partitions, toAdd)
	}

	return partitions
}

func generateSecrets(alphabet string, n int, wg *sync.WaitGroup, done chan struct{}) <-chan string {
	if n <= 0 {
		return nil
	}

	c := make(chan string)

	wg.Add(1)
	go func() {
		defer close(c)
		var helper func(string)
		helper = func(input string) {
			if len(input) == n {
				return
			}
			select {
			case <-done:
				return
			default:
			}
			for _, char := range alphabet {
				s := input + string(char)
				c <- s
				helper(s)
			}
		}
		helper("")
		wg.Done()
	}()

	return c
}

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	token := flag.String("token", "", "The full HS256 jwt token to crack")
	alphabet := flag.String("alphabet", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "The alphabet to use for the brute force")
	// prefix := flag.String("prefix", "", "A string that is always prefixed to the secret")
	// suffix := flag.String("suffix", "", "A string that is always suffixed to the secret")
	maxLength := flag.Int("maxlen", 12, "The max length of the string generated during the brute force")
	jobs := flag.Int("jobs", runtime.NumCPU(), "Max concurrent goroutines to run crack")

	flag.Parse()

	if *token == "" {
		fmt.Println("Parameter token is empty\n")
		flag.Usage()
		return
	}
	if *alphabet == "" {
		fmt.Println("Parameter alphabet is empty\n")
		flag.Usage()
		return
	}
	if *maxLength == 0 {
		fmt.Println("Parameter maxlen is 0\n")
		flag.Usage()
		return
	}

	if *jobs < 1 {
		fmt.Println("Parameter jobs is lower than 1")
		flag.Usage()
		return
	}

	err, parsed := jwt.ParseJWT(*token)
	if err != nil {
		fmt.Printf("Could not parse JWT: %v\n", err)
		return
	}

	fmt.Printf("Parsed JWT:\n- Algorithm: %s\n- Type: %s\n- Payload: %s\n- Signature (hex): %s\n\n",
		parsed.Header.Algorithm,
		parsed.Header.Type,
		parsed.Payload,
		hex.EncodeToString(parsed.Signature))

	if strings.ToUpper(parsed.Header.Algorithm) != "HS256" {
		fmt.Println("Unsupported algorithm")
		return
	}

	generateSecretPartitions(*alphabet, *maxLength, *jobs)

	// fmt.Println(partitions)

	// combinations := big.NewInt(0)
	// for i := 1; i <= *maxLength; i++ {
	// 	alen, mlen := big.NewInt(int64(len(*alphabet))), big.NewInt(int64(i))
	// 	combinations.Add(combinations, alen.Exp(alen, mlen, nil))
	// }
	// fmt.Printf("There are %s combinations to attempt\nCracking JWT secret...\n", combinations.String())

	// done := make(chan struct{})
	// wg := &sync.WaitGroup{}
	// var found bool
	// var attempts uint64
	// startTime := time.Now()

	// for secret := range generateSecrets(*alphabet, *maxLength, wg, done) {
	// 	wg.Add(1)
	// 	go func(s string, i uint64) {
	// 		select {
	// 		case <-done:
	// 			wg.Done()
	// 			return
	// 		default:
	// 		}
	// 		if bytes.Equal(parsed.Signature, jwt.GenerateSignature(parsed.Message, []byte(*prefix+s+*suffix))) {
	// 			fmt.Printf("Found secret in %d attempts: %s\n", attempts, *prefix+s+*suffix)
	// 			found = true
	// 			close(done)
	// 		}
	// 		wg.Done()
	// 	}(secret, attempts)

	// 	attempts++
	// 	if attempts%100000 == 0 {
	// 		elapsedSeconds := time.Since(startTime).Seconds()
	// 		startTime = time.Now()
	// 		fmt.Printf("Attempts: %d /// APS: %.3f a/s\n", attempts, 100000/elapsedSeconds)
	// 	}
	// }

	// wg.Wait()
	// if !found {
	// 	fmt.Printf("No secret found in %d attempts\n", attempts)
	// }
}
