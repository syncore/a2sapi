wget http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz
gzip -d GeoLite2-City.mmdb.gz
mkdir -p ../bin/db/
mv GeoLite2-City.mmdb ../bin/db/