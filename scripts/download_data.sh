#!/bin/bash

DATE=2025/01/06
BASE_URL=https://ton.twimg.com/birdwatch-public-data

mkdir -p input/ratings

# download notes file
file=${BASE_URL}/${DATE}/notes/notes-00000.tsv
echo "Downloading ${file}"
curl -L --compressed ${file} -o input/notes-00000.tsv

# download 16 ratings files
for i in {0..15}; do
    num=$(printf "%05d" $i)
    file=${BASE_URL}/${DATE}/noteRatings/ratings-${num}.tsv
    echo "Downloading ${file}"
    curl -L --compressed ${file} -o input/ratings/ratings-${num}.tsv
done

#download other files
file=${BASE_URL}/${DATE}/noteStatusHistory/noteStatusHistory-00000.tsv
echo "Downloading ${file}"
curl -L --compressed ${file} -o input/noteStatusHistory-00000.tsv

file=${BASE_URL}/${DATE}/userEnrollment/userEnrollment-00000.tsv
echo "Downloading ${file}"
curl -L --compressed ${file} -o input/userEnrollment-00000.tsv
