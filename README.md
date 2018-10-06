10 0000
=======

Implementation of the "10 000" concept (see https://10ktblog.wordpress.com/a-propos/ for a description of the concept)

10 000 is a web server

**Running the web server**

```
cd runtime_server
go run runtime_server.go -sourceCountry fra -sourceCountryNbBodiesPtr 697529 -sourceCountryStep 4723 -targetCountry hti -targetCountryNbBodiesPtr 927787 -targetCountryStep 8564
```

a vscode configuration is availble to run and debug the server.

**Running the web client**

launch your browser at http://localhost:8001/tkv-client.html

On the top panel, zoom to your place of interest (it is currently limited to france). Left click. You terrritory appears as well as the matching territory in Haiti.
