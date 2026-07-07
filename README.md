# taskdispatcher

A concurrent job dispatcher in Go, built as a focused exercise in goroutines,
channels, and `select`. It exposes a TCP server for submitting jobs and an
HTTP server for polling their status, backed by a channel-based semaphore
gatekeeper with priority-based wait timeouts across three job classes.

## How it works

- Jobs belong to one of three classes, each with a different tolerance for
  waiting on a free worker slot before being rejected as busy:

  | Class  | Wait tolerance | Simulated work duration |
  |--------|-----------------|--------------------------|
  | class1 | indefinite      | 5s                       |
  | class2 | 15s             | 10s                      |
  | class3 | 5s              | 15s                      |

- A fixed pool of worker slots (default: 5) is enforced by a buffered channel
  acting as a semaphore.
- When all slots are busy, an incoming job is queued by class rather than
  rejected outright. Whenever a slot frees up, the dispatcher checks queues
  in priority order (class3 → class2 → class1) and dispatches the next
  eligible job — so a higher-priority job can jump ahead of jobs that arrived
  earlier.
- If a queued job's wait tolerance expires before a slot becomes available,
  it's marked `worker_busy` instead of running.
- All coordination happens through channels owned by a single dispatcher
  goroutine — no locks are needed on the job queues themselves.

State is in-memory only and does not persist across restarts; this is
intentional, since the exercise is about concurrency primitives, not storage.

## Running it

```
go run .
```

This starts two servers in the same process:
- TCP server on `:9000` — for submitting jobs.
- HTTP server on `:8080` — for polling job status.

## Submitting a job

Send a single line over TCP: `<class> <name>`.

```
echo "class1 mytask" | nc localhost 9000
```

Response:
```
task_id: 1
```

## Checking job status

```
curl localhost:8080/status/1
```

Response:
```json
{"id":1,"task_name":"mytask","class":1,"status":"running"}
```

`status` will be one of: `queued`, `running`, `done`, `failed`, `worker_busy`.

## Trying the priority/timeout behavior

Saturate all workers with class1 jobs, then submit a class3 job and watch it
either jump the queue once a slot frees, or time out to `worker_busy` if none
frees within 5 seconds:

```
for i in 1 2 3 4 5; do echo "class1 task$i" | nc localhost 9000; done
echo "class3 urgent" | nc localhost 9000
```

Poll the returned task IDs over the next several seconds to watch the state
transitions.

## Project structure

```
taskdispatcher/
├── main.go                       # wiring: starts dispatcher, TCP, HTTP
├── models/
│   └── models.go                 # Task, JobClass, Status types
├── internal/
│   ├── gatekeeper/
│   │   └── gatekeeper.go         # buffered-channel semaphore
│   └── dispatcher/
│       └── dispatcher.go         # queues, priority dispatch, timeouts
└── server/
    ├── tcp.go                    # job submission
    └── http.go                   # status polling
```