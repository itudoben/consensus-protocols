package main

import (
	"flag"
	"hash/fnv"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

// Number of philosophers is simply the length of this list.
var ph = []string{"Mark", "Russell", "Rocky", "Haris", "Root"}

const hunger = 3                // Number of times each philosopher eats
const think = time.Second / 100 // Mean think time
const eat = time.Second / 100   // Mean eat time

var fmt = log.New(os.Stdout, "", 0)

var dining sync.WaitGroup

func diningProblem(phName string, dominantHand, otherHand *sync.Mutex) {
	fmt.Println(phName, "Seated")
	h := fnv.New64a()
	h.Write([]byte(phName))
	rg := rand.New(rand.NewSource(int64(h.Sum64())))
	rSleep := func(t time.Duration) {
		time.Sleep(t/2 + time.Duration(rg.Int63n(int64(t))))
	}
	for h := hunger; h > 0; h-- {
		fmt.Println(phName, "Hungry")
		dominantHand.Lock() // pick up forks
		otherHand.Lock()
		fmt.Println(phName, "Eating")
		rSleep(eat)
		dominantHand.Unlock() // put down forks
		otherHand.Unlock()
		fmt.Println(phName, "Thinking")
		rSleep(think)
	}
	fmt.Println(phName, "Satisfied")
	dining.Done()
	fmt.Println(phName, "Left the table")
}

func diningPhilosophers() {
	fmt.Println("Table empty")
	dining.Add(5)
	fork0 := &sync.Mutex{}
	forkLeft := fork0
	for i := 1; i < len(ph); i++ {
		forkRight := &sync.Mutex{}
		go diningProblem(ph[i], forkLeft, forkRight)
		forkLeft = forkRight
	}
	go diningProblem(ph[0], fork0, forkLeft)
	dining.Wait() // wait for philosphers to finish
	fmt.Println("Table empty")
}

// Worker example

func worker(part string) {
	log.Println(part, "worker begins part")
	time.Sleep(time.Duration(rand.Int63n(1e6)))
	log.Println(part, "worker completes part")
	wg.Done()
}

var (
	partList    = []string{"A", "B", "C", "D"}
	nAssemblies = 3
	wg          sync.WaitGroup
)

func partWorkers() {
	rand.Seed(time.Now().UnixNano())
	for c := 1; c <= nAssemblies; c++ {
		log.Println("begin assembly cycle", c)
		wg.Add(len(partList))
		for _, part := range partList {
			go worker(part)
		}
		wg.Wait()
		log.Println("assemble.  cycle", c, "complete")
	}
}

// Producer/Consumer example

type Consumer struct {
	msgs *chan int
}

// NewConsumer creates a Consumer
func NewConsumer(msgs *chan int) *Consumer {
	return &Consumer{msgs: msgs}
}

// consume reads the msgs channel
func (c *Consumer) consume() {
	fmt.Println("consume: Started")
	for {
		msg := <-*c.msgs
		fmt.Println("consume: Received:", msg)
	}
}

// Producer definition
type Producer struct {
	msgs *chan int
	done *chan bool
}

// NewProducer creates a Producer
func NewProducer(msgs *chan int, done *chan bool) *Producer {
	return &Producer{msgs: msgs, done: done}
}

// produce creates and sends the message through msgs channel
func (p *Producer) produce(max int) {
	fmt.Println("produce: Started")
	for i := 0; i < max; i++ {
		fmt.Println("produce: Sending ", i)
		*p.msgs <- i
	}
	*p.done <- true // signal when done
	fmt.Println("produce: Done")
}

func producerConsumer() {
	// profile flags
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile := flag.String("memprofile", "", "write memory profile to `file`")

	// get the maximum number of messages from flags
	max := flag.Int("n", 5, "defines the number of messages")

	flag.Parse()

	// utilize the max num of cores available
	runtime.GOMAXPROCS(runtime.NumCPU())

	// CPU Profile
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	var msgs = make(chan int)  // channel to send messages
	var done = make(chan bool) // channel to control when production is done

	// Start a goroutine for Produce.produce
	go NewProducer(&msgs, &done).produce(*max)

	// Start a goroutine for Consumer.consume
	go NewConsumer(&msgs).consume()

	// Finish the program when the production is done
	<-done

	// Memory Profile
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}

func main() {
	diningPhilosophers()
	partWorkers()
	producerConsumer()
}
