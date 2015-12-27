wget http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz
gzip -d GeoLite2-City.mmdb.gz
mkdir db
mv GeoLite2-City.mmdb db/