#!/bin/bash
GO_VERSION="1.16.5"
ARCH=$(dpkg --print-architecture)
case $ARCH in
	armhf)
		GOLANG_ARCH_NAME="linux-armv6l"
		ARCHIVE_NAME=$(echo "go$GO_VERSION.$GOLANG_ARCH_NAME.tar.gz")
		ARCHIVE_URL=$(echo "https://golang.org/dl/$ARCHIVE_NAME")
		;;
	arm64)
		GOLANG_ARCH_NAME="linux-arm64"
		ARCHIVE_NAME=$(echo "go$GO_VERSION.$GOLANG_ARCH_NAME.tar.gz")
		ARCHIVE_URL=$(echo "https://golang.org/dl/$ARCHIVE_NAME")
		;;
	amd64)
		GOLANG_ARCH_NAME="linux-amd64"
		ARCHIVE_NAME=$(echo "go$GO_VERSION.$GOLANG_ARCH_NAME.tar.gz")
		ARCHIVE_URL=$(echo "https://golang.org/dl/$ARCHIVE_NAME")
		;;
	*)
		echo "unknown arch"
		#GOLANG_ARCH_NAME="windows-amd64.zip"
		;;
esac
echo "Downloading $ARCHIVE_URL"
wget $ARCHIVE_URL --progress=bar -O go.tar.gz
echo "Extracting To : /usr/local/go"
# sudo tar -C /usr/local -xzf $ARCHIVE_NAME
#sudo ln -s /usr/lib/go/bin/go /usr/local/go
#echo "PATH=$PATH:/usr/local/go/bin" | sudo tee -a /etc/profile
# https://askubuntu.com/questions/217570/bc-set-number-of-digits-after-decimal-point
export GO_TAR_KILOBYTES=$(printf "%.3f\n" $(echo "$(stat --format="%s" go.tar.gz) / 1000" | bc -l))
echo $GO_TAR_KILOBYTES
# GO_TAR_KILOBYTES=$(printf "%.3f\n" $(echo "$(stat --format="%s" go.tar.gz) / 1000" | bc -l));
# sudo tar --checkpoint=1 --checkpoint-action=exec='/bin/bash -c "cmd=$(echo R09fVEFSX0tJTE9CWVRFUz0kKHByaW50ZiAiJS4zZlxuIiAkKGVjaG8gIiQoc3RhdCAtLWZvcm1hdD0iJXMiIC9ob21lL21vcnBocy9nby50YXIuZ3opIC8gMTAwMCIgfCBiYyAtbCkpOw== | base64 -d ; echo); eval $cmd; echo [$TAR_CHECKPOINT] of $GO_TAR_KILOBYTES kilobytes"' -C /usr/local -xzf /home/$USERNAME/go.tar.gz
mkdir -p ~/go