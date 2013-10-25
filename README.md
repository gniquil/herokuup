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

Note you should probably set `sendonfail` to true.