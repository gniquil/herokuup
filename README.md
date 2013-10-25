# HerokuUp

This is a really simple program used to ping heroku to keep it alive. If any of the specified url does not return 200 status code, it will send a consolidated email to those specified in the config file, listing all the failed urls. 

```json
{
  "urls": [
    "http://example.herokuapp.com",
    "http://example2.herokuapp.com"
  ],
  "from": "admin@example.com",
  "tos": [
    "admin@example.com"
  ],
  "sendonfail": false,
  "serveraddr": "localhost:25"
}
```
## Install

To install, you just run the following in command line:

```bash
go get github.com/gniquil/herokuup

cp ~/go/src/github.com/gniquil/herokuup/config.json.example herokuup.config.json

vi herokuup.config.json
```
Modify the config file accordingly. Note you should probably set `sendonfail` to true.

Now run it

```bash
herokuup herokuup.config.json
```

Everything should work as it is. Finally you can put it in your cron. Following is a simple cron entry, checking the list urls every 5 minutes

```bash
*/5 * * * * /home/{user}/go/bin/herokuup /home/{user}/herokuup.config.json > /home/{user}/herokuup.log 2>&1
```