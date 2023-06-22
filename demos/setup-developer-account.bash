#!/bin/bash
if [ "$USER" = "" ]; then
	USER="admin-test"
fi
# Make sure we have a database in PostgreSQL 15 for our username
# and a database to match
if [ "$HOSTNAME" = "" ]; then
	echo 'Missing HOSTNANE in the environment, aborting'
	exit 1
fi
if [ "$USER" = "" ]; then
	echo 'Missing USER in the environment, aborting'
	exit 1
fi
cat <<EOT
Run the following command if you need to create a PostgreSQL
support user account and database needed for this demo.

    createuser --no-password --createdb --superuser \\
         --host "localhost" "${USER}"
    createdb "${USER}"

Depending on how your Postgres is setup it maybe easier to
create a "superuser" account via SQL run from the postgres
account via psql. If so

    sudo -u postgres psql -c "CREATE USER ${USER} SUPERUSER"
    createdb "${USER}"
    
EOT
