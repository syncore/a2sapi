# a2sapi

a2sapi is a RESTful API for receiving [**master server information**](https://developer.valvesoftware.com/wiki/Master_Server_Query_Protocol) and for querying [**A2S information**](https://developer.valvesoftware.com/wiki/Server_queries) from servers running on the Steam (Source) platform.

This back end service was written to provide information to a number of sites (for example, [here](http://reflex.syncore.org) and [here](http://ql.syncore.org) for which I needed this specific information.

*Please note, this is the first project that I have written in the Go programming language.* :scream: Pull requests are welcome!


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

### Configuration (binaries and source)
The configuration is handled interactively by passing the `--config` flag to the a2sapi executable. The configuration file will be stored in the `conf` directory. Any existing configuration will be overwritten.


### Launching: Binaries
  - Linux/OSX: Launch with: `./a2sapi`
  - Windows: Launch by running the `a2sapi.exe` executable.
  - You can pass the `--h` flag to the executable to see a few command-line options.


### Build from Source

- Alternatively, you can build from source. This assumes that you have a working Go environment. If not, check out the [Golang Getting Started guide](https://golang.org/doc/install).
- Extract the archive.
- Change directory to `build/nix` if you're on Linux/OSX or `build\win` if you're on Windows and launch the appropriate `build.sh` or `build.bat` script.
- Change back to the root directory, then change directory to `getfiles` and run the appropriate `get_countrydb` script to get the geolocation database file, which is the GeoLite2 City free database [provided by MaxMind](http://dev.maxmind.com/geoip/geoip2/geolite2/).
  - Note: if you're on Windows you'll need `wget` and `gzip`. For more info, see the discussion above for the binary installation.
  - Updates for this geolocation database are provided by MaxMind on the first Tuesday of every month, so you can run the script again at that time to get the updates.

### Launching: Source
- After building, the resulting executable will be located in the `bin` directory.
- On first run, you will need to generate the configuration file by passing the `--config` flag to the executable.
- After generating the configuration file, launch the application by running the `a2sapi` executable in the `bin` directory. If you'd like to see a few command-line options, then pass the `--h` flag to the executable.




# Usage
:book: For interactive documentation and more detail, see the a2sapi Swagger UI documentation in use [on one of my pages that uses this API](http://reflex.syncore.org/apidoc/) or you can use the included a2sapi-swagger files with Swagger UI/Editor.

The API ships with three endpoints:
- /servers
- /serverIDs
- /query


### `GET: /servers`
The `servers` endpoint provides a list of the most recent servers returned from the Valve master server. Data from this endpoint is only available if the application has been configured to retrieve servers from Valve's master server. This list can be filtered by specifying one or more of the filter parameters below. Separate multiple parameter values with commas. Multiple filters can be combined after the first filter by using the & character before any additional filters, for example: `/servers?countries=US,SE&maps=overkill&hasPlayers=true&serverOS=Linux`

### String parameters (filters):
- ***countries***
  - Filter by 2-letter ISO 3166-1 country code.
  - `/servers?countries=US,SE,NL`
- ***regions***
  - Filter by region. Possible regions: `Africa, Antarctica, Asia, Europe, Oceania, North America, South America`
  - `/servers?regions=North America,Oceania`
- ***states***
  - Filter by 2-letter US state. United States of America only.
  - `/servers?states=NY,TX`
- ***serverNames***
  - Filter by server name. Results are loosely matched.
  - `/servers?serverNames=Newbies,practice,fun server`
- ***maps***
  - Filter by map. Results are loosely matched.
  - `/servers?maps=bdm3,cpm22,dp6`
- ***games***
  - Filter by game.
  - `/servers?games=Reflex`
- ***gametypes***
  - Filter by gametype.
  - `/servers?gametypes=CA,CTF`
- ***serverTypes***
  - Filter by server types. Possible types: `dedicated, listen`
  - `/servers?serverTypes=dedicated`
- ***serverOS***
  - Filter by server operating system. Possible types: `Linux, Windows, Mac`
  - `/servers?serverOS=Linux`
- ***serverVersions***
  - Filter by server version.
  - `/servers?serverVersions=1.33,1.66,2.02`
- ***serverKeywords***
  - Filter by server keywords. Results are loosely matched.
  - `/servers?serverKeywords=minqlx,clanarena,stats`

### Boolean parameters (filters):
- ***hasPlayers***
  - Filter by whether server has players (true) or is empty (false).
  - `/servers?hasPlayers=true`
- ***hasBots***
  - Filter by whether server has bots (true) or not (false).
  - `/servers?hasBots=false`
- ***hasPassword***
  - Filter by whether server has a password (true) or not (false).
  - `/servers?hasPassword=false`
- ***hasAntiCheat***
  - Filter by whether server is secured by anti-cheat (true) or not (false).
  - `/servers?hasAntiCheat=true`
- ***isNotFull***
  - Filter by whether server is full (true) or not (false).
  - `/servers?isNotFull=true`

### `GET: /serverIDs`
The `serverIDs` endpoint retrieves servers' internal ID numbers. The ID number(s) will be used with the `ids` parameter of the `query` endpoint to retrieve a server's real-time information. Separate multiple parameter values with commas.

### Parameters:
- ***hosts***
  - The host in the format of IP:port to retrieve the ID for. Multiple IP:ports can separated with commas.
  - `/serverIDs?hosts=54.93.46.254:25801,46.101.8.188:27960`


### `GET: /query`
The `query` endpoint retrieves servers' real-time information. Depending on how the application is configured, this can be done via server ID numbers (retrieved via the `serverIDs` endpoint) and/or directly from IP addresses and ports. Separate multiple parameter values with commas.

### Parameters for querying by server ID:
- ***ids***
  - The server ID(s) whose information should be retrieved.
  - `/query?ids=123,456,999,10340`

### Parameters for directly querying by address:
- ***hosts***
  - The host in the format of IP:port whose information should be retrieved. :warning: Note, address queries might be disabled, depending on the application configuration.
  - `/query?hosts=54.93.46.254:25801,46.101.8.188:27960`


# Quick Examples
**`/servers` endpoint:**

These examples for the `/servers` endpoint will assume that a2api is configured to retrieve Quake Live servers from Valve's master server at timed intervals (that is, data is available for the `/servers` endpoint):

*Get all clan arena servers in North America that have players, do not have bots, and are running the minqlx server extension:*

- `http://some-webserver.com/servers?gametypes=CA&regions=North America&hasPlayers=true&hasBots=false&serverKeywords=minqlx`

*Get all servers in Germany that contain the word "fun" in their name that are running the map overkill or thunderstruck:*

- `http://some-webserver.com/servers?countries=DE&serverNames=fun&maps=overkill,thunderstruck`

**`/serverIDs` endpoint:**

*Get the ID for the server with address: 127.0.0.1:27960*
- `http://some-webserver.com/serverIDs?hosts=127.0.0.1:27960`

*Get the IDs for the servers with addresses: 127.0.0.1:27960, 10.0.0.1:27597, and 172.16.0.1:27015*
- `http://some-webserver.com/serverIDs?hosts=127.0.0.1:27960,10.0.0.1:27597,172.16.0.1:27015`

**`/query` endpoint:**

*Get the real-time information for the servers with IDs 100, 200, 300, and 400*
- `http://some-webserver.com/query?ids=100,200,300,400`

*Get the real-time information for the servers with addresses 127.0.0.1:27960, 10.0.0.1:27597, and 172.16.0.1:27015*
- `http://some-webserver.com/query?hosts=127.0.0.1:27960,10.0.0.1:27597,172.16.0.1:27015`
- :warning: The API administrator may have direct server address queries disabled, in which case this would not work!





# Issues

The preferable method of contact would be for you to open up an [issue on Github](https://github.com/syncore/a2sapi/issues).
Alternatively, I usually can be found under the name "syncore" on QuakeNet IRC - irc.quakenet.org


# License
----
See [LICENSE.md]

[LICENSE.md]:https://github.com/syncore/a2sapi/blob/master/LICENSE.md

