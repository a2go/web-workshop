# service-training

This project is the training material for the [`service`][service] repo broken
out into steps.

## Requirements

This project was designed against Go 1.12. It should work for 1.11 but 1.12 is
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
git clone https://github.com/ardanlabs/service-training.git
mkdir garagesale
cd garagesale
git init .
```

---

You must also use `go mod init` to set the import path for this project. Doing
this exactly as shown will allow you to copy and paste code without a need to
modify import paths.

```sh
go mod init github.com/ardanlabs/garagesale
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

## Maintenance of this project

**This section is only relevant to maintainers of this repository.**

Many of the commands in this section rely on a shell feature that converts
expressions like `**/*.go` to a list of all `.go` files under the current
directory, recursively. Test this by running `echo **/*.md`. If it shows a list
of all Markdown files then your shell supports it. If it shows the literal
string `**/*.md` then the feature is not enabled. Add `shopt -s globstar` to
your shell init file and try again.

### Applying changes

Whenever a change is made to [`service`][service] it needs to be backported to
this repository. Find the step folder where the affected features were
introduced, make the changes, then copy them forward through to the last step.

Before making changes I recommend using the [`md5c`](https://github.com/jcbwlkr/md5c)
tool to visualize the different steps of evolution a particular file has
through the project. For example run `md5c **/web/web.go` to see where the
`web.go` file is introduced and which steps modify it.

---

If you modify a file in an early step and multiple subsequent steps are
identical you can copy the file forward without manually applying the changes.
For example if you change `web.go` in step 13 and want to apply the same
changes to steps 14 through 23 you can run this

```sh
for f in {14..23}*/**/web.go; do cp 13*/**/web.go $f; done
```

---

When applying changes between files that **are not** identical you must
manually copy only the relevant lines. This may involve merging two different
lines together if a newer step already changed the affected line. I recommend
using a diffing program. I use `vimdiff` but you can use any diffing tool that
feels comfortable.

```sh
vimdiff {32..33}*/internal/user/user.go
```

---

After bringing changes all the way forward you should rerun the `md5c` and
compare it with the original output. The output will be different but in most
cases the pattern of colors should be the same.

### Regenerating Docs

Be sure to regenerate the File Lists of the READMEs after every change.
Reviewing these changes to the READMEs can help validate that you didn't miss a
step.

```sh
./scripts/file_lists.sh
```

### Adding New Folders

If needed, you can insert a whole new step folder in the middle of these
materials. For example if you wanted to add a "ready check" section after
folder `22-health` you would need to bump the folder numbers for steps 23+ then
copy folder 22 to become the new 23.

There is a script to assist in this so for the example case I would do this:

```sh
for dir in {23..36}*; do ./scripts/bump.sh $dir; done
cp -a 22-health 23-readycheck
```

But really please don't insert folders in the middle. If it is reasonable just
add a new section to the end by copying the newest folder.

```sh
cp -a 36-self-shutdown 37-new-thing
```

---

### Synchronizing Changes

Use `git dirdiff` against the [`service`][service] repo periodically to look
for major changes. To reduce noise you can change all of the imports in this
project to match the other. Be sure to have a clean `git status` before you do
the diff so you can quickly revert the changes.

```sh
sed -i "" -e "s/garagesale/service/" **/*.go
git dirdiff {service,.}/cmd
git checkout -- **/*.go
```

[service]: https://github.com/ardanlabs/service
