package queue

// implements FIFO string collection
type Queue []string

// returns first element of queue and deletes it from queue, panics if queue is empty
func (q *Queue) Pop() string {
	var res string = (*q)[0]
	*q = (*q)[1:]
	return res
}

// adds value to end of queue
func (q *Queue) Push(val string) {
	*q = append(*q, val)
}
