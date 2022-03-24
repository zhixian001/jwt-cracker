package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"runtime"
	"strings"
	"sync"
	"time"

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

	jobsPerEachPartition := calc.Div(numberOfCombinations, big.NewInt(int64(partitionCount)))

	partitions := make([]*SecretPartition, 0)

	// Calculate Partition
	cntr.ToMinValue()

	accumStart := cntr.ToBigInt()
	accumEnd := big.NewInt(0).Set(jobsPerEachPartition)

	for i := 0; i < partitionCount; i++ {
		fmt.Printf("partition: %d // ", i)

		fmt.Print(accumStart)
		cntr.LoadBigInt(accumStart)
		startSecret := cntr.ToString()
		fmt.Printf(" startSecret: %s // ", startSecret)

		// to end
		fmt.Print(accumEnd)
		if i == partitionCount-1 {
			accumEnd.Set(numberOfCombinations)
		}
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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	token := flag.String("token", "", "The full HS256 jwt token to crack")
	alphabet := flag.String("alphabet", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "The alphabet to use for the brute force")
	maxLength := flag.Int("maxlen", 12, "The max length of the string generated during the brute force")
	jobs := flag.Int("jobs", runtime.NumCPU(), "Max concurrent goroutines to run crack")
	reportInterval := flag.Int("report-interval", 10, "Running status report interval (in seconds)")

	flag.Parse()

	if *token == "" {
		fmt.Println("Parameter token is empty")
		flag.Usage()
		return
	}
	if *alphabet == "" {
		fmt.Println("Parameter alphabet is empty")
		flag.Usage()
		return
	}
	if *maxLength == 0 {
		fmt.Println("Parameter maxlen is 0")
		flag.Usage()
		return
	}

	if *jobs < 1 {
		fmt.Println("Parameter jobs is lower than 1")
		flag.Usage()
		return
	}

	if *reportInterval < 1 {
		fmt.Println("Report Interval is lower than 1")
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

	partitions := generateSecretPartitions(*alphabet, *maxLength, *jobs)

	fmt.Println()

	done := make(chan struct{})
	wg := &sync.WaitGroup{}
	var found bool
	var foundSecret string
	// var attempts uint64
	startTime := time.Now()

	for partitionIdx, partition := range partitions {
		wg.Add(1)
		go func(startSecret, endSecret, alphabet string, maxLength, pIdx int, wg *sync.WaitGroup, done chan struct{}) {
			partitionCounter := counter.MakeCounter(alphabet, maxLength)

			partitionCounter.LoadString(startSecret)

			currentSecret := partitionCounter.ToString()

			reportIntervalBegin := time.Now()

			var reportIntervalElapsed float64

			for {
				select {
				case <-done:
					wg.Done()
					return
				default:
					if bytes.Equal(parsed.Signature, jwt.GenerateSignature(parsed.Message, []byte(currentSecret))) {
						foundSecret = currentSecret

						found = true
						close(done)
						wg.Done()
						return
					}

					// update counter logic
					if currentSecret == endSecret {
						wg.Done()
						return
					}
					partitionCounter.Increase()
					currentSecret = partitionCounter.ToString()

					// report
					reportIntervalElapsed = time.Since(reportIntervalBegin).Seconds()
					if reportIntervalElapsed > float64(*reportInterval) {
						fmt.Printf("(Partition: %d) Running: %s / %s\n", pIdx, currentSecret, endSecret)
						reportIntervalBegin = time.Now()
					}
				}
			}

		}(partition.startSecret, partition.endSecret, partition.alphabets, partition.maxLength, partitionIdx, wg, done)
	}

	wg.Wait()

	elapsedSeconds := time.Since(startTime).Seconds()

	if !found {
		fmt.Printf("\nNo secret found (in %f seconds )\n", elapsedSeconds)
	} else {
		fmt.Printf("\nFound Secret (in %f seconds ): %s\n", elapsedSeconds, foundSecret)
	}
}
