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
## 2018, october the 11th

Trying to reproduce gophers-1 for serving the 10000 web site (only the map).

The definition of the app.yaml is cryptic. 
```
# All URLs are handled by the Go application script
handlers:
- url: /.*
  script: _go_app
```
Why "/.*" to describe all files ? It is described in https://cloud.google.com/appengine/docs/standard/go/config/appref
This is for the description of ONE service. There, we will need only ONE app.yaml file.
The above definition defines one "Handler" instruction. It states that all requests should be served by the go application.
I make the hypothesis that this will be enough for us. 
```
- url: /tkv-client.html
  static_files: tkv-client.html
  upload: tkv-client.html

- url: /.*
  script: _go_app

```
à essayer

```
cd goroot/src/github.com/thomaspeugeot/tkv/gae_tkv
```

http://localhost:8080/tkv-client.html

does not work. At least, a configuration seems to work with dev_appserver.py
handlers:
```
- url: /tkv-client.html
  static_files: tkv-client.html
  upload: tkv-client.html

- url: /css
  static_dir: css
  
- url: /js
  static_dir: js
```
Let's try to see if works when the path is outside the root directory
```
handlers:
- url: /tkv-client.html
  static_files: tkv-client.html
  upload: tkv-client.html

- url: /css
  static_dir: ../tkv-client/css
  
- url: /js
  static_dir: ../tkv-client/js
```
it does work localy but does work with gcloud app deploy

## 2018, october the 13th

lets try if the path below root directory works with gcloud deploy. It does not. Let's check logs
```
gcloud --quiet app deploy
gcloud app logs tail -s default
```
the
```
--quiet
```
option is cool for not having to avoid to interactive confirmation.

It does not work anymore on the deployment configuration.
There is a trace
```
2018-10-13 06:16:06 default[20181013t081528]  "GET /tkv-client.html HTTP/1.1" 200
```
meaning the tkv-client.html is correctly loaded but all the rest is not.
or it does not GET the tkv-client.html at all
```
2018-10-13 06:18:28 default[20181013t081744]  "GET /js/leaflet.js HTTP/1.1" 404
2018-10-13 06:18:28 default[20181013t081744]  "GET /css/angular-material.css HTTP/1.1" 404
2018-10-13 06:18:28 default[20181013t081744]  "GET /js/angular.js HTTP/1.1" 404
2018-10-13 06:18:28 default[20181013t081744]  "GET /js/tkv-client.js HTTP/1.1" 404
2018-10-13 06:18:28 default[20181013t081744]  "GET /css/leaflet.css HTTP/1.1" 404
2018-10-13 06:18:28 default[20181013t081744]  "GET /js/angular-leaflet-directive.js HTTP/1.1" 404
2018-10-13 06:18:28 default[20181013t081744]  "GET /js/leaflet.js HTTP/1.1" 404
2018-10-13 06:18:28 default[20181013t081744]  "GET /js/angular-leaflet-directive.js HTTP/1.1" 404
2018-10-13 06:18:28 default[20181013t081744]  "GET /js/tkv-client.js HTTP/1.1" 404
```
It works localy but not on remote, Let's try gcloud --verbosity=info.

no info.
Let's try from the console with the thomas.peugeot@10kt.org account
https://console.cloud.google.com/logs/viewer?project=tenktorg&authuser=1&organizationId=174259221484&resource=gae_app%2Fmodule_id%2Fdefault&minLogLevel=0&expandAll=false&timestamp=2018-10-13T06:36:14.521000000Z&customFacets=&limitCustomFacetWidth=true&dateRangeStart=2018-10-13T05:36:14.775Z&dateRangeEnd=2018-10-13T06:36:14.775Z&interval=PT1H&logName=projects%2Ftenktorg%2Flogs%2Fappengine.googleapis.com%252Frequest_log&scrollTimestamp=2018-10-13T06:28:03.555834000Z

Now, we have confirmation that the directory path OUTSIDE the root directory does work locally but NOT with the deployment.

```
ERROR: (gcloud.app.deploy) INVALID_ARGUMENT: Your app may not have more than 15 versions. Please delete one of the existing versions before trying to create a new version.
```
https://console.cloud.google.com/appengine/versions?authuser=1&organizationId=174259221484&project=tenktorg&serviceId=default&versionssize=50
to remove versions.
**Let's try to work with the services**
From now, we have only dealt with serving the html/css/js files. Let's try to have handlers perform services tasks.

## 2018, october the 14th
Today, we try to have the service working.
First find, no need to have the http server running in main.go
```
func main() {
	appengine.Main()
}
```
is all you need to serve your file
If you want to add a service, it seems quite simple.
```
func main() {
	http.HandleFunc("/checkEnv", checkEnv)
	appengine.Main()
}
func checkEnv(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "IsDevAppServer: %v", appengine.IsDevAppServer())
}
```
If I want to do something more complex, 
```
func main() {
	http.HandleFunc("/translateLatLngInSourceCountryToLatLngInTargetCountry",
		handler.TranslateLatLngInSourceCountryToLatLngInTargetCountry)
	http.HandleFunc("/checkEnv", checkEnv)
	appengine.Main()
}
func checkEnv(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "IsDevAppServer: %v", appengine.IsDevAppServer())
}
```
The server answer to the request http://localhost:8080/translateLatLngInSourceCountryToLatLngInTargetCountry
is
```
...
2018/10/14 09:23:37 http: panic serving 127.0.0.1:59274: runtime error: invalid memory address or nil pointer dereference
...
github.com/thomaspeugeot/tkv/translation.(*Country).ClosestBodyInOriginalPosition(0x171ae80, 0x0, 0x0, 0x11e05ee, 0xc420086280, 0x2, 0xc4200722a0, 0x54, 0x16d8e60, 0xc420192a20, ...)
	/Users/thomaspeugeot/goroot/src/github.com/thomaspeugeot/tkv/translation/country.go:193 +0x3e
github.com/thomaspeugeot/tkv/translation.(*Translation).ClosestBodyInOriginalPosition(0x171ae60, 0x0, 0x0, 0xc420040b28, 0x2, 0x2, 0x14b5fe0, 0xc420040a48, 0xc42008ea40, 0xc420056600, ...)
	/Users/thomaspeugeot/goroot/src/github.com/thomaspeugeot/tkv/translation/translation.go:36 +0x4d
github.com/thomaspeugeot/tkv/handler.TranslateLatLngInSourceCountryToLatLngInTargetCountry(0x16e1460, 0xc4201360e0, 0xc42011a400)
...

``` 
The crash is normal since the translation has not been inited. Let's do this init in the translation file through a singloton pattern.
IT WORKS .... on the local development server

On the gcloud server https://tenktorg.appspot.com/translateLatLngInSourceCountryToLatLngInTargetCountry, we now have an issue, not a surprise since 
# - the client javascript asks for localhost to get service 
- the memory footprint is above the standard free quota

The log on the application is pretty explicit.
```
2018-10-14 12:19:27.205 CEST
Exceeded soft memory limit of 128 MB with 205 MB after servicing 0 requests total. Consider setting a larger instance class in app.yaml.
2018-10-14 12:19:27.205 CEST
This request caused a new process to be started for your application, and thus caused your application code to be loaded for the first time. This request may thus take longer and use more CPU than a typical request for your application.
2018-10-14 12:19:27.205 CEST
While handling this request, the process that handled this request was found to be using too much memory and was terminated. This is likely to cause a new process to be used for the next request to your application. If you see this message frequently, you may have a memory leak in your application or may be using an instance with insufficient memory. Consider setting a larger instance class in app.yaml.
```
**upgrading the google application engine class**

apparently, the application consummes more than 128 MB (377), in top
```
9790  _go_app      0.0   00:19.89 11    0    35    219M   0B     377M   55433 55436 sleeping *0[1]          0.00000 0.00000    501  197249    319
5
```
in the application monitor on osx, it is rather above 596 MB
```
 _go_app	0.0	20.29	11	0	59790	thomaspeugeot	219.5 MB	596.9 MB	0 bytes		0 bytes	0 bytes	64 bit	0 bytes	0 bytes	0	0		-	No	No	No	0 bytes	0 bytes	No	No		0 bytes	No	
 ```
 According to the go documentation, https://cloud.google.com/appengine/docs/standard/#instance_classes, 
 we need to be in the F4_1G class
 ```
 Instance Class 	Memory Limit 	CPU Limit 	Supported Scaling Types
F1 (default) 	128 MB 	600 MHz 	automatic
F2 	256 MB 	1.2 GHz 	automatic
F4 	512 MB 	2.4 GHz 	automatic
F4_1G 	1024 MB 	2.4 GHz 	automatic
....
```
This is set in the app.yaml file.

Activation of the 300$ free trial.

https://tenktorg.appspot.com/translateLatLngInSourceCountryToLatLngInTargetCountry
is responding correctly. VERY GOOD NEWS.

Modification of the hostname in tkv-client.js
```
var hostname = "https://tenktorg.appspot.com/"
```
!!! one needs to have the hostname be automaticaly computed --> window.location.hostname
```
		hostname = window.location.hostname
		protocol = window.location.protocol
		port = window.location.port
		targetService = protocol + "//"+ hostname + ":" + port + "/"
        ...
		$http.post( targetService +'translateLatLngInSourceCountryToLatLngInTargetCountry', jsonLatLng )...
```
that works

PUTAIN CA MARCHE !!!!!

Except it works from the firefox on my mac. Not from safari on the ipad or the iphone.

