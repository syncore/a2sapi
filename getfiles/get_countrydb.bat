:: wget for windows can be downloaded at http://nebm.ist.utl.pt/~glopes/wget/
:: gizp for windows can be downloaded at http://www.gzip.org
@echo off
wget http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz
gzip -d GeoLite2-City.mmdb.gz
mkdir %cd%\..\..\bin\db
move /Y GeoLite2-City.mmdb %cd%\..\..\bin\db\