# service-training

The training material for the service repo broken out into steps.

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

## Contributing

After making changes run the `scrips/file_lists.sh` script to update the File
Changes and Dependency lists in each section's README.
