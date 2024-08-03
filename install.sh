#!/usr/bin/env bash

set -ex

go build .

if [ ! -d ~/bin/.executables ] ; then
    mkdir -p ~/bin/.executables
fi

cp ./shenme ~/bin/.executables/shenme

cat <<EOF > ~/bin/shenme
#!/usr/bin/env bash
~/bin/.executables/shenme "\$@" 2>/dev/null
EOF

chmod +x ~/bin/shenme
