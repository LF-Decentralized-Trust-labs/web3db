#!/bin/bash

# Wait for hive-metastore to be ready
while ! nc -z hive-metastore 9083; do   
  echo "Waiting for hive-metastore to be ready..."
  sleep 2
done

# Now that hive-metastore is ready, run hiveserver2
echo "Starting hiveserver2..."
./hiveserver2
