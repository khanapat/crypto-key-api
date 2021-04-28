# /bin/sh

function printHelp() {
  echo "Usage: This script can remove, build and push image from crypto-key-api on utility project'"
  echo "       This script need 1 argument for processing including"
  echo "       1. Latest version of image"
  echo "                                                            "
  echo "    ./build.sh <version> "
  echo "                                                            "
  echo "Example :"
  echo "	./build.sh 1.0.7"
}

if [ -z "$1" ]
then
	printHelp
	exit 1
fi

version=$1

echo "build image $version"

docker build -t kcskbcnd93.kcs:5000/utility/crypto-key-api:"$version" .

echo "push image $version"

docker push kcskbcnd93.kcs:5000/utility/crypto-key-api:"$version"

echo "remove image in local"

docker rmi -f kcskbcnd93.kcs:5000/utility/crypto-key-api:"$version"

