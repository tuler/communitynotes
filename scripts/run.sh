#!/bin/sh

docker run \
    -v $(PWD)/input:/input:ro \
    -v $(PWD)/output:/output \
    communitynotes \
    -e /input/userEnrollment-00000.tsv \
    -n /input/notes-00000.tsv \
    -r /input/ratings \
    -s /input/noteStatusHistory-00000.tsv \
    -o /output \
