
# About the project
This is a task manager application that I made as my first project in Golang, and an attempt at utilising three tier architecture (Backend, Middleware, Interface)


# How to run (v1.0)
on the current version (1.0) you need to do some set-up to get it working properly. The three-tier architecture approach means this is more suited to run on servers (one for the Database, one for the API, and one for the webpage).
So getting it to run on a local machine can be a teeny bit tough if you have no experience with programming or APIs. (This guide assumes you've already installed [Go](https://go.dev)



## Start the API server
first you have to CD into the api/ directory then run the command `go run .`

## Create a Database
Use an API testing software such as Postman to send a post request to `http://0.0.0.0:8080` which creates a new database.

## Start Webpage server
now cd into the project directory (taskManager) and run the command `go run .`

*congratulations! you've now got my task manager project running!*

You don't need to create a new database when you restart either servers.
