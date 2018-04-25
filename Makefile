PLATFORMS := windows/amd64 darwin/amd64 linux/amd64
package = cmd/pre-commit/main.go
binary = build/pre-commit

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

compile: clean ## Build binaries for all platforms defined in $PLATFORMS
	@ for platform in $(PLATFORMS); do 										  \
		platform_split=($${platform//\// }); 								  \
		GOOS=$${platform_split[0]}; 										  \
		GOARCH=$${platform_split[1]}; 										  \
		output_name=$(binary)'_'$$GOOS'-'$$GOARCH;							  \
		if [ $$GOOS = "windows" ]; then 									  \
			output_name+='.exe'; 											  \
		fi;	 																  \
		echo "Build for $$platform -> $$output_name";						  \
		env GOOS=$$GOOS GOARCH=$$GOARCH go build -o $$output_name $(package); \
	done

clean: ## Cleanup build artifacts and remove build/
	rm -rf build/