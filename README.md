# A2Go web-workshop
A series of progressive workshops to build webservices

We are going to shift gears and adapt the excellent Ardan Labs service training to make a TodoMVC application.
Currently, the master branch of this repository is in transition to this goal.
### Our story so far

We have an app that says hello world! It.starts up, and gracefully shuts down.

### Client

Go to <a href="./client">The client directory</a> for instructions on how to run the client!

### Get started with Go

The quickest way to get started with Go, is install [VS Code](https://github.com/microsoft/vscode) (Jet Brains GoLand is nice, but not free). If you go to the [vscode-config](./vscode-config) directory, there's useful tips and settings for you to get started with.

## Why

A2 Go, the developer meetup, is a forum for people working in or around Ann Arbor with Go to discuss ideas, issues and share solutions. One of the topics of interest that we heard from our community was "Why do I hear so much about web services and Go?". These workshops are an attempt to teach, from ground zero, the how and why of web services in Go. In this case we are interpreting web services fairly broadly and we think that is good. It gives plenty of fun and interesting rabbit holes to fall into while learning.

## How to use this web-workshop

While each workshop is written with a user group meeting in mind, it could be used in another format. The git repository is structured such that you change branches for each workshop. The zero branch is a very basic hello world web application. Each branch has its own README which walks the attendee through steps to build whatever is being built in that workshop. Subsequent branches pick up where the previous workshop left off. Our goal is for someone with little to no Go experience to be able to start this workshop at step zero.

To get started at the beginning, clone this repo and change to the [workshop_0](https://github.com/a2go/web-workshop/tree/workshop_0) branch.
To start with a basic webserver, clone this repo and change to the [workshop_1](https://github.com/a2go/web-workshop/tree/workshop_1) branch.
# service-training

This project is the training material for the [`service`][service] repo broken
out into steps.

## Requirements

This project was designed against Go 1.13. It should work for 1.12 but 1.13 is
recommended.

Supporting services like the database are hosted in Docker. If you cannot
install Docker on your machine you can still follow most of this material by
hosting a database elsewhere and modifying the connection information to your
needs.

## Setup

Clone this repository somewhere on your computer. The location does not
especially matter but if it is outside of your `$GOPATH` then the Go modules
features will work automatically.

In a separate folder make a directory where you will be building your API. We
recommend you initialize that folder as a Git repository to track your work.


```sh
mkdir ~/training
cd ~/training
git clone https://github.com/a2go/web-workshop.git
mkdir garagesale
cd garagesale
git init .
```

---

You must also use `go mod init` to set the import path for this project. Doing
this exactly as shown will allow you to copy and paste code without a need to
modify import paths.

```sh
go mod init github.com/a2go/garagesale
git add go.mod
git commit -m "Initial commit"
```

## Postman API Client

For the class we will be building up a REST API. You may use any HTTP client
you prefer to make requests but we recommend [Postman](https://www.getpostman.com/).
For convenience you may use the import button in the top left to import the
included `postman_environment.json` and `postman_collection.json` files to get
a client up and running quickly. Be sure to select the "Garage Sale Service"
environment in the top right.

## Diffing Folders

Reviewing the differences between the successive steps helps to reinforce the
ideas each change is about. This is made easier by running the following
command to define a git alias called `dirdiff`:

```sh
git config --global alias.dirdiff 'diff -p --stat -w --no-index'
```

With that alias in place, run this command from the top level folder to see the
differences between the `01-startup` directory and the `02-shutdown` directory.

```sh
git dirdiff 01-startup 02-shutdown`
```

---
---
---