platforms = windows/amd64 darwin/amd64 linux/amd64
package = cmd/pre-commit/main.go
binary = build/pre-commit

help:
	@ echo "See Makefile for options"

compile:
	@ for platform in $(platforms); do 										  \
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
	