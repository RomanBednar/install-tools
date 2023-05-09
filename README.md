# How to use this tool

1. Configure desired values. The priority order is following:
   * Defaults in the code - do not use these for custom settings, they're used for some reasonable defaults
   * Config files - two locations are searched (./config and ~/.install-tools) for a file named "config.toml". The one in homedir has higher priority.
   * Command line arguments - these have the highest priority


2. Install dependencies 
   * podman
   * oc

# Obtaining all pull secrets

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

4. Run it - a working example:
  
    ```$go run main.go --action create --cloud aws --image quay.io/openshift-release-dev/ocp-release:4.10.0-rc.2-x86_64 --outputdir /tmp/installdir```


## TODO

1) correct destroy action so it does not extract tools from payload - it's not needed