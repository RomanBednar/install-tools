# How to use this tool

1. Install required dependencies
   * podman
   * oc

2. Configure installer tool

   1. There are multiple places where the installer tool looks for configuration:
      1. Defaults in the code - lowest priority
      2. Config files - two locations are searched (./config and ~/.install-tools) for a file named "conf.env", homedir has higher priority.
      3. Environment variables - higher priority than config files
      4. Command line arguments - highest priority
   2. Create a config file in one of the two locations mentioned above, e.g. ~/.install-tools/conf.env, example: 
   
   ```
   cloud=aws
   clusterName=mycluster-02
   userName=jdoe
   outputDir=./output
   cloudRegion=eu-central-1
   sshPublicKeyFile=${HOME}/.ssh/id_rsa.pub
   pullSecretFile=${HOME}/.docker/config.json
   imageTag=install-tools:latest
   engine=docker
   ```

4. Start the installation:

```
go run main.go --action create  --cloud aws --image registry.ci.openshift.org/ocp/release:4.17.0-0.ci-2024-07-25-020703 -o ~/openshift/clusters/aws/cluster-01 
```

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

4. Set `pullSecretFile` in your config file to point to the secrets file you created.

## Known issues & future work

* add cli tool to prompt user for required values interactively and save them to config (can be done by GUI instead)
* add scraper for image payloads so users do not have to copy/paste it manually: https://amd64.ocp.releases.ci.openshift.org/
* sometimes docker/podman adds `"quay.io":{}` into `config.json` which will break openshift-install if this lands in `pullSecret`
* vSphere installations are currently supported in CLI only due to being slightly more complex with preflight checks (VPN and password)
* each cloud has different regions, need a better way to handle this than hardcoding to templates (also ccoctl has region flag too)
