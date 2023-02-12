## Challenge

In this exercise we ask you to write a command line REPL (read-eval-print loop) that drives a simple in-memory key/value storage system. This system should also allow for nested transactions. A transaction can then be committed or aborted.

Please see instructions at the bottom of this doc about how to submit your code.

## EXAMPLE RUN

```
$ my-program
> WRITE a hello
> READ a
hello
> START
> READ a
hello
> WRITE a hello-again
> READ a
hello-again
> START
> DELETE a
> READ a
Key not found: a
> COMMIT
> READ a
Key not found: a
> WRITE a once-more
> READ a
once-more
> ABORT
> READ a
hello
> QUIT
Exiting...
```

## COMMANDS

* `READ <key>` Reads and prints, to stdout, the val associated with key. If the value is not present an error is printed to stderr.
* `WRITE <key> <val>` Stores val in key.
* `DELETE <key>` Removes a key from the store. Future READ commands on that key will return an error.
* `START` Start a transaction.
* `COMMIT` Commit a transaction. All actions in the current transaction are committed to the parent transaction or the root store. If there is no current transaction an error is output to stderr.
* `ABORT` Abort a transaction. All actions in the current transaction are discarded.
* `QUIT` Exit the REPL cleanly. A message to stderr may be output.

## OTHER DETAILS

For simplicity, all keys and values are simple ASCII strings delimited by whitespace. No quoting is needed.
All errors are output to stderr.
Commands are case-insensitive.
As this is a simple command line program with no networking, there is only one “client” at a time. There is no need for locking or multiple threads.



# Approach
In terms of approach I chose to use a Linked-List data structure
for connecting transactions together. The most recent/head transaction
is placed at the store root so that it can be easily accessed for read/writes.

When a transaction is then committed the global state of that transaction is either
passed down the chain to the next transaction or if there is no other transactions
it will change the stores global state.

**See `simple-kv-store` for main implementation**
## Assumptions
Keys are case-sensitive. Meaning that you can store a key of `anApple` and `ANAPPLE`.
I did not see value in adding complexity of varying cases in the kv,store. Also comparing
to other kv stores, like REDIS keys are case-insensitive.

Keys and values also must not be blank, i.e just be whitespace.

## Improvements
* Provide further testcases for the integration/CLI part. Due to time I focussed on writing
  unit tests on the store and then covering the base scenario from the README. I chose to use a
  `testdata` folder to allow a text file for `stdin` to be passed in and assert the `stdout` response.
  I contemplated writing commands in the test but felt `testdata` was easier to run through multiple commands.
* Depending on how this tool could be used Docker could be introduced to package the application.
* Could add CI/CD  to this to ensure all tests pass/build compiles, etc.
* For simplicity I'm ignoring error in the CLI for `fmt.Fprintf` this can error if there are any write issues so it should be handled.
* CLI Tool
  * Due to the requirement of commands being `case-sensitive` I went with the approach
    of using a basic buf scanner. This then means I have to `strings.Upper` incoming
    requests and add a switch for each command. Using something like Go's built-in flags
    or a library like [Cobra](https://github.com/spf13/cobra) could be a good way of creating it/
    generating the CLI as it has a lot of useful features. However, both are case-sensitive when it comes to commands.
  * I added a very basic help command so that you don't need to switch to README to see commands. This could be extended
    so that for any time a command is entered incorrectly a help menu for that specific command is shown.


# Usage
In order to use the CLI Tool there is a pre-requisite of [Go](https://github.com/golang/go) being installed.
Execute `make run` to run the program and from there you are able to send in commands.

## Makefile
* `make lint` - run linter across project
* `make run` - to run the CLI (requires Go to be installed)
* `make test` - executes tests
