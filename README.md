# How to set up

1. Have a polygon node ipc
2. 
    - Create pgsql user "bot" with a password.
    - CREATE USER bot WITH ENCRYPTED PASSWORD '';
    - Create database: CREATE DATABASE arbitrage
    - log in to bot user with psql -U bot -d arbitrage -W.
     If failed, than find the "pg_hba.conf" file (find / -name pg_hba.conf) and add the line "local | all | bot | (leave empty) | md5.


# Generate bindings like this

abigen  --abi univ2pair.abi --pkg univ2pair --type UniV2Pair --out univ2pair.go
