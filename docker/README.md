**Docker installation**

The files in this directory are set up to be a quick installation of bulbistry within a VERY minimal, preconfigured docker container. 

They are not guaranteed to work yet, and are untested.

To build the bulbistry container:
   
    docker build -t bulbistry-dev:latest -t bulbistry-dev -f docker/Dockerfile .
