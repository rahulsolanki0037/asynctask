# Async Task - Go Notes

## Step 1 - HTTP Server
- `net/http` provides HTTP Server functionality in Go.
- `http/HandleFunc()` register a handler for a route.
- `http.ListenAndServe()` starts the HTTP server.
- `http.ResponseWriter` --> used to send HTTP response.
- `*http.Request` --> contains the incoming HTTP request.

Flow: Client --> HTTP Server --> Handler --> Response

## Step 2 - HTTP Handler

Handler Signature:

func handler(w http.ResponseWriter, r *http.Request)
- `r.Method` --> HTTP Method
- `r.Body` --> request body
- `r.URL` --> request URL
- `w.Header` --> response headers
- `w.WriteHeader` --> response status code
- `http.Error` --> response error

Example:
GET /health --> 200 OK
Health endpoint: {"status": "UP"}

## Step 3 - Create Job (POST /jobs)

Accepts information from the client

Request: 
POST /jobs
{
    "Type" : "Generate Report",
    "Payload" : "Monthly Sales Report"
}

Concepts:
- json.NewDecoder(r.Body).Decode(&request) --> Decodes the request and converts to Go struct
- json.NewEncoder(w).Encode(job) --> Encodes the response and converts it to JSON response.
- Request model & job model are differnt.

Create Job Request:
- Type
- Payload

Job:
- ID
- Type
- Payload
- Status

## Step 4 — Repository

Separate data storage logic from the rest of the application.
Initially used an in-memory map: map[int]model.Job

Example:

1 → Job 1
2 → Job 2
3 → Job 3

`nextID` generates unique IDs.

Repository operations:

Create()
GetAll()
GetByID()

Repository is currently in-memory, so data is lost when the application stops.

## Step 5 — GET APIs

GET /jobs → returns all jobs.
GET /jobs/{id} → returns one job.

`strconv.Atoi()` --> converts string → integer.
Example: "2" → 2

Repository lookup:

job, exists := repository.GetByID(id)

Go commonly uses the second boolean value to indicate whether a value exists.

HTTP status codes used:

200 → successful GET
201 → resource created
400 → bad request
404 → resource not found
405 → method not allowed

## Step 6 — Service Layer

Before: Handler → Repository
After: Handler → Service → Repository

Responsibilities:

Handler:
- HTTP request/response
- Decode request
- Encode response

Service:
- Business/application logic
- Coordinates operations

Repository:
- Data storage/retrieval

Why:
- Separation of concerns
- Easier maintenance
- Easier testing
- Easier to extend

Dependency wiring:

Repository
    ↓
Service
    ↓
Handler

## Step 7 - Interfaces

- Interfaces define a set of required methods
- Go interfaces are satisfied implicitly.
- A type satisfies an interface when the required methods are implemented.
- Service depends on repository interface instead of a concrete repository.
- Loose coupling, easier testing, easy to replace implementations.


Flow: Handler --> Service --> Repository Interface (PostgreSQL / In Memory)

## Step 8 - Concurrency-Safe Repository

- Go maps are not safe for concurrent writes.
- Multiple goroutines accessing the shared data can cause data race
- sync.Mutex allows mutual exclusion
- Lock() allows one goroutine to access the shared data at a time
- Unlock() releases the lock
- defer Unlock() ensures the lock is released when the function exits.

## Step 9 - Job Queue

- A channel can be used as a queue between goroutines
- make (chan model.Job, size) creates a buffered channel
- chan <- model.Job sends the required job in the channel
- <- chan model.Job exposes it as a receiver type channel
- Service stores the job & puts it in the Queue.

We expose a receive-only channel because Workers only need to take jobs from the queue, one by one, and process them. Workers don't need to send jobs back into the queue.

We use a buffered channel because it acts as a queue that can temporarily hold multiple jobs while workers are busy.

## Step 10 - Worker

- Worker consumes job from the queue
- `for job := range channel` receives job one by one
- Worker waits when the channel has no jobs
- This workers runs on its own goroutine.
- `go worker.Start()` allows the worker to run concurrently with the actual server.

Flow : Handler --> Service --> Queue --> Job --> Worker

## Step 11 - Worker Pool

- A worker pool consists of multiple workers processing on jobs from the same queue
- Multiple jobs get processed concurrently
- Each worker runs in its own goroutine
- Worker count controls the maximum processing concurrency.
- Worker pools provide controlled concurrency instead of creation one goroutine across a job.

## Step 12 - Job Status Lifecycle

Lifecycle : QUEUED --> PROCESSING --> COMPLETED

- QUEUED --> When you create a Job
- PROCESSING --> Worker has picked the job
- COMPLETED --> Worker has processed the job

UpdateStatus() updated the status for the job & repository access is protected using mutex.

If we try to call GetById(jobId) in UpdateStatus(), updateStatus looks first & then call GetById() which locks again and it goes it waiting state forever because Go's `sync.Mutex` is not reentrant. The same goroutine cannot lock the same mutex twice without unlocking it first.
 
Note: `Don't call another repository method that acquires the same mutex while you are already holding that mutex`.