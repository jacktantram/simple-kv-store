# Cryptowatch Golang Hiring Test: K/V REPL with nested transactions

Thank you for applying to work on Cryptowatch at Kraken!

We'd like the opportunity to see how you write code and solve problems, so we created this little challenge.

We don't limit the time you spend on the solution, you're free to spend as much as you want. Quality-wise, consider it a real task (albeit a little toyish) given to you by your manager.

## The Goal

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

## Submitting Your Solution

This is a git repo. Please keep a commit history as you work on your solution; it's helpful to see.

You should work locally. If you want to push to a remote to back up your work, please make sure to use a _private_ repo (not a public one).

When you're ready to submit your solution, run `prepare.sh` and email the resulting `.bundle.gz` file back to your Kraken recruiter.

```
bash prepare.sh
```

After running the script, please ensure that the resulting archive isn't empty or corrupted before sending the file to your recruiter.

You can check that by running the following:

```bash
gunzip backend-go-Name-LastName.bundle.gz # or gzip -d
git clone backend-go-Name-LastName.bundle myarchive
cd myarchive
```

You should be able to see the folder populated with your submission. If not, repeat the process.


# Usage
In order to use the CLI Tool there is a pre-requisite of [Go](https://github.com/golang/go) being installed.
Execute `make run` to run the program and from there you are able to send in commands. 
# Assumptions
* Keys
  * They are case-sensitive. Meaning that you can store a key of anApple and ANAPPLE. 
  * Keys and values must not be blank, i.e just be whitespace

# Improvements
* Provide further testcases for the integration/CLI part. Due to time I focussed on writing
unit tests on the store and then covering the base scenario from the README. I chose to use a 
`testdata` folder to allow a text file for `stdin` to be passed in and assert the `stdout` response.
* Depending on how this tool could be used Docker could be introduced to package the application