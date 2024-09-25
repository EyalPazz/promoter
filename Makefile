INSTALL_PATH = /usr/local/bin/promoter

install:
	@echo "Installing promoter binary..."
	./install.sh

uninstall:
	@echo "Uninstalling promoter binary..."
	rm -f $(INSTALL_PATH)

reinstall: uninstall install

.PHONY: uninstall install reinstall

	
