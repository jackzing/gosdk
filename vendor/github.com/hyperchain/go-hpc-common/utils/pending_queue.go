package utils

import (
	"errors"
	"math"
	"sync"
)

// mark is the inner event, represents the Done() op or the Wait() op, hold the seqNo as resource
type mark struct {
	seqNo  uint64 // related seqNo
	waiter chan struct{}
}

type resetEvent struct {
	seqNo  uint64 // related seqNo
	waiter chan error
}

// PendingQueueImpl is the instance for a following produce-consume model:
//   the main resource in the model is the uint64 type seqNo
//   For produce - call Done() interface, outer logic should ensure that
//                 the `seqNo` resource must be consecutively incremented(allow repeated),
//                 but the Done() interface can be called concurrently,
//                 which means the `seqNo` arriving sequence can be disordered.
//   For consume - call Wait() interface, and the Wait() interface will return only if
//                 all the prev `seqNo` resource has been continuously released by inner logic.
// the PendingQueueImpl instance cannot be re-used after Stop() interface is called,
// and all newly created instance will work after Start() interface is called
type PendingQueueImpl struct {
	id          string
	expectSeqNo uint64          // expectSeqNo is the next to-received sequential seqNo
	markCh      chan mark       // major channel for in-coming events(mark)
	resetCh     chan resetEvent // only for reset op
	// for safe close
	closeCh    chan struct{}
	closeAckCh chan struct{}
	closed     bool
	closeLock  sync.RWMutex
	log        Logger
}

// NewPendingQueue creates a PendingQueue implementation instance
func NewPendingQueue(buffer int, id string, log Logger) *PendingQueueImpl {
	return &PendingQueueImpl{
		id:      id,
		markCh:  make(chan mark, buffer),
		resetCh: make(chan resetEvent),
		closed:  true, // init as true, current instance can perform nothing
		log:     log,
	}
}

func (q *PendingQueueImpl) process() {
	// variables within the process function context

	// pending stores the received but not in-sequential seqNo,
	// each time a new `mark` comes, should iterate from q.expectSeqNo till the largest sequential one
	pending := make(map[uint64]bool)
	// waiters can be regard as the subscription of a specific seqNo,
	// when the seqNo can be released, all the subscriber should be notified
	waiters := make(map[uint64][]chan struct{})

	// pending can ensure it will receive all the sequential in-coming seqNo event
	// waiters registration is determined by outer user

	for {
		select {
		case ev := <-q.resetCh:
			if ev.seqNo >= q.expectSeqNo {
				q.log.Errorf("reset seqNo [%d] is larger than current expectedSeqNo [%d], illegal", ev.seqNo, q.expectSeqNo)
				ev.waiter <- errors.New("seqNo illegal")
				continue
			}
			q.expectSeqNo = ev.seqNo
			// re-init `pending` map, leave `waiters` map
			pending = make(map[uint64]bool)
			ev.waiter <- nil
		case ev := <-q.markCh:
			// q.Wait and q.Done ops will all be handled here

			// process waiter
			if ev.waiter != nil {
				// q.expectSeqNo is the currently not arrived one
				if ev.seqNo == uint64(math.MaxUint64) || q.expectSeqNo > ev.seqNo {
					// directly return
					ev.waiter <- struct{}{}
					continue
				}
				// set into waiters
				ws, ok := waiters[ev.seqNo]
				if !ok {
					waiters[ev.seqNo] = make([]chan struct{}, 0)
				}
				waiters[ev.seqNo] = append(ws, ev.waiter)
				continue
			}
			// process Done
			if ev.seqNo == uint64(math.MaxUint64) || q.expectSeqNo > ev.seqNo {
				q.log.Warningf("[%s] PendingQueue Done op for %d has already be triggered before", q.id, ev.seqNo)
				// directly return
				continue
			}
			pending[ev.seqNo] = true
			for {
				processing := q.expectSeqNo
				_, ok := pending[processing]
				if !ok {
					break
				}
				delete(pending, processing)
				for _, ch := range waiters[processing] {
					ch <- struct{}{}
				}
				delete(waiters, processing)
				q.expectSeqNo++
			}
		case <-q.closeCh:
			if len(pending) != 0 || len(waiters) != 0 {
				q.log.Warningf("[%s] PendingQueue instance is closing, "+
					"left with %d elements in `pending` and %d elements in `waiters`", q.id, len(pending), len(waiters))
			}
		DrainLoop:
			// drain network and executor messages, no more messages can enter.
			for {
				select {
				case <-q.markCh:
				default:
					break DrainLoop
				}
			}
			close(q.closeAckCh)
			return
		}
	}
}

// Start will change PendingQueue instance inner status and actually start working
func (q *PendingQueueImpl) Start(seqNo uint64) error {
	q.closeLock.Lock()
	defer q.closeLock.Unlock()
	if !q.closed {
		return errors.New("PendingQueue instance already started")
	}
	q.closed = false
	q.expectSeqNo = seqNo
	q.closeCh = make(chan struct{})
	q.closeAckCh = make(chan struct{})
	go q.process()
	return nil
}

// Stop will sync return after main goroutine being released
func (q *PendingQueueImpl) Stop() {
	q.closeLock.Lock()
	defer q.closeLock.Unlock()
	if q.closed {
		return
	}
	close(q.closeCh)
	<-q.closeAckCh
	q.closed = true
}

// Reset will set a new init seqNo for running session, param must be less than current running session
func (q *PendingQueueImpl) Reset(seqNo uint64) error {
	q.closeLock.Lock()
	defer q.closeLock.Unlock()
	if q.closed {
		return nil
	}
	waitCh := make(chan error)
	q.resetCh <- resetEvent{seqNo: seqNo, waiter: waitCh}
	return <-waitCh
}

// Wait triggers a wait op registration, and will sync return after param seqNo is released
func (q *PendingQueueImpl) Wait(seqNo uint64) {
	q.closeLock.RLock()
	defer q.closeLock.RUnlock()
	if q.closed {
		return
	}
	waitCh := make(chan struct{})
	q.markCh <- mark{seqNo: seqNo, waiter: waitCh}
	<-waitCh
}

// Done triggers a release op for param seqNo, directly return, may cause some Wait ops to be notified
func (q *PendingQueueImpl) Done(seqNo uint64) {
	q.closeLock.RLock()
	defer q.closeLock.RUnlock()
	if q.closed {
		return
	}
	q.markCh <- mark{seqNo: seqNo}
}
