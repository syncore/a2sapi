# a2sapi

a2sapi is a REST API for receiving [**master server information**](https://developer.valvesoftware.com/wiki/Master_Server_Query_Protocol) and for querying [**A2S information**](https://developer.valvesoftware.com/wiki/Server_queries) from servers running on the Steam (Source) platform.

This back end service was written to provide information to a number of sites (for example, [here](http://reflex.syncore.org) and [here](http://ql.syncore.org) for which I needed this specific information.

*Please note, this is the first project that I have written in the Go programming language.* Pull requests are welcome!


----------

# Installation

### Installation: Binaries
- Grab one of the binaries for your platform from [releases](https://github.com/syncore/a2sapi/releases).
  - Extract the archive.
  - Change directory to `getfiles` and run the appropriate `get_countrydb` script to grab the geolocation database.
  	- This is the GeoLite2 City free database [provided by MaxMind](http://dev.maxmind.com/geoip/geoip2/geolite2/).
  	- MaxMind updates this database on the first Tuesday of every month, so you can run the script again at that time, if you'd like.
  	- If you are on Windows, you will need [wget](http://nebm.ist.utl.pt/~glopes/wget/) and [gzip](http://www.gzip.org) to use the `get_countrydb` script. Or, alternatively, simply download the database [here](http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz), extract the `GeoLite2-City.mmdb` file into a directory called `db` in the same location as the `a2sapi` executable.
  - Change back to the directory where the `a2sapi` executable is located.
  - Linux/OSX users: you must generate the configuration file with `./a2sapi --config`
  - Windows users: you must generate the configuration file with `a2sapi.exe --config`
  - The configuration file will be stored in the `conf` directory.

  ### Launching
  - Linux/OSX: Launch with: `./a2sapi`
  - Windows: Launch by running the `a2sapi.exe` executable.
  - You can pass the `--h` flag to the executable to see a few command-line options.


### Installation: From Source

- Alternatively, you can build from source. This assumes that you have a working Go environment. If not, check out the [Golang Getting Started guide](https://golang.org/doc/install).
- Extract the archive.
- Change directory to `build/nix` if you're on Linux/OSX or `build\win` if you're on Windows and launch the appropriate `build.sh` or `build.bat` script.
- Change back to the root directory, then change directory to `getfiles` and run the appropriate `get_countrydb` script to get the geolocation database file, which is the GeoLite2 City free database [provided by MaxMind](http://dev.maxmind.com/geoip/geoip2/geolite2/).
	- Note: if you're on Windows you'll need `wget` and `gzip`. For more info, see the discussion above for the binary installation.
	- Updates for this geolocation database are provided by MaxMind on the first Tuesday of every month, so you can run the script again at that time to get the updates.

### Launching
- The resulting executable will be located in the `bin` directory.
- On first run, you will need to generate the configuration file by passing the `--config` flag to the executable.
  - The configuration file will be stored in the `conf` directory.
- After generating the configuration file, launch the application by running the `a2sapi` executable in the `bin` directory.
	- If you'd like to see a few command-line options, then pass the `--h` flag to the executable.



### Usage (TODO)


### Issues

The preferable method of contact would be for you to open up an [issue on Github](https://github.com/syncore/a2sapi/issues).
Alternatively, I usually can be found under the name "syncore" on QuakeNet IRC - irc.quakenet.org


License
----
See [LICENSE.md]

[LICENSE.md]:https://github.com/syncore/a2sapi/blob/master/LICENSE.md

