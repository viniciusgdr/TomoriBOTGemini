#!/bin/bash

echo "Go - Installing Dependencies..."
go get

echo "Go - Building..."
go build

echo "NodeJS - Installing Dependencies..."
cd ./nodejs
npm install
cd ..

echo "Python - Installing Dependencies..."
cd ./python
if ! [ -x "$(command -v pip)" ]; then
  echo 'Error: pip is not installed.' >&2
  echo "Installing pip..."
  # check distro and install pip
  if [ -x "$(command -v apt-get)" ]; then
    sudo apt-get install python3-pip
  elif [ -x "$(command -v yum)" ]; then
    sudo yum install python3-pip
  elif [ -x "$(command -v pacman)" ]; then
    sudo pacman -S python-pip
  else
    echo "Error: Unsupported distro"
  fi
  sudo pip3 install --upgrade pip
else
  echo "pip is already installed"
fi


pip install -r requirements.txt --break-system-packages
cd ..

echo "Checking if ffmpeg is installed..."
if ! [ -x "$(command -v ffmpeg)" ]; then
  echo 'Error: ffmpeg is not installed.' >&2
  echo "Installing ffmpeg..."
  # check distro and install ffmpeg
  if [ -x "$(command -v apt-get)" ]; then
    sudo apt-get install ffmpeg
  elif [ -x "$(command -v yum)" ]; then
    sudo yum install ffmpeg
  elif [ -x "$(command -v pacman)" ]; then
    sudo pacman -S ffmpeg
  else
    echo "Error: Unsupported distro"
  fi
else
  echo "ffmpeg is already installed"
fi

echo "Done!
 
 The  installer.sh  script is a simple bash script that installs the dependencies for the Go and NodeJS applications. It then builds the Go application and installs the NodeJS dependencies. 
 The  Dockerfile  is used to build the Docker image for the application. 
 # Path: Dockerfile"