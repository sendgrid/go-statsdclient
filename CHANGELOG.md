Changelog
=========

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
