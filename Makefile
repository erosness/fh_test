version="0.0.1"
version_file=VERSION
working_dir=$(shell pwd)
arch="armhf"
remote_host = "fh@cube.local"

clean:
	-rm thingsplex_service_template

build-go:
	cd ./src;go build -o thingsplex_service_template service.go;cd ../

build-go-arm:
	cd ./src;GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o thingsplex_service_template service.go;cd ../

build-go-amd:
	cd ./src;GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o thingsplex_service_template service.go;cd ../


configure-arm:
	python ./scripts/config_env.py prod $(version) armhf

configure-amd64:
	python ./scripts/config_env.py prod $(version) amd64

package-tar:
	tar cvzf thingsplex_service_template_$(version).tar.gz thingsplex_service_template VERSION

clean-deb:
	find package/debian_tp -name ".DS_Store" -delete
	find package/debian_tp -name "delete_me" -delete
	find package/debian_fh -name ".DS_Store" -delete
	find package/debian_fh -name "delete_me" -delete

package-deb-doc-tp:clean-deb
	@echo "Packaging application using Thingsplex debian package layout"
	chmod a+x package/debian_tp/DEBIAN/*
	cp ./src/thingsplex_service_template package/debian_tp/opt/thingsplex/thingsplex_service_template
	cp VERSION package/debian_tp/opt/thingsplex/thingsplex_service_template
	docker run --rm -v ${working_dir}:/build -w /build --name debuild debian dpkg-deb --build package/debian_tp
	@echo "Done"

package-deb-doc-fh:clean-deb
	@echo "Packaging application using Futurehome debian package layout"
	chmod a+x package/debian_fh/DEBIAN/*
	cp ./src/thingsplex_service_template package/debian_fh/usr/bin/thingsplex_service_template
	cp VERSION package/debian_fh/var/lib/futurehome/thingsplex_service_template
	docker run --rm -v ${working_dir}:/build -w /build --name debuild debian dpkg-deb --build package/debian_fh
	@echo "Done"

deb-arm-fh : clean configure-arm build-go-arm package-deb-doc-fh
	@echo "Building Futurehome ARM package"
	mv package/debian_fh.deb package/build/thingsplex_service_template_$(version)_armhf.deb

deb-arm-tp : clean configure-arm build-go-arm package-deb-doc-tp
	@echo "Building Thingsplex ARM package"
	mv package/debian_tp.deb package/build/thingsplex_service_template_$(version)_armhf.deb

deb-amd : configure-amd64 build-go-amd package-deb-doc-tp
	@echo "Building Thingsplex AMD package"
	mv package/debian_tp.deb thingsplex_service_template_$(version)_amd64.deb

upload :
	@echo "Uploading the package to remote host"
	scp package/build/thingsplex_service_template_$(version)_armhf.deb $(remote_host):~/

remote-install : upload
	@echo "Uploading and installing the package on remote host"
	ssh -t $(remote_host) "sudo dpkg -i thingsplex_service_template_$(version)_armhf.deb"

deb-tp-remote-install : deb-arm-tp remote-install
	@echo "Package was built and installed on remote host"

deb-fh-remote-install : deb-arm-fh remote-install
	@echo "Package was built and installed on remote host"

run :
	cd ./src; go run service.go -c testdata/config.json;cd ../


.phony : clean
