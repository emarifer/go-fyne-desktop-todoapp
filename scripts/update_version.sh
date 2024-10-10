# Script to update prod version on configs
version=$(git rev-parse --short HEAD)
sed -i'.bak' -e "s/PROD_VERSION/$version/g" configs/conf_prod.go