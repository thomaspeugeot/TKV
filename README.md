# TKV
Ten Kilo Villages

todo:

DEBUG OF DYNAMIC
- allow for pin-pointing the mid point between the closest bodies
- allow for understanding the dynamic of the min distance

PRODUCTION FOR ALL COUNTRIES

INTEGRATION OF EN USER CLIENT

1/ provide a web service getVillages() which takes an area an input (top left and bottom right coords) and returns
the coordinates of the barycenters of the village in the area  

- optional paramters (cutoff for the number of villages barycenters to get)

- create a struct "villageLayer" which is iinitialized with the data of all villages 
	
	init function
	- for a set of countries
		- load initial ".bods"
		- load final ".bods"
	- from the final ".bods", compute the village barycenters

	getVillages functions
	- given the area, parse all villages and pick up the ones inside the area
