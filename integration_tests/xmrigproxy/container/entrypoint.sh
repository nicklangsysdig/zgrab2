#!/bin/sh


# This should do the same work as the init.d start script, but the 
# process should run in the foreground.
# In some cases, it may be appropriate to exit (and hence stop the 
# container) if the daemon process ends, but in others that may be
# expected behavior.

set -x
tar -xvzf xmrig-proxy-6.15.1-linux-static-x64.tar.gz
./xmrig-proxy-6.15.1/xmrig-proxy -o xmr-us-east1.nanopool.org:14444 --coin XMR
 
