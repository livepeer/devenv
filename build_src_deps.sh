#!/usr/bin/env bash

set -ex

BASEDIR="$HOME/src"

##
# Note that we use shared libraries here to improve build times.
# Release builds should use static libraries throughout.
##

function build_nasm {
  cd "$BASEDIR"
  rm -rf "$BASEDIR/nasm"
  git clone -b nasm-2.14.02 https://repo.or.cz/nasm.git "$BASEDIR/nasm"
  cd "$BASEDIR/nasm"
  ./autogen.sh
  ./configure
  make
}

function build_x264 {
  cd "$BASEDIR"
  rm -rf "$BASEDIR/x264"
  git clone http://git.videolan.org/git/x264.git "$BASEDIR/x264"
  cd "$BASEDIR/x264"
  # git master as of this writing
  git checkout 545de2ffec6ae9a80738de1b2c8cf820249a2530
  ./configure --enable-pic --enable-shared
  make
}

function build_ffmpeg {
  cd "$BASEDIR"
  rm -rf "$BASEDIR/ffmpeg"
  git clone -b n4.1 https://git.ffmpeg.org/ffmpeg.git "$BASEDIR/ffmpeg"
  cd "$BASEDIR/ffmpeg"
  ./configure --enable-shared --disable-static --enable-gpl --enable-libx264 --enable-gnutls
  make
}

if [ -z "$1" ]; then
  build_nasm
  build_x264
  build_ffmpeg
else
  build_"$1"
fi
