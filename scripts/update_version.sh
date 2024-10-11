# Script to update prod version on configs
version=$(git describe --abbrev=0 --tags)
sed -i'.bak' -e "s/PROD_VERSION/$version/g" configs/conf_prod.go