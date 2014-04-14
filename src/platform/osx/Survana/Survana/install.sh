#!/bin/sh

#  install.sh
#  Survana
#
#  Created by Victor Petrov on 4/12/14.
#  Copyright (c) 2014 The Neuroinformatics Research Group at Harvard University. All rights reserved.

SERVICES_DIR="services"
MONGODB_DIR="${SERVICES_DIR}/mongodb"
MONGODB_DL_URL="http://fastdl.mongodb.org/osx/mongodb-osx-x86_64-2.6.0.tgz"
MONGODB_DL_PATH=$SERVICES_DIR/mongodb.tgz
MONGODB_EXTRACTED_DIR="$SERVICES_DIR/mongodb-osx-x86_64-2.6.0"
MONGODB_BINDIR="$MONGODB_DIR/bin"

if [ -d $MONGODB_BINDIR ]; then
    exit 0;
fi

#download mongodb archive
curl $MONGODB_DL_URL -o $MONGODB_DL_PATH || exit 1

#extract mongodb archive
tar -zxf $MONGODB_DL_PATH -C $SERVICES_DIR || exit 1

#rename folder to "mongodb"
mv $MONGODB_EXTRACTED_DIR $MONGODB_DIR || exit 1

#remove mongodb archive
rm -f $MONGODB_DL_PATH

#done
echo "OK"