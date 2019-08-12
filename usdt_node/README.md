# Test of https://medium.com/sthima-insights/we-just-got-a-new-super-power-runtime-usdt-comes-to-linux-814dc47e909f

## Install from package manager (***didn’t work***)

```
sudo add-apt-repository ppa:sthima/oss
```

Edit `/etc/apt/sources.list.d/sthima-ubuntu-oss-bionic.list`. To turn off gpg checks, add `[trusted=yes]` after `deb ` to get something like:

```
deb [trusted=yes] http://ppa.launchpad.net/sthima/oss/ubuntu bionic main
```

Then
```
sudo apt-get update
sudo apt-get install libstapsdt0 libstapsdt-dev
```

## Install from source

Clone https://github.com/sthima/libstapsdt and 

```
make
sudo make install
sudo ldconfig
```

## Final step


```
npm install usdt
npm index.js
```

