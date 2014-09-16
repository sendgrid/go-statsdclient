Changelog
=========
# 3.0.0
- using a map to track stats that have been logged. This way, we can also tell how many times a given stat has been logged

# 2.1.0
- SetPrefix ensures there is only a single "." delimeter at the end. Will remove extraneous ones if present and add one if not present.

# 2.0.0
- Dial returns an interface, not a concrete type
  (https://github.com/sendgrid/go-statsdclient/pull/5)
- Added MockClient type for testing
  (https://github.com/sendgrid/go-statsdclient/pull/5)
- Removed MakePrefix method since Ops says the entire prefix should be
  passed in via config, not constructed by the process
  (https://github.com/sendgrid/go-statsdclient/pull/7)
  

# 1.1.1
- Add benchmarks to tests
  (https://github.com/sendgrid/go-statsdclient/pull/3)

# 1.1.0
- Add SetPrefix and MakePrefix functions
