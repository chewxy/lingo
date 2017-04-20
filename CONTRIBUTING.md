# Contributing #

Contributors are welcome! We want to make contributing as easy as possible, and the process is very Github-centric. [Github Issues](https://github.com/chewxy/lingo/issues) are used to manage any contributions and changes. If you don't have a github account, please feel free to email me (my  user name [at] gmail.com), and I'll gladly open an issue on your behalf.

# Process #

Say you have a change you want to make, this is the process:

1. Open an issue.
2. I'll have a brief discussion with you. If you don't feel comfortable with a public discussion, I'm okay to email. 
3. Fork this project on Github, and clone it to your local machine.
4. Make your changes
5. Make sure you have tests. If you foresee breaking any API, it is vital that it be discussed beforehand.
6. Make sure your tests pass.
7. `gofmt` your code
8. Send a Pull Request.

Say you instead saw one of the [many issues](https://github.com/chewxy/lingo/issues) and want to solve one of them. This is the process:

1. Comment on the issue saying you'll pick it up. (Alternatively, email me)
2. Fork the project on Github, clone to your local drive.
3. Fork this project on Github, and clone it to your local machine.
4. Make your changes
5. Make sure you have tests. If you foresee breaking any API, it is vital that it be discussed beforehand.
6. Make sure your tests pass.
7. `gofmt` your code
8. Send a Pull Request.

## Pull Requests ##

I'll review every pull request. I may request some changes, or delve into further discussions. After that, once I'm satisfied everything passes, I'll merge the pull request. Then I'll add your name into the CONTRIBUTORS list.

# Debugging #

This package comes with a debug tag option. Most subpackages will have a `debug.go` which contain a `logf` function for logging any traces you wish to trace. 