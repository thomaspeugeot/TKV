Working note on deploying 10000 on google application engine
======================================================

## 2018, october the 9th

**choosing standard or flexible environment**
According to https://cloud.google.com/appengine/docs, there are two kind of environment :
- standard
- flexible

Standard provides free environment as soon as you don't overtake some limits (1GB storage). Go version must be lower or equal to 1.9.

Flexible seems a very different approach. It is based on google compute engine (docker).

https://cloud.google.com/appengine/docs/the-appengine-environments details the difference. For standard, it states " Intended to run for free or at very low cost, where you pay only for what you need and when you need it. For example, your application can scale to 0 instances when there is no traffic. "

Let's go for standard. First thing, we have to downgrade to 1.9

**let's learn standard environment**

https://cloud.google.com/appengine/docs/standard/go/building-app/

https://cloud.google.com/appengine/docs/standard/go/building-app/creating-your-application

Beware, this is the tutorial, not the quick start (https://cloud.google.com/appengine/docs/standard/go/quickstart)

First step is to create your environment on the Google Compute Platform (GCP). 

From there, either you work fully in the cloud (with a cloud shell) or you work locally on your computer (we do the second).

```
cd /Users/thomaspeugeot/goroot/src/go-app
dev_appserver.py app.yaml
```
An error occured
```
    with open(configuration_path) as f:
IOError: [Errno 2] No such file or directory: 'app.yaml'
```
quite strange since there is such a file
## 2018, october the 10th
Let's try to see  if this is the go v1.11 that causes the problem. 
40' to go back to go 1.9.7 (!)
Problem still present
upgrading gcloud to version 220.0.0
it seems to be a known problem
https://stackoverflow.com/questions/52653776/why-is-dev-appserver-py-reporting-no-such-file-or-directory-after-todays-gclo
It is an known issue and it has been solved with 220.0.0

```
INFO     2018-10-10 17:55:55,019 devappserver2.py:278] Skipping SDK update check.
WARNING  2018-10-10 17:55:56,209 simple_search_stub.py:1196] Could not read search indexes from /var/folders/9z/9jhjxqj95k5dvpr_tbqjtvpr0000gn/T/appengine.None.thomaspeugeot/search_indexes
INFO     2018-10-10 17:55:56,252 api_server.py:275] Starting API server at: http://localhost:49577
INFO     2018-10-10 17:55:56,637 dispatcher.py:270] Starting module "default" running at: http://localhost:8080
INFO     2018-10-10 17:55:56,693 admin_server.py:152] Starting admin server at: http://localhost:8000
/Users/thomaspeugeot/dev/google-cloud-sdk/platform/google_appengine/google/appengine/tools/devappserver2/mtime_file_watcher.py:182: UserWarning: There are too many files in your application for changes in all of them to be monitored. You may have to restart the development server to see some changes to your files.
  'There are too many files in your application for '
INFO     2018-10-10 17:56:42,139 instance.py:294] Instance PID: 12581
INFO     2018-10-10 18:09:56,162 module.py:880] default: "GET / HTTP/1.1" 200 23
INFO     2018-10-10 18:09:56,321 module.py:880] default: "GET /favicon.ico HTTP/1.1" 302 24
INFO     2018-10-10 18:09:56,413 module.py:880] default: "GET /favicon.ico HTTP/1.1" 302 24
INFO     2018-10-10 18:09:56,420 module.py:880] default: "GET / HTTP/1.1" 200 23
```

**testing the deployment**
```
gcloud app deploy
```

https://tenktorg.appspot.com

it works fine.

**serving static files**
```
cd $GOPATH/src/github.com/GoogleCloudPlatform/golang-samples/appengine/gophers/gophers-1

```


