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

## Step 13 - Job Failure Handling

- `processJob()` handles the error if the processing fails.
- If it returns nil, then processing is successful.
- A failed job is marked as FAILED
- If a job is failed, we continue the process for next job instead of stopping with help of `continue`

Flow: 
        QUEUED
           ↓
        PROCESSING
           ↓
    ┌─────────────┐
    ↓             ↓
COMPLETED       FAILED
    
## Step 14: Retry Failed Jobs

- If a job is failed, retry mechanism should handle it instead of keeping the status as `FAILED`
- Retry Count should be tracked for the particular job to check how many times it has been retried
- Max Retires Count decide how many times the retry mechanism should work instead of infinite tries
- When the retry mechanism gets executed the status gets changed to QUEUED & the job is passed to the Queue
- if the retry mechanism reached the max retries count then status gets set to FAILED.

Flow: 

      PROCESSING
           ↓
         FAILED
           ↓
RetryCount < MaxRetries?
        /       \
      YES       NO
       ↓         ↓
     QUEUED    FAILED
       ↓
     Queue
       ↓
     Worker

## Step 15: Graceful Shutdown

- Graceful shutdown allows the application to stop cleanly instead of abruptly terminating running goroutines and services.
- OS sends termination signals to the application when the process needs to stop.
- `os` → provides functionality for interacting with the operating system
- `syscall` → provides access to low-level operating-system functionality.
- `os.Interrupt` → interrupt signal commonly generated when pressing Ctrl+C.
- `syscall.SIGTERM` → Signal Terminate; commonly sent by the OS, Docker, Kubernetes, or process managers to request application termination.
- `os/signal` package allows the Go application to listen for OS signals.
- `signalChan := make(chan os.Signal, 1)` creates a buffered channel used to receive OS signals. The 1 means the channel can temporarily hold one signal without requiring an immediate receiver.
- `signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)` tells Go to send the specified OS signals to signalChan.
- `<-signalChan` receives the signal from the channel and blocks the main goroutine until a signal is received.
- Receiving a value from a channel does not mean that the channel should be closed. A channel is generally closed by the sender when no more values will be sent.

# Context Cancellation
- `context.WithCancel()` creates a cancellable context and returns a cancel() function.
- The context is passed to workers so that cancellation can be propagated to them.
- cancel() triggers the context's Done() channel.
- Workers listen for cancellation using `case <-ctx.Done():`.
- When cancellation occurs, workers return and exit their goroutines.
- `defer cancel()` ensures that the context is cancelled when the function exits.

Example:

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

Worker:

case <-ctx.Done():
    return

# HTTP Server Shutdown
- `http.ListenAndServe()` is blocking, so if it runs directly in main, main cannot continue to wait for the shutdown signal.
- The HTTP server is therefore started in a goroutine.
- Using http.Server allows the application to explicitly shut down the HTTP server.

Example:

server := &http.Server{Addr: ":8080"}
go server.ListenAndServe()

- `server.Shutdown()` gracefully stops the HTTP server.
- Graceful shutdown allows existing HTTP requests to finish while preventing new requests from being accepted.

# Shutdown Flow
Application Starts
       ↓
Create Context
       ↓
Create Signal Channel
       ↓
Register OS Signals
       ↓
Start Workers
       ↓
Start HTTP Server
       ↓
Application Running
       ↓
Ctrl+C / SIGTERM
       ↓
OS sends Signal
       ↓
signalChan receives Signal
       ↓
<-signalChan continues execution
       ↓
cancel()
       ↓
ctx.Done() triggered
       ↓
Workers Exit
       ↓
server.Shutdown()
       ↓
HTTP Server Stops
       ↓
Application Exits Cleanly

## 16. PostgreSQL Persistence

### 16.1 Database & Migration
- Created the `jobs` PostgreSQL table using `golang-migrate`.
- Added `id` as a PostgreSQL identity column so job IDs are generated by the database.
- Added migration files for schema creation and identity configuration.
- PostgreSQL database used for local development: `async_tasks`.

### 16.2 PostgreSQL Connection
- Added `pgx/v5` for PostgreSQL connectivity.
- Created a `database` package using `pgxpool.Pool`.
- Database configuration is loaded through environment variables.
- Added `.env` support using `godotenv` for local development.
- `.env` is kept outside Git using `.gitignore`.
- PostgreSQL connection URL safely handles special characters in credentials.

### 16.3 PostgreSQL Repository
- Created `PostgresJobRepository` without removing the existing in-memory repository code.
- PostgreSQL repository implements:
  - `CreateJob`
  - `GetAll`
  - `GetJobById`
  - `UpdateStatus`
  - `RetryJob`
- Repository methods accept `context.Context`.
- Database errors are returned from repository methods instead of being hidden.

### 16.4 CreateJob
- Job ID is generated by PostgreSQL using the identity column.
- `INSERT ... RETURNING` is used to insert the job and retrieve the generated job ID and other fields.
- New jobs are persisted to PostgreSQL before being added to the in-memory queue.

### 16.5 GetJob / GetAll
- `QueryRow` is used when retrieving a single job.
- `Query` is used when retrieving multiple jobs.
- `rows.Close()` and `rows.Err()` are handled when iterating over query results.
- `pgx.ErrNoRows` is returned when a requested job does not exist.

### 16.6 UpdateStatus
- Job status is updated directly in PostgreSQL.
- `updated_at` is updated using `NOW()`.
- `RowsAffected()` is used to determine whether the requested job existed.
- SQL placeholder arguments must be passed in the same order as `$1`, `$2`, etc.

### 16.7 RetryJob
- Retry count is incremented atomically in PostgreSQL:
  `retry_count = retry_count + 1`
- Retried jobs are moved back to `QUEUED`.
- `updated_at` is updated during retry.
- `RETURNING` retrieves the updated job so it can be placed back into the queue.
- SQL string values use single quotes (`'QUEUED'`); double quotes are reserved for identifiers.

### 16.8 Service Layer
- Updated the `JobRepository` interface to return errors from database operations.
- Added `context.Context` to repository and service calls.
- Service remains independent of the PostgreSQL implementation by depending on the `JobRepository` interface.
- `PostgresJobRepository` is injected into `JobService`.

### 16.9 Handler Layer
- HTTP request context (`r.Context()`) is passed to service methods.
- Database/service errors are converted into appropriate HTTP responses.
- Handlers return immediately after calling `http.Error()` to avoid writing multiple responses.

### 16.10 Worker Integration
- Workers receive jobs from the in-memory queue and update their status in PostgreSQL.
- Worker lifecycle:
  `QUEUED → PROCESSING → COMPLETED`
- Failed jobs follow:
  `QUEUED → PROCESSING → FAILED → QUEUED → retry`
- Retry count is persisted in PostgreSQL.
- Worker uses the same `JobService` and PostgreSQL repository as the HTTP handlers.
- Worker processing was tested with both successful and intentionally failing jobs.

### 16.11 Application Wiring
- PostgreSQL repository is now used at runtime instead of the in-memory repository.
- A shared `JobQueue` is created and passed to the service and worker pool.
- Multiple workers process jobs concurrently.
- Database pool is closed during application shutdown.

### 16.12 Testing Completed
- Tested job creation and PostgreSQL persistence.
- Tested retrieving all jobs.
- Tested retrieving a job by ID.
- Tested updating job status.
- Tested retry functionality.
- Tested persistence across application restart.
- Tested worker + PostgreSQL integration.
- Verified successful jobs reach `COMPLETED`.
- Verified failed jobs are retried according to the configured retry limit.

### Key Architecture Decision

PostgreSQL is the persistent source of truth for jobs.

The in-memory queue is responsible only for delivering jobs to workers.

Jobs that already exist in PostgreSQL are not automatically placed back into the in-memory queue when the application starts. A future startup/recovery mechanism may be required to load pending `QUEUED` jobs from PostgreSQL into the queue.