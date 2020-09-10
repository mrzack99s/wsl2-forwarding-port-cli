# wsl2-forwarding-port-cli

>	Author: Chatdanai Phakaket <br>
>	Email: zchatdanai@gmail.com 

WSL2-forwarding-port-cli is command line tools for wsl2 forwarding port configure


## How to install

1. Open WSL2
2. Download the binary with the command 
```
    curl -LO https://github.com/mrzack99s/wsl2-forwarding-port-cli/releases/download/v1.1.2/wfp-cli
```
3. Make the kubectl binary executable.
```
    chmod +x wfp-cli
```
4. Move the binary in to PATH.
```
    sudo mv ./wfp-cli /usr/local/bin/wfp-cli
```

Let's enjoy !!!!

## How to use
# Use in WSL2
1. Open WSL2 with an administrator
2. Use the command
```
    wfp-cli <command>
```

# Use in Powershell or CMD
1. Open Powershell or CMD with an administrator
2. Use the command
```
    wsl wfp-cli <command>
```

## License

Copyright (c) 2020 - Chatdanai Phakaket

	

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)