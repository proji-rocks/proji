#!/usr/bin/env bash

package="./cmd/proji/"
package_split=(${package//\// })
package_name="${package_split[-1]}"
output_path="bin/"
platforms=("linux/amd64" "linux/386")

for platform in "${platforms[@]}"; do
    platform_split=(${platform//\// })
    GOOS="${platform_split[0]}"
    GOARCH="${platform_split[1]}"
    output_name="${output_path}${GOOS}/${GOARCH}/${package_name}"

    env CGO_ENABLED=1 GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags "-s -w" -o "$output_name" "$package"
    if [ $? -ne 0 ]; then
        echo "An error has occurred! Aborting the script execution..."
        exit 1
    fi
done
