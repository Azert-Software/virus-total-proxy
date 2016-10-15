# Virus Proxy #

This is a Virus Proxy that uses Virus Total to scan a file for the presence of viruses.

This was created for use with Hmail server but could be used elsewhere of course.

The proxy is written in GoLang and an executable is located in the dist folder allowing you to download an run with a minimum of effort.

### Use ###

This takes a single parameter which is the location of the file to be checked.

```
.\virus-proxy.exe -file %FILEPATH%
```
A response code is then returned, 0 if virus 1 if no virus.

Logs will be written to log.txt at the same location as the exe.

### Setup ###

* Download exe
* Open config.json and add your Virus-Total apikey
* That's it
