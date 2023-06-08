# How to use this tool

1. Install dependencies
   * docker
   * oc

1. Configure installer tool

   1. There are multiple places where the installer tool looks for configuration:
      1. Defaults in the code - do not use these for custom settings, they're used for some reasonable defaults
      1. Config files - two locations are searched (./config and ~/.install-tools) for a file named "config.toml". The one in homedir has higher priority.
      1. Command line arguments - these have the highest priority
   1. Example config file: 
   
   ```
   clusterName = "mycluster-01"
   userName = "user-01"
   outputDir = "./output"
   cloudRegion = "eu-central-1"
   vmwarePassword = "xxxxxx"
   sshPublicKeyFile = "secrets/id_rsa.pub"
   pullSecretFile = "secrets/config.json"
   ```
   
1. Gather all secrets, this will copy ssh keys and pull secret under ./secrets
```
make get-secrets
```

1. Start the installation:

```
go run main.go --action create --cloud aws --image quay.io/openshift-release-dev/ocp-release:4.10.0-rc.2-x86_64 --outputdir /tmp/installdir
```

> You might need to enable go modules explicitly if your environment has it disabled by appending `GO111MODULE=on` to the command above.

# Obtaining pull secrets

1. Visit installer web page

    https://console.redhat.com/openshift/install/aws/installer-provisioned

    This one is for AWS, there are others as well but the pull secret should be the same. Get the pull secret and save it to a file.

2. Visit the following links, if prompted for authentication use GitHub, if not available use SSO. Some links don't redirect to the page where token is displayed, if so click your name in the top right corner and choose "Copy login command". Continue until you see "Display Token" button on a blank page.

      * app.ci (registry.ci.openshift.org)
      
      https://oauth-openshift.apps.ci.l2s4.p1.openshiftapps.com/oauth/authorize?client_id=openshift-browser-client&redirect_uri=https%3A%2F%2Foauth-openshift.apps.ci.l2s4.p1.openshiftapps.com%2Foauth%2Ftoken%2Fdisplay&response_type=code
   
      * arm01 (registry.arm-build01.arm-build.devcluster.openshift.com)
      
      https://oauth-openshift.apps.arm-build01.arm-build.devcluster.openshift.com/oauth/authorize?client_id=openshift-browser-client&redirect_uri=https%3A%2F%2Foauth-openshift.apps.arm-build01.arm-build.devcluster.openshift.com%2Foauth%2Ftoken%2Fdisplay&response_type=code

      * build01 (registry.build01.ci.openshift.org)
      
      https://oauth-openshift.apps.build01.ci.devcluster.openshift.com/oauth/authorize?client_id=openshift-browser-client&redirect_uri=https%3A%2F%2Foauth-openshift.apps.build01.ci.devcluster.openshift.com%2Foauth%2Ftoken%2Fdisplay&response_type=code

      * build02 (registry.build02.ci.openshift.org)
      
      https://oauth-openshift.apps.build02.gcp.ci.openshift.org/oauth/authorize?client_id=console&redirect_uri=https%3A%2F%2Fconsole.build02.ci.openshift.org%2Fauth%2Fcallback&response_type=code&scope=user%3Afull&state=139b3493

      * build03 (registry.build03.ci.openshift.org)
      
      https://oauth-openshift.apps.build03.ky4t.p1.openshiftapps.com/oauth/authorize?client_id=openshift-browser-client&redirect_uri=https%3A%2F%2Foauth-openshift.apps.build03.ky4t.p1.openshiftapps.com%2Foauth%2Ftoken%2Fdisplay&response_type=code

      * build04 (registry.build04.ci.openshift.org)
      
      https://oauth-openshift.apps.build04.34d2.p2.openshiftapps.com/oauth/authorize?client_id=console&redirect_uri=https%3A%2F%2Fconsole-openshift-console.apps.build04.34d2.p2.openshiftapps.com%2Fauth%2Fcallback&response_type=code&scope=user%3Afull&state=84e04d6b

      * vsphere (registry.apps.build01-us-west-2.vmc.ci.openshift.org)
      
      https://oauth-openshift.apps.build01-us-west-2.vmc.ci.openshift.org/oauth/authorize?client_id=console&redirect_uri=https%3A%2F%2Fconsole-openshift-console.apps.build01-us-west-2.vmc.ci.openshift.org%2Fauth%2Fcallback&response_type=code&scope=user%3Afull&state=05f61d3d
   
3. Run full 'oc login...' as displayed on each of the pages followed by 'oc 
   registry 
   login 
   --to=<your_secrets_file>'. Repeat the steps for all registries or for subset of those if you're absolutely sure which ones you need. Each 'oc registry' command will *append* to your secrets file from Step 1.



## Possible pitfalls when handling secrets

1. You can have two different secrets for quay.io, one from 

## TODO

1) correct destroy action, so it does not extract tools from payload - it's not needed
2) create config file interactively - ask for values and save them to config/config.env or config/config.toml
3) add scraper image payloads so users do not have to copy/paste it manually and just specify or search versions in CLI: https://amd64.ocp.releases.ci.openshift.org/
4) handle pull secret file better - parse it directly from docker conf.json to install template, so it does not have to be copied to temporary file just to get rid of spaces
5) for some reason docker can store `"quay.io":{}` in its config.json which will break openshift-install if this lands in `pullSecret` - handle this case
6) explore how to run everything in docker, dind seems to work fine when socket is mounted: docker run -it -v /var/run/docker.sock:/var/run/docker.sock docker:dind sh
7) 