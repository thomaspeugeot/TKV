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

First step is to create your environment on the Google Compute Platform (GCP). 

From there, either you work fully in the cloud (with a cloud shell) or you work locally on your computer (we do the second).

```
cd /Users/thomaspeugeot/goroot/src/go-app
```




