# FancyJobRunner

Run fan-out jobs like downloading, AI inference, and more on large groups of workers using only SSH and Docker.

Automatically schedules inputs (e.g. a link to download) and resources (e.g. a proxy url) across worker nodes.

Captures stdout and stderr from all workers, and provide simple metrics of Job and Task execution

## Why

FancyJobRunner was built to fulfill many of the data collection and processing jobs required for modern large AI tasks. For example downloading millions of files from datasets to your own S3 in parallel across many machines, so each can have their network cards full saturated.

Or downloading tar files, untaring them, and running them through AI models across many GPU workers in parallel.

Having a system that could fairly schedule download links, lock proxy urls to individual workers, and queue task executions across workers is essential to spending less time collecting data and more time acting on it.

## Entity Model

### Job

A Job is a definition of a task with many inputs.

For example, a job might be to download thousands of files from a provided list of URLs, where each execution of that download (a task) is the same.

### Task

A task is a single execution instance of a Job.

A task is given an input (e.g. a URL to download), possibly a Resource (e.g. a proxy url), and 

### Worker Resource

A Worker Resource is a pool of resources that are automatically given to each worker for the duration of a Job. Each item in the Worker Resource pool is guaranteed to only be given to a single worker, and no more workers are recruited for task execution than items in the pool.

This is useful for things like proxy URLs, where you want to give each worker its own proxy, and ensure that no other worker uses the same URL, nor that any worker executes tasks without one.

### Task Resource

A Task Resource, like a Worker Resource, is a pool of resources that is given out per task. Unlike Worker Resources, Task Resources are released after task execution, and are available to be provided to the next

## Writing a Task Executing Docker Container

FancyJobRunner was designed to allow your code to be very agnostic to its environment. The only requirement is that you can provide a docker container that your worker nodes can pull, and 

### Docker

You can configure a private docker registry so that your worker nodes will log in with docker before pulling the image for task execution.

## Configuration

The `fjr.y(a)ml` file follows the (LINK) JSON schema.

The configuration file uses go templates, which is given the [`JobRun` struct](). Anything in this can be templated into the configuration.

### worker_resources

An array of named worker resources.

The resources will be available on the worker as `WORKER_RESOURCE_{NAME}`.

⚠️ **WARNING:**

The total number of workers used can be limited by the smallest amount of available worker resources. For example, if you have 10 workers, 8 `PROXY` resources, no more than 8 task instances will be executed concurrently, leaving 2 workers excluded. If multiple resources are used, the total number of concurrent workers is equivalent to the fewest resources. effectively this is calculated as `Min(workers, resource_a, resource_b, ...)`.

Example:
```yaml
worker_resources:
  - PROXY: # each worker gets a proxy from this list
      file:
        path: proxys.txt
        format: eachLine
```

### Functions

There are a few helper functions provided through go templates as well.

#### Exec

Returns the output of some terminal function. Not multi-line safe (will break YAML parsing).

Example:

```yaml
docker:
  registry:
    password: {{ Exec "aws ecr get-login-password" }} # replaces with the output of this command
```

#### Env

Returns the environment variable.

Example:
```yaml
docker:
  environment:
    AWS_ACCESS_KEY_ID: {{ Env "AWS_ACCESS_KEY_ID" }} # replaces with the env var
    AWS_SECRET_ACCESS_KEY: {{ Env "AWS_SECRET_ACCESS_KEY" }} # replaces with the env var
```