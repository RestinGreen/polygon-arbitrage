#!/bin/bash


initdb -D /var/lib/postgresql/data -A md5 --username=postgres --pwfile=<(echo $POSTGRES_PASSWORD )


pg_ctl -D /var/lib/postgresql/data -w start

echo $POSTGRES_PASSWORD | psql -c "CREATE USER $PGSQL_USER WITH ENCRYPTED PASSWORD '$PGSQL_PASSWORD';" -U postgres -W
echo $POSTGRES_PASSWORD | psql -c "CREATE DATABASE $PGSQL_DBNAME WITH OWNER $PGSQL_USER;" -U postgres -W

echo $POSTGRES_PASSWORD | psql -c "GRANT ALL PRIVILEGES ON DATABASE $PGSQL_DBNAME TO $PGSQL_USER;" -U postgres -W

echo $PGSQL_PASSWORD | psql -d $PGSQL_DBNAME -f /init/sql/dexs.sql -U $PGSQL_USER -W
echo $PGSQL_PASSWORD | psql -d $PGSQL_DBNAME -f /init/sql/tokens.sql -U $PGSQL_USER -W
echo $PGSQL_PASSWORD | psql -d $PGSQL_DBNAME -f /init/sql/pairs.sql -U $PGSQL_USER -W


pg_ctl -V ON_ERROR_STOP=1 
pg_ctl -D "/var/lib/postgresql/data" -m fast -w stop

exec "$@"