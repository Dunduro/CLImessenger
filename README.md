# CLImessenger

flags

| Flag | excepted inputs | description |
| --- | ---| --- |
| `--v` | | toggles on verbose mode for debugging !not secure for production use! |
| `--mode` | client, server, test | default operation mode is set to client |
| | |client: runs the application as a client connecting to a server |
| | |server: runs the application as host, |
| | |test: preformances functionality purely for development purposes |
| `--user` | username | user handle to be used to reference the user with in chat. this functionality is used purely to directly login as a specific user from commandline without being prompted for it

commands can be found by calling the `/help` command running in client mode