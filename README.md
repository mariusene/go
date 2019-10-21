# Caserola: Order lunch for me
## What will order?
One random appetizer, main and desert from your favorite restaurant.

Future versions will have the option to configure selection algorithm. eg. repete last, check my history, etc.

## What favorite restaurant?
You will have to configure your favorite one.
```
{
    "email":"marius.ene@uipath.com",
    "pwd":"*****",
    "restaurant":"rawdiacorp",
    "utcH":7,
    "utcT":30
}
```
Restaurants available for selection:

- rawdiacorp
- saladrevolution

New restaurants will be added in the future versions.

## When will order?
Every day, except weekends, at 10:30 if you did not order already

## How to run?
Get the **bin** and run the **lunch.exe** after you change the **config.json** with your data.

Don't close it. After ordering will sleep until tomorrow. (future version will run in system tray)

You should also add it to Startup.

***If you don't trust the .exe, download [Go](https://golang.org/dl/) and build the src.***
