#How to use this tool#

1. Configure desired values. The priority order is following:
   * defaults in the code - do not use these for custom settings, they're used for some reasonable defaults
   * config files - two locations are searched (./config and ~/.install-tools) for a file named "config.toml" the one in homedir has higher priority
   * command line arguments - these have the highest priority

#Obtaining pull secrets#

1. Visit installer web page

    https://console.redhat.com/openshift/install/aws/installer-provisioned

    This one is for AWS, there are others as well but the pull secret should be the same.

2. 


app.ci
registry.ci.openshift.org
the authoritative, central CI registry
https://oauth-openshift.apps.ci.l2s4.p1.openshiftapps.com/oauth/authorize?client_id=openshift-browser-client&redirect_uri=https%3A%2F%2Foauth-openshift.apps.ci.l2s4.p1.openshiftapps.com%2Foauth%2Ftoken%2Fdisplay&response_type=code

arm01
registry.arm-build01.arm-build.devcluster.openshift.com
TOKEN: https://oauth-openshift.apps.arm-build01.arm-build.devcluster.openshift.com/oauth/authorize?client_id=openshift-browser-client&redirect_uri=https%3A%2F%2Foauth-openshift.apps.arm-build01.arm-build.devcluster.openshift.com%2Foauth%2Ftoken%2Fdisplay&response_type=code

build01
registry.build01.ci.openshift.org
https://oauth-openshift.apps.build01.ci.devcluster.openshift.com/oauth/authorize?client_id=openshift-browser-client&redirect_uri=https%3A%2F%2Foauth-openshift.apps.build01.ci.devcluster.openshift.com%2Foauth%2Ftoken%2Fdisplay&response_type=code

build02
registry.build02.ci.openshift.org
https://oauth-openshift.apps.build02.gcp.ci.openshift.org/oauth/authorize?client_id=console&redirect_uri=https%3A%2F%2Fconsole.build02.ci.openshift.org%2Fauth%2Fcallback&response_type=code&scope=user%3Afull&state=139b3493

build03
registry.build03.ci.openshift.org
https://oauth-openshift.apps.build03.ky4t.p1.openshiftapps.com/oauth/authorize?client_id=openshift-browser-client&redirect_uri=https%3A%2F%2Foauth-openshift.apps.build03.ky4t.p1.openshiftapps.com%2Foauth%2Ftoken%2Fdisplay&response_type=code

build04
registry.build04.ci.openshift.org
https://oauth-openshift.apps.build04.34d2.p2.openshiftapps.com/oauth/authorize?client_id=console&redirect_uri=https%3A%2F%2Fconsole-openshift-console.apps.build04.34d2.p2.openshiftapps.com%2Fauth%2Fcallback&response_type=code&scope=user%3Afull&state=84e04d6b

vsphere
registry.apps.build01-us-west-2.vmc.ci.openshift.org
https://oauth-openshift.apps.build01-us-west-2.vmc.ci.openshift.org/oauth/authorize?client_id=console&redirect_uri=https%3A%2F%2Fconsole-openshift-console.apps.build01-us-west-2.vmc.ci.openshift.org%2Fauth%2Fcallback&response_type=code&scope=user%3Afull&state=05f61d3d