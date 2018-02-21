// Copyright © 2013 Felipe Rodrigues <lfrs.web@gmail.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"os"
	"time"

	"github.com/thr27/go-mailer/conf"
	"github.com/thr27/go-mailer/queue"
)

// Check queue and proccess each file in queue
// in a different goroutine.
func run() {
	q := queue.New()

	hasQueue, err := q.HasQueue()
	if err != nil {
		log.Println(err)
		return
	}

	if hasQueue {
		log.Printf("%d files in queue.", len(q.Files))

		for _, file := range q.Files {
			go q.Process(file)
		}
	} else {
		log.Println("Queue is empty")
	}
}

func main() {
	// Prepare log file
	f, err := os.OpenFile("log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Println("Starting")
	log.Printf("Using config %s\n", conf.Path)
	log.Println(conf.String())

	wait, _ := time.ParseDuration(conf.WaitFor())

	// Infinit loop
	for {
		run()

		log.Printf("Waiting %s\n", conf.WaitFor())
		time.Sleep(wait)
	}
}
