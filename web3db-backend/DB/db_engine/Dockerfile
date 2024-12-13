# Step 1: AdoptOpenJDK for Java 11
FROM adoptopenjdk:11-jre-hotspot as java_base

# Step 2: Start from Python 3.8 base image
FROM python:3.8-slim

# Set the working directory in the container
WORKDIR /app

# Copy Java Runtime from the Java base image
COPY --from=java_base /opt/java/openjdk /opt/java/openjdk

# Set JAVA_HOME environment variable
ENV JAVA_HOME=/opt/java/openjdk
ENV PATH="$JAVA_HOME/bin:$PATH"

# Install necessary build dependencies for compiling Python packages
RUN apt-get update -y \
    && apt-get install -y gcc g++ libc-dev libffi-dev libsasl2-dev libssl-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy your application code to the container
COPY . /app

# Install Python dependencies
RUN pip install --trusted-host pypi.python.org \
    Flask==2.0.1 \
    flask-restx==0.5.1 \
    Werkzeug==2.0.1 \
    requests \
    pyspark==3.0.0 \
    pyhive \
    thrift \
    thrift_sasl \
    --upgrade

# Expose port 5000 for the application
EXPOSE 3001

# Command to run the application
CMD ["python", "db_engine/app.py"]
