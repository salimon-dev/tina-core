#!/bin/bash

# Check if username and password arguments are provided
if [ $# -ne 2 ]; then
    echo "Usage: $0 <username> <password>"
    exit 1
fi

# Assign arguments to variables
username="$1"
password="$2"

# Generate a random UUID for id
random_id=$(uuidgen)

# Generate MD5 hash of the provided password
password_hash=$(echo -n "$password" | openssl dgst -md5 | awk '{print $2}')

# Generate a random secret key (e.g., 16-character alphanumeric string)
secretkey=$(openssl rand -hex 8)

cat << EOF 
run this query in nexus database

INSERT INTO users (id, username, password, invitation_id, credit, role, status, secret_key, registered_at, updated_at) VALUES ('$random_id','$username','$password_hash',null,9999,1,1,'$secretkey',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP);

EOF
