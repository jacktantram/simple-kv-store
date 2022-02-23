set -e

echo "What's your name? (Default: $(whoami))"
printf ">"

read candidate

candidate=$(echo $candidate | sed "s/ /-/g")

if [ -z $candidate ]
then
	candidate=$(whoami)
fi

testname=$(basename $(git rev-parse --show-toplevel))

bundlename="$testname-$candidate.bundle"

if [ -f $bundlename ]; then rm $bundlename; fi
if [ -f $bundlename.gz ]; then rm $bundlename.gz; fi

git bundle create $bundlename --all 2> /dev/null
gzip $bundlename

printf "\n\033[0;32mCreated an archive containing this git repo: $bundlename.gz\n\nIf you're ready to submit please email this file to Kraken recruiting.\033[0m\n\n"
