#!/bin/sh

docker run \
    -v $(pwd)/input:/input:ro \
    -v $(pwd)/output:/output \
    communitynotes \
    -e /input/userEnrollment-00000.tsv \
    -n /input/notes-00000.tsv \
    -r /input/ratings \
    -s /input/noteStatusHistory-00000.tsv \
    -o /output
