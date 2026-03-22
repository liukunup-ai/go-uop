package parallel

import (
	"fmt"
	"sync"
)

type Task func() error

type Result struct {
	Index int
	Value interface{}
	Err   error
}

type Executor struct {
	maxWorkers int
	results    chan Result
}

func NewExecutor(maxWorkers int) *Executor {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	return &Executor{
		maxWorkers: maxWorkers,
		results:    make(chan Result, maxWorkers),
	}
}

func (e *Executor) Execute(tasks []Task) []Result {
	if len(tasks) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, e.maxWorkers)

	for i, task := range tasks {
		wg.Add(1)
		go func(index int, t Task) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			var result Result
			result.Index = index
			result.Err = t()
			e.results <- result
		}(i, task)
	}

	go func() {
		wg.Wait()
		close(e.results)
	}()

	var results []Result
	for r := range e.results {
		results = append(results, r)
	}

	return results
}

func (e *Executor) ExecuteWithValues(inputs []interface{}, fn func(interface{}) error) []Result {
	tasks := make([]Task, len(inputs))
	for i, input := range inputs {
		inp := input
		tasks[i] = func() error {
			return fn(inp)
		}
	}
	return e.Execute(tasks)
}

type Pool struct {
	tasks   chan Task
	results chan Result
	wg      sync.WaitGroup
}

func NewPool(workers int, queueSize int) *Pool {
	if queueSize <= 0 {
		queueSize = 100
	}
	p := &Pool{
		tasks:   make(chan Task, queueSize),
		results: make(chan Result, queueSize),
	}

	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}

	return p
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()
	for task := range p.tasks {
		result := Result{Index: id}
		result.Err = task()
		p.results <- result
	}
}

func (p *Pool) Submit(task Task) {
	p.tasks <- task
}

func (p *Pool) SubmitAndWait(task Task) Result {
	done := make(chan Result, 1)
	p.tasks <- func() error {
		result := Result{}
		result.Err = task()
		done <- result
		return result.Err
	}
	return <-done
}

func (p *Pool) Close() {
	close(p.tasks)
	p.wg.Wait()
	close(p.results)
}

func CollectResults(results []Result) error {
	for _, r := range results {
		if r.Err != nil {
			return fmt.Errorf("task %d: %w", r.Index, r.Err)
		}
	}
	return nil
}
